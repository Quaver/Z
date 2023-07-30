package db

import (
	"example.com/Quaver/Z/config"
	"log"
	"testing"
)

func TestGetRandomSongMap(t *testing.T) {
	_ = config.Load("../config.json")

	if config.Instance == nil {
		return
	}

	InitializeSQL()

	song, err := GetRandomSongMap(0, 100)

	if err != nil {
		t.Fatal(err)
	}

	log.Println(song)
}
