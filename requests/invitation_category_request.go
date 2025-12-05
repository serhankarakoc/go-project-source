package requests

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type InvitationCategoryRequest struct {
	Name     string  `form:"name" validate:"required,min=2"`
	Icon     string  `form:"icon" validate:"required"`
	Template string  `form:"template" validate:"required"`
	IsActive string  `form:"is_active" validate:"required,oneof=true false"`
	Price    float64 `form:"price" validate:"required"`
}

func ParseAndValidateInvitationCategoryRequest(c *fiber.Ctx) (InvitationCategoryRequest, error) {
	var req InvitationCategoryRequest

	if err := c.BodyParser(&req); err != nil {
		return req, errors.New("geçersiz istek formatı")
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		field := validationErrors[0].Field()
		tag := validationErrors[0].Tag()
		errorMessages := map[string]string{
			"Name_required":     "Kategori adı zorunludur.",
			"Name_min":          "Kategori adı en az 2 karakter olmalıdır.",
			"Icon_required":     "İkon zorunludur.",
			"Template_required": "Şablon zorunludur.",
			"IsActive_required": "Durum (Aktif/Pasif) seçilmelidir.",
			"IsActive_oneof":    "Durum için geçersiz bir değer seçildi.",
			"Price_required":    "Fiyat zorunludur.",
		}
		if msg, ok := errorMessages[field+"_"+tag]; ok {
			return req, errors.New(msg)
		}
		return req, errors.New("lütfen formdaki hataları düzeltin")
	}
	return req, nil
}
