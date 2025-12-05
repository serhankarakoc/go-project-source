package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"os"

	"zatrano/configs/logconfig"
	"zatrano/configs/sessionconfig"
	"zatrano/models"
	"zatrano/pkg/flashmessages"
	"zatrano/pkg/renderer"
	"zatrano/requests"
	"zatrano/services"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type AuthHandler struct {
	service    services.IAuthService
	mailSender services.IMailService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		service:    services.NewAuthService(),
		mailSender: services.NewMailService(),
	}
}

func (h *AuthHandler) handleError(c *fiber.Ctx, err error, userID uint, email string, action string) error {
	var (
		errMsg         string
		flashKey       = flashmessages.FlashErrorKey
		redirectTarget = "/auth/login"
		logoutUser     = false
	)

	errorMessages := map[error]struct {
		message       string
		redirect      string
		shouldLogout  bool
		logAdditional []zap.Field
	}{
		services.ErrInvalidCredentials:       {message: "Kullanıcı adı veya şifre hatalı."},
		services.ErrUserInactive:             {message: "Hesabınız aktif değil. Lütfen yöneticinizle iletişime geçin."},
		services.ErrUserNotFound:             {message: "Kullanıcı bulunamadı, lütfen tekrar giriş yapın.", shouldLogout: true, logAdditional: []zap.Field{zap.Uint("user_id", userID)}},
		services.ErrCurrentPasswordIncorrect: {message: "Mevcut şifreniz hatalı.", redirect: "/auth/profile"},
		services.ErrPasswordTooShort:         {message: "Şifre çok kısa.", redirect: "/auth/profile"},
		services.ErrPasswordSameAsOld:        {message: "Yeni şifre eski şifre ile aynı olamaz.", redirect: "/auth/profile"},
	}

	if customErr, ok := errorMessages[err]; ok {
		errMsg = customErr.message
		if customErr.redirect != "" {
			redirectTarget = customErr.redirect
		}
		logoutUser = customErr.shouldLogout
	} else {
		errMsg = "İşlem sırasında bir sorun oluştu. Lütfen tekrar deneyin."
		logconfig.Log.Error(action+": Beklenmeyen hata",
			zap.Uint("user_id", userID),
			zap.String("email", email),
			zap.Error(err))
	}

	if logoutUser {
		h.destroySession(c)
	}

	_ = flashmessages.SetFlashMessage(c, flashKey, errMsg)
	return c.Redirect(redirectTarget, fiber.StatusSeeOther)
}

func (h *AuthHandler) getSessionUser(c *fiber.Ctx) (uint, error) {
	if userID, ok := c.Locals("userID").(uint); ok {
		return userID, nil
	}
	sess, err := sessionconfig.SessionStart(c)
	if err != nil {
		return 0, err
	}
	switch v := sess.Get("user_id").(type) {
	case uint:
		return v, nil
	case int:
		return uint(v), nil
	case float64:
		return uint(v), nil
	default:
		return 0, fiber.ErrUnauthorized
	}
}

func (h *AuthHandler) destroySession(c *fiber.Ctx) {
	sess, err := sessionconfig.SessionStart(c)
	if err != nil {
		logconfig.Log.Warn("Oturum yok edilemedi (zaten yok olabilir)", zap.Error(err))
		return
	}
	_ = sess.Destroy()
}

