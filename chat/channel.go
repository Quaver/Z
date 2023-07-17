package chat

import (
	"example.com/Quaver/Z/db"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
	"example.com/Quaver/Z/webhooks"
	"github.com/disgoorg/disgo/webhook"
	"log"
	"sync"
)

type Channel struct {
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	AdminOnly      bool                   `json:"admin_only"`
	AutoJoin       bool                   `json:"auto_join"`
	DiscordWebhook string                 `json:"discord_webhook"`
	WebhookClient  webhook.Client         `json:"-"`
	Participants   map[int]*sessions.User `json:"-"`
	IsMultiplayer  bool                   `json:"-"`
	mutex          *sync.Mutex
}

// NewChannel Creates a new chat channel instance
func NewChannel(name string, description string, adminOnly bool, autoJoin bool, isMultiplayer bool, discordWebhook string) *Channel {
	channel := Channel{
		Name:           name,
		Description:    description,
		AdminOnly:      adminOnly,
		AutoJoin:       autoJoin,
		DiscordWebhook: discordWebhook,
		WebhookClient:  nil,
		Participants:   map[int]*sessions.User{},
		IsMultiplayer:  isMultiplayer,
		mutex:          &sync.Mutex{},
	}

	channel.initializeWebhook()
	return &channel
}

// Initializes the Discord webhook for the channel
func (channel *Channel) initializeWebhook() {
	if channel.WebhookClient != nil {
		return
	}

	if channel.DiscordWebhook == "" {
		log.Printf("Empty webhook url for chat channel: %v\n", channel.Name)
		return
	}

	var err error

	channel.WebhookClient, err = webhook.NewWithURL(channel.DiscordWebhook)

	if err != nil {
		panic(err)
	}

	log.Printf("Initialized webhook for channel: %v (%v)\n", channel.Name, channel.WebhookClient.ID().String())
}

// AddUser Adds a user to the channel
func (channel *Channel) AddUser(user *sessions.User) {
	channel.mutex.Lock()
	defer channel.mutex.Unlock()

	if channel.AdminOnly && !isChatModerator(user.Info.UserGroups) {
		sessions.SendPacketToUser(packets.NewServerFailedToJoinChatChannel(channel.Name), user)
		return
	}

	channel.Participants[user.Info.Id] = user
	sessions.SendPacketToUser(packets.NewServerJoinedChatChannel(channel.Name), user)
}

// RemoveUser Removes a user from the channel
func (channel *Channel) RemoveUser(user *sessions.User) {
	channel.mutex.Lock()
	defer channel.mutex.Unlock()

	if _, ok := channel.Participants[user.Info.Id]; ok {
		delete(channel.Participants, user.Info.Id)
	}

	sessions.SendPacketToUser(packets.NewServerLeftChatChannel(channel.Name), user)
}

// SendMessage Sends a message to all the users in the channel
func (channel *Channel) SendMessage(sender *sessions.User, message string) {
	channel.mutex.Lock()
	defer channel.mutex.Unlock()

	packet := packets.NewServerChatMessage(sender.Info.Id, sender.Info.Username, channel.Name, message)

	for _, user := range channel.Participants {
		if user == sender {
			continue
		}

		sessions.SendPacketToUser(packet, user)
	}

	err := db.InsertPublicChatMessage(sender.Info.Id, channel.Name, message)

	if err != nil {
		log.Printf("Failed to insert chat message to DB: %v\n", err)
	}
}

// Sends a discord webhook
func (channel *Channel) sendWebhook(sender *sessions.User, message string) {
	webhooks.SendChatMessage(channel.WebhookClient, sender.Info.Username, sender.Info.GetProfileUrl(), sender.Info.AvatarUrl, channel.Name, message)
}

// Removes all users from the channel
func (channel *Channel) removeAllUsers() {
	channel.mutex.Lock()
	defer channel.mutex.Unlock()

	for _, user := range channel.Participants {
		channel.RemoveUser(user)
	}
}
