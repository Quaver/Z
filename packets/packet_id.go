package packets

type PacketId int

const (
	PacketIdInvalid = iota // UNUSED
	PacketIdServerPing
	PacketIdClientPong
	PacketIdServerLoginReply
	PacketIdServerUserDisconnected
	PacketIdServerUserConnected
	PacketIdServerAvailableChatChannel
	PacketIdServerJoinedChatChannel
	PacketIdServerChatMessage
	PacketIdClientChatMessage
	PacketIdClientRequestLeaveChatChannel
	PacketIdServerLeftChatChannelPacket
	PacketIdClientRequestJoinChatChannel
	PacketIdServerFailedToJoinChannelPacket
	PacketIdServerMuteEndTimePacket // UNUSED
	PacketIdServerNotification
	PacketIdClientStatusUpdate
	PacketIdServerUsersOnline
	PacketIdClientRequestUserInfo
	PacketIdServerUserInfo
	PacketIdClientRequestUserStatus
	PacketIdServerUserStatus
	PacketIdServerFailedToLogin // UNUSED
	PacketIdServerChooseUsername
	PacketIdClientLobbyJoin
	PacketIdClientLobbyLeave
	PacketIdClientCreateGame
	PacketIdServerMultiplayerGameInfo
	PacketIdServerJoinGame
	PacketIdServerChangeGameHost
	PacketIdClientLeaveGame
	PacketIdServerGameDisbanded
	PacketIdClientJoinGame
	PacketIdServerJoinGameFailed
	PacketIdServerUserJoinedGame
	PacketIdServerUserLeftGame
	PacketIdClientChangeGameMap
	PacketIdServerGameMapChanged
	PacketIdClientGamePlayerNoMap
	PacketIdServerGamePlayerNoMap
	PacketIdClientGamePlayerHasMap
	PacketIdServerGamePlayerHasMap
	PacketIdServerGameStart
	PacketIdClientPlayerFinished
	PacketIdServerGameEnded
	PacketIdClientGameJudgements
	PacketIdServerGameJudgements
	PacketIdClientGameScreenLoaded
	PacketIdServerAllPlayersLoaded
	PacketIdClientGameSongSkipRequest
	PacketIdServerGameAllPlayersSkipped
	PacketIdClientGamePlayerReady
	PacketIdServerGamePlayerReady
	PacketIdClientGamePlayerNotReady
	PacketIdServerGamePlayerNotReady
	PacketIdClientGameStartCountdown
	PacketIdServerGameStartCountdown
	PacketIdClientGameStopCountdown
	PacketIdServerGameStopCountdown
	PacketIdServerGameDifficultyRangeChanged
	PacketIdServerGameMaxSongLengthChanged
	PacketIdServerGameAllowedModesChanged
	PacketIdClientGameChangeModifiers
	PacketIdServerGameChangeModifiers
	PacketIdServerGameFreeModTypeChanged
	PacketIdClientGamePlayerChangeModifiers
	PacketIdServerGamePlayerChangeModifiers
	PacketIdServerGameKicked
	PacketIdServerGameNameChanged
	PacketIdServerGameInvite
	PacketIdClientGameAcceptInvite
	PacketIdServerGameHealthTypeChanged // UNUSED
	PacketIdServerGameLivesChanged      // UNUSED
	PacketIdServerGameHostRotationChanged
	PacketIdServerGamePlayerTeamChanged // UNUSED
	PacketIdClientGamePlayerTeamChanged // UNUSED
	PacketIdServerGameRulesetChanged    // UNUSED
	PacketIdServerGameLongNotePercentageChanged
	PacketIdServerGameMaxPlayersChanged
	PacketIdServerGameMinimumRateChanged // UNUSED
	PacketIdServerGameTeamWinCount       // UNUSED
	PacketIdServerGamePlayerWinCount
	PacketIdClientRequestUserStats
	PacketIdServerUserStats
	PacketIdServerGamePlayerBattleRoyaleEliminated // UNUSED
	PacketIdClientGameKickPlayer
	PacketIdClientGameTransferHost
	PacketIdClientGameChangeOtherPlayerTeam // UNUSED
	PacketIdClientGameChangeRuleset         // UNUSED
	PacketIdClientGameChangeMaxPlayers
	PacketIdClientGameChangeAutoHostRotation
	PacketIdClientGameChangeHealthType // UNUSED
	PacketIdClientGameChangeLivesCount // UNUSED
	PacketIdClientGameChangeFreeModType
	PacketIdClientGameHostSelectingMap
	PacketIdServerGameHostSelectingMap
	PacketIdServerGameSetReferee
	PacketIdClientStartSpectatePlayer
	PacketIdClientStopSpectatePlayer
	PacketIdServerStartSpectatePlayer
	PacketIdServerStopSpectatePlayer
	PacketIdServerSpectatorJoined
	PacketIdServerSpectatorLeft
	PacketIdClientSpectatorReplayFrames
	PacketIdServerSpectatorReplayFrames
	PacketIdServerListeningPartyJoined          // UNUSED
	PacketIdServerListeningPartyLeft            // UNUSED
	PacketIdClientListeningPartyStateUpdate     // UNUSED
	PacketIdServerListeningPartyStateUpdate     // UNUSED
	PacketIdServerListeningPartyFellowJoined    // UNUSED
	PacketIdServerListeningPartyFellowLeft      // UNUSED
	PacketIdClientListeningPartyChangeHost      // UNUSED
	PacketIdServerListeningPartyChangeHost      // UNUSED
	PacketIdClientListeningPartyKickUser        // UNUSED
	PacketIdClientListeningPartyUserMissingSong // UNUSED
	PacketIdServerListeningPartyUserMissingSong // UNUSED
	PacketIdClientListeningPartyUserHasSong     // UNUSED
	PacketIdServerListeningPartyUserHasSong     // UNUSED
	PacketIdServerUserFriendsList
	PacketIdClientFriendship
	PacketIdClientJoinListeningParty // UNUSED
	PacketIdClientInviteToGame
	PacketIdServerSongRequest
	PacketIdServerTwitchConnection
	PacketIdClientTwitchUnlink
	PacketIdServerGameMapsetShared
	PacketIdClientPacketChangeGameName
	PacketIdClientPacketChangeGamePassword
	PacketIdClientSpectateMultiplayerGame
	PacketIdServerSpectateMultiplayerGame
	PacketIdServerGameTournamentMode
	PacketIdServerGameNeedDifficultyRatings
	PacketIdClientGameDifficultyRatings
	PacketIdClientGameAutoHost
	PacketIdServerGameAutoHost
	PacketIdClientLogout
	PacketIdServerClearSpectateeReplayFrames
)
