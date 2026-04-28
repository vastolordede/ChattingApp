package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv              string
	AppPort             string
	DBHost              string
	DBPort              string
	DBName              string
	DBUser              string
	DBPassword          string
	DBSSLMode           string
	JWTSecret           string
	JWTExpiresHours     int
	RefreshExpiresHours int
}

func Load() (*Config, error) {
	_ = godotenv.Load("../../.env")

	cfg := &Config{
		AppEnv:     getEnv("APP_ENV", "development"),
		AppPort:    getEnv("APP_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "chatting_app"),
		DBUser:     getEnv("DB_USER", "chatting_app"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		JWTSecret:  getEnv("JWT_SECRET", "change_me_to_a_long_random_secret_key"),
	}

	jwtHoursStr := getEnv("JWT_EXPIRES_HOURS", "72")
	jwtHours, err := strconv.Atoi(jwtHoursStr)
	if err != nil {
		return nil, errors.New("JWT_EXPIRES_HOURS phải là số nguyên")
	}
	cfg.JWTExpiresHours = jwtHours

	refreshHoursStr := getEnv("REFRESH_EXPIRES_HOURS", "168")
	refreshHours, err := strconv.Atoi(refreshHoursStr)
	if err != nil {
		return nil, errors.New("REFRESH_EXPIRES_HOURS phải là số nguyên")
	}
	cfg.RefreshExpiresHours = refreshHours

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.DBHost == "" {
		return errors.New("thiếu DB_HOST")
	}
	if c.DBPort == "" {
		return errors.New("thiếu DB_PORT")
	}
	if c.DBName == "" {
		return errors.New("thiếu DB_NAME")
	}
	if c.DBUser == "" {
		return errors.New("thiếu DB_USER")
	}
	if c.JWTSecret == "" {
		return errors.New("thiếu JWT_SECRET")
	}
	if c.AppPort == "" {
		return errors.New("thiếu APP_PORT")
	}
	return nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
