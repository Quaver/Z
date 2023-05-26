package chat

import (
	"example.com/Quaver/Z/config"
	"log"
	"sync"
)

var (
	channels map[string]*Channel
	mutex    *sync.Mutex
)

// Initialize Initializes the chat channels
func Initialize() {
	channels = make(map[string]*Channel)
	mutex = &sync.Mutex{}

	for _, channel := range config.Instance.ChatChannels {
		addChannel(NewChannel(channel.Name, channel.Description, channel.AdminOnly, channel.AutoJoin, channel.DiscordWebhook))
	}
}

// Adds a channel to channels
func addChannel(channel *Channel) {
	mutex.Lock()
	defer mutex.Unlock()

	channels[channel.Name] = channel
	log.Printf("Initialized chat channel: %v\n", channel.Name)
}

// Removes a channel from channels
func removeChannel(channel *Channel) {
	mutex.Lock()
	defer mutex.Unlock()

	// TODO: Remove all users from the channel
	delete(channels, channel.Name)
}
