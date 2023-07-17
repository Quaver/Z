package handlers

import (
	"database/sql"
	"example.com/Quaver/Z/db"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
	"log"
)

// Handles when the client requests to unlink their twitch account
func handleClientUnlinkTwitch(user *sessions.User, packet *packets.ClientUnlinkTwitch) {
	if packet == nil {
		return
	}

	user.Info.TwitchUsername = sql.NullString{}

	err := db.UnlinkUserTwitch(user.Info.Id)

	if err != nil {
		log.Printf("Error unlinking user #%v's twitch account - %v\n", user.Info.Id, err)
		return
	}

	sessions.SendPacketToUser(packets.NewServerTwitchConnection(""), user)
}
