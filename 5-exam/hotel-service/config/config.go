package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Postgres PostgresConfig

	UserServiceHost string
	UserServicePort string

	ServiceHost string
	ServicePort string
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func Load(path string) Config {
	err := godotenv.Load(path + "/.env")
	if err != nil {
		log.Fatalf("Error loading.env file: %s", err)
	}
	conf := viper.New()
	conf.AutomaticEnv()

	cfg := Config{
		Postgres: PostgresConfig{
			Host:     conf.GetString("POSTGRES_HOST"),
			Port:     conf.GetString("POSTGRES_PORT"),
			User:     conf.GetString("POSTGRES_USER"),
			Password: conf.GetString("POSTGRES_PASSWORD"),
			Database: conf.GetString("POSTGRES_DATABASE"),
		},


		ServiceHost: conf.GetString("SERVICE_HOST"),
		ServicePort: conf.GetString("SERVICE_PORT"),
	}

	return cfg
}
