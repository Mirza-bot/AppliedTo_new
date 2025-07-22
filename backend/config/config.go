package config

import (
	"appliedTo/controllers"
	"appliedTo/models"
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Load() {
    viper.AutomaticEnv()
    // viper.SetDefault("APP_PORT", "8080")
}

func ConnectDatabase() {
    var err error

    dsn := fmt.Sprintf("host=postgres user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=UTC",
        viper.GetString("DB_USER"),
        viper.GetString("DB_PASSWORD"),
        viper.GetString("DB_NAME"),
    )

    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Unable to connect to database %v \n", err)
    }

    fmt.Println("Connected to PostgreSQL!")

    Migrate(DB)
    controllers.SetDB(DB)
}



func Migrate(db *gorm.DB) {
    models := []any{
        &models.User{},
    }

    for _, model := range models {
        if err := db.AutoMigrate(model); err != nil {
            log.Printf("Failed to migrate model %T: %v", model, err)
        } else {
            log.Printf("Successfully migrated model %T", model)
        }
    }
}


