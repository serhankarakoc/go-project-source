package middlewares

import (
	"zatrano/configs/sessionconfig"
	"zatrano/pkg/currentuser"
	"zatrano/services"

	"github.com/gofiber/fiber/v2"
)

// GuestMiddleware — oturum açıkken giriş sayfasına erişimi engeller.
func GuestMiddleware(c *fiber.Ctx) error {
	// Session kontrolü
	userID, err := sessionconfig.GetUserIDFromSession(c)
	if err != nil || userID == 0 {
		return c.Next() // Oturum yoksa devam et
	}

	// Kullanıcı doğrulaması
	authService := services.NewAuthService()
	user, err := authService.GetUserProfile(userID)
	if err != nil {
		_ = sessionconfig.DestroySession(c)
		return c.Next()
	}

	// Fiber locals ve context güncellemesi
	c.Locals("userID", userID)
	c.Locals("userEmail", user.Email)
	c.Locals("userTypeID", user.UserTypeID)

	ctx := currentuser.SetToContext(c.Context(), currentuser.CurrentUser{
		ID:         userID,
		Email:      user.Email,
		UserTypeID: user.UserTypeID,
	})
	c.SetUserContext(ctx)

	// Kullanıcı türüne göre yönlendir
	if user.UserTypeID == 1 { // Admin
		return c.Redirect("/dashboard/home", fiber.StatusSeeOther)
	}
	return c.Redirect("/panel/anasayfa", fiber.StatusSeeOther)
}
