package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type ServerConfig struct {
	Env        string        `yaml:"env" env-required:"true"`
	Secret     string        `yaml:"secret" env-required:"true"`
	User       User          `yaml:"user" env-required:"true"`
	TokenTTL   time.Duration `yaml:"token_ttl" env-required:"true"`
	HTTPServer HTTPServer    `yaml:"http_server" env-required:"true"`
}

type HTTPServer struct {
	Port        int           `env:"SERVER_PORT" env-default:"8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	Secure      bool          `yaml:"secure" env-default:"true"`
}

type DatabaseConfig struct {
	Host     string `env:"DATABASE_HOST" env-required:"true"`
	Port     string `env:"DATABASE_PORT" env-required:"true"`
	Name     string `env:"DATABASE_NAME" env-required:"true"`
	Username string `env:"DATABASE_USER" env-required:"true"`
	Password string `env:"DATABASE_PASSWORD" env-required:"true"`
}

type User struct {
	DefaultBalance int `yaml:"default_balance" env-required:"true"`
}

func LoadServerConfig() *ServerConfig {
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	_, err := os.Stat(cfgPath)
	if os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", cfgPath)
	}

	var cfg ServerConfig
	err = cleanenv.ReadConfig(cfgPath, &cfg)
	if err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}

func LoadDatabaseConfig() *DatabaseConfig {
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	_, err := os.Stat(cfgPath)
	if os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", cfgPath)
	}

	var cfg DatabaseConfig
	err = cleanenv.ReadConfig(cfgPath, &cfg)
	if err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
