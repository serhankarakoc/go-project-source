package models

type Service struct {
	BaseModel

	Name string `gorm:"type:varchar(100);not null;uniqueIndex"`
}

func (Service) TableName() string {
	return "services"
}
