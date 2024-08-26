package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env         string        `yaml:"evn" env-default:"local"`
	StoragePath string        `yaml:"db_path" env-required:"true"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC        GPRCConfig    `yaml:"grpc"`
}

type GPRCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoadConfig() *Config {
	configPath := parseConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadConfigByPath(configPath)
}

func parseConfigPath() string {
	var result string

	flag.StringVar(&result, "config", "config", "path to config file")
	flag.Parse()

	if result == "" {
		result = os.Getenv("CONFIG_PATH")
	}

	return result
}

func MustLoadConfigByPath(configPath string) *Config {
	if _, err := os.Stat(configPath); err != nil {
		panic("config file not found" + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("config file read error: " + err.Error())
	}

	return &cfg
}
