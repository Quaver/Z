package chat

import (
	"example.com/Quaver/Z/common"
	"example.com/Quaver/Z/config"
	"example.com/Quaver/Z/db"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
	"example.com/Quaver/Z/webhooks"
	"github.com/disgoorg/disgo/webhook"
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

// SendMessage Sends a message to a given a receiver
func SendMessage(sender *sessions.User, receiver string, message string) {
	chatMutex.Lock()
	defer chatMutex.Unlock()

	if receiver == "" || message == "" {
		return
	}

	if sender.IsMuted() {
		return
	}

	// TODO: Track spam messages and mute user for spamming
	// TODO: Censor message

	var discordWebhook webhook.Client

	if receiver[0] == '#' {
		channel := channels[receiver]

		if channel == nil {
			return
		}

		sendPublicMessage(sender, channel, message)
		discordWebhook = channel.WebhookClient
	} else {
		receivingUser := sessions.GetUserByUsername(receiver)

		if receivingUser == nil {
			return
		}

		sendPrivateMessage(sender, receivingUser, message)
		discordWebhook = webhooks.PrivateChat
	}

	if discordWebhook != nil {
		webhooks.SendChatMessage(discordWebhook, sender.Info.Username, sender.Info.GetProfileUrl(), sender.Info.AvatarUrl, receiver, message)
	}
}

// Sends a message to a public chat channel
func sendPublicMessage(sender *sessions.User, channel *Channel, message string) {
	channel.SendMessage(sender, message)
}

// Sends a private message to a user
func sendPrivateMessage(sender *sessions.User, receiver *sessions.User, message string) {
	sessions.SendPacketToUser(packets.NewServerChatMessage(sender.Info.Id, sender.Info.Username, receiver.Info.Username, message), receiver)

	err := db.InsertPrivateChatMessage(sender.Info.Id, receiver.Info.Id, receiver.Info.Username, message)

	if err != nil {
		log.Printf("Error inserting private chat into DB: %v\n", err)
		return
	}
}

// RemoveUserFromAllChannels Removes a user from every single channel if they are in them
func RemoveUserFromAllChannels(user *sessions.User) {
	chatMutex.Lock()
	defer chatMutex.Unlock()

	for _, channel := range channels {
		channel.RemoveUser(user)
	}
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
