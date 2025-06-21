package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

var Cfg Config

type Config struct {
	LogLevel   string           `split_words:"true" required:"true" default:"INFO"`
	HTTPServer HTTPServerConfig `split_words:"true" required:"true"`
	HLS        HLSConfig        `split_words:"true" required:"true"`
}

type HTTPServerConfig struct {
	Port        int      `split_words:"true" required:"true" default:"8080"`
	CORSOrigin  string   `split_words:"true" required:"true" `
	CORSMethods []string `split_words:"true" required:"true" `
}

type HLSConfig struct {
	OutputDir string `split_words:"true" required:"true"`
}

func Load() {
	if err := envconfig.Process("", &Cfg); err != nil {
		log.Fatal(err)
	}
}
