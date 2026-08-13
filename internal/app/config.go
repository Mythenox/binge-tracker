package app

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const configFileName = ".btconfig.json"

type Config struct {
	DBPath     string `json:"db_path"`
	SocketPath string `json:"socket_path"`
	ScriptPath string `json:"script_path"`
}

func ReadConfig() (Config, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	jsonFile, err := os.Open(configFilePath)
	if err != nil {
		return Config{}, err
	}
	defer jsonFile.Close()

	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("An error has occurred while attempting to read JSON data from the config")
		return Config{}, err
	}

	var cfg Config
	err = json.Unmarshal(jsonData, &cfg)
	if err != nil {
		fmt.Println("An error has occurred while attempting to unmarshal JSON data")
		return Config{}, err
	}

	return cfg, nil
}

func getConfigFilePath() (string, error) {
	/*homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configFilePath := filepath.Join(homedir, configFileName)
	if len(configFilePath) == 0 {
		fmt.Println("An error has occurred while attempting to construct the config filepath")
		return "", fmt.Errorf("An error has occurred while attempting to construct the config filepath")
	}

	return configFilePath, nil*/
	configFilePath := filepath.Join(".", configFileName)
	return configFilePath, nil
}

func writeConfig(cfg Config) error {
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		fmt.Println("An error has occurred while attempting unmarshal JSON data")
		return err
	}

	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	err = os.WriteFile(configFilePath, jsonData, 0644)
	if err != nil {
		fmt.Printf("An error has occurred while attempting to write to the config file: %v\n", err)
		return err
	}

	return nil
}
