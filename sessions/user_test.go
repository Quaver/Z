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
	db.InitializeRedis()

	user := NewUser(nil, &db.User{Id: 1})

	err := user.SetStats()

	if err != nil {
		t.Fatal(err)
	}

	stats := user.GetStats()

	if len(stats) != int(common.ModeEnumMaxValue)-1 {
		t.Fatalf("expected (%v) mode stats. only fetched %v", int(common.ModeEnumMaxValue)-1, len(stats))
	}

	db.CloseSQLConnection()
}
