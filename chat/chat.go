package chat

import (
	"example.com/Quaver/Z/common"
	"example.com/Quaver/Z/config"
	"example.com/Quaver/Z/db"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
	"example.com/Quaver/Z/utils"
	"example.com/Quaver/Z/webhooks"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

var (
	channels               map[string]*Channel
	chatMutex              *sync.Mutex
	publicMessageHandlers  []func(user *sessions.User, channel *Channel, args []string) string
	privateMessageHandlers []func(user *sessions.User, receiver *sessions.User, args []string) string
)

// Initialize Initializes the chat channels
func Initialize() {
	channels = make(map[string]*Channel)
	chatMutex = &sync.Mutex{}
	publicMessageHandlers = []func(user *sessions.User, channel *Channel, args []string) string{}
	privateMessageHandlers = []func(user *sessions.User, receiver *sessions.User, args []string) string{}

	for _, channel := range config.Instance.ChatChannels {
		addChannel(NewChannel(ChannelNormal, channel.Name, channel.Description, channel.AdminOnly, channel.AutoJoin, channel.LimitedChat, channel.DiscordWebhook))
	}

	_ = sessions.AddUser(Bot)
	addBotChatHandlers()
	addSpectatorHandlers()
}

// GetAvailableChannels Returns the available channels that the user is able to join
func GetAvailableChannels(userGroups common.UserGroups) []*Channel {
	chatMutex.Lock()
	defer chatMutex.Unlock()

	var availableChannels []*Channel

	for _, channel := range channels {
		if (channel.Type != ChannelTypeMultiplayer && channel.Type != ChannelTypeSpectator && !channel.AdminOnly) || (channel.AdminOnly && isChatModerator(userGroups)) {
			availableChannels = append(availableChannels, channel)
		}
	}

	return availableChannels
}

// GetChannelByName Gets a chat channel by its name
func GetChannelByName(name string) *Channel {
	chatMutex.Lock()
	defer chatMutex.Unlock()

	if _, ok := channels[name]; ok {
		return channels[name]
	}

	return nil
}

// SendMessage Sends a message to a given a receiver
func SendMessage(sender *sessions.User, receiver string, message string) {
	if receiver == "" || message == "" {
		return
	}

	if sender.IsMuted() {
		return
	}

	sender.IncrementSpammedMessagesCount()

	if sender.GetSpammedMessagesCount() >= 10 && !isChatModerator(sender.Info.UserGroups) {
		_ = sender.MuteUser(time.Minute * 30)
		return
	}

	message = utils.TruncateString(message, 500)

	if censored, err := utils.CensorString(message); err == nil {
		message = censored
	} else {
		log.Printf("Error censoring chat message string - %v - %v\n", message, err)
	}

	if receiver[0] == '#' {
		channel := GetChannelByName(receiver)

		if channel == nil {
			return
		}

		if channel.LimitedChat && !isChatModerator(sender.Info.UserGroups) {
			return
		}

		channel.SendMessage(sender, message)
		webhooks.SendChatMessage(channel.WebhookClient, sender.Info.Username, sender.Info.GetProfileUrl(), sender.Info.AvatarUrl.String, receiver, message)
		runPublicMessageHandlers(sender, channel, message)
	} else {
		receivingUser := sessions.GetUserByUsername(receiver)

		if receivingUser == nil {
			return
		}

		sendPrivateMessage(sender, receivingUser, message)
		webhooks.SendChatMessage(webhooks.PrivateChat, sender.Info.Username, sender.Info.GetProfileUrl(), sender.Info.AvatarUrl.String, receiver, message)
		runPrivateMessageHandlers(sender, receivingUser, message)
	}
}

// AddMultiplayerChannel Adds a multiplayer channel.
func AddMultiplayerChannel(id string) *Channel {
	channel := NewChannel(ChannelTypeMultiplayer, fmt.Sprintf("#multiplayer_%v", id), "", false, false,
		false, config.Instance.DiscordWebhooks.Multiplayer)

	addChannel(channel)
	return channel
}

// RemoveMultiplayerChannel Removes a multiplayer channel
func RemoveMultiplayerChannel(id string) {
	channel := GetChannelByName(fmt.Sprintf("#multiplayer_%v", id))

	if channel == nil {
		return
	}

	removeChannel(channel)
}

// GetSpectatorChannel Returns a user's spectator channel
func GetSpectatorChannel(userId int) *Channel {
	return GetChannelByName(getSpectatorChannelName(userId))
}

// AddSpectatorChannel Adds a spectator channel for a user
func AddSpectatorChannel(userId int) *Channel {
	channel := NewChannel(ChannelTypeSpectator, getSpectatorChannelName(userId), "", false, false,
		false, config.Instance.DiscordWebhooks.Spectator)

	addChannel(channel)

	return channel
}

// RemoveSpectatorChannel Removes a user's spectator channel
func RemoveSpectatorChannel(userId int) {
	channel := GetChannelByName(getSpectatorChannelName(userId))

	if channel == nil {
		return
	}

	removeChannel(channel)
}

// AddPublicMessageHandler Adds a message handler for public chat channels
func AddPublicMessageHandler(f func(user *sessions.User, channel *Channel, args []string) string) {
	chatMutex.Lock()
	defer chatMutex.Unlock()

	publicMessageHandlers = append(publicMessageHandlers, f)
}

// AddPrivateMessageHandler Adds a message handler for private chats
func AddPrivateMessageHandler(f func(user *sessions.User, receivingUser *sessions.User, args []string) string) {
	chatMutex.Lock()
	defer chatMutex.Unlock()

	privateMessageHandlers = append(privateMessageHandlers, f)
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

// Runs all the chat message handlers for a given public channel message.
// Ran in separate goroutine due to possible chat deadlocks.
func runPublicMessageHandlers(sender *sessions.User, channel *Channel, message string) {
	go func() {
		for _, handler := range publicMessageHandlers {
			responseMsg := handler(sender, channel, strings.Split(message, " "))

			if responseMsg != "" {
				channel.SendMessage(Bot, responseMsg)
			}
		}
	}()
}

// Runs all the private message handlers for a given private channel message.
// Ran in separate goroutine due to separate chat deadlocks
func runPrivateMessageHandlers(sender *sessions.User, receivingUser *sessions.User, message string) {
	go func() {
		for _, handler := range privateMessageHandlers {
			responseMsg := handler(sender, receivingUser, strings.Split(message, " "))

			if responseMsg != "" {
				sendPrivateMessage(Bot, sender, responseMsg)
			}
		}
	}()
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
	log.Printf("Uninitialized chat channel: %v\n", channel.Name)
}

// Adds handlers for when a user begins/stops spectating someone
func addSpectatorHandlers() {
	// Create spectator channel and add users
	sessions.AddSpectatorAddedHandler(func(user *sessions.User, spectator *sessions.User) {
		if len(user.GetSpectators()) == 1 {
			AddSpectatorChannel(user.Info.Id)
		}

		channel := GetSpectatorChannel(user.Info.Id)

		if channel == nil {
			return
		}

		if !channel.isUserInChannel(user) {
			channel.AddUser(user)
		}

		channel.AddUser(spectator)
	})

	// Remove spectators from channel and delete channel
	sessions.AddSpectatorLeftHandler(func(user *sessions.User, spectator *sessions.User) {
		channel := GetSpectatorChannel(user.Info.Id)

		if channel == nil {
			return
		}

		channel.RemoveUser(spectator)

		if len(user.GetSpectators()) == 0 {
			RemoveSpectatorChannel(user.Info.Id)
		}
	})
}

// Returns if the user is a moderator of the chat
func isChatModerator(userGroups common.UserGroups) bool {
	return common.HasAnyUserGroup(userGroups, []common.UserGroups{
		common.UserGroupSwan,
		common.UserGroupDeveloper,
		common.UserGroupAdmin,
		common.UserGroupModerator,
		common.UserGroupBot,
	})
}

// Returns a user's spectator channel name
func getSpectatorChannelName(userId int) string {
	return fmt.Sprintf("#spectator_%v", userId)
}
