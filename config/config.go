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

	BypassSteamLogin bool `json:"bypass_steam_login"`

	SQL struct {
		Host     string `json:"host"`
		Username string `json:"username"`
		Password string `json:"password"`
		Database string `json:"database"`
	} `json:"sql"`

	Redis struct {
		Host     string `json:"host"`
		Password string `json:"password"`
		Database int    `json:"database"`
	}

	Steam struct {
		AppId        int    `json:"app_id"`
		PublisherKey string `json:"publisher_key"`
	} `json:"steam"`

	DiscordWebhooks struct {
		AntiCheat   string `json:"anti_cheat"`
		PrivateChat string `json:"private_chat"`
		Multiplayer string `json:"multiplayer"`
		Spectator   string `json:"spectator"`
	} `json:"discord_webhooks"`

	ChatChannels []struct {
		Name           string `json:"name"`
		Description    string `json:"description"`
		AdminOnly      bool   `json:"admin_only"`
		AutoJoin       bool   `json:"auto_join"`
		DiscordWebhook string `json:"discord_webhook"`
		LimitedChat    bool   `json:"limited_chat"`
	} `json:"chat_channels"`
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
