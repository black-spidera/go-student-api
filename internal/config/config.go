package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServerConfig struct {
	Addr string `yaml:"addr" env:"HTTP_SERVER_ADDR" env-required:"true" env-default:"localhost:8080"`
}
type Config struct {
	Env         string           `yaml:"env" env:"ENV" env-required:"true" env-default:"production"`
	StoragePath string           `yaml:"storage_path" env:"STORAGE_PATH" env-required:"true" env-default:"storage/storage.db"`
	HTTPServer  HTTPServerConfig `yaml:"http_server" env-prefix:"HTTP_SERVER_" env-required:"true"`
}

func ConfigLoader() *Config {
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "config/local.yml", "Path to the configuration file")
		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatal("CONFIG_PATH environment variable or --config flag must be set")
		}

	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Configuration file does not exist at path: %s", configPath)
	}

	config := &Config{}
	err := cleanenv.ReadConfig(configPath, config)
	if err != nil {
		log.Fatalf("Failed to read configuration file: %v", err)
	}

	return config
}
