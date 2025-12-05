package models

type Gallery struct {
	BaseModel

	BusinessID uint   `gorm:"index;not null"`             // Hangi işletmeye ait
	Image      string `gorm:"type:varchar(255);not null"` // Görselin yolu

	Type string `gorm:"type:varchar(50);not null"` // "logo", "banner", "gallery" gibi

	Business *Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (Gallery) TableName() string {
	return "galleries"
}
