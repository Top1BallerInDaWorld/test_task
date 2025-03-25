package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log/slog"
)

type ConfigDatabase struct {
	ConnectionString string `yaml:"connection_string" env:"CONNECTION_STRING" env-default:"postgres://user:password@localhost:5432/postgres?sslmode=disable"`
	Port             string `yaml:"port" env:"PORT" env-default:":8080"`
}

func GetConfig() ConfigDatabase {
	var cfg ConfigDatabase
	err := cleanenv.ReadConfig("config.yml", &cfg)
	if err != nil {
		slog.Error("Error reading config file", "error", err)
	}
	return cfg
}
