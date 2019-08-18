package settings

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/novikovSU/rocketchat-desktop-native/utils"
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
	Conf         *Config
	settingsFile = "settings.json"
)

func GetConfig(params ...string) (*Config, error) {
	var config Config

	if len(params) > 0 {
		settingsFile = params[0]
	}

	var _, err = os.Stat(settingsFile)
	if os.IsNotExist(err) {
		return nil, errors.New("Config file does not exists")
	}

	b, err := ioutil.ReadFile(settingsFile)
	utils.AssertErrMsg(err, "Failed to read config file: %s\n")

	err = json.Unmarshal(b, &config)
	utils.AssertErrMsg(err, "Failed to unmarshal file: %s\n")

	return &config, nil
}

func CreateDefaultConfig() *Config {
	return &Config{UseTLS: true}
}

func StoreConfig(config *Config) error {
	confContent, err := json.MarshalIndent(config, "", " ")
	if err == nil {
		err = ioutil.WriteFile(settingsFile, confContent, 0644)
	}

	return err
}
