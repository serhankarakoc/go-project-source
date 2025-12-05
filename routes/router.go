package routes

import (
	"zatrano/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// global limiter & middleware’lar
	app.Use(middlewares.GlobalRateLimit())
	app.Use(middlewares.FormPostRateLimit())

	app.Use(middlewares.SessionMiddleware())
	app.Use(middlewares.ZapLogger())

	// önce auth alanı
	registerAuthRoutes(app)
	registerDashboardRoutes(app)
	registerPanelRoutes(app)

	// web/public area
	registerWebsiteRoutes(app)
}
