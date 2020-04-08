package config

import (
	_ "github.com/joho/godotenv/autoload"
	"os"
)

type AppConfig struct {
	AppHost    string
	DBHost     string
	DBPassword string
	DBName     string
	DBUser     string
	DBDriver   string
	Dialect    string
	DBPort     string
	SSLMode    string
	JwtSecret  string
}

func LoadConfig() *AppConfig {
	return &AppConfig{
		AppHost:    os.Getenv("APP_PORT"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBUser:     os.Getenv("DB_USER"),
		DBDriver:   os.Getenv("DB_DRIVER"),
		Dialect:    os.Getenv("DB_DRIVER"),
		DBPort:     os.Getenv("DB_PORT"),
		SSLMode:    os.Getenv("SSL_MODE"),
		JwtSecret:  os.Getenv("JWT_SECRET_KEY"),
	}
}
