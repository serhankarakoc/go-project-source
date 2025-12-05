package models

type Country struct {
	BaseModel

	Name string `gorm:"type:varchar(100);not null;uniqueIndex"`

	Cities []City `gorm:"foreignKey:CountryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (Country) TableName() string {
	return "countries"
}
