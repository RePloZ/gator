package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (cfg *Config) Read() error {
	configFilePath, err := getConfigFilePath()

	if err != nil {
		return fmt.Errorf("cannot find home directory : %w", err)
	}

	fileData, err := os.ReadFile(configFilePath)
	if err != nil {
		return fmt.Errorf("cannot read file : %w", err)
	}

	return json.Unmarshal(fileData, cfg)
}

func (cfg Config) SetUser(current_user_name string) error {
	cfg.CurrentUserName = current_user_name

	if err := write(cfg); err != nil {
		return fmt.Errorf("cannot write config : %w", err)
	}
	return nil
}

func getConfigFilePath() (path string, err error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}
	path = fmt.Sprintf("%v/%v", homeDir, configFileName)
	return
}

func write(cfg Config) error {
	configFilePath, err := getConfigFilePath()

	if err != nil {
		return fmt.Errorf("cannot find home directory : %w", err)
	}

	config, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("cannot find home directory : %w", err)
	}

	if err := os.WriteFile(configFilePath, config, 0666); err != nil {
		return fmt.Errorf("cannot write file : %w", err)
	}
	return nil
}
