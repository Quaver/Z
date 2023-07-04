package multiplayer

import (
	"example.com/Quaver/Z/objects"
	"example.com/Quaver/Z/utils"
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

	game.Data.SetDefaults()
	game.setupCreationPassword()
	game.validateSettings()

	err := game.insertIntoDatabase()

	if err != nil {
		return nil, err
	}

	return &game, nil
}
