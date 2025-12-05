package seeders

import (
	"fmt"

	"zatrano/configs/logconfig"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SeedAll fonksiyonu, tüm seed işlemlerini doğru sırayla çalıştırır.
func SeedAll(db *gorm.DB) error {
	logconfig.SLog.Info("Tüm seed işlemleri başlatılıyor...")

	// Sıralama: foreign key bağımlılıklarına göre belirlenmiştir
	seedSteps := []struct {
		Name string
		Func func(*gorm.DB) error
	}{
		{"UserTypes", SeedUserTypes},
		{"SystemUser", SeedSystemUser},
		{"Countries", SeedCountries},
		{"Cities", SeedCities},
		{"Districts", SeedDistricts},
		{"BusinessTypes", SeedBusinessTypes},
		{"Services", SeedServices},
	}

	for _, step := range seedSteps {
		logconfig.SLog.Info(fmt.Sprintf("%s verileri seed ediliyor...", step.Name))
		if err := step.Func(db); err != nil {
			logconfig.Log.Error("Seed işlemi başarısız",
				zap.String("step", step.Name),
				zap.Error(err),
			)
			return err
		}
		logconfig.SLog.Info(fmt.Sprintf("%s verileri başarıyla eklendi.", step.Name))
	}

	logconfig.SLog.Info("Tüm seed işlemleri başarıyla tamamlandı.")
	return nil
}
