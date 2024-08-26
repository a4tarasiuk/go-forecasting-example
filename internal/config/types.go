package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type (
	dbConfig struct {
		Host string `required:"true"`
		Port int    `required:"true"`

		User     string `required:"true"`
		Password string `required:"true"`

		Name string `required:"true"`

		Schema string `required:"false" default:"postgresql"`
	}

	AppConfig struct {
		DB dbConfig
	}
)

func InitConfig() (cfg AppConfig, err error) {
	err = godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file", err)
		return AppConfig{}, err
	}

	port, _err := strconv.Atoi(os.Getenv("DB_PORT"))
	if _err != nil {
		return AppConfig{}, _err
	}

	cfg = AppConfig{
		DB: dbConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     port,
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			Schema:   os.Getenv("DB_SCHEMA"),
		},
	}

	return cfg, nil
}
