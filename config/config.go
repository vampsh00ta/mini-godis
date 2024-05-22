package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type (
	// Config -.
	Config struct {
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		Log  `yaml:"logger"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}
)

func New() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}
	currPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	filePath := currPath + os.Getenv("path") + "/" + os.Getenv("env") + ".yml"
	fmt.Println(filePath)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	var cfg *Config

	if err := d.Decode(&cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
