package config

import (
	"context"
	"fmt"
	slogGorm "github.com/orandin/slog-gorm"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"prb_care_api/internal/entity"
	"time"
)

func NewDatabase(config *viper.Viper) *gorm.DB {
	username := config.GetString("db.username")
	password := config.GetString("db.password")
	host := config.GetString("db.host")
	port := config.GetInt("db.port")
	database := config.GetString("db.name")

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable&lock_timeout=5000", username, password, host, port, database)

	logger := slogGorm.New(
		slogGorm.WithRecordNotFoundError(),
		slogGorm.WithSlowThreshold(500*time.Millisecond),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger})
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Context(context.Background())
	if err := Migrate(ctx, db); err != nil {
		log.Fatalln(err)
	}

	return db
}

func Migrate(ctx context.Context, db *gorm.DB) error {
	tx := db.WithContext(ctx).Begin()
	defer tx.Rollback()

	enumQueries := []string{
		"DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_pasien_enum') THEN CREATE TYPE status_pasien_enum AS ENUM ('aktif', 'selesai'); END IF; END $$;",
		"DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_pengambilan_obat_enum') THEN CREATE TYPE status_pengambilan_obat_enum AS ENUM ('menunggu', 'diambil', 'batal'); END IF; END $$;",
		"DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_kontrol_balik_enum') THEN CREATE TYPE status_kontrol_balik_enum AS ENUM ('menunggu', 'selesai', 'batal'); END IF; END $$;",
	}

	for _, query := range enumQueries {
		if err := tx.Exec(query).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	entities := []interface{}{
		&entity.AdminSuper{},
		&entity.AdminPuskesmas{},
		&entity.AdminApotek{},
		&entity.Pengguna{},
		&entity.Pasien{},
		&entity.Obat{},
		&entity.KontrolBalik{},
		&entity.PengambilanObat{},
		&entity.Artikel{},
		&entity.File{},
	}

	for _, e := range entities {
		if err := tx.AutoMigrate(e); err != nil {
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
