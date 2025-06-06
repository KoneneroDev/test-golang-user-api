package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string     `yaml:"env" env-default:"local"`
	Data       Data       `yaml:"data" env-required:"true"`
	HTTPServer HTTPServer `yaml:"http_server" env-required:"true"`
}

type Data struct {
	Postgres Postgres `yaml:"postgres" env-required:"true"`
}

type Postgres struct {
	Host     string `yaml:"host" env:"POSTGRES_HOST" env-required:"true"`
	Port     string `yaml:"port" env:"POSTGRES_PORT" env-required:"true"`
	User     string `yaml:"user" env:"POSTGRES_USER" env-required:"true"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-required:"true"`
	Dbname   string `yaml:"dbname" env:"POSTGRES_DBNAME" env-required:"true"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env:"ADDRESS" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func LoadConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist %s", configPath)
	}

	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("failed to read config: %s", err)
	}

	if err := cleanenv.ReadEnv(&config); err != nil {
		log.Fatalf("failed to override config with env: %s", err)
	}

	return &config
}
