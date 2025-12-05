package sessionconfig

import (
	"encoding/gob"
	"strings"
	"time"

	"zatrano/configs/envconfig"
	"zatrano/configs/logconfig"
	"zatrano/configs/redisconfig"
	"zatrano/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis/v3"
)

var Store *session.Store

// Basit yardımcı fonksiyonlar
func isProd() bool { return envconfig.IsProd() }

// SameSite modunu .env’den oku
func sessionSameSite() string {
	def := "Strict"
	if !isProd() {
		def = "Lax" // dev ortamında en uyumlu
	}
	v := strings.ToLower(strings.TrimSpace(envconfig.String("SESSION_COOKIE_SAMESITE", "")))
	switch v {
	case "strict":
		return "Strict"
	case "lax":
		return "Lax"
	case "none":
		return "None"
	default:
		return def
	}
}

// InitSession
func InitSession() {
	Store = createSessionStore()
	registerGobTypes()
	logconfig.SLog.Info("Oturum (session) sistemi başlatıldı.")
}

func SetupSession() *session.Store {
	if Store == nil {
		logconfig.SLog.Warn("Session store henüz başlatılmamış, şimdi başlatılıyor.")
		InitSession()
	}
	return Store
}

// Asıl session store oluşturma
func createSessionStore() *session.Store {
	cookieName := envconfig.String("SESSION_COOKIE_NAME", "session_id")
	expirationHours := envconfig.Int("SESSION_TTL_HOURS", 24)
	sameSite := sessionSameSite()

	secure := isProd()
	cookieDomain := ""

	// Development ortamında secure cookie kullanma
	if !isProd() {
		secure = false
	}

	// Production ortamında SameSite=None ise Secure zorunlu
	if isProd() && sameSite == "None" {
		secure = true
		cookieDomain = envconfig.String("COOKIE_DOMAIN", "")
	}

	// Redis Store (fiber/storage/redis v3)
	redisStore := redis.NewFromConnection(redisconfig.GetClient())

	store := session.New(session.Config{
		KeyLookup:      "cookie:" + cookieName,
		CookieHTTPOnly: true,
		CookieSecure:   secure,
		CookieSameSite: sameSite,
		CookieDomain:   cookieDomain,
		Expiration:     time.Duration(expirationHours) * time.Hour,
		Storage:        redisStore,
	})

	logconfig.SLog.Infow("Session store Redis ile yapılandırıldı",
		"cookie_name", cookieName,
		"cookie_http_only", true,
		"cookie_secure", secure,
		"same_site", sameSite,
		"domain", cookieDomain,
		"expiration_hours", expirationHours,
	)

	return store
}

// Gob kayıtları
func registerGobTypes() {
	gob.Register(&models.User{})
	logconfig.SLog.Debug("Session için gob türleri kaydedildi: *models.User")
}

// Session başlatma
func SessionStart(c *fiber.Ctx) (*session.Session, error) {
	if Store == nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "session store not initialized")
	}
	return Store.Get(c)
}

// Session yok etme
func DestroySession(c *fiber.Ctx) error {
	sess, err := SessionStart(c)
	if err != nil {
		return err
	}
	return sess.Destroy()
}

// Basit yardımcılar
func GetUserID(c *fiber.Ctx) (uint, error) {
	return GetUserIDFromSession(c)
}

func SetValue(c *fiber.Ctx, key string, value interface{}) error {
	return SetSessionValue(c, key, value)
}

// Orijinal prefix destekli fonksiyonlar
func GetUserIDFromSession(c *fiber.Ctx) (uint, error) {
	sess, err := SessionStart(c)
	if err != nil {
		return 0, err
	}
	key := "user_id"
	switch v := sess.Get(key).(type) {
	case uint:
		return v, nil
	case int:
		return uint(v), nil
	case int64:
		return uint(v), nil
	case float64:
		if v < 0 {
			return 0, fiber.ErrUnauthorized
		}
		return uint(v), nil
	default:
		return 0, fiber.ErrUnauthorized
	}
}

func SetSessionValue(c *fiber.Ctx, key string, value interface{}) error {
	sess, err := SessionStart(c)
	if err != nil {
		return err
	}
	prefixedKey := key
	sess.Set(prefixedKey, value)
	return sess.Save()
}
