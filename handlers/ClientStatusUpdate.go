package handlers

import (
	"example.com/Quaver/Z/objects"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
	"example.com/Quaver/Z/spectator"
)

// Handles when a user's client sends a status update
func handleClientStatusUpdate(user *sessions.User, packet *packets.ClientStatusUpdate) {
	if packet == nil {
		return
	}

	if packet.Status == (objects.ClientStatus{}) {
		return
	}

	user.SetClientStatus(&packet.Status)
	spectator.GetUser(user).SendClientStatusToSpectators()
}
