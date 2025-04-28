package config

import (
	"encoding/json"
	"errors"
	"os"
	"path"
)

type Config struct {
	Db_url            string // json: db_url
	Current_user_name string // json: current_user_name
}

func Read() (Config, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	configPath := path.Join(homePath, CONFIG_FILE_NAME)

	configContent, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(configContent, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func (c Config) SetUser() error {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := path.Join(homePath, CONFIG_FILE_NAME)

	jsonBytes, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return errors.New("failed to serialize configs to json bytes")
	}

	os.WriteFile(configPath, jsonBytes, 0644)

	return nil
}
