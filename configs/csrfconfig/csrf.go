package csrfconfig

import (
	"strings"
	"time"

	"zatrano/configs/envconfig"
	"zatrano/configs/logconfig"
	"zatrano/pkg/flashmessages"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/utils"
	"go.uber.org/zap"
)

func isProd() bool { return envconfig.IsProd() }

// Gerekirse muaf yollar
var csrfExemptPaths = []string{"/healthz", "/readyz"}

// SameSite değerini env'den oku: CSRF_COOKIE_SAMESITE (Strict|Lax|None)
// Prod varsayılan: Strict, Dev varsayılan: Lax
func csrfSameSite() string {
	def := "Strict"
	if !isProd() {
		def = "Lax"
	}
	v := strings.ToLower(strings.TrimSpace(envconfig.String("CSRF_COOKIE_SAMESITE", "")))
	switch v {
	case "strict":
		return "Strict"
	case "lax":
		return "Lax"
	case "none":
		return "None" // DİKKAT: None => Secure zorunlu
	default:
		return def
	}
}

func SetupCSRF() fiber.Handler {
	sameSite := csrfSameSite()

	// None seçildiyse tarayıcılar cookie'yi yalnız Secure=true iken kabul eder
	secure := isProd()
	if sameSite == "None" {
		secure = true
	}

	cookieDomain := ""
	if isProd() {
		cookieDomain = envconfig.String("COOKIE_DOMAIN", "") // örn: .zatrano (lokalde boş bırak)
	}

	cfg := csrf.Config{
		KeyLookup:      "header:X-CSRF-Token",
		CookieName:     "csrf_token",
		CookieHTTPOnly: true,
		CookieSecure:   secure,
		CookieSameSite: sameSite,
		CookieDomain:   cookieDomain,
		Expiration:     1 * time.Hour,
		KeyGenerator:   utils.UUID,
		ContextKey:     "csrf",

		ErrorHandler: func(c *fiber.Ctx, err error) error {
			logconfig.Log.Warn("CSRF validation failed",
				zap.Error(err),
				zap.String("ip", c.IP()),
				zap.String("path", c.Path()),
				zap.String("method", c.Method()),
			)
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey,
				"Güvenlik doğrulaması başarısız oldu. Lütfen sayfayı yenileyip tekrar deneyin.")
			return c.Redirect("/auth/login", fiber.StatusSeeOther)
		},

		// GET/HEAD/OPTIONS’u SKIP etmiyoruz ki token üretilebilsin.
		// Yalnız muaf yolları atla; formdan gelen token’ı header’a köprüle.
		Next: func(c *fiber.Ctx) bool {
			if c.Get("X-CSRF-Token") == "" {
				if t := c.FormValue("csrf_token"); t != "" {
					c.Request().Header.Set("X-CSRF-Token", t)
				}
			}
			path := c.Path()
			for _, p := range csrfExemptPaths {
				if strings.HasPrefix(path, p) {
					return true
				}
			}
			return false
		},
	}

	logconfig.SLog.Infow("CSRF middleware yapılandırıldı",
		"exempt_paths", csrfExemptPaths,
		"secure", cfg.CookieSecure,
		"samesite", cfg.CookieSameSite,
		"domain", cfg.CookieDomain,
	)
	return csrf.New(cfg)
}
