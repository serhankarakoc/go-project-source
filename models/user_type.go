package models

type UserType struct {
	BaseModel
	Name string `gorm:"size:50;unique;not null;index"`
}
