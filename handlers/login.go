package login

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

type Data struct {
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
	data, err := parseLoginData(r)

	if err != nil {
		return fmt.Errorf("[%v] login failed - %v", conn.RemoteAddr(), err)
	}

	fmt.Println(data)

	err = authenticateSteamTicket(data)

	if err != nil {
		return fmt.Errorf("[%v] login failed - %v", conn.RemoteAddr(), err)
	}

	return nil
}

// Parses the raw data into a LoginData struct
func parseLoginData(r *http.Request) (*Data, error) {
	data := r.URL.Query().Get("login")

	if data == "" {
		return nil, fmt.Errorf("empty login data")
	}

	decoded, err := base64.StdEncoding.DecodeString(data)

	if err != nil {
		return nil, fmt.Errorf("failed to decode login data - %v", err)
	}

	var parsed Data
	err = json.Unmarshal(decoded, &parsed)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal login data - %v", err)
	}

	return &parsed, nil
}
