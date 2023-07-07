package handlers

import (
	"encoding/json"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
	"fmt"
	"log"
	"net"
)

// HandleIncomingPackets Handles incoming messages from clients
func HandleIncomingPackets(conn net.Conn, msg string) {
	user := sessions.GetUserByConnection(conn)

	if user == nil {
		log.Printf("[%v] Received packet while not logged in: %v\n", conn.RemoteAddr(), msg)
		return
	}

	var p packets.Packet

	if err := json.Unmarshal([]byte(msg), &p); err != nil {
		log.Println(err)
		return
	}

	switch p.Id {
	case packets.PacketIdClientPong:
		handleClientPong(user, unmarshalPacket[packets.ClientPong](msg))
	case packets.PacketIdClientChatMessage:
		handleClientChatMessage(user, unmarshalPacket[packets.ClientChatMessage](msg))
	case packets.PacketIdClientStatusUpdate:
		handleClientStatusUpdate(user, unmarshalPacket[packets.ClientStatusUpdate](msg))
	case packets.PacketIdClientRequestUserInfo:
		handleClientRequestUserInfo(user, unmarshalPacket[packets.ClientRequestUserInfo](msg))
	case packets.PacketIdClientRequestLeaveChatChannel:
		handleClientRequestLeaveChatChannel(user, unmarshalPacket[packets.ClientRequestLeaveChatChannel](msg))
	case packets.PacketIdClientRequestJoinChatChannel:
		handleClientRequestJoinChatChannel(user, unmarshalPacket[packets.ClientRequestJoinChatChannel](msg))
	case packets.PacketIdClientRequestUserStatus:
		handleClientRequestUserStatus(user, unmarshalPacket[packets.ClientRequestUserStatus](msg))
	case packets.PacketIdClientLobbyJoin:
		handleClientLobbyJoin(user, unmarshalPacket[packets.ClientLobbyJoin](msg))
	case packets.PacketIdClientLobbyLeave:
		handleClientLobbyLeave(user, unmarshalPacket[packets.ClientLobbyLeave](msg))
	case packets.PacketIdClientCreateGame:
		handleClientCreateGame(user, unmarshalPacket[packets.ClientCreateGame](msg))
	case packets.PacketIdClientLeaveGame:
		handleClientLeaveGame(user, unmarshalPacket[packets.ClientLeaveGame](msg))
	case packets.PacketIdClientJoinGame:
		handleClientJoinGame(user, unmarshalPacket[packets.ClientJoinGame](msg))
	case packets.PacketIdClientChangeGameMap:
		handleClientChangeGameMap(user, unmarshalPacket[packets.ClientChangeGameMap](msg))
	case packets.PacketIdClientGamePlayerNoMap:
		handleClientGamePlayerNoMap(user, unmarshalPacket[packets.ClientGamePlayerNoMap](msg))
	case packets.PacketIdClientGamePlayerHasMap:
		handleClientGamePlayerHasMap(user, unmarshalPacket[packets.ClientGamePlayerHasMap](msg))
	case packets.PacketIdClientGamePlayerReady:
		handleClientGamePlayerReady(user, unmarshalPacket[packets.ClientGamePlayerReady](msg))
	case packets.PacketIdClientGamePlayerNotReady:
		handleClientGamePlayerNotReady(user, unmarshalPacket[packets.ClientGamePlayerNotReady](msg))
	case packets.PacketIdClientGameStartCountdown:
		handleClientGameStartCountdown(user, unmarshalPacket[packets.ClientGameStartCountdown](msg))
	case packets.PacketIdClientGameStopCountdown:
		handleClientGameStopCountdown(user, unmarshalPacket[packets.ClientGameStopCountdown](msg))
	case packets.PacketIdClientPacketChangeGameName:
		handleClientChangeGameName(user, unmarshalPacket[packets.ClientChangeGameName](msg))
	default:
		log.Println(fmt.Errorf("unknown packet: %v", msg))
	}
}

// unmarshalPacket Unmarshal a packet of a specified type
func unmarshalPacket[T any](packet string) *T {
	var data T

	if err := json.Unmarshal([]byte(packet), &data); err != nil {
		log.Println(err)
		return nil
	}

	return &data
}
