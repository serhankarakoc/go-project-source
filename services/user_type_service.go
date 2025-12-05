package services

import (
	"context"
	"errors"

	"zatrano/configs/logconfig"
	"zatrano/models"
	"zatrano/repositories"
	"zatrano/requests"

	"go.uber.org/zap"
)

type IUserTypeService interface {
	GetAllUserTypes(ctx context.Context, params requests.UserTypeListParams) (*requests.PaginatedResult, error)
	GetUserTypeByID(ctx context.Context, id uint) (*models.UserType, error)
	CreateUserType(ctx context.Context, req requests.UserTypeRequest) error
	UpdateUserType(ctx context.Context, id uint, req requests.UserTypeRequest) error
	DeleteUserType(ctx context.Context, id uint) error
}

type UserTypeService struct {
	repo repositories.IUserTypeRepository
}

func NewUserTypeService() IUserTypeService {
	return &UserTypeService{
		repo: repositories.NewUserTypeRepository(),
	}
}

func (s *UserTypeService) GetAllUserTypes(ctx context.Context, params requests.UserTypeListParams) (*requests.PaginatedResult, error) {
	// Repository'yi çağır
	userTypes, totalCount, err := s.repo.GetAllUserTypes(ctx, params)
	if err != nil {
		// Hata log'la
		return nil, err
	}

	// PaginatedResult oluştur
	return requests.CreatePaginatedResult(userTypes, totalCount, params.Page, params.PerPage), nil
}

func (s *UserTypeService) GetUserTypeByID(ctx context.Context, id uint) (*models.UserType, error) {
	userType, err := s.repo.GetUserTypeByID(ctx, id)
	if err != nil {
		logconfig.Log.Warn("Kullanıcı Tipi bulunamadı", zap.Uint("user_type_id", id), zap.Error(err))
		return nil, errors.New("kullanıcı tipi bulunamadı")
	}
	return userType, nil
}

func (s *UserTypeService) CreateUserType(ctx context.Context, req requests.UserTypeRequest) error {
	// Request'i convert et
	converted := req.BaseUserTypeRequest.Convert()

	// Model oluştur
	userType := &models.UserType{
		BaseModel: models.BaseModel{IsActive: false}, // default
		Name:      converted.Name,
	}

	// IsActive nil kontrolü
	if converted.IsActive != nil {
		userType.BaseModel.IsActive = *converted.IsActive
	}

	// Repository'e kaydet
	return s.repo.CreateUserType(ctx, userType)
}

func (s *UserTypeService) UpdateUserType(ctx context.Context, id uint, req requests.UserTypeRequest) error {
	// Mevcut user type'ı kontrol et
	_, err := s.repo.GetUserTypeByID(ctx, id)
	if err != nil {
		return errors.New("kullanıcı tipi bulunamadı")
	}

	// Request'i convert et
	converted := req.BaseUserTypeRequest.Convert()

	// Update data hazırla
	updateData := map[string]interface{}{
		"name":      converted.Name,
		"is_active": converted.IsActive != nil && *converted.IsActive,
	}

	// Repository'de güncelle
	return s.repo.UpdateUserType(ctx, id, updateData)
}

func (s *UserTypeService) DeleteUserType(ctx context.Context, id uint) error {
	return s.repo.DeleteUserType(ctx, id)
}
