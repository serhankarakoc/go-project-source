package middlewares

import (
	"zatrano/configs/logconfig"

	"github.com/gofiber/fiber/v2"
)

// ZapLogger — sade, seviyelere göre düzenlenmiş logger middleware
func ZapLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		// IP ve path
		ip := c.IP()
		path := c.Path()

		// Sadece logla, caller ve stack trace yok
		logconfig.SLog.Desugar().Sugar().Infow("request",
			"ip", ip,
			"path", path,
		)

		return err
	}
}
