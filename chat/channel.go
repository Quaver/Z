package chat

import (
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
	"github.com/disgoorg/disgo/webhook"
	"log"
)

type Channel struct {
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	AdminOnly      bool                   `json:"admin_only"`
	AutoJoin       bool                   `json:"auto_join"`
	DiscordWebhook string                 `json:"discord_webhook"`
	WebhookClient  webhook.Client         `json:"-"`
	Participants   map[int]*sessions.User `json:"-"`
}

// NewChannel Creates a new chat channel instance
func NewChannel(name string, description string, adminOnly bool, autoJoin bool, discordWebhook string) *Channel {
	channel := Channel{
		Name:           name,
		Description:    description,
		AdminOnly:      adminOnly,
		AutoJoin:       autoJoin,
		DiscordWebhook: discordWebhook,
		WebhookClient:  nil,
		Participants:   map[int]*sessions.User{},
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
	mutex.Lock()
	defer mutex.Unlock()

	channel.Participants[user.Info.Id] = user
	sessions.SendPacketToUser(packets.NewServerJoinedChatChannel(channel.Name), user)
}

// RemoveUser Removes a user from the channel
func (channel *Channel) RemoveUser(user *sessions.User) {
	mutex.Lock()
	defer mutex.Unlock()

	delete(channel.Participants, user.Info.Id)
	sessions.SendPacketToUser(packets.NewServerPing(), user)
}
