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

type IUserService interface {
	GetAllUsers(ctx context.Context, params requests.UserListParams) (*requests.PaginatedResult, error)
	GetUserByID(ctx context.Context, id uint) (*models.User, error)
	CreateUser(ctx context.Context, req requests.CreateUserRequest) error
	UpdateUser(ctx context.Context, id uint, req requests.UpdateUserRequest) error
	DeleteUser(ctx context.Context, id uint) error
	GetUserCount(ctx context.Context) (int64, error)
}

type UserService struct {
	repo repositories.IUserRepository
}

func NewUserService() IUserService {
	return &UserService{repo: repositories.NewUserRepository()}
}

func (s *UserService) GetAllUsers(ctx context.Context, params requests.UserListParams) (*requests.PaginatedResult, error) {
	// Repository'yi çağır (UserListParams tipinde)
	users, totalCount, err := s.repo.GetAllUsers(ctx, params)
	if err != nil {
		logconfig.Log.Error("Kullanıcılar alınamadı", zap.Error(err))
		return nil, errors.New("kullanıcılar getirilirken bir hata oluştu")
	}

	// UserType'daki gibi requests.CreatePaginatedResult kullan
	return requests.CreatePaginatedResult(users, totalCount, params.Page, params.PerPage), nil
}

func (s *UserService) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		logconfig.Log.Warn("Kullanıcı bulunamadı", zap.Uint("user_id", id), zap.Error(err))
		return nil, errors.New("kullanıcı bulunamadı")
	}
	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, req requests.CreateUserRequest) error {
	// Request'i convert et
	converted := req.BaseUserRequest.Convert()

	// UserType seçimi zorunlu
	if converted.UserTypeID == nil {
		return errors.New("kullanıcı tipi seçilmelidir")
	}

	// Model oluştur
	user := &models.User{
		BaseModel: models.BaseModel{
			IsActive: converted.IsActive != nil && *converted.IsActive,
		},
		Name:              converted.Name,
		Email:             converted.Email,
		Password:          req.Password,
		UserTypeID:        *converted.UserTypeID,
		ResetToken:        converted.ResetToken,
		EmailVerified:     converted.EmailVerified != nil && *converted.EmailVerified,
		VerificationToken: converted.VerificationToken,
		Provider:          converted.Provider,
		ProviderID:        converted.ProviderID,
	}

	// Şifre kontrolü ve hash'leme
	if user.Password == "" {
		return errors.New("şifre alanı boş olamaz")
	}
	if err := user.SetPassword(user.Password); err != nil {
		logconfig.Log.Error("Şifre oluşturulamadı", zap.Error(err))
		return errors.New("şifre oluşturulurken hata oluştu")
	}

	// Repository'e kaydet
	return s.repo.CreateUser(ctx, user)
}

func (s *UserService) UpdateUser(ctx context.Context, id uint, req requests.UpdateUserRequest) error {
	// Mevcut kullanıcıyı kontrol et
	_, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return errors.New("kullanıcı bulunamadı")
	}

	// Request'i convert et
	converted := req.BaseUserRequest.Convert()

	// UserType seçimi zorunlu
	if converted.UserTypeID == nil {
		return errors.New("kullanıcı tipi seçilmelidir")
	}

	// Update data hazırla
	updateData := map[string]interface{}{
		"name":               converted.Name,
		"email":              converted.Email,
		"is_active":          converted.IsActive != nil && *converted.IsActive,
		"user_type_id":       *converted.UserTypeID,
		"reset_token":        converted.ResetToken,
		"email_verified":     converted.EmailVerified != nil && *converted.EmailVerified,
		"verification_token": converted.VerificationToken,
		"provider":           converted.Provider,
		"provider_id":        converted.ProviderID,
	}

	// Şifre değişikliği (optional)
	if req.Password != "" {
		hasher := models.User{}
		if err := hasher.SetPassword(req.Password); err != nil {
			return errors.New("şifre oluşturulurken hata oluştu")
		}
		updateData["password"] = hasher.Password
	}

	// Repository'de güncelle
	return s.repo.UpdateUser(ctx, id, updateData)
}

func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	return s.repo.DeleteUser(ctx, id)
}

func (s *UserService) GetUserCount(ctx context.Context) (int64, error) {
	return s.repo.GetUserCount(ctx)
}
