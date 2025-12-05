package seeders

import (
	"time"

	"zatrano/configs/logconfig"
	"zatrano/models"

	"gorm.io/gorm"
)

func SeedBusinessTypes(db *gorm.DB) error {
	businessTypes := []models.BusinessType{
		{Name: "Erkek Kuaförü"},
		{Name: "Kadın Kuaförü"},
		{Name: "Unisex Kuaför"},
		{Name: "Güzellik Merkezi"},
		{Name: "Nail Art Stüdyosu"},
		{Name: "Manikür & Pedikür Salonu"},
		{Name: "Epilasyon Merkezi"},
	}

	logconfig.SLog.Info("İşletme türleri yükleniyor...")

	for _, bt := range businessTypes {
		bt.BaseModel = models.BaseModel{
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			CreatedBy: 1,
			UpdatedBy: 1,
		}

		if err := db.Create(&bt).Error; err != nil {
			logconfig.SLog.Error("İşletme türleri eklenirken hata: "+bt.Name, err)
			return err
		}
	}

	logconfig.SLog.Info("İşletme türleri yükleme tamamlandı.")
	return nil
}
