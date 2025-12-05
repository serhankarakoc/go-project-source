package handlers

import (
	"net/http"
	"strings"

	"zatrano/models"
	"zatrano/pkg/flashmessages"
	"zatrano/pkg/formflash"
	"zatrano/pkg/queryparams"
	"zatrano/pkg/renderer"
	"zatrano/requests"
	"zatrano/services"

	"github.com/gofiber/fiber/v2"
)

type DashboardUserHandler struct {
	userService services.IUserService
}

func NewDashboardUserHandler() *DashboardUserHandler {
	return &DashboardUserHandler{
		userService: services.NewUserService(),
	}
}

func (h *DashboardUserHandler) ListUsers(c *fiber.Ctx) error {
	params, fieldErrors, err := requests.ParseAndValidateUserList(c)
	if err != nil {
		renderData := fiber.Map{
			"Title":            "Kullanıcılar",
			"ValidationErrors": fieldErrors,
			"Params": fiber.Map{
				"Name":       params.Name,
				"Email":      params.Email,
				"IsActive":   params.IsActive,
				"UserTypeID": params.UserTypeID,
				"SortBy":     params.SortBy,
				"OrderBy":    params.OrderBy,
				"Page":       params.Page,
				"PerPage":    params.PerPage,
			},
			"Result": &queryparams.PaginatedResult{
				Data: []models.User{},
				Meta: queryparams.PaginationMeta{
					CurrentPage: params.Page,
					PerPage:     params.PerPage,
					TotalItems:  0,
					TotalPages:  0,
				},
			},
		}
		return renderer.Render(c, "dashboard/users/list", "layouts/app", renderData, http.StatusBadRequest)
	}

	paginatedResult, err := h.userService.GetAllUsers(c.UserContext(), params)

	renderData := fiber.Map{
		"Title":  "Kullanıcılar",
		"Result": paginatedResult,
		"Params": fiber.Map{
			"Name":       params.Name,
			"Email":      params.Email,
			"IsActive":   params.IsActive,
			"UserTypeID": params.UserTypeID,
			"SortBy":     params.SortBy,
			"OrderBy":    params.OrderBy,
			"Page":       params.Page,
			"PerPage":    params.PerPage,
		},
	}

	if err != nil {
		renderData[renderer.FlashErrorKeyView] = "Kullanıcılar getirilirken bir hata oluştu."
		renderData["Result"] = &queryparams.PaginatedResult{
			Data: []models.User{},
			Meta: queryparams.PaginationMeta{
				CurrentPage: params.Page,
				PerPage:     params.PerPage,
				TotalItems:  0,
				TotalPages:  0,
			},
		}
	}

	return renderer.Render(c, "dashboard/users/list", "layouts/app", renderData, http.StatusOK)
}

func (h *DashboardUserHandler) ShowCreateUser(c *fiber.Ctx) error {
	return renderer.Render(c, "dashboard/users/create", "layouts/app", fiber.Map{
		"Title": "Yeni Kullanıcı Ekle",
	})
}

func (h *DashboardUserHandler) CreateUser(c *fiber.Ctx) error {
	// Form verilerini map olarak al
	formData := make(map[string]string)
	args := c.Request().PostArgs()
	args.VisitAll(func(key, value []byte) {
		formData[string(key)] = string(value)
	})

	// Create için özel validation
	req, fieldErrors, err := requests.ParseAndValidateCreateRequest(c)
	if err != nil {
		// Form verilerini kaydet
		formflash.SetData(c, formData)

		// Field-specific hataları kaydet
		formflash.SetValidationErrors(c, fieldErrors)

		// Genel hata mesajı
		flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, err.Error())

		return c.Redirect("/dashboard/users/create")
	}

	// YENİ: CreateFromRequest kullan (CreateUserRequest tipinde)
	if err := h.userService.CreateUser(c.UserContext(), req); err != nil {
		// Servis hatası - form verilerini koru
		formflash.SetData(c, formData)
		flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Kullanıcı oluşturulamadı: "+err.Error())
		return c.Redirect("/dashboard/users/create")
	}

	// BAŞARILI - form verilerini temizle
	formflash.ClearData(c)
	flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Kullanıcı başarıyla oluşturuldu.")
	return c.Redirect("/dashboard/users", fiber.StatusFound)
}

func (h *DashboardUserHandler) ShowUpdateUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("Geçersiz kullanıcı ID")
	}

	user, err := h.userService.GetUserByID(c.UserContext(), uint(id))
	if err != nil {
		flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Kullanıcı bulunamadı.")
		return c.Redirect("/dashboard/users", fiber.StatusSeeOther)
	}

	return renderer.Render(c, "dashboard/users/update", "layouts/app", fiber.Map{
		"Title": "Kullanıcı Düzenle",
		"User":  user,
	})
}

func (h *DashboardUserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Geçersiz kullanıcı ID")
	}

	// Form verilerini map olarak al
	formData := make(map[string]string)
	args := c.Request().PostArgs()
	args.VisitAll(func(key, value []byte) {
		formData[string(key)] = string(value)
	})

	// Update için özel validation
	req, fieldErrors, err := requests.ParseAndValidateUpdateRequest(c)
	if err != nil {
		// Form verilerini kaydet
		formflash.SetData(c, formData)

		// Field-specific hataları kaydet
		formflash.SetValidationErrors(c, fieldErrors)

		// Genel hata mesajı
		flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, err.Error())

		return c.Redirect("/dashboard/users/update/" + c.Params("id"))
	}

	// YENİ: UpdateFromRequest kullan (UpdateUserRequest tipinde)
	if err := h.userService.UpdateUser(c.UserContext(), uint(id), req); err != nil {
		// Servis hatası - form verilerini koru
		formflash.SetData(c, formData)
		flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Kullanıcı güncellenemedi: "+err.Error())
		return c.Redirect("/dashboard/users/update/" + c.Params("id"))
	}

	// BAŞARILI - form verilerini temizle
	formflash.ClearData(c)
	flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Kullanıcı başarıyla güncellendi.")
	return c.Redirect("/dashboard/users", fiber.StatusFound)
}

func (h *DashboardUserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Geçersiz kullanıcı ID")
	}

	if err := h.userService.DeleteUser(c.UserContext(), uint(id)); err != nil {
		errMsg := "Kullanıcı silinemedi: " + err.Error()
		if strings.Contains(c.Get("Accept"), "application/json") {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": errMsg})
		}
		flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, errMsg)
		return c.Redirect("/dashboard/users", fiber.StatusSeeOther)
	}

	if strings.Contains(c.Get("Accept"), "application/json") {
		return c.JSON(fiber.Map{"message": "Kullanıcı başarıyla silindi."})
	}

	flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Kullanıcı başarıyla silindi.")
	return c.Redirect("/dashboard/users", fiber.StatusFound)
}
