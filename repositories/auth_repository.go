package repositories

import (
	"context"

	"zatrano/configs/databaseconfig"
	"zatrano/configs/logconfig"
	"zatrano/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IAuthRepository interface {
	FindUserByEmail(email string) (*models.User, error)
	FindUserByID(id uint) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	CreateUser(ctx context.Context, user *models.User) error
	FindUserByResetToken(token string) (*models.User, error)
	FindUserByVerificationToken(token string) (*models.User, error)
	FindByProviderAndID(provider, providerID string) (*models.User, error)
}

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository() IAuthRepository {
	return &AuthRepository{db: databaseconfig.GetDB()}
}

func (r *AuthRepository) executeQuery(query *gorm.DB, operation string, fields ...zap.Field) error {
	if err := query.Error; err != nil {
		fields = append(fields, zap.Error(err))
		logconfig.Log.Error(operation+" hatası", fields...)
		return err
	}
	return nil
}

func (r *AuthRepository) findUser(query *gorm.DB, operation string, fields ...zap.Field) (*models.User, error) {
	var user models.User
	err := r.executeQuery(query.First(&user), operation, fields...)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) FindUserByEmail(email string) (*models.User, error) {
	return r.findUser(
		r.db.Where("email = ?", email),
		"Kullanıcı sorgulama (email)",
		zap.String("email", email),
	)
}

func (r *AuthRepository) FindUserByID(id uint) (*models.User, error) {
	return r.findUser(
		r.db.Where("id = ?", id),
		"Kullanıcı sorgulama (ID)",
		zap.Uint("user_id", id),
	)
}

func (r *AuthRepository) UpdateUser(ctx context.Context, user *models.User) error {
	return r.executeQuery(
		r.db.WithContext(ctx).Save(user),
		"Kullanıcı güncelleme",
		zap.Uint("user_id", user.ID),
		zap.String("email", user.Email),
	)
}

func (r *AuthRepository) CreateUser(ctx context.Context, user *models.User) error {
	return r.executeQuery(
		r.db.WithContext(ctx).Create(user),
		"Kullanıcı oluşturma",
		zap.String("email", user.Email),
	)
}

func (r *AuthRepository) FindUserByResetToken(token string) (*models.User, error) {
	return r.findUser(
		r.db.Where("reset_token = ?", token),
		"Kullanıcı sorgulama (reset token)",
		zap.String("reset_token", token),
	)
}

func (r *AuthRepository) FindUserByVerificationToken(token string) (*models.User, error) {
	return r.findUser(
		r.db.Where("verification_token = ?", token),
		"Kullanıcı sorgulama (verification token)",
		zap.String("verification_token", token),
	)
}

func (r *AuthRepository) FindByProviderAndID(provider, providerID string) (*models.User, error) {
	return r.findUser(
		r.db.Where("provider = ? AND provider_id = ?", provider, providerID),
		"Kullanıcı sorgulama (provider ve ID)",
		zap.String("provider", provider),
		zap.String("provider_id", providerID),
	)
}

var _ IAuthRepository = (*AuthRepository)(nil)
