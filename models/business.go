package models

type Business struct {
	BaseModel

	UserID         uint   `gorm:"index;not null"`
	BusinessTypeID uint   `gorm:"index;not null"`
	Capacity       uint   `gorm:"not null;default:0"`
	Slug           string `gorm:"type:varchar(150);uniqueIndex;not null"`
	Title          string `gorm:"type:varchar(255)"`
	Description    string `gorm:"type:text"`
	Gsm            string `gorm:"type:varchar(20)"`
	Telephone      string `gorm:"type:varchar(20)"`
	Email          string `gorm:"type:varchar(100)"`
	Website        string `gorm:"type:varchar(255)"`

	TaxOffice  string `gorm:"type:varchar(255)"`
	TaxNumber  string `gorm:"type:varchar(50)"`
	KEPAddress string `gorm:"type:varchar(255)"`
	MersisNo   string `gorm:"type:varchar(50)"`
	IbanNo     string `gorm:"type:varchar(34)"`

	AddressID uint   `gorm:"index;not null"`
	Map       string `gorm:"type:text"`

	Logo   string `gorm:"type:varchar(255)"`
	Banner string `gorm:"type:varchar(255)"`
	Video  string `gorm:"type:varchar(255)"`

	Whatapp   string `gorm:"type:varchar(20)"`
	Instagram string `gorm:"type:varchar(255)"`
	Facebook  string `gorm:"type:varchar(255)"`
	Twitter   string `gorm:"type:varchar(255)"`
	Linkedin  string `gorm:"type:varchar(255)"`
	Youtube   string `gorm:"type:varchar(255)"`
	Tiktok    string `gorm:"type:varchar(255)"`

	Galleries            []Gallery             `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Professionals        []Professional        `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ProfessionalServices []ProfessionalService `gorm:"-"`

	User         *User         `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	BusinessType *BusinessType `gorm:"foreignKey:BusinessTypeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Address      *Address      `gorm:"foreignKey:AddressID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (Business) TableName() string {
	return "businesses"
}
