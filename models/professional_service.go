package models

type ProfessionalService struct {
	BaseModel

	ProfessionalID uint    `gorm:"index;not null"`
	ServiceID      uint    `gorm:"index;not null"`
	Duration       uint    `gorm:"not null"`
	Price          float64 `gorm:"not null"`

	Professional *Professional `gorm:"foreignKey:ProfessionalID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Service      *Service      `gorm:"foreignKey:ServiceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (ProfessionalService) TableName() string {
	return "professional_services"
}
