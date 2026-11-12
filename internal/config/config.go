package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	User     string
	Password string
	DBName   string
	Host     string
	Port     string
	SSLMode  string
}


func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func Load() (*Config, error) {
	cfg := &Config{
		User:     getEnv("DB_USERNAME"),
		Password: getEnv("PASSWORD"),
		DBName:   getEnv("DB"),
		Host:     getEnv("HOST"),
		Port:     getEnv("PORT"),
		SSLMode:  getEnv("SSLMODE"),
	}

	return cfg, nil
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		return ""
	}
	return value
}
