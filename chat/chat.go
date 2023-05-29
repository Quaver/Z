package chat

import (
	"example.com/Quaver/Z/common"
	"example.com/Quaver/Z/config"
	"example.com/Quaver/Z/sessions"
	"log"
	"sync"
)

var (
	channels  map[string]*Channel
	chatMutex *sync.Mutex
)

// Initialize Initializes the chat channels
func Initialize() {
	channels = make(map[string]*Channel)
	chatMutex = &sync.Mutex{}

	for _, channel := range config.Instance.ChatChannels {
		addChannel(NewChannel(channel.Name, channel.Description, channel.AdminOnly, channel.AutoJoin, channel.DiscordWebhook))
	}
}

// GetAvailableChannels Returns the available channels that the user is able to join
func GetAvailableChannels(userGroups common.UserGroups) []*Channel {
	chatMutex.Lock()
	defer chatMutex.Unlock()

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

// SendPublicMessage Sends a message to a public chat channel
func SendPublicMessage(sender *sessions.User, channel *Channel, message string) {
	channel.SendMessage(sender, message)
}

// SendPrivateMessage Sends a private message to a user
func SendPrivateMessage(sender *sessions.User, receiver *sessions.User, message string) {
	// TODO: Send Discord Webhook
	// TODO: Log In Database
}

// Adds a channel to channels
func addChannel(channel *Channel) {
	chatMutex.Lock()
	defer chatMutex.Unlock()

	channels[channel.Name] = channel
	log.Printf("Initialized chat channel: %v\n", channel.Name)
}

// Removes a channel from channels
func removeChannel(channel *Channel) {
	chatMutex.Lock()
	defer chatMutex.Unlock()

	channel.removeAllUsers()
	delete(channels, channel.Name)
}

// RemoveUserFromAllChannels Removes a user from every single channel if they are in them
func RemoveUserFromAllChannels(user *sessions.User) {
	chatMutex.Lock()
	defer chatMutex.Unlock()

	for _, channel := range channels {
		channel.RemoveUser(user)
	}
}
