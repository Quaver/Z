package handlers

import (
	"example.com/Quaver/Z/db"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
	"example.com/Quaver/Z/webhooks"
	"golang.org/x/exp/slices"
	"log"
)

// Handles when a user sends a pong packet
func handleClientPong(user *sessions.User, packet *packets.ClientPong) {
	if packet == nil {
		return
	}

	user.SetLastPongTimestamp()

	parsed := packet.Parse()

	checkProcesses(user, parsed.Processes)
	// checkLibraries(user, parsed.Libraries)
}

func checkProcesses(user *sessions.User, processes []packets.Process) {
	if processes == nil || len(processes) == 0 {
		return
	}

	dbProcesses, err := db.FetchProcesses()

	if err != nil {
		log.Printf("Failed to fetch process from database - %v\n", err)
		return
	}

	detected := detectProcesses(dbProcesses, processes)

	if len(detected) == 0 {
		user.SetLastDetectedProcesses([]string{})
		return
	}

	if slices.Equal(detected, user.GetLastDetectedProcesses()) {
		return
	}

	user.SetLastDetectedProcesses(detected)

	webhooks.SendAntiCheatProcessLog(user.Info.Username, user.Info.Id, user.Info.GetProfileUrl(), user.Info.AvatarUrl.String, detected)

	log.Printf("[%v - #%v] Detected %v flagged processes \n", user.Info.Username, user.Info.Id, len(detected))
}

func checkLibraries(user *sessions.User, libraries []string) {
	if libraries == nil || len(libraries) == 0 {
		return
	}

	if slices.Equal(libraries, user.GetLastLibraries()) {
		return
	}

	user.SetLastLibraries(libraries)

	webhooks.SendAntiCheatLibraries(user.Info.Username, user.Info.Id, user.Info.GetProfileUrl(),
		user.Info.AvatarUrl.String, libraries)

	log.Printf("[%v - #%v] Detected %v libraries \n", user.Info.Username, user.Info.Id, len(libraries))
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
