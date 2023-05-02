package db

import (
	"example.com/Quaver/Z/config"
	"fmt"
	"testing"
)

func TestFetchProcesses(t *testing.T) {
	_ = config.Load("../config.json")

	if config.Instance == nil {
		return
	}

	InitializeSQL()

	processes, err := FetchProcesses()

	if err != nil {
		t.Fatal(err)
	}

	if len(processes) == 0 {
		t.Fatal("expected more than zero processes")
	}

	for i, p := range processes {
		fmt.Printf("[#%v] %v | %v\n", i, p.Name, p.GetMD5())
	}

	CloseSQLConnection()
}
