package sessions

import (
	"example.com/Quaver/Z/db"
	"example.com/Quaver/Z/utils"
	"net"
)

type User struct {
	// The connection for the user
	Conn net.Conn

	// The token used to identify the user for requests.
	Token string

	// All user table information from the database
	Info *db.User
}

// NewUser Creates a new user session struct object
func NewUser(conn net.Conn, user *db.User) *User {
	return &User{
		Conn:  conn,
		Token: utils.GenerateRandomString(64),
		Info:  user,
	}
}
