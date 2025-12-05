package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"zatrano/configs/sessionconfig"
	"zatrano/pkg/flashmessages"
	"zatrano/services"

	"github.com/gofiber/fiber/v2"
)

func GoogleLogin(c *fiber.Ctx) error {
	cfg := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URI"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}

	sess, err := sessionconfig.SessionStart(c)
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Oturum başlatılamadı.")
		return c.Redirect("/auth/login", fiber.StatusSeeOther)
	}

	stateToken, err := generateToken()
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "State token oluşturulamadı.")
		return c.Redirect("/auth/login", fiber.StatusSeeOther)
	}

	sess.Set("oauth_state", stateToken)
	_ = sess.Save()

	return c.Redirect(cfg.AuthCodeURL(stateToken, oauth2.AccessTypeOffline), http.StatusTemporaryRedirect)
}

func GoogleCallback(c *fiber.Ctx) error {
	cfg := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URI"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}

	state := c.Query("state")
	if state == "" {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "State parametresi eksik.")
		return c.Redirect("/auth/login", fiber.StatusSeeOther)
	}

	sess, err := sessionconfig.SessionStart(c)
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Oturum başlatılamadı.")
		return c.Redirect("/auth/login", fiber.StatusSeeOther)
	}
	if sess.Get("oauth_state") != state {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz state token.")
		return c.Redirect("/auth/login", fiber.StatusSeeOther)
	}

	code := c.Query("code")
	if code == "" {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Code parametresi eksik.")
		return c.Redirect("/auth/login", fiber.StatusSeeOther)
	}

	tok, err := cfg.Exchange(context.Background(), code)
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Token değişimi başarısız.")
		return c.Redirect("/auth/login", fiber.StatusSeeOther)
	}

	client := cfg.Client(context.Background(), tok)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Kullanıcı bilgileri alınamadı.")
		return c.Redirect("/auth/login", fiber.StatusSeeOther)
	}
	defer resp.Body.Close()

	var userInfo struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Kullanıcı bilgileri parse edilemedi.")
		return c.Redirect("/auth/login", fiber.StatusSeeOther)
	}

	authService := services.NewAuthService()
	// UserTypeID 2 → normal kullanıcı
	user, err := authService.FindOrCreateOAuthUser(userInfo.ID, userInfo.Email, userInfo.Name)
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Kullanıcı oluşturulamadı veya giriş yapılamadı.")
		return c.Redirect("/auth/login", fiber.StatusSeeOther)
	}

	sess.Set("user_id", user.ID)
	sess.Set("user_type_id", user.UserTypeID)
	sess.Set("is_active", user.IsActive)
	if err = sess.Save(); err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Oturum kaydedilemedi.")
		return c.Redirect("/auth/login", fiber.StatusSeeOther)
	}

	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Google ile giriş başarılı.")
	if user.UserTypeID == 1 {
		return c.Redirect("/dashboard/home", fiber.StatusSeeOther)
	}
	return c.Redirect("/panel/anasayfa", fiber.StatusSeeOther)
}
