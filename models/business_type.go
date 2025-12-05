package models

type BusinessType struct {
	BaseModel
	Name string `gorm:"size:50;unique;not null;index"`
	Icon string `gorm:"type:varchar(50);not null"`

	Businesses []Business `gorm:"foreignKey:BusinessTypeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
