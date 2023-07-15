package handlers

import (
	"example.com/Quaver/Z/db"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
	"log"
)

// Handles when the client requests to add/remove a user
func handleClientFriendship(user *sessions.User, packet *packets.ClientFriendship) {
	if packet == nil {
		return
	}

	if packet.UserId == user.Info.Id {
		return
	}

	relationship, err := db.GetUserRelationship(user.Info.Id, packet.UserId)

	if err != nil {
		log.Printf("Failed to get used relationship (#%v -> #%v) - %v\n", user.Info.Id, packet.UserId, err)
	}

	switch packet.Action {
	case packets.FriendsListActionAdd:
		if relationship != nil {
			return
		}

		err = db.AddFriend(user.Info.Id, packet.UserId)

		if err != nil {
			log.Printf("Failed to add friend (#%v -> #%v) - %v\n", user.Info.Id, packet.UserId, err)
		}
	case packets.FriendsListActionRemove:
		if relationship == nil {
			return
		}

		err = db.RemoveFriend(user.Info.Id, packet.UserId)

		if err != nil {
			log.Printf("Failed to remove friend (#%v -> #%v) - #%v\n", user.Info.Id, packet.UserId, err)
		}
	}
}
