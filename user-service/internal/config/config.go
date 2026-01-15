package config
import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	App      AppConfig
	DB       DBConfig
	JWT      JWTConfig
}

type AppConfig struct {
	Port string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessTTL     int
	RefreshTTL    int 
}

func Load() *Config {
	return &Config{
		App: AppConfig{
			Port: getEnv("APP_PORT", "8080"),
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "user_service"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			AccessSecret:  mustGetEnv("JWT_ACCESS_SECRET"),
			RefreshSecret: mustGetEnv("JWT_REFRESH_SECRET"),
			AccessTTL:     getEnvAsInt("JWT_ACCESS_TTL", 15),
			RefreshTTL:    getEnvAsInt("JWT_REFRESH_TTL", 7),
		},
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("ENV %s is required", key)
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Fatalf("ENV %s must be integer", key)
	}

	return value
}
