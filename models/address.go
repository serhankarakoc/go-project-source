package models

type Address struct {
	BaseModel

	CountryID  uint   `gorm:"index"`
	CityID     uint   `gorm:"index"`
	DistrictID uint   `gorm:"index"`
	Address    string `gorm:"type:varchar(255)"`

	Country  *Country  `gorm:"foreignKey:CountryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	City     *City     `gorm:"foreignKey:CityID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	District *District `gorm:"foreignKey:DistrictID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (Address) TableName() string {
	return "addresses"
}
