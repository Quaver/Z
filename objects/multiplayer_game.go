package objects

import "example.com/Quaver/Z/common"

type MultiplayerGame struct {
	Id                        int                         `json:"gid"`   // The id of the game in the database
	GameId                    string                      `json:"id"`    // A unique identifier for the game
	Name                      string                      `json:"n"`     // The name of the game
	CreationPassword          string                      `json:"pw"`    // The password of the game during creation
	HasPassword               bool                        `json:"hp"`    // If the game has a password on it
	MaxPlayers                int                         `json:"mp"`    // The maximum amount of players allowed in the game
	MapMD5                    string                      `json:"md5"`   // The MD5 hash of the currently played map
	MapMD5Alternative         string                      `json:"amd5"`  // An alternative md5 hash for the current played map. Usually osu! beatmap md5
	MapId                     int                         `json:"mid"`   // The id of the map in the database
	MapsetId                  int                         `json:"msid"`  // The id of the mapset in the database
	MapName                   string                      `json:"map"`   // The full name of the map
	MapGameMode               common.Mode                 `json:"gm"`    // The game mode for the currently selected map
	MapJudgementCount         int                         `json:"jc"`    // The amount of judgements possible in the map
	MapDifficultyRating       float64                     `json:"d"`     // The difficulty rating of the currently selected map
	MapAllDifficultyRatings   []float64                   `json:"adr"`   // The difficulty rating for all rates of the map. Host provides this for scoring on unsubmitted maps
	Ruleset                   MultiplayerGameRuleset      `json:"r"`     // The rules of the match (free-for-all, team, etc)
	IsHostRotation            bool                        `json:"hr"`    // Whether the server will control host rotation for the game
	InProgress                bool                        `json:"inp"`   // IF the match is currently in progress
	HostId                    int                         `json:"h"`     // The id of the host
	RefereeId                 int                         `json:"ref"`   // The id of the referee of the game
	PlayerIds                 []int                       `json:"ps"`    // The ids of the players in the game
	PlayersWithoutMap         []int                       `json:"pwm"`   // The players in the match that do not have the currently selected map
	PlayersReady              []int                       `json:"pri"`   // The players in the match that are readied up
	PlayerModifiers           []MultiplayerGamePlayerMods `json:"pm"`    // The modifiers that each player is using
	PlayersRedTeam            []int                       `json:"rtp"`   // The players that are on the red team
	PlayersBlueTeam           []int                       `json:"btp"`   // The players that are on the blue team
	PlayerWins                []MultiplayerGamePlayerWins `json:"plw"`   // The amount of wins each player has
	MatchCountdownTimestamp   int64                       `json:"cst"`   // A unix timestamp of the time the match countdown has started
	GlobalModifiers           int64                       `json:"md"`    // The modifiers that are used globally for every player in the match
	FreeModType               MultiplayerGameFreeMod      `json:"fm"`    // The type of free mod that is active for the match.
	TeamRedWins               int                         `json:"rtw"`   // The amount of wins the red team has
	TeamBlueWins              int                         `json:"btw"`   // The amount of wins the blue team has
	IsHostSelectingMap        bool                        `json:"hsm"`   // If the host is currently selecting a map
	IsMapsetShared            bool                        `json:"ims"`   // If the mapset is temporarily uploaded and shared by the host
	IsTournamentMode          bool                        `json:"trn"`   // If the game is currently in tournament mode
	FilterMinDifficultyRating float32                     `json:"mind"`  // The minimum difficulty rating allowed for maps in the game
	FilterMaxDifficultyRating float32                     `json:"maxd"`  // The maximum difficulty rating allowed for maps in the game
	FilterMaxSongLength       int                         `json:"maxl"`  // The maximum length allowed for maps in the lobby
	FilterAllowedGameModes    []common.Mode               `json:"ag"`    // The game modes that are allowed to be selected in the game
	FilterMinLongNotePercent  int                         `json:"lnmin"` // The minimum long note percentage for the map
	FilterMaxLongNotePercent  int                         `json:"lnmax"` // The maximum long note percentage for the map
	FilterMinAudioRate        float64                     `json:"mr"`    // The minimum audio rate allowed for free mod
}

func (mg *MultiplayerGame) SetDefaults() {
	mg.PlayerIds = []int{}
	mg.PlayersWithoutMap = []int{}
	mg.PlayersReady = []int{}
	mg.PlayerModifiers = []MultiplayerGamePlayerMods{}
	mg.PlayersRedTeam = []int{}
	mg.PlayersBlueTeam = []int{}
	mg.PlayerWins = []MultiplayerGamePlayerWins{}
	mg.FilterAllowedGameModes = []common.Mode{}

	mg.FilterMaxDifficultyRating = 999999999
	mg.FilterMaxSongLength = 999999999
	mg.FilterMaxLongNotePercent = 100
	mg.FilterMinAudioRate = 0.5
}
