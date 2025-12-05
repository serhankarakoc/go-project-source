package middlewares

import (
	"zatrano/configs/sessionconfig"

	"github.com/gofiber/fiber/v2"
)

// SessionMiddleware her istek için session başlatır
// ve c.Locals içine session objesini ekler.
func SessionMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {

		sess, err := sessionconfig.SessionStart(c)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Session başlatılamadı")
		}

		// Locals içine session objesi + prefix koy
		c.Locals("session", sess)

		return c.Next()
	}
}
