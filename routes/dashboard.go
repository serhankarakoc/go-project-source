package routes

import (
	handlers "zatrano/handlers/dashboard"
	"zatrano/middlewares"

	"github.com/gofiber/fiber/v2"
)

func registerDashboardRoutes(app *fiber.App) {
	dashboardGroup := app.Group("/dashboard")
	dashboardGroup.Use(
		middlewares.AuthMiddleware,
		middlewares.UserTypeMiddleware(1),
	)

	// Dashboard anasayfa
	dashboardHomeHandler := handlers.NewDashboardHomeHandler()
	dashboardGroup.Get("/", dashboardHomeHandler.HomePage)
	dashboardGroup.Get("/home", dashboardHomeHandler.HomePage)

	// Kullanıcı türleri yönetimi
	userTypeHandler := handlers.NewDashboardUserTypeHandler()
	dashboardGroup.Get("/user-types", userTypeHandler.ListUserTypes)
	dashboardGroup.Get("/user-types/create", userTypeHandler.ShowCreateUserType)
	dashboardGroup.Post("/user-types/create", userTypeHandler.CreateUserType)
	dashboardGroup.Get("/user-types/update/:id", userTypeHandler.ShowUpdateUserType)
	dashboardGroup.Post("/user-types/update/:id", userTypeHandler.UpdateUserType)
	dashboardGroup.Delete("/user-types/delete/:id", userTypeHandler.DeleteUserType)

	// Kullanıcı yönetimi
	userHandler := handlers.NewDashboardUserHandler()
	dashboardGroup.Get("/users", userHandler.ListUsers)
	dashboardGroup.Get("/users/create", userHandler.ShowCreateUser)
	dashboardGroup.Post("/users/create", userHandler.CreateUser)
	dashboardGroup.Get("/users/update/:id", userHandler.ShowUpdateUser)
	dashboardGroup.Post("/users/update/:id", userHandler.UpdateUser)
	dashboardGroup.Delete("/users/delete/:id", userHandler.DeleteUser)
}
