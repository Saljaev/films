package config

import (
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type (
	Config struct {
		StoragePath string `yaml:"storage_path" env-required:"true"`
		HTTPServer  `yaml:"http_server"`
		Secret      string        `yaml:"secret"`
		TokenTTL    time.Duration `yaml:"token_ttl"`
		SessionTTL  time.Duration `yaml:"session_ttl"`
	}

	HTTPServer struct {
		Address     string        `yaml:"address" env-default:"localhost:8080"`
		Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
		IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	}
)

func ConfgiLoad() *Config {
	yamlFile, err := os.ReadFile("config/config.yaml")

	if err != nil {
		log.Fatal("failed to read config.yaml")
	}

	var cfg Config

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		log.Fatalf("unmarshal failed with error: %w", err)
	}

	return &cfg
}
