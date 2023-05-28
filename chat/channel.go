package chat

import (
	"example.com/Quaver/Z/db"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
	"example.com/Quaver/Z/webhooks"
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
	"log"
	"sync"
	"time"
)

type Channel struct {
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	AdminOnly      bool                   `json:"admin_only"`
	AutoJoin       bool                   `json:"auto_join"`
	DiscordWebhook string                 `json:"discord_webhook"`
	WebhookClient  webhook.Client         `json:"-"`
	Participants   map[int]*sessions.User `json:"-"`
	mutex          *sync.Mutex
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

	channel.Participants[user.Info.Id] = user
	sessions.SendPacketToUser(packets.NewServerJoinedChatChannel(channel.Name), user)
}

// RemoveUser Removes a user from the channel
func (channel *Channel) RemoveUser(user *sessions.User) {
	channel.mutex.Lock()
	defer channel.mutex.Unlock()

	delete(channel.Participants, user.Info.Id)
	sessions.SendPacketToUser(packets.NewServerPing(), user)
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

	channel.sendWebhook(sender, message)

	err := db.InsertPublicChatMessage(sender.Info.Id, channel.Name, message)

	if err != nil {
		log.Printf("Failed to insert chat message to DB: %v\n", err)
	}
}

// Sends a webhook to Discord
func (channel *Channel) sendWebhook(sender *sessions.User, message string) {
	if channel.WebhookClient == nil {
		return
	}

	embed := discord.NewEmbedBuilder().
		SetAuthor(fmt.Sprintf("%v → %v", sender.Info.Username, channel.Name), sender.Info.GetProfileUrl(), sender.Info.AvatarUrl).
		SetDescription(message).
		SetFooter("Quaver", webhooks.QuaverLogo).
		SetTimestamp(time.Now()).
		SetColor(0x00FFFF).
		Build()

	_, err := channel.WebhookClient.CreateEmbeds([]discord.Embed{embed})

	if err != nil {
		log.Printf("Failed to send webhook to channel: %v - %v\n", channel.Name, err)
	}
}
