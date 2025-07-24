package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DBDsn    string `env:"DB_DSN" env-required:"true"`
	HTTPPort string `env:"HTTP_PORT" env-default:"8080"`
}

func LoadConfig() Config {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("config load: %v", err)
	}
	return cfg
}
