package webhooks

import (
	"example.com/Quaver/Z/config"
	"testing"
)

func TestInitialize(t *testing.T) {
	_ = config.Load("../config.json")

	if config.Instance == nil {
		return
	}

	Initialize()
}

func TestSendAntiCheat(t *testing.T) {
	_ = config.Load("../config.json")

	if config.Instance == nil {
		return
	}

	Initialize()

	SendAntiCheatProcessLog("Quaver", 11, "https://quavergame.com/user/Quaver", quaverLogo, []string{
		"L33t H4x",
		"h3ll0 1f ur r34d1ng",
	})
}
