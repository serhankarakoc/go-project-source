package seeders

import (
	"time"

	"zatrano/configs/logconfig"
	"zatrano/models"

	"gorm.io/gorm"
)

func SeedServices(db *gorm.DB) error {
	logconfig.SLog.Info("Servisler yükleniyor...")

	services := []models.Service{
		{
			Name: "Türkiye",
		},
	}

	for _, service := range services {
		service.BaseModel = models.BaseModel{
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			CreatedBy: 1,
			UpdatedBy: 1,
		}

		if err := db.Create(&service).Error; err != nil {
			logconfig.SLog.Error("Servisler eklenirken hata: "+service.Name, err)
			return err
		}

	}

	logconfig.SLog.Info("Servisler yükleme tamamlandı.")
	return nil
}
