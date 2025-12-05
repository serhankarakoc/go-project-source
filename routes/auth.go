package routes

import (
	handlers "zatrano/handlers/auth"
	"zatrano/middlewares"
	"zatrano/requests"

	"github.com/gofiber/fiber/v2"
)

func registerAuthRoutes(app *fiber.App) {
	authHandler := handlers.NewAuthHandler()
	authGroup := app.Group("/auth")

	authGroup.Get("/login", middlewares.GuestMiddleware, authHandler.ShowLogin)
	// Login'de brute-force limiter ekli
	authGroup.Post("/login",
		middlewares.GuestMiddleware,
		middlewares.LoginRateLimit(),
		requests.ValidateLoginRequest,
		authHandler.Login,
	)

	authGroup.Get("/logout", middlewares.AuthMiddleware, authHandler.Logout)
	authGroup.Get("/profile", middlewares.AuthMiddleware, authHandler.Profile)
	authGroup.Post("/profile/update-password", middlewares.AuthMiddleware, requests.ValidateUpdatePasswordRequest, authHandler.UpdatePassword)
	authGroup.Post("/profile/update-info", middlewares.AuthMiddleware, requests.ValidateUpdateInfoRequest, authHandler.UpdateInfo)

	authGroup.Get("/register", middlewares.GuestMiddleware, authHandler.ShowRegister)
	authGroup.Post("/register", middlewares.GuestMiddleware, requests.ValidateRegisterRequest, authHandler.Register)

	authGroup.Get("/forgot-password", middlewares.GuestMiddleware, authHandler.ShowForgotPassword)
	authGroup.Post("/forgot-password", middlewares.GuestMiddleware, requests.ValidateForgotPasswordRequest, authHandler.ForgotPassword)

	authGroup.Get("/reset-password", middlewares.GuestMiddleware, authHandler.ShowResetPassword)
	authGroup.Post("/reset-password", middlewares.GuestMiddleware, requests.ValidateResetPasswordRequest, authHandler.ResetPassword)

	authGroup.Get("/verify-email", middlewares.GuestMiddleware, authHandler.VerifyEmail)
	authGroup.Get("/resend-verification", middlewares.GuestMiddleware, authHandler.ShowResendVerification)
	authGroup.Post("/resend-verification", middlewares.GuestMiddleware, requests.ValidateResendVerificationRequest, authHandler.ResendVerification)

	authGroup.Get("/google/login", middlewares.GuestMiddleware, handlers.GoogleLogin)
	authGroup.Get("/google/callback", middlewares.GuestMiddleware, handlers.GoogleCallback)
}
