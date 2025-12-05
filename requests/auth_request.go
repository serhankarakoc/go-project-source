package requests

import (
	"zatrano/pkg/flashmessages"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type (
	LoginRequest struct {
		Email    string `form:"email" validate:"required,email"`
		Password string `form:"password" validate:"required,min=6"`
	}

	UpdatePasswordRequest struct {
		CurrentPassword string `form:"current_password" validate:"required,min=6"`
		NewPassword     string `form:"new_password" validate:"required,min=8,nefield=CurrentPassword"`
		ConfirmPassword string `form:"confirm_password" validate:"required,eqfield=NewPassword"`
	}

	RegisterRequest struct {
		Name            string `form:"name" validate:"required,min=3"`
		Email           string `form:"email" validate:"required,email"`
		Password        string `form:"password" validate:"required,min=8"`
		ConfirmPassword string `form:"confirm_password" validate:"required,eqfield=Password"`
	}

	ForgotPasswordRequest struct {
		Email string `form:"email" validate:"required,email"`
	}

	ResetPasswordRequest struct {
		Token           string `form:"token" validate:"required"`
		NewPassword     string `form:"new_password" validate:"required,min=8"`
		ConfirmPassword string `form:"confirm_password" validate:"required,eqfield=NewPassword"`
	}

	ResendVerificationRequest struct {
		Email string `form:"email" validate:"required,email"`
	}

	UpdateInfoRequest struct {
		Name  string `form:"name" validate:"required,min=3"`
		Email string `form:"email" validate:"required,email"`
	}

	VerifyEmailRequest struct {
		Token string `form:"token" validate:"required"`
	}
)

func validateRequest(c *fiber.Ctx, req interface{}, errorMessages map[string]string, redirectPath string) error {
	if err := c.BodyParser(req); err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz istek formatı")
		return c.Redirect(redirectPath, fiber.StatusSeeOther)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		err := err.(validator.ValidationErrors)[0]
		if msg, ok := errorMessages[err.Field()+"_"+err.Tag()]; ok {
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, msg)
		} else {
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz giriş bilgileri")
		}
		return c.Redirect(redirectPath, fiber.StatusSeeOther)
	}

	return nil
}

func ValidateLoginRequest(c *fiber.Ctx) error {
	var req LoginRequest
	errorMessages := map[string]string{
		"Email_required":    "E-posta adresi zorunludur",
		"Email_email":       "Geçerli bir e-posta adresi giriniz",
		"Password_required": "Şifre zorunludur",
		"Password_min":      "Şifre en az 6 karakter olmalıdır",
	}

	if err := validateRequest(c, &req, errorMessages, "/auth/login"); err != nil {
		return err
	}

	c.Locals("loginRequest", req)
	return c.Next()
}

func ValidateUpdatePasswordRequest(c *fiber.Ctx) error {
	var req UpdatePasswordRequest
	errorMessages := map[string]string{
		"CurrentPassword_required": "Mevcut şifre zorunludur",
		"CurrentPassword_min":      "Mevcut şifre en az 6 karakter olmalıdır",
		"NewPassword_required":     "Yeni şifre zorunludur",
		"NewPassword_min":          "Yeni şifre en az 8 karakter olmalıdır",
		"NewPassword_nefield":      "Yeni şifre mevcut şifreden farklı olmalıdır",
		"ConfirmPassword_required": "Şifre tekrarı zorunludur",
		"ConfirmPassword_eqfield":  "Yeni şifreler uyuşmuyor",
	}

	if err := validateRequest(c, &req, errorMessages, "/auth/update-password"); err != nil {
		return err
	}

	c.Locals("updatePasswordRequest", req)
	return c.Next()
}

func ValidateRegisterRequest(c *fiber.Ctx) error {
	var req RegisterRequest
	errorMessages := map[string]string{
		"Name_required":            "İsim zorunludur",
		"Name_min":                 "İsim en az 3 karakter olmalıdır",
		"Email_required":           "E-posta zorunludur",
		"Email_email":              "Geçerli bir e-posta adresi giriniz",
		"Password_required":        "Şifre zorunludur",
		"Password_min":             "Şifre en az 8 karakter olmalıdır",
		"ConfirmPassword_required": "Şifre tekrarı zorunludur",
		"ConfirmPassword_eqfield":  "Şifreler eşleşmiyor",
	}

	if err := validateRequest(c, &req, errorMessages, "/auth/register"); err != nil {
		return err
	}

	c.Locals("registerRequest", req)
	return c.Next()
}

func ValidateForgotPasswordRequest(c *fiber.Ctx) error {
	var req ForgotPasswordRequest
	errorMessages := map[string]string{
		"Email_required": "E-posta zorunludur",
		"Email_email":    "Geçerli bir e-posta adresi giriniz",
	}

	if err := validateRequest(c, &req, errorMessages, "/auth/forgot-password"); err != nil {
		return err
	}

	c.Locals("forgotPasswordRequest", req)
	return c.Next()
}

func ValidateResetPasswordRequest(c *fiber.Ctx) error {
	var req ResetPasswordRequest
	errorMessages := map[string]string{
		"Token_required":           "Token zorunludur",
		"NewPassword_required":     "Yeni şifre zorunludur",
		"NewPassword_min":          "Yeni şifre en az 8 karakter olmalıdır",
		"ConfirmPassword_required": "Şifre onayı zorunludur",
		"ConfirmPassword_eqfield":  "Şifreler eşleşmiyor",
	}

	if err := validateRequest(c, &req, errorMessages, "/auth/reset-password"); err != nil {
		return err
	}

	c.Locals("resetPasswordRequest", req)
	return c.Next()
}

func ValidateResendVerificationRequest(c *fiber.Ctx) error {
	var req ResendVerificationRequest
	errorMessages := map[string]string{
		"Email_required": "E-posta zorunludur",
		"Email_email":    "Geçerli bir e-posta adresi giriniz",
	}

	if err := validateRequest(c, &req, errorMessages, "/auth/resend-verification"); err != nil {
		return err
	}

	c.Locals("resendVerificationRequest", req)
	return c.Next()
}

func ValidateUpdateInfoRequest(c *fiber.Ctx) error {
	var req UpdateInfoRequest
	errorMessages := map[string]string{
		"Name_required":  "İsim zorunludur",
		"Name_min":       "İsim en az 3 karakter olmalıdır",
		"Email_required": "E-posta zorunludur",
		"Email_email":    "Geçerli bir e-posta adresi giriniz",
	}

	if err := validateRequest(c, &req, errorMessages, "/auth/profile/update-info"); err != nil {
		return err
	}

	c.Locals("updateInfoRequest", req)
	return c.Next()
}

func ValidateVerifyEmailRequest(c *fiber.Ctx) error {
	var req VerifyEmailRequest
	errorMessages := map[string]string{
		"Token_required": "Doğrulama token'ı zorunludur",
	}

	if err := validateRequest(c, &req, errorMessages, "/auth/login"); err != nil {
		return err
	}

	c.Locals("verifyEmailRequest", req)
	return c.Next()
}
