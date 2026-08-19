package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mythenox/bingetracker/internal/database"
	_ "modernc.org/sqlite"
)

const configFileName = ".btconfig.json"

type State struct {
	DB  *sql.DB
	Q   *database.Queries
	Cfg *Config
}

type Config struct {
	DBPath               string `json:"db_path"`
	SocketPath           string `json:"socket_path"`
	PipePath             string `json:"pipe_path"`
	ScriptPath           string `json:"script_path"`
	CountPartialProgress bool   `json:"count_partial_progress"`
	VideoPlayer          string `json:"video_player"`
}

func (s *State) LoadConfig() error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	jsonFile, err := os.Open(configFilePath)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("An error has occurred while attempting to read JSON data from the config")
		return err
	}

	var cfg Config
	err = json.Unmarshal(jsonData, &cfg)
	if err != nil {
		fmt.Println("An error has occurred while attempting to unmarshal JSON data")
		return err
	}

	s.Cfg = &cfg
	return nil
}

func (s *State) ConnectDB() error {
	db, err := sql.Open("sqlite", s.Cfg.DBPath)
	if err != nil {
		return err
	}

	s.DB = db
	s.Q = database.New(db)

	return nil
}

func getConfigFilePath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configFilePath := filepath.Join(homedir, configFileName)
	if len(configFilePath) == 0 {
		fmt.Println("An error has occurred while attempting to construct the config filepath")
		return "", fmt.Errorf("An error has occurred while attempting to construct the config filepath")
	}

	return configFilePath, nil
}

func (s *State) WriteConfig() error {
	jsonData, err := json.Marshal(s.Cfg)
	if err != nil {
		fmt.Println("An error has occurred while attempting to unmarshal JSON data")
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
