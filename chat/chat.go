package chat

import (
	"example.com/Quaver/Z/common"
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

// GetAvailableChannels Returns the available channels that the user is able to join
func GetAvailableChannels(userGroups common.UserGroups) []*Channel {
	mutex.Lock()
	defer mutex.Unlock()

	hasAdminChannelAccess := common.HasAnyUserGroup(userGroups, []common.UserGroups{
		common.UserGroupSwan,
		common.UserGroupDeveloper,
		common.UserGroupAdmin,
		common.UserGroupModerator,
	})

	var availableChannels []*Channel

	for _, channel := range channels {
		if !channel.AdminOnly || (channel.AdminOnly && hasAdminChannelAccess) {
			availableChannels = append(availableChannels, channel)
		}
	}

	return availableChannels
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
