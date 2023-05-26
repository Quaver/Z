package chat

import (
	"example.com/Quaver/Z/sessions"
	"github.com/disgoorg/disgo/webhook"
	"log"
)

type Channel struct {
	Name           string           `json:"name"`
	Description    string           `json:"description"`
	AdminOnly      bool             `json:"admin_only"`
	AutoJoin       bool             `json:"auto_join"`
	DiscordWebhook string           `json:"discord_webhook"`
	Participants   []*sessions.User `json:"-"`
	WebhookClient  webhook.Client   `json:"-"`
}

// NewChannel Creates a new chat channel instance
func NewChannel(name string, description string, adminOnly bool, autoJoin bool, discordWebhook string) *Channel {
	channel := Channel{
		Name:           name,
		Description:    description,
		AdminOnly:      adminOnly,
		AutoJoin:       autoJoin,
		DiscordWebhook: discordWebhook,
		Participants:   []*sessions.User{},
		WebhookClient:  nil,
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
