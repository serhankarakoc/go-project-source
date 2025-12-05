package main

import (
	"flag"
	"log"

	"zatrano/configs/databaseconfig"
	"zatrano/configs/envconfig"
	"zatrano/configs/logconfig"
	"zatrano/database"
)

func main() {
	// ⬇️ Dev ortamdaysa .env dosyasını otomatik yükle
	envconfig.LoadIfDev()

	// Logger başlat
	logconfig.InitLogger()
	defer logconfig.SyncLogger()

	// Komut satırı parametreleri
	migrateFlag := flag.Bool("migrate", false, "Veritabanı migrasyonlarını çalıştır")
	seedFlag := flag.Bool("seed", false, "Veritabanı seederlarını çalıştır")
	flag.Parse()

	// DB başlat
	databaseconfig.InitDB()
	defer func() {
		if err := databaseconfig.CloseDB(); err != nil {
			log.Println("Database kapanırken hata:", err)
		}
	}()

	db := databaseconfig.GetDB()

	logconfig.SLog.Infow("Veritabanı başlatma işlemi çalıştırılıyor",
		"migrate", *migrateFlag,
		"seed", *seedFlag,
	)

	// Migrasyon ve seed işlemleri
	database.Initialize(db, *migrateFlag, *seedFlag)

	logconfig.SLog.Info("Veritabanı başlatma işlemi tamamlandı.")
}
