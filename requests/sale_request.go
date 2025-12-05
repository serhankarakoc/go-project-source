package requests

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type SaleRequest struct {
	UserID        uint    `form:"user_id" validate:"required,gt=0"`
	CategoryID    uint    `form:"category_id" validate:"required,gt=0"`
	InvitationID  uint    `form:"invitation_id" validate:"required,gt=0"`
	TransactionID uint    `form:"transaction_id" validate:"required,gt=0"`
	Amount        float64 `form:"amount" validate:"required,gt=0"`
}

func ParseAndValidateSaleRequest(c *fiber.Ctx) (SaleRequest, error) {
	var req SaleRequest

	if err := c.BodyParser(&req); err != nil {
		return req, errors.New("geçersiz istek formatı")
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		field := validationErrors[0].Field()
		tag := validationErrors[0].Tag()
		errorMessages := map[string]string{
			"UserID_required":        "Kullanıcı zorunludur.",
			"UserID_gt":              "Geçerli bir kullanıcı seçiniz.",
			"CategoryID_required":    "Kategori zorunludur.",
			"CategoryID_gt":          "Geçerli bir kategori seçiniz.",
			"InvitationID_required":  "Davetiye zorunludur.",
			"InvitationID_gt":        "Geçerli bir davetiye seçiniz.",
			"TransactionID_required": "İşlem zorunludur.",
			"TransactionID_gt":       "Geçerli bir işlem seçiniz.",
			"Amount_required":        "Tutar zorunludur.",
			"Amount_gt":              "Tutar 0'dan büyük olmalıdır.",
		}
		if msg, ok := errorMessages[field+"_"+tag]; ok {
			return req, errors.New(msg)
		}
		return req, errors.New("lütfen formdaki hataları düzeltin")
	}

	return req, nil
}
