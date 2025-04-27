package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Db_url            string // json: db_url
	Current_user_name string // json: current_user_name
}

func Read() (Config, error) {
	configContent, err := os.ReadFile(CONFIG_FILE_NAME)
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
	jsonBytes, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}

	os.WriteFile(".gatorconfig.json", jsonBytes, 0644)

	return nil
}
