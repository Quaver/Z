package webhooks

import (
	"example.com/Quaver/Z/config"
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
	"log"
	"time"
)

var (
	AntiCheat   webhook.Client
	PrivateChat webhook.Client
)

const QuaverLogo string = "https://i.imgur.com/DkJhqvT.jpg"
const antiCheatDescription string = "**❌ Anti-cheat Triggered!**"

func Initialize() {
	hooks := []string{
		config.Instance.DiscordWebhooks.AntiCheat,
		config.Instance.DiscordWebhooks.PrivateChat,
	}

	for _, hook := range hooks {
		if hook == "" {
			continue
		}

		var err error

		switch hook {
		case config.Instance.DiscordWebhooks.AntiCheat:
			AntiCheat, err = webhook.NewWithURL(hook)
		case config.Instance.DiscordWebhooks.PrivateChat:
			PrivateChat, err = webhook.NewWithURL(hook)
		}

		if err != nil {
			panic(err)
		}

		log.Printf("Initialized Webhook: %v\n", hook)
	}
}

func SendAntiCheat(username string, userId int, url string, icon string, reason string, text string) {
	if AntiCheat == nil {
		log.Printf("Cannot send nil Anti-Cheat webhook\n")
		return
	}

	viewProfile := fmt.Sprintf("[View Profile](%v)", url)
	banUser := fmt.Sprintf("[Ban User](https://a.quavergame.com/ban/%v)", userId)
	editUser := fmt.Sprintf("[Edit User](https://a.quavergame.com/edituser/%v)", userId)

	embed := discord.NewEmbedBuilder().
		SetAuthor(username, url, icon).
		SetDescription(antiCheatDescription).
		SetFields(discord.EmbedField{
			Name:  reason,
			Value: text,
		}, discord.EmbedField{
			Name:   "Admin Actions",
			Value:  fmt.Sprintf("%v | %v | %v", viewProfile, banUser, editUser),
			Inline: nil,
		}).
		SetThumbnail(QuaverLogo).
		SetFooter("Quaver", QuaverLogo).
		SetTimestamp(time.Now()).
		SetColor(0xFF0000).
		Build()

	_, err := AntiCheat.CreateEmbeds([]discord.Embed{embed})

	if err != nil {
		log.Printf("Failed to send anti-cheat webhook: %v\n", err)
	}
}

func SendAntiCheatProcessLog(username string, userId int, url string, icon string, processes []string) {
	formatted := ""

	for i, proc := range processes {
		formatted += fmt.Sprintf("**%v. %v**\n", i+1, proc)
	}

	SendAntiCheat(username, userId, url, icon, "Detected Processes", formatted)
}

func SendAntiCheatLibraries(username string, userId int, url string, icon string, libraries []string) {
	formatted := ""

	for i, library := range libraries {
		formatted += fmt.Sprintf("**%v. %v**\n", i+1, library)
	}

	SendAntiCheat(username, userId, url, icon, "Detected Libraries", formatted)
}

// SendChatMessage Sends a chat message webhook to Discord
func SendChatMessage(webhook webhook.Client, senderUsername string, senderProfileUrl string, senderAvatarUrl, receiverName string, message string) {
	if webhook == nil {
		return
	}

	embed := discord.NewEmbedBuilder().
		SetAuthor(fmt.Sprintf("%v → %v", senderUsername, receiverName), senderProfileUrl, senderAvatarUrl).
		SetDescription(message).
		SetFooter("Quaver", QuaverLogo).
		SetTimestamp(time.Now()).
		SetColor(0x00FFFF).
		Build()

	_, err := webhook.CreateEmbeds([]discord.Embed{embed})

	if err != nil {
		log.Printf("Failed to send webhook to channel: %v - %v\n", receiverName, err)
	}
}
