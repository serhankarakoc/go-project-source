package handlers

import (
	"net/http"
	"zatrano/pkg/renderer"

	"github.com/gofiber/fiber/v2"
)

type DashboardHomeHandler struct {
}

func NewDashboardHomeHandler() *DashboardHomeHandler {
	return &DashboardHomeHandler{}
}

func (h *DashboardHomeHandler) HomePage(c *fiber.Ctx) error {

	mapData := fiber.Map{
		"Title": "Dashboard",
	}
	return renderer.Render(c, "dashboard/home/home", "layouts/app", mapData, http.StatusOK)
}
