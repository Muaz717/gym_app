package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string        `yaml:"env" env-default:"local"`
	TokenTTL   time.Duration `yaml:"token_ttl" env-required:"true"`
	HTTPServer `yaml:"http_server"`
	DB         `yaml:"db"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"local"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type DB struct {
	Host       string `yaml:"host" env-required:"true"`
	DBPort     string `yaml:"port" env-required:"true"`
	Username   string `yaml:"username" env-required:"true"`
	DBName     string `yaml:"dbname" env-required:"true"`
	DBPassword string `yaml:"dbpassword" env-required:"true" env:"DB_PASSWORD"`
}

func MustLoad() *Config {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load .env file: %s", err.Error())
	}

	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		log.Fatal("config path is empty")
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exists: %s", err)
	}

	var cfg Config

	err := cleanenv.ReadConfig(cfgPath, &cfg)
	if err != nil {
		log.Fatalf("failed to read config: %s", err.Error())
	}

	return &cfg
}
