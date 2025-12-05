package models

type Professional struct {
	BaseModel

	UserID      uint   `gorm:"index;not null"`
	BusinessID  uint   `gorm:"index;not null"`
	Title       string `gorm:"type:varchar(255)"`
	Image       string `gorm:"type:varchar(255)"`
	Description string `gorm:"type:text"`

	// İlişkiler
	User     *User     `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Business *Business `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (Professional) TableName() string {
	return "professionals"
}
