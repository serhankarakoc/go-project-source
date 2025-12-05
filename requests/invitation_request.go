package requests

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type InvitationRequest struct {
	Image                 string                  `form:"image" validate:"-"`
	InvitationKey         string                  `form:"invitation_key" validate:"-"`
	CategoryID            uint                    `form:"category_id" validate:"required,gt=0"`
	IsConfirmed           string                  `form:"is_confirmed" validate:"required,oneof=true false"`
	IsParticipant         string                  `form:"is_participant" validate:"-"`
	IsMultipleParticipant string                  `form:"is_multiple_participant" validate:"-"`
	IsFree                string                  `form:"is_free" validate:"-"`
	Description           string                  `form:"description" validate:"-"`
	Venue                 string                  `form:"venue" validate:"-"`
	Address               string                  `form:"address" validate:"-"`
	Location              string                  `form:"location" validate:"-"`
	Link                  string                  `form:"link" validate:"-"`
	Telephone             string                  `form:"telephone" validate:"-"`
	Note                  string                  `form:"note" validate:"-"`
	Date                  string                  `form:"date" validate:"-"`
	Time                  string                  `form:"time" validate:"-"`
	Detail                InvitationDetailRequest `form:"detail" validate:"-"`
}

type InvitationDetailRequest struct {
	Title              string `form:"title" validate:"-"`
	BrideName          string `form:"bride_name" validate:"-"`
	BrideSurname       string `form:"bride_surname" validate:"-"`
	BrideMotherName    string `form:"bride_mother_name" validate:"-"`
	BrideMotherSurname string `form:"bride_mother_surname" validate:"-"`
	BrideFatherName    string `form:"bride_father_name" validate:"-"`
	BrideFatherSurname string `form:"bride_father_surname" validate:"-"`
	GroomName          string `form:"groom_name" validate:"-"`
	GroomSurname       string `form:"groom_surname" validate:"-"`
	GroomMotherName    string `form:"groom_mother_name" validate:"-"`
	GroomMotherSurname string `form:"groom_mother_surname" validate:"-"`
	GroomFatherName    string `form:"groom_father_name" validate:"-"`
	GroomFatherSurname string `form:"groom_father_surname" validate:"-"`
	Person             string `form:"person" validate:"-"`
	MotherName         string `form:"mother_name" validate:"-"`
	MotherSurname      string `form:"mother_surname" validate:"-"`
	FatherName         string `form:"father_name" validate:"-"`
	FatherSurname      string `form:"father_surname" validate:"-"`
	IsMotherLive       string `form:"is_mother_live" validate:"-"`
	IsFatherLive       string `form:"is_father_live" validate:"-"`
	IsBrideMotherLive  string `form:"is_bride_mother_live" validate:"-"`
	IsBrideFatherLive  string `form:"is_bride_father_live" validate:"-"`
	IsGroomMotherLive  string `form:"is_groom_mother_live" validate:"-"`
	IsGroomFatherLive  string `form:"is_groom_father_live" validate:"-"`
}

func ParseAndValidateInvitationRequest(c *fiber.Ctx) (InvitationRequest, error) {
	var req InvitationRequest

	if err := c.BodyParser(&req); err != nil {
		return req, errors.New("geçersiz istek formatı")
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		field := validationErrors[0].Field()
		tag := validationErrors[0].Tag()
		errorMessages := map[string]string{
			"CategoryID_required":  "Kategori seçilmelidir.",
			"IsConfirmed_required": "Onay durumu seçilmelidir.",
		}
		if msg, ok := errorMessages[field+"_"+tag]; ok {
			return req, errors.New(msg)
		}
		return req, errors.New("lütfen formdaki hataları düzeltin")
	}
	return req, nil
}
