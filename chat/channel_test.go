package chat

import (
	"example.com/Quaver/Z/db"
	"example.com/Quaver/Z/sessions"
	"example.com/Quaver/Z/webhooks"
	"testing"
)

func TestChannelSendWebhook(t *testing.T) {
	// Testing Webhook
	hook := "https://discord.com/api/webhooks/1112399851762745415/z72i_7LVi9ShbVnRcOjUj4VDyv4-jTqKSsif2c7yi52qWCHrLUuLymmWhYJXWDaa11h3"

	channel := NewChannel("#test", "", false, false, hook)

	user := &sessions.User{Info: &db.User{
		Id:        2,
		Username:  "Quaver Bot",
		AvatarUrl: webhooks.QuaverLogo,
	}}

	channel.sendWebhook(user, "This is a test")
}
