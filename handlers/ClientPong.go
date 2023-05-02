package handlers

import (
	"example.com/Quaver/Z/db"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
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
		user.LastDetectedProcesses = []int{}
		return
	}

	if slices.Equal(detected, user.LastDetectedProcesses) {
		return
	}

	user.LastDetectedProcesses = detected
	log.Printf("%v %v provided: %v processes | %v detected\n", user.Info.Username, user.Info.Id, len(packetProcs), len(detected))
}

// Goes through both the db processes and packet processes and checks if any are found
func detectProcesses(dbProcesses []*db.Process, packetProcesses []packets.Process) []int {
	detectedIds := make([]int, 0)

	for _, dbProcess := range dbProcesses {
		md5 := dbProcess.GetMD5()

		for _, packetProcess := range packetProcesses {
			if md5 != packetProcess.Name && md5 != packetProcess.WindowTitle && md5 != packetProcess.FileName {
				continue
			}

			if !slices.Contains(detectedIds, dbProcess.Id) {
				detectedIds = append(detectedIds, dbProcess.Id)
				break
			}
		}
	}

	return detectedIds
}
