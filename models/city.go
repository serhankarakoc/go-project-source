package models

type City struct {
	BaseModel

	CountryID uint   `gorm:"index;not null"`
	Name      string `gorm:"type:varchar(100);not null;index"`

	Country   *Country   `gorm:"foreignKey:CountryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Districts []District `gorm:"foreignKey:CityID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (City) TableName() string {
	return "cities"
}
