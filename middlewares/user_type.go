package middlewares

import (
	"zatrano/configs/sessionconfig"
	"zatrano/pkg/flashmessages"

	"github.com/gofiber/fiber/v2"
)

// UserTypeMiddleware — Kullanıcının belirli UserTypeID'lere sahip olduğunu kontrol eder.
func UserTypeMiddleware(allowedTypes ...uint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		val := c.Locals("authUser")
		if val == nil {
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Oturum bulunamadı")
			return c.Redirect("/auth/login", fiber.StatusSeeOther)
		}

		user, ok := val.(AuthUser)
		if !ok {
			_ = sessionconfig.DestroySession(c)
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Oturum bilgileri geçersiz")
			return c.Redirect("/auth/login", fiber.StatusSeeOther)
		}

		// İzin verilen tipler arasında mı?
		allowed := false
		for _, t := range allowedTypes {
			if user.UserTypeID == t {
				allowed = true
				break
			}
		}

		if !allowed {
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey,
				"Bu sayfaya erişim yetkiniz bulunmamaktadır.")

			// Kullanıcı tipine göre ana sayfaya yönlendir
			switch user.UserTypeID {
			case 1:
				return c.Redirect("/dashboard/home", fiber.StatusSeeOther)
			case 2:
				return c.Redirect("/panel/anasayfa", fiber.StatusSeeOther)
			default:
				// Tanımlı değilse oturumu sonlandır
				_ = sessionconfig.DestroySession(c)
				return c.Redirect("/auth/login", fiber.StatusSeeOther)
			}
		}

		return c.Next()
	}
}
