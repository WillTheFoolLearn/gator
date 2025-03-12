package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

const (
	configFileName = ".gatorconfig.json"
)

func Read() (Config, error) {
	dir, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	jsonFile, err := os.Open(dir)

	if err != nil {
		fmt.Println("Error with opening json file")
		return Config{}, err
	}
	defer jsonFile.Close()

	data, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error with ReadAll file")
		return Config{}, err
	}

	var cfg Config

	json.Unmarshal(data, &cfg)

	return cfg, nil
}

func getConfigFilePath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("can't get directory")
	}

	dir += "/" + configFileName

	return dir, nil
}

func (cfg *Config) SetUser(name string) error {
	cfg.CurrentUserName = name
	return write(*cfg)
}

func write(cfg Config) error {
	jsonString, _ := json.Marshal((cfg))

	dir, _ := getConfigFilePath()

	err := os.WriteFile(dir, jsonString, 0777)
	if err != nil {
		return errors.New("couldn't write to file")
	}

	return nil
}

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}
