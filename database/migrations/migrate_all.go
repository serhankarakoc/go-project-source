package migrations

import (
	"fmt"
	"strings"

	"zatrano/configs/logconfig"
	"zatrano/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// MigrateAll fonksiyonu tüm tabloları doğru sırayla migrate eder.
func MigrateAll(db *gorm.DB) error {
	logconfig.SLog.Info("Tüm migrasyon işlemleri başlatılıyor...")

	// Migrasyon sırası (foreign key ilişkilerine göre)
	modelsToMigrate := []interface{}{
		&models.UserType{},
		&models.User{},
		&models.Country{},
		&models.City{},
		&models.District{},
		&models.Address{},
		&models.BusinessType{},
		&models.Business{},
		&models.Account{},
		&models.Professional{},
		&models.ProfessionalService{},
		&models.Gallery{},
		&models.Service{},
	}

	for _, model := range modelsToMigrate {
		tableName := modelName(model)
		logconfig.SLog.Info(fmt.Sprintf("%s tablosu migrate ediliyor...", tableName))

		if err := db.AutoMigrate(model); err != nil {
			logconfig.Log.Error("Migrasyon hatası",
				zap.String("model", tableName),
				zap.Error(err),
			)
			return err
		}

		logconfig.SLog.Info(fmt.Sprintf("%s tablosu migrate edildi.", tableName))
	}

	logconfig.SLog.Info("Tüm migrasyon işlemleri başarıyla tamamlandı.")
	return nil
}

// modelName fonksiyonu struct tipinin adını çözer
func modelName(m interface{}) string {
	typeName := fmt.Sprintf("%T", m)
	if typeName == "" {
		return "Unknown"
	}
	parts := strings.Split(typeName, ".")
	return parts[len(parts)-1]
}
