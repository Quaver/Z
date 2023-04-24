package sessions

import (
	"example.com/Quaver/Z/common"
	"example.com/Quaver/Z/config"
	"example.com/Quaver/Z/db"
	"testing"
)

func TestPopulateUserStats(t *testing.T) {
	_ = config.Load("../config.json")

	if config.Instance == nil {
		return
	}

	db.InitializeSQL()

	user := NewUser(nil, &db.User{Id: 1})

	err := user.UpdateStats()

	if err != nil {
		t.Fatal(err)
	}

	if len(user.Stats) != int(common.ModeEnumMaxValue)-1 {
		t.Fatalf("expected (%v) mode stats. only fetched %v", int(common.ModeEnumMaxValue)-1, len(user.Stats))
	}

	db.CloseSQLConnection()
}
