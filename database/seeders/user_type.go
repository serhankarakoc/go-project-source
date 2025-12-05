package seeders

import (
	"time"

	"zatrano/configs/logconfig"
	"zatrano/models"

	"gorm.io/gorm"
)

func SeedUserTypes(db *gorm.DB) error {
	userTypes := []models.UserType{
		{Name: "Admin"},
		{Name: "User"},
		{Name: "Business"},
		{Name: "Specialist"},
	}

	logconfig.SLog.Info("Kullanıcı tipleri yükleniyor...")

	for _, ut := range userTypes {
		ut.BaseModel = models.BaseModel{
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			CreatedBy: 1,
			UpdatedBy: 1,
		}

		if err := db.Create(&ut).Error; err != nil {
			logconfig.SLog.Error("Kullanıcı tipi eklenirken hata: "+ut.Name, err)
			return err
		}
	}

	logconfig.SLog.Info("Kullanıcı tipleri yükleme tamamlandı.")
	return nil
}
