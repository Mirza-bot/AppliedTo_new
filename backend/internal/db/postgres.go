package db

import (
	"fmt"
	"log"

	"appliedTo/models"
	"appliedTo/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode, cfg.DBTimeZone,
	)
	g, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := g.DB()
	if err == nil {
		if cfg.DBMaxOpen > 0 {
			sqlDB.SetMaxOpenConns(cfg.DBMaxOpen)
		}
		if cfg.DBMaxIdle > 0 {
			sqlDB.SetMaxIdleConns(cfg.DBMaxIdle)
		}
		if cfg.DBMaxLife > 0 {
			sqlDB.SetConnMaxLifetime(cfg.DBMaxLife)
		}
	}

	return g, nil
}

func Migrate(g *gorm.DB) error {
	entities := []any{
		&models.User{},
		&models.JobApplication{},
		&models.Employment{},
		&models.SalaryRange{},
	}
	for _, e := range entities {
		if err := g.AutoMigrate(e); err != nil {
			return fmt.Errorf("migrate %T: %w", e, err)
		}
		log.Printf("migrated %T", e)
	}
	return nil
}
