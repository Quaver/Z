package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

type LoginData struct {
	// The Steam ID of the user
	Id string `json:"id"`

	// Steam authentication PTicket
	PTicket string `json:"ticket"`

	// Steam authentication PCBTicket
	PcbTicket byte `json:"pcb"`

	// Game Client file signatures
	Client string `json:"client"`
}

// HandleLogin Handles the login of a client
func HandleLogin(conn net.Conn, r *http.Request) error {
	data, err := parseLoginData(conn, r)

	if err != nil {
		return fmt.Errorf("[%v] failed to login - %v", conn.RemoteAddr(), err)
	}

	fmt.Println(data)
	return nil
}

// Parses the raw data into a LoginData struct
func parseLoginData(conn net.Conn, r *http.Request) (*LoginData, error) {
	data := r.URL.Query().Get("login")

	if data == "" {
		return nil, fmt.Errorf("[%v] failed to login - empty login data", conn.RemoteAddr())
	}

	decoded, err := base64.StdEncoding.DecodeString(data)

	if err != nil {
		return nil, fmt.Errorf("[%v] failed to decode login data - %v", conn.RemoteAddr(), err)
	}

	var parsed LoginData

	err = json.Unmarshal(decoded, &parsed)

	if err != nil {
		return nil, fmt.Errorf("[%v] failed to unmarshal login data - %v", conn.RemoteAddr(), err)
	}

	return &parsed, nil
}
