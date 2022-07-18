package config

import (
	"github.com/yosuke-furukawa/json5/encoding/json5"
	"log"
	"os"
)

func LoadConfig() *Config {
	configFile, err := os.Open("private/config.json5")
	if err != nil {
		log.Panicln(err)
	}
	configDecoder := json5.NewDecoder(configFile)
	cfg := new(Config)
	err = configDecoder.Decode(&cfg)
	if err != nil {
		log.Panicln(err)
	}
	return cfg
}

// TODO parse into usable struct
