package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func Read() (Config, error) {
	cfgFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	file, err := os.Open(cfgFilePath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	cfg := Config{}
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (cfg *Config) SetUser(name string) error {
	cfg.CurrentUserName = name
	err := write(*cfg)
	if err != nil {
		return err
	}
	return nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	cfgFilepath := filepath.Join(homeDir, configFileName)

	return cfgFilepath, nil
}

func write(cfg Config) error {
	cfgFilepath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(cfgFilepath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}

	return nil
}
