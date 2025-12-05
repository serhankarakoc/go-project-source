package middlewares

import (
	"slices"
	"time"

	"zatrano/configs/envconfig"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// Ortak whitelist IP’leri (örn: local, internal)
var whitelistIPs = []string{
	"127.0.0.1",
	"::1",
}

// Ortak kontrol: Whitelist + Dev ortamında limit devre dışı
func shouldSkipLimit(c *fiber.Ctx) bool {
	if slices.Contains(whitelistIPs, c.IP()) {
		return true
	}
	if !envconfig.IsProd() { // Development veya staging ortamında limiti devre dışı bırak
		return true
	}
	return false
}

// 🌍 Global limiter — tüm uygulama için IP başına genel limit
func GlobalRateLimit() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        envconfig.Int("GLOBAL_RATE_MAX", 300), // 1 dakikada 300 istek (IP başına)
		Expiration: time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return "global:" + c.IP()
		},
		Next: func(c *fiber.Ctx) bool {
			if c.Path() == "/healthz" || c.Path() == "/readyz" {
				return true // health endpointleri hariç
			}
			return shouldSkipLimit(c)
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).
				SendString("Çok fazla istek gönderildi. Lütfen kısa bir süre sonra tekrar deneyin.")
		},
	})
}

// 🧾 Form POST limiter — örn: iletişim, kayıt, form gönderimleri
func FormPostRateLimit() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        envconfig.Int("FORM_POST_RATE_MAX", 30),
		Expiration: time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return "form:" + c.IP() + ":" + c.Path()
		},
		Next: func(c *fiber.Ctx) bool {
			if c.Method() != fiber.MethodPost {
				return true // sadece POST istekleri için uygula
			}
			return shouldSkipLimit(c)
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).
				SendString("Çok fazla form isteği gönderildi. Lütfen biraz bekleyin.")
		},
	})
}

// 🔐 Login özel limiter — brute force engelleme
func LoginRateLimit() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        envconfig.Int("LOGIN_RATE_MAX", 5),
		Expiration: time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return "login:" + c.IP()
		},
		Next: func(c *fiber.Ctx) bool {
			if !(c.Method() == fiber.MethodPost && c.Path() == "/auth/login") {
				return true
			}
			return shouldSkipLimit(c)
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).
				SendString("Çok fazla giriş denemesi yaptınız. Lütfen 1 dakika sonra tekrar deneyin.")
		},
	})
}
