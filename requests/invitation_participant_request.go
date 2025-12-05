package requests

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type InvitationParticipantRequest struct {
	InvitationID     uint   `form:"invitation_id" validate:"required"`
	Name             string `form:"name" validate:"required,min=2"`
	Telephone        string `form:"telephone" validate:"required,min=10"`
	ParticipantCount int    `form:"participant_count" validate:"required"`
}

func ParseAndValidateInvitationParticipantRequest(c *fiber.Ctx) (InvitationParticipantRequest, error) {
	var req InvitationParticipantRequest

	if err := c.BodyParser(&req); err != nil {
		return req, errors.New("geçersiz istek formatı")
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		field := validationErrors[0].Field()
		tag := validationErrors[0].Tag()
		errorMessages := map[string]string{
			"InvitationID_required":     "Davet ID zorunludur.",
			"Name_required":             "Ad Soyad zorunludur.",
			"Name_min":                  "Ad Soyad en az 2 karakter olmalıdır.",
			"Telephone_required":        "Telefon numarası zorunludur.",
			"Telephone_min":             "Telefon numarası en az 10 karakter olmalıdır.",
			"ParticipantCount_required": "Kişi sayısı zorunludur.",
		}
		if msg, ok := errorMessages[field+"_"+tag]; ok {
			return req, errors.New(msg)
		}
		return req, errors.New("lütfen formdaki hataları düzeltin")
	}
	return req, nil
}
