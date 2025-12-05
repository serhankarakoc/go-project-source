package middlewares

import (
	"zatrano/configs/sessionconfig"
	"zatrano/pkg/currentuser"
	"zatrano/pkg/flashmessages"
	"zatrano/services"

	"github.com/gofiber/fiber/v2"
)

// AuthUser — Locals içinde tutulacak kullanıcı bilgisi
type AuthUser struct {
	ID            uint
	Email         string
	UserTypeID    uint
	IsActive      bool
	EmailVerified bool
}

// AuthMiddleware — Sırayla tüm authentication kontrollerini yapar:
// 1. Oturum kontrolü
// 2. Email doğrulanmış mı?
// 3. Kullanıcı aktif mi?
func AuthMiddleware(c *fiber.Ctx) error {
	// 1. OTURUM KONTROLÜ
	userID, err := sessionconfig.GetUserIDFromSession(c)
	if err != nil || userID == 0 {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Oturum süresi dolmuş veya geçersiz.")
		return c.Redirect("/auth/login", fiber.StatusSeeOther)
	}

	// Kullanıcı bilgilerini getir
	authService := services.NewAuthService()
	user, err := authService.GetUserProfile(userID)
	if err != nil {
		// Oturumu temizle
		_ = sessionconfig.DestroySession(c)
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Kullanıcı bulunamadı, lütfen tekrar giriş yapın.")
		return c.Redirect("/auth/login", fiber.StatusSeeOther)
	}

	// 2. EMAIL DOĞRULAMA KONTROLÜ
	if !user.EmailVerified {
		// Oturumu temizle
		_ = sessionconfig.DestroySession(c)

		// Kullanıcıyı email doğrulama sayfasına yönlendir
		sess, _ := sessionconfig.SessionStart(c)
		sess.Set("pending_verification", true)
		sess.Set("user_email", user.Email)
		_ = sess.Save()

		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey,
			"Lütfen e-posta adresinizi doğrulayınız. Doğrulama linki e-postanıza gönderilmiştir.")

		return c.Redirect("/auth/login", fiber.StatusSeeOther)
	}

	// 3. AKTİF Mİ KONTROLÜ
	if !user.IsActive {
		// Oturumu temizle
		_ = sessionconfig.DestroySession(c)

		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey,
			"Hesabınız pasif durumda. Lütfen yöneticinizle iletişime geçin.")

		return c.Redirect("/auth/login", fiber.StatusSeeOther)
	}

	// Tüm kontroller başarılı, kullanıcı bilgilerini kaydet
	authUser := AuthUser{
		ID:            user.ID,
		Email:         user.Email,
		UserTypeID:    user.UserTypeID,
		IsActive:      user.IsActive,
		EmailVerified: user.EmailVerified,
	}

	// Locals'a kaydet
	c.Locals("authUser", authUser)

	// Context'e kaydet
	ctx := currentuser.SetToContext(c.UserContext(), currentuser.CurrentUser{
		ID:         user.ID,
		Email:      user.Email,
		UserTypeID: user.UserTypeID,
	})
	c.SetUserContext(ctx)

	return c.Next()
}
