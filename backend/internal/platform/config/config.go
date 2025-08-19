package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort    string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	DBTimeZone string
	BcryptCost int
	DBMaxOpen  int
	DBMaxIdle  int
	DBMaxLife  time.Duration
}

func Load() Config {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// Defaults
	v.SetDefault("APP_PORT", "8080")
	v.SetDefault("DB_HOST", "db")
	v.SetDefault("DB_PORT", "5432")
	v.SetDefault("DB_SSLMODE", "disable")
	v.SetDefault("DB_TIMEZONE", "UTC")
	v.SetDefault("BCRYPT_COST", 12)
	v.SetDefault("DB_MAXOPEN", 20)
	v.SetDefault("DB_MAXIDLE", 10)
	v.SetDefault("DB_MAXLIFE", "1h")

	dur, err := time.ParseDuration(v.GetString("DB_MAXLIFE"))
	if err != nil {
		dur = time.Hour
	}

	return Config{
		AppPort:    v.GetString("APP_PORT"),
		DBHost:     v.GetString("DB_HOST"),
		DBPort:     v.GetString("DB_PORT"),
		DBUser:     v.GetString("DB_USER"),
		DBPassword: v.GetString("DB_PASSWORD"),
		DBName:     v.GetString("DB_NAME"),
		DBSSLMode:  v.GetString("DB_SSLMODE"),
		DBTimeZone: v.GetString("DB_TIMEZONE"),
		BcryptCost: v.GetInt("BCRYPT_COST"),
		DBMaxOpen:  v.GetInt("DB_MAXOPEN"),
		DBMaxIdle:  v.GetInt("DB_MAXIDLE"),
		DBMaxLife:  dur,
	}
}

