package handlers

import (
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
	"example.com/Quaver/Z/spectator"
)

// Handles when a client provides spectator replay frames
func handleClientSpectatorReplayFrames(user *sessions.User, packet *packets.ClientSpectatorReplayFrames) {
	if packet == nil {
		return
	}

	spectator.GetUser(user).HandleIncomingFrames(packet)
}
