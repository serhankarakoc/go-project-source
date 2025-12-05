package models

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	BaseModel
	Name              string `gorm:"size:100;not null;index"`
	Email             string `gorm:"size:100;unique;not null"`
	Password          string `gorm:"size:255;not null"`
	UserTypeID        uint   `gorm:"index"`
	ResetToken        string `gorm:"size:255;index"`
	EmailVerified     bool   `gorm:"default:false;index"`
	VerificationToken string `gorm:"size:255;index"`
	Provider          string `gorm:"size:50;index"`
	ProviderID        string `gorm:"size:100;index"`

	UserType UserType `gorm:"foreignKey:UserTypeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// SetPassword - Kullanıcının şifresini hashler ve set eder
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword - Kullanıcının şifresini doğrular
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
