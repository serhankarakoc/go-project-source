package models

type District struct {
	BaseModel

	CityID uint   `gorm:"index;not null"`
	Name   string `gorm:"type:varchar(100);not null;index"`

	City *City `gorm:"foreignKey:CityID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (District) TableName() string {
	return "districts"
}
