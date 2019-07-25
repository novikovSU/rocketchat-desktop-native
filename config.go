package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Config -- AAA
type Config struct {
	Server   string `json:"server"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Email    string `json:"email"`
	Password string `json:"password"`
	UseTLS   bool   `json:"use_tls,omitempty"`
	Debug    bool   `json:"debug,omitempty"`
}

var (
	config *Config
)

func getConfig(params ...string) (*Config, error) {
	var config Config

	file := "settings.json"

	if len(params) > 0 {
		file = params[0]
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Failed to read config file: %s\n", err)
	}

	err = json.Unmarshal(b, &config)
	if err != nil {
		log.Fatalf("Failed to unmarshal file: %s\n", err)
	}

	return &config, nil
}
