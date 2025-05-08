package main

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"os"
	"path/filepath"
)

type Config struct {
	AppPort string `required:"true" envconfig:"APP_PORT"`

	TmdbApiUrl      string `required:"true" envconfig:"TMDB_API_URL"`
	TmdbFilesApiUrl string `required:"true" envconfig:"TMDB_FILES_API_URL"`
	TmdbImageApiUrl string `required:"true" envconfig:"TMDB_IMAGE_API_URL"`
	TmdbApiKey      string `required:"true" envconfig:"TMDB_API_KEY"`
}

func New() *Config {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	envPath := filepath.Join(wd, ".env")

	var newCfg Config
	_ = godotenv.Load(envPath)

	if err := envconfig.Process("", &newCfg); err != nil {
		panic(err)
	}

	return &newCfg
}
