package handlers

import (
	"net/http"

	"zatrano/pkg/renderer"

	"github.com/gofiber/fiber/v2"
)

type WebsiteHandler struct {
}

func NewWebsiteHandler() *WebsiteHandler {
	return &WebsiteHandler{}
}

func (h *WebsiteHandler) HomePage(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/home", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) KullanimSartlari(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/kullanim-sartlari", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) Acilis(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-acilis-davetiyesi", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) AfterParty(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-after-party-davetiyesi", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) AnitkabirZiyareti(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-anitkabir-ziyareti-davetiyesi", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) AskerEglencesi(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-asker-eglencesi-davetiyesi", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) BabyShower(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-baby-shower-davetiyesi", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) Balo(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-balo-davetiyesi", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) BekarligaVeda(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-bekarliga-veda-davetiyesi", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) CinsiyetPartisi(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-cinsiyet-partisi-davetiyesi", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) Defile(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-defile-davetiyesi", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) DiniToren(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-dini-toren-davetiyesi", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) DogumGunu(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-dogum-gunu-davetiyesi", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) Dugun(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-dugun-davetiyesi", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) Egitim(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-egitim-davetiyesi", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) Festival(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-festival-davetiyesi", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) FilmGalasi(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-film-galasi-davetiyesi", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) Fuar(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-fuar-davetiyesi", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) GelinHamami(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-gelin-hamami-davetiyesi", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) Gezi(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-gezi-davetiyesi", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) Kamp(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-kamp-davetiyesi", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) KinaGecesi(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-kina-gecesi-davetiyesi", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) Konferans(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-konferans-davetiyesi", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) Kongre(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-kongre-davetiyesi", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) Konser(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-konser-davetiyesi", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) Lansman(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-lansman-davetiyesi", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) Mezuniyet(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-mezuniyet-davetiyesi", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) NikahToreni(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-nikah-toreni-davetiyesi", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) Nisan(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-nisan-davetiyesi", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) OnlineEtkinlik(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-online-etkinlik-davetiyesi", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) Seminer(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-seminer-davetiyesi", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) Sergi(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-sergi-davetiyesi", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) SporMusabakasi(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-spor-musabakasi-davetiyesi", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) SunnetDugunu(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-sunnet-dugunu-davetiyesi", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) Tanitim(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-tanitim-davetiyesi", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) TiyatroGosterisi(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-tiyatro-gosterisi-davetiyesi", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) Toplanti(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-toplanti-davetiyesi", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) VedaPartisi(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-veda-partisi-davetiyesi", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) YilbasiPartisi(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-yilbasi-partisi-davetiyesi", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) Yildonumu(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/dijital-yildonumu-davetiyesi", "layouts/website", mapData, http.StatusOK)
}
