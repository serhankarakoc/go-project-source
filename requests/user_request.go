package requests

import (
	"errors"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type BaseUserRequest struct {
	Name              string `form:"name" validate:"required,min=3"`
	Email             string `form:"email" validate:"required,email"`
	IsActive          string `form:"is_active" validate:"required,oneof=true false"`
	UserTypeID        string `form:"user_type_id" validate:"required"`
	EmailVerified     string `form:"email_verified" validate:"required,oneof=true false"`
	ResetToken        string `form:"reset_token"`
	VerificationToken string `form:"verification_token"`
	Provider          string `form:"provider"`
	ProviderID        string `form:"provider_id"`
}

type ConvertedBaseUserRequest struct {
	Name              string
	Email             string
	IsActive          *bool
	UserTypeID        *uint
	EmailVerified     *bool
	ResetToken        string
	VerificationToken string
	Provider          string
	ProviderID        string
}

func (r *BaseUserRequest) Convert() ConvertedBaseUserRequest {
	var isActivePtr *bool
	if r.IsActive != "" {
		val := r.IsActive == "true"
		isActivePtr = &val
	}

	var emailVerifiedPtr *bool
	if r.EmailVerified != "" {
		val := r.EmailVerified == "true"
		emailVerifiedPtr = &val
	}

	var userTypeIDPtr *uint
	if r.UserTypeID != "" {
		if val, err := strconv.ParseUint(r.UserTypeID, 10, 32); err == nil {
			uintVal := uint(val)
			userTypeIDPtr = &uintVal
		}
	}

	return ConvertedBaseUserRequest{
		Name:              r.Name,
		Email:             r.Email,
		IsActive:          isActivePtr,
		UserTypeID:        userTypeIDPtr,
		EmailVerified:     emailVerifiedPtr,
		ResetToken:        r.ResetToken,
		VerificationToken: r.VerificationToken,
		Provider:          r.Provider,
		ProviderID:        r.ProviderID,
	}
}

type CreateUserRequest struct {
	BaseUserRequest
	Password        string `form:"password" validate:"required,min=6"`
	ConfirmPassword string `form:"confirmPassword" validate:"required,eqfield=Password"`
}

type UpdateUserRequest struct {
	BaseUserRequest
	Password string `form:"password" validate:"omitempty,min=6"`
}

type UserListRequest struct {
	Name       string `query:"name"`
	Email      string `query:"email"`
	IsActive   string `query:"is_active" validate:"omitempty,oneof=true false"`
	UserTypeID string `query:"user_type_id" validate:"omitempty,numeric"`
	SortBy     string `query:"sortBy" validate:"omitempty,oneof=id name email created_at"`
	OrderBy    string `query:"orderBy" validate:"omitempty,oneof=asc desc"`
	Page       string `query:"page" validate:"omitempty,numeric,min=1"`
	PerPage    string `query:"perPage" validate:"omitempty,numeric,min=1,max=200"`
	Search     string `query:"search"`
}

type UserListParams struct {
	Name       string
	Email      string
	IsActive   string
	UserTypeID *uint
	SortBy     string
	OrderBy    string
	Page       int
	PerPage    int
	Search     string
}

func (r *UserListRequest) ToServiceParams() UserListParams {
	params := UserListParams{
		Name:     strings.TrimSpace(r.Name),
		Email:    strings.TrimSpace(r.Email),
		IsActive: strings.TrimSpace(r.IsActive),
		SortBy:   strings.TrimSpace(r.SortBy),
		OrderBy:  strings.TrimSpace(r.OrderBy),
		Search:   strings.TrimSpace(r.Search),
	}

	if r.UserTypeID != "" {
		if val, err := strconv.ParseUint(r.UserTypeID, 10, 32); err == nil {
			uintVal := uint(val)
			params.UserTypeID = &uintVal
		}
	}

	if r.Page != "" {
		if page, err := strconv.Atoi(r.Page); err == nil && page > 0 {
			params.Page = page
		}
	}

	if r.PerPage != "" {
		if perPage, err := strconv.Atoi(r.PerPage); err == nil && perPage > 0 {
			params.PerPage = perPage
		}
	}

	params.applyDefaults()

	return params
}

func (p *UserListParams) applyDefaults() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PerPage <= 0 {
		p.PerPage = 20
	}
	if p.SortBy == "" {
		p.SortBy = "name"
	}
	if p.OrderBy == "" {
		p.OrderBy = "asc"
	}
}

