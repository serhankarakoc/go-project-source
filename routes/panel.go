package routes

import (
	handlers "zatrano/handlers/panel"
	"zatrano/middlewares"

	"github.com/gofiber/fiber/v2"
)

func registerPanelRoutes(app *fiber.App) {
	panelGroup := app.Group("/panel")
	panelGroup.Use(
		middlewares.AuthMiddleware,
		middlewares.UserTypeMiddleware(2),
	)

	panelHomeHandler := handlers.NewPanelHomeHandler()
	panelGroup.Get("/", panelHomeHandler.HomePage)
	panelGroup.Get("/anasayfa", panelHomeHandler.HomePage)
}