func (h *AuthHandler) ShowLogin(c *fiber.Ctx) error {
	sess, err := sessionconfig.SessionStart(c)

	var pendingVerification bool
	var userEmail string

	if err == nil {
		// 1. Önce email_not_verified'i kontrol et (middleware'den geliyor)
		if notVerified := sess.Get("email_not_verified"); notVerified != nil {
			if b, ok := notVerified.(bool); ok && b {
				pendingVerification = true
				if email := sess.Get("user_email"); email != nil {
					userEmail = email.(string)
				}
			}
			// Temizle
			sess.Delete("email_not_verified")
			sess.Delete("user_email")
		}

		// 2. Sonra pending_verification'i kontrol et (eski kod ve yeni middleware)
		if v := sess.Get("pending_verification"); v != nil {
			if b, ok := v.(bool); ok && b {
				pendingVerification = true
			}
			sess.Delete("pending_verification")
		}

		// 3. user_email'i al (middleware'den geliyor)
		if email := sess.Get("user_email"); email != nil {
			if e, ok := email.(string); ok {
				userEmail = e
			}
			sess.Delete("user_email")
		}

		_ = sess.Save()
	}

	return renderer.Render(c, "auth/login", "layouts/auth", fiber.Map{
		"Title":               "Giriş Yap",
		"PendingVerification": pendingVerification,
		"UserEmail":           userEmail,
	}, http.StatusOK)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	req, ok := c.Locals("loginRequest").(requests.LoginRequest)
	if !ok {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz istek formatı")
		return c.Redirect("/auth/login", fiber.StatusSeeOther)
	}

	user, err := h.service.Authenticate(req.Email, req.Password)
	if err != nil {
		return h.handleError(c, err, 0, req.Email, "Login")
	}

	// EMAİL DOĞRULAMA KONTROLÜ EKLE!
	if !user.EmailVerified {
		// Email doğrulanmamışsa direkt login sayfasına mesajla yönlendir
		sess, _ := sessionconfig.SessionStart(c)
		sess.Set("pending_verification", true)
		sess.Set("user_email", user.Email)
		_ = sess.Save()

		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey,
			"Lütfen e-posta adresinizi doğrulayınız. Doğrulama linki e-postanıza gönderilmiştir.")
		return c.Redirect("/auth/login", fiber.StatusSeeOther)
	}

	// Normal oturum açma işlemleri...
	sess, err := sessionconfig.SessionStart(c)
	if err != nil {
		logconfig.Log.Error("Oturum başlatılamadı",
			zap.Uint("user_id", user.ID),
			zap.String("email", user.Email),
			zap.Error(err))
		return h.handleError(c, fiber.ErrInternalServerError, user.ID, user.Email, "Login")
	}

	sess.Set("user_id", user.ID)
	sess.Set("user_type_id", user.UserTypeID)
	sess.Set("is_active", user.IsActive)
	if err := sess.Save(); err != nil {
		logconfig.Log.Error("Oturum kaydedilemedi",
			zap.Uint("user_id", user.ID),
			zap.String("email", user.Email),
			zap.Error(err))
		return h.handleError(c, fiber.ErrInternalServerError, user.ID, user.Email, "Login")
	}

	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Başarıyla giriş yapıldı")
	if user.UserTypeID == 1 {
		return c.Redirect("/dashboard/home", fiber.StatusFound)
	}
	return c.Redirect("/panel/anasayfa", fiber.StatusFound)
}

func (h *AuthHandler) Profile(c *fiber.Ctx) error {
	userID, err := h.getSessionUser(c)
	if err != nil {
		logconfig.Log.Warn("Profil: Geçersiz oturum", zap.Error(err))
		h.destroySession(c)
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz oturum, lütfen tekrar giriş yapın.")
		return c.Redirect("/auth/login", fiber.StatusSeeOther)
	}

	user, err := h.service.GetUserProfile(userID)
	if err != nil {
		return h.handleError(c, err, userID, "", "Profil")
	}

	return renderer.Render(c, "auth/profile", "layouts/auth", fiber.Map{
		"Title": "Profilim",
		"User":  user,
	}, http.StatusOK)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	h.destroySession(c)
	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Başarıyla çıkış yapıldı.")
	return c.Redirect("/auth/login", fiber.StatusFound)
}

func (h *AuthHandler) UpdatePassword(c *fiber.Ctx) error {
	userID, err := h.getSessionUser(c)
	if err != nil {
		h.destroySession(c)
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz oturum bilgisi, lütfen tekrar giriş yapın.")
		return c.Redirect("/auth/login", fiber.StatusSeeOther)
	}

	req, ok := c.Locals("updatePasswordRequest").(requests.UpdatePasswordRequest)
	if !ok {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz istek formatı.")
		return c.Redirect("/auth/profile", fiber.StatusSeeOther)
	}

	if err := h.service.UpdatePassword(c.UserContext(), userID, req.CurrentPassword, req.NewPassword); err != nil {
		return h.handleError(c, err, userID, "", "Parola Güncelleme")
	}

	h.destroySession(c)
	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Şifre başarıyla güncellendi. Lütfen yeni şifrenizle tekrar giriş yapın.")
	return c.Redirect("/auth/login", fiber.StatusFound)
}

func (h *AuthHandler) ShowRegister(c *fiber.Ctx) error {
	return renderer.Render(c, "auth/register", "layouts/auth", fiber.Map{
		"Title": "Kayıt Ol",
	}, http.StatusOK)
}

func generateToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	req, ok := c.Locals("registerRequest").(requests.RegisterRequest)
	if !ok {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz kayıt isteği")
		return c.Redirect("/auth/register", fiber.StatusSeeOther)
	}

	user := &models.User{
		Name:              req.Name,
		Email:             req.Email,
		Password:          req.Password,
		UserTypeID:        ptrUint(2),
		ResetToken:        "",
		EmailVerified:     false,
		VerificationToken: "",
		Provider:          "",
		ProviderID:        "",
	}

	if token, err := generateToken(); err == nil {
		user.ResetToken = token
	}
	if vtoken, err := generateToken(); err == nil {
		user.VerificationToken = vtoken
	}

	if err := h.service.CreateUser(c.UserContext(), user); err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Kullanıcı oluşturulamadı. Lütfen tekrar deneyin.")
		return c.Redirect("/auth/register", fiber.StatusSeeOther)
	}

	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Kayıt işlemi tamamlandı. Lütfen email adresinizi doğrulayın.")
	baseURL := os.Getenv("APP_BASE_URL")
	verificationLink := baseURL + "/auth/verify-email?token=" + user.VerificationToken
	_ = h.mailSender.SendMail(user.Email, "Email Doğrulama", "Lütfen doğrulamak için tıklayın: "+verificationLink)

	return renderer.Render(c, "auth/verify_email_notice", "layouts/auth", fiber.Map{
		"Title": "Email Doğrulama",
	}, http.StatusOK)
}

