package tests

import (
	"example.com/Quaver/Z/config"
	"example.com/Quaver/Z/db"
	"sync"
	"testing"
)

// Tests a 100 users connecting to the server at one time
func TestLogin100Clients(t *testing.T) {
	if err := config.Load("../../config.json"); err != nil {
		t.Fatal(err)
	}

	db.InitializeSQL()

	wg := sync.WaitGroup{}

	for i := 3; i < 103; i++ {
		wg.Add(1)

		i := i

		go func() {
			defer wg.Done()

			user, err := db.GetUserById(i)

			if err != nil {
				t.Error(err)
				return
			}

			client := newClient(user)
			client.login()
		}()
	}

	wg.Wait()
	db.CloseSQLConnection()
}
