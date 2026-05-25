package config

import (
	"encoding/json"
	"errors"
	"os"
)

const CONFIG_PATH string = "bootdotdev-go-build-a-blog-aggregator/gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (self *Config) SetUser(userName string) error {
	if userName == "" {
		return errors.New("userName must not be empty")
	}

	self.CurrentUserName = userName
	bytes, err := json.Marshal(self)
	if err != nil {
		return err
	}

	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func Read() (Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return Config{}, err
	}

	bytes, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}

	var ret Config
	if err = json.Unmarshal(bytes, &ret); err != nil {
		return Config{}, err
	}

	return ret, nil
}

func getConfigPath() (string, error) {
	ret := os.Getenv("XDG_CONFIG_HOME")
	if ret == "" {
		return "", errors.New("Please set your XDG_CONFIG_HOME environment variable.")
	}

	ret += "/" + CONFIG_PATH

	return ret, nil
}
