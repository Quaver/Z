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
	case packets.PacketIdClientGameHostSelectingMap:
		handleClientGameHostSelectingMap(user, unmarshalPacket[packets.ClientGameHostSelectingMap](msg))
	case packets.PacketIdClientPacketChangeGamePassword:
		handleClientChangeGamePassword(user, unmarshalPacket[packets.ClientChangeGamePassword](msg))
	case packets.PacketIdClientGameChangeModifiers:
		handleClientGameChangeModifiers(user, unmarshalPacket[packets.ClientGameChangeModifiers](msg))
	case packets.PacketIdClientGameChangeFreeModType:
		handleClientGameChangeFreeMod(user, unmarshalPacket[packets.ClientGameFreeModTypeChanged](msg))
	case packets.PacketIdClientGamePlayerChangeModifiers:
		handleClientGameChangePlayerModifiers(user, unmarshalPacket[packets.ClientGameChangePlayerModifiers](msg))
	case packets.PacketIdClientGameChangeAutoHostRotation:
		handleClientGameHostRotation(user, unmarshalPacket[packets.ClientGameHostRotation](msg))
	case packets.PacketIdClientGameChangeMaxPlayers:
		handleClientGameChangeMaxPlayers(user, unmarshalPacket[packets.ClientGameChangeMaxPlayers](msg))
	case packets.PacketIdClientGameAcceptInvite:
		handleClientGameAcceptInvite(user, unmarshalPacket[packets.ClientGameAcceptInvite](msg))
	case packets.PacketIdClientRequestUserStats:
		handleClientRequestUserStats(user, unmarshalPacket[packets.ClientRequestUserStats](msg))
	case packets.PacketIdClientGameKickPlayer:
		handleClientGameKickPlayer(user, unmarshalPacket[packets.ClientGameKickPlayer](msg))
	case packets.PacketIdClientGameTransferHost:
		handleClientGameTransferHost(user, unmarshalPacket[packets.ClientGameTransferHost](msg))
	case packets.PacketIdClientInviteToGame:
		handleClientGameInvite(user, unmarshalPacket[packets.ClientGameInvite](msg))
	case packets.PacketIdClientGameScreenLoaded:
		handleClientGameScreenLoaded(user, unmarshalPacket[packets.ClientGameScreenLoaded](msg))
	case packets.PacketIdClientPlayerFinished:
		handleClientGamePlayerFinished(user, unmarshalPacket[packets.ClientGamePlayerFinished](msg))
	case packets.PacketIdClientGameSongSkipRequest:
		handleClientGamePlayerSkipSong(user, unmarshalPacket[packets.ClientGamePlayerSkipSong](msg))
	case packets.PacketIdClientGameJudgements:
		handleClientGameJudgements(user, unmarshalPacket[packets.ClientGameJudgements](msg))
	case packets.PacketIdClientFriendship:
		handleClientFriendship(user, unmarshalPacket[packets.ClientFriendship](msg))
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
