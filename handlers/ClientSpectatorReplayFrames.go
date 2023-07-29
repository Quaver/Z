package handlers

import (
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
)

// Handles when a client provides spectator replay frames
func handleClientSpectatorReplayFrames(user *sessions.User, packet *packets.ClientSpectatorReplayFrames) {
	if packet == nil {
		return
	}

	user.HandleNewSpectatorFrames(packet)
}