func ptrUint(v uint) uint {
	return v
}

func (h *AuthHandler) ShowForgotPassword(c *fiber.Ctx) error {
	return renderer.Render(c, "auth/forgot_password", "layouts/auth", fiber.Map{
		"Title": "Şifremi Unuttum",
	}, http.StatusOK)
}

func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	req, ok := c.Locals("forgotPasswordRequest").(requests.ForgotPasswordRequest)
	if !ok {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz istek")
		return c.Redirect("/auth/forgot-password", fiber.StatusSeeOther)
	}
	if err := h.service.SendPasswordResetLink(req.Email); err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Şifre sıfırlama bağlantısı gönderilemedi.")
		return c.Redirect("/auth/forgot-password", fiber.StatusSeeOther)
	}
	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Şifre sıfırlama bağlantısı gönderildi.")
	return c.Redirect("/auth/login", fiber.StatusSeeOther)
}

func (h *AuthHandler) ShowResetPassword(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz veya eksik token.")
		return c.Redirect("/auth/forgot-password", fiber.StatusSeeOther)
	}
	return renderer.Render(c, "auth/reset_password", "layouts/auth", fiber.Map{
		"Title": "Şifre Sıfırla",
		"Token": token,
	}, http.StatusOK)
}

func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	req, ok := c.Locals("resetPasswordRequest").(requests.ResetPasswordRequest)
	if !ok || req.Token == "" {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz veya eksik token.")
		return c.Redirect("/auth/forgot-password", fiber.StatusSeeOther)
	}
	if err := h.service.ResetPassword(req.Token, req.NewPassword); err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Şifre sıfırlama başarısız.")
		return c.Redirect("/auth/reset-password", fiber.StatusSeeOther)
	}
	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Şifre sıfırlandı. Lütfen giriş yapın.")
	return c.Redirect("/auth/login", fiber.StatusSeeOther)
}

func (h *AuthHandler) VerifyEmail(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Doğrulama tokeni eksik veya geçersiz.")
		return c.Redirect("/auth/forgot-password", fiber.StatusSeeOther)
	}
	if err := h.service.VerifyEmail(token); err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Email doğrulama başarısız.")
		return c.Redirect("/auth/forgot-password", fiber.StatusSeeOther)
	}
	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Email başarıyla doğrulandı.")
	return c.Redirect("/auth/login", fiber.StatusSeeOther)
}

func (h *AuthHandler) ShowResendVerification(c *fiber.Ctx) error {
	// URL'den email parametresini al
	email := c.Query("email")

	return renderer.Render(c, "auth/resend_verification", "layouts/auth", fiber.Map{
		"Title": "Email Doğrulama Linkini Yeniden Gönder",
		"Email": email, // Email parametresini template'e gönder
	}, http.StatusOK)
}

func (h *AuthHandler) ResendVerification(c *fiber.Ctx) error {
	req, ok := c.Locals("resendVerificationRequest").(requests.ResendVerificationRequest)
	if !ok {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz istek")
		return c.Redirect("/auth/resend-verification", fiber.StatusSeeOther)
	}
	if err := h.service.ResendVerificationLink(req.Email); err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Doğrulama linki gönderilemedi.")
		return c.Redirect("/auth/resend-verification", fiber.StatusSeeOther)
	}
	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Doğrulama linki e-posta adresinize gönderildi.")
	return c.Redirect("/auth/login", fiber.StatusSeeOther)
}

func (h *AuthHandler) UpdateInfo(c *fiber.Ctx) error {
	userID, err := h.getSessionUser(c)
	if err != nil {
		h.destroySession(c)
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz oturum bilgisi, lütfen tekrar giriş yapın.")
		return c.Redirect("/auth/login", fiber.StatusSeeOther)
	}

	req, ok := c.Locals("updateInfoRequest").(requests.UpdateInfoRequest)
	if !ok {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz istek formatı.")
		return c.Redirect("/auth/profile", fiber.StatusSeeOther)
	}

	if err := h.service.UpdateUserInfo(c.UserContext(), userID, req.Name, req.Email); err != nil {
		return h.handleError(c, err, userID, req.Email, "Profil Bilgileri Güncelleme")
	}

	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Profil bilgileri güncellendi.")
	return c.Redirect("/auth/profile", fiber.StatusSeeOther)
}
