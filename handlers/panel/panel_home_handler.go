package handlers

import (
	"net/http"

	"zatrano/pkg/renderer"

	"github.com/gofiber/fiber/v2"
)

type PanelHomeHandler struct {
}

func NewPanelHomeHandler() *PanelHomeHandler {
	return &PanelHomeHandler{}
}

func (h *PanelHomeHandler) HomePage(c *fiber.Ctx) error {

	mapData := fiber.Map{
		"Title": "Panel",
	}
	return renderer.Render(c, "panel/home/home", "layouts/panel", mapData, http.StatusOK)
}
