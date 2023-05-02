package handlers

import (
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
	"log"
	"time"
)

// Handles when a user sends a pong packet
func handleClientPong(user *sessions.User, packet *packets.ClientPong) {
	if packet == nil {
		return
	}

	user.LastPingTimestamp = time.Now().UnixMilli()

	processes := packet.ParseProcessList()

	// TODO: No process list provided. Ban them if they can't bypass this.
	if processes == nil || len(processes) == 0 {
		log.Printf("[%v - %v] Sent a Pong packet without any process list\n", user.Info.Id, user.Info.Username)
		return
	}

	log.Printf("%v %v provided: %v processes\n", user.Info.Username, user.Info.Id, len(processes))
}
