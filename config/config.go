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
	}

	Steam struct {
		AppId        int    `json:"app_id"`
		PublisherKey string `json:"publisher_key"`
	} `json:"steam"`
}

var Instance *Configuration

// Load Parses the config file into Instance
func Load() error {
	if Instance != nil {
		return fmt.Errorf("config already loaded")
	}

	data, err := os.ReadFile("./config.json")

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
