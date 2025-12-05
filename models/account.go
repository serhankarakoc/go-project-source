package models

type Account struct {
	BaseModel

	UserID    uint   `gorm:"index;not null"`
	Gsm       string `gorm:"type:varchar(20)"`
	TCKN      string `gorm:"type:varchar(11)"`
	Photo     string `gorm:"type:varchar(255)"`
	AddressID uint   `gorm:"index;not null"`

	User    *User    `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Address *Address `gorm:"foreignKey:AddressID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (Account) TableName() string {
	return "accounts"
}