func (p *UserListParams) CalculateOffset() int {
	if p.Page <= 0 {
		return 0
	}
	return (p.Page - 1) * p.PerPage
}

func ParseAndValidateCreateRequest(c *fiber.Ctx) (CreateUserRequest, map[string]string, error) {
	var req CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return req, make(map[string]string), errors.New("geçersiz istek formatı")
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		validationErrors := GetUserValidationErrors(err)
		return req, validationErrors, errors.New("lütfen formdaki hataları düzeltin")
	}

	return req, make(map[string]string), nil
}

func ParseAndValidateUpdateRequest(c *fiber.Ctx) (UpdateUserRequest, map[string]string, error) {
	var req UpdateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return req, make(map[string]string), errors.New("geçersiz istek formatı")
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		validationErrors := GetUserValidationErrors(err)
		return req, validationErrors, errors.New("lütfen formdaki hataları düzeltin")
	}

	return req, make(map[string]string), nil
}

func ParseAndValidateUserList(c *fiber.Ctx) (UserListParams, map[string]string, error) {
	var req UserListRequest

	if err := c.QueryParser(&req); err != nil {
		return UserListParams{}, make(map[string]string), errors.New("geçersiz sorgu parametreleri")
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		validationErrors := GetUserListValidationErrors(err)
		return UserListParams{}, validationErrors, errors.New("lütfen filtreleri kontrol edin")
	}

	return req.ToServiceParams(), make(map[string]string), nil
}

func GetUserValidationErrors(err error) map[string]string {
	errorMessages := map[string]string{
		"Name_required":           "Kullanıcı adı zorunludur.",
		"Name_min":                "Kullanıcı adı en az 3 karakter olmalıdır.",
		"Email_required":          "E-posta adresi zorunludur.",
		"Email_email":             "Geçerli bir e-posta adresi giriniz.",
		"IsActive_required":       "Kullanıcı durumu seçilmelidir.",
		"IsActive_oneof":          "Geçerli bir durum seçiniz (Aktif/Pasif).",
		"EmailVerified_required":  "E-posta doğrulama durumu seçilmelidir.",
		"EmailVerified_oneof":     "Geçerli bir e-posta doğrulama durumu seçiniz.",
		"UserTypeID_required":     "Kullanıcı tipi seçilmelidir.",
		"Password_required":       "Şifre zorunludur.",
		"Password_min":            "Şifre en az 6 karakter olmalıdır.",
		"ConfirmPassword_eqfield": "Şifreler eşleşmiyor.",
	}

	return CommonValidationErrors(err, errorMessages)
}

func GetUserListValidationErrors(err error) map[string]string {
	errorMessages := map[string]string{
		"IsActive_oneof":     "Durum sadece 'true' veya 'false' olabilir.",
		"UserTypeID_numeric": "Kullanıcı tipi ID'si sayı olmalıdır.",
		"SortBy_oneof":       "Sıralama alanı sadece 'id', 'name', 'email' veya 'created_at' olabilir.",
		"OrderBy_oneof":      "Sıralama yönü sadece 'asc' veya 'desc' olabilir.",
		"Page_numeric":       "Sayfa numarası sayı olmalıdır.",
		"Page_min":           "Sayfa numarası en az 1 olmalıdır.",
		"PerPage_numeric":    "Sayfa başı kayıt sayısı sayı olmalıdır.",
		"PerPage_min":        "Sayfa başı kayıt en az 1 olmalıdır.",
		"PerPage_max":        "Sayfa başı kayıt en fazla 200 olmalıdır.",
	}

	return CommonValidationErrors(err, errorMessages)
}
