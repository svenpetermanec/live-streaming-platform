package config

import (
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"
)

var Cfg Config

type Config struct {
	LogLevel string    `split_words:"true" required:"true" default:"INFO"`
	SRT      SRTConfig `split_words:"true" required:"true"`
	HLS      HLSConfig `split_words:"true" required:"true"`
}

type SRTConfig struct {
	TransType          string `split_words:"true" required:"true" default:"live"`
	Port               int    `split_words:"true" required:"true" default:"5270"`
	BacklogConnections int    `split_words:"true" required:"true" default:"100"`
}

type HLSConfig struct {
	OutputDir   string `split_words:"true" required:"true"`
	Resolutions []Resolution
}

type Resolution struct {
	Name      string `split_words:"true" required:"true"`
	Width     int    `split_words:"true" required:"true"`
	Height    int    `split_words:"true" required:"true"`
	Framerate int    `split_words:"true" required:"true"`
	Bitrate   string `split_words:"true" required:"true"`
}

func Load() {
	if err := envconfig.Process("", &Cfg); err != nil {
		log.Fatal(err)
	}

	var resolutions []Resolution
	for i := 0; ; i++ {
		var r Resolution
		prefix := fmt.Sprintf("HLS_RESOLUTION_%d", i)
		err := envconfig.Process(prefix, &r)
		if err != nil {
			fmt.Println(err)
			break
		}
		if r.Name == "" {
			break
		}
		resolutions = append(resolutions, r)
	}
	Cfg.HLS.Resolutions = resolutions
}
