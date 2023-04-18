package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Configuration struct {
	Server struct {
		Port int `json:"port"`
	} `json:"server"`

	SQL struct {
		Host     string `json:"host"`
		Username string `json:"username"`
		Password string `json:"password"`
		Database string `json:"database"`
	} `json:"sql"`

	Steam struct {
		AppId        int    `json:"app_id"`
		PublisherKey string `json:"publisher_key"`
	} `json:"steam"`
}

var Instance *Configuration

// Load Parses the config file into Instance
func Load(path string) error {
	if Instance != nil {
		return fmt.Errorf("config already loaded")
	}

	data, err := os.ReadFile(path)

	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &Instance)

	if err != nil {
		return err
	}

	log.Println("Config file has been successfully read")
	return nil
}
