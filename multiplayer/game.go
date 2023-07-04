package multiplayer

import (
	"example.com/Quaver/Z/common"
	"example.com/Quaver/Z/db"
	"example.com/Quaver/Z/objects"
	"example.com/Quaver/Z/utils"
	"math"
)

type Game struct {
	Data      *objects.MultiplayerGame
	Password  string // The password for the game. This is different from Data.CreationPassword, as it is hidden from users.
	CreatorId int    // The id of the user who created the game
}

// NewGame Creates a new multiplayer game from a game
func NewGame(gameData *objects.MultiplayerGame, creatorId int) (*Game, error) {
	game := Game{Data: gameData}
	game.Data.GameId = utils.GenerateRandomString(32)
	game.CreatorId = creatorId

	// We don't want the password to be exposed in the JSON of the multiplayer game, so we are using another property to hide it.
	if game.Data.CreationPassword != "" {
		game.Password = game.Data.CreationPassword
		game.Data.CreationPassword = ""
	}

	var err error
	game.Data.Id, err = db.InsertMultiplayerGame(game.Data.Name, game.Data.GameId)

	if err != nil {
		return nil, err
	}

	game.Data.SetDefaults()
	game.validateSettings()

	return &game, nil
}

// validateSettings Checks the multiplayer settings to see if they are in an acceptable range
func (game *Game) validateSettings() {
	data := game.Data

	data.Name = utils.TruncateString(data.Name, 50)
	data.HasPassword = game.Password != ""
	data.MaxPlayers = utils.Clamp(data.MaxPlayers, 2, 16)
	data.Ruleset = utils.Clamp(data.Ruleset, objects.MultiplayerGameRulesetFreeForAll, objects.MultiplayerGameRulesetTeam)
	data.FreeModType = utils.Clamp(data.FreeModType, objects.MultiplayerGameFreeModNone, objects.MultiplayerGameFreeModRegular|objects.MultiplayerGameFreeModRate)

	data.MapMD5 = utils.TruncateString(data.MapMD5, 32)
	data.MapMD5Alternative = utils.TruncateString(data.MapMD5Alternative, 32)
	data.MapName = utils.TruncateString(data.MapName, 250)
	data.MapGameMode = utils.Clamp(data.MapGameMode, common.ModeKeys4, common.ModeKeys7)

	data.FilterMinDifficultyRating = utils.Clamp(data.FilterMinDifficultyRating, 0, math.MaxInt32)
	data.FilterMaxDifficultyRating = utils.Clamp(data.FilterMaxDifficultyRating, 0, math.MaxInt32)
	data.FilterMaxSongLength = utils.Clamp(data.FilterMaxSongLength, 0, math.MaxInt32)
	data.FilterMinLongNotePercent = utils.Clamp(data.FilterMinLongNotePercent, 0, 100)
	data.FilterMaxLongNotePercent = utils.Clamp(data.FilterMaxLongNotePercent, 0, 100)
	data.FilterMinAudioRate = utils.Clamp(data.FilterMinAudioRate, 0.5, 2.0)

	// There is a maximum of 21 rates allowed in the game. So if we don't have all of them, then just clear it.
	if len(data.MapAllDifficultyRatings) < 21 {
		data.MapAllDifficultyRatings = []float64{}
	}
}
