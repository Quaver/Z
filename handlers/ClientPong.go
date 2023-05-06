package handlers

import (
	"example.com/Quaver/Z/db"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
	"example.com/Quaver/Z/webhooks"
	"golang.org/x/exp/slices"
	"log"
	"time"
)

// Handles when a user sends a pong packet
func handleClientPong(user *sessions.User, packet *packets.ClientPong) {
	if packet == nil {
		return
	}

	user.LastPongTimestamp = time.Now().UnixMilli()

	packetProcs := packet.ParseProcessList()

	if packetProcs == nil || len(packetProcs) == 0 {
		webhooks.SendAntiCheatProcessLog(user.Info.Username, user.Info.GetProfileUrl(), user.Info.AvatarUrl, []string{"NO PROCESSES PROVIDED"})
		log.Printf("[%v - %v] Sent a Pong packet without any process list\n", user.Info.Id, user.Info.Username)
		return
	}

	dbProcs, err := db.FetchProcesses()

	if err != nil {
		log.Printf("Failed to fetch process from database - %v\n", err)
		return
	}

	detected := detectProcesses(dbProcs, packetProcs)

	if len(detected) == 0 {
		user.LastDetectedProcesses = []string{}
		return
	}

	if slices.Equal(detected, user.LastDetectedProcesses) {
		return
	}

	user.LastDetectedProcesses = detected
	webhooks.SendAntiCheatProcessLog(user.Info.Username, user.Info.GetProfileUrl(), user.Info.AvatarUrl, user.LastDetectedProcesses)

	log.Printf("[%v - #%v] Detected %v flagged processes \n", user.Info.Username, user.Info.Id, len(detected))
}

// Goes through both the db processes and packet processes and checks if any are found
func detectProcesses(dbProcesses []*db.Process, packetProcesses []packets.Process) []string {
	detected := make([]string, 0)

	for _, dbProcess := range dbProcesses {
		md5 := dbProcess.GetMD5()

		for _, packetProcess := range packetProcesses {
			if md5 != packetProcess.Name && md5 != packetProcess.WindowTitle && md5 != packetProcess.FileName {
				continue
			}

			if !slices.Contains(detected, dbProcess.Name) {
				detected = append(detected, dbProcess.Name)
				break
			}
		}
	}

	return detected
}
