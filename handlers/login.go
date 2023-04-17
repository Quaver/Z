package handlers

import (
	"encoding/base64"
	"encoding/json"
	"example.com/Quaver/Z/config"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net"
	"net/http"
	"strconv"
	"strings"
)

// LoginData The data that the user sends to log in
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
	data, err := parseLoginData(r)

	if err != nil {
		return fmt.Errorf("[%v] login failed - %v", conn.RemoteAddr(), err)
	}

	err = authenticateSteamTicket(data)

	if err != nil {
		return fmt.Errorf("[%v] login failed - %v", conn.RemoteAddr(), err)
	}

	return nil
}

// Parses the raw data into a LoginData struct
func parseLoginData(r *http.Request) (*LoginData, error) {
	data := r.URL.Query().Get("login")

	if data == "" {
		return nil, fmt.Errorf("empty login data")
	}

	decoded, err := base64.StdEncoding.DecodeString(data)

	if err != nil {
		return nil, fmt.Errorf("failed to decode login data - %v", err)
	}

	var parsed LoginData
	err = json.Unmarshal(decoded, &parsed)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal login data - %v", err)
	}

	return &parsed, nil
}

// Authenticates the user via Steam. Makes sure the user has a valid id and ticket
func authenticateSteamTicket(data *LoginData) error {
	resp, err := resty.New().R().
		SetQueryParams(map[string]string{
			"key":    config.Instance.Steam.PublisherKey,
			"appid":  strconv.Itoa(config.Instance.Steam.AppId),
			"ticket": strings.Replace(data.PTicket, "-", "", -1),
		}).
		Get("https://api.steampowered.com/ISteamUserAuth/AuthenticateUserTicket/v1/")

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("%v failed to authenticate steam ticket - %v", resp.StatusCode(), string(resp.Body()))
	}

	type authenticateSteamTicketResponse struct {
		Response struct {
			Params struct {
				Result          string `json:"result,omitempty"`
				SteamId         string `json:"steamid,omitempty"`
				OwnerSteamId    string `json:"ownersteamid,omitempty"`
				VacBanned       bool   `json:"vacbanned,omitempty"`
				PublisherBanned bool   `json:"publisherbanned,omitempty"`
			} `json:"params"`
			Error interface{} `json:"error,omitempty"`
		} `json:"response"`
	}

	var parsed authenticateSteamTicketResponse
	err = json.Unmarshal(resp.Body(), &parsed)

	if err != nil {
		return fmt.Errorf("failed to authenticate steam ticket - json unmarshal - %v - %v", err, string(resp.Body()))
	}

	if parsed.Response.Error != nil || parsed.Response.Params.Result != "OK" {
		return fmt.Errorf("failed to authenticate steam ticket - invalid response result - %v", string(resp.Body()))
	}

	const failed string = "failed to authenticate steam ticket"

	if parsed.Response.Params.VacBanned {
		return fmt.Errorf("%v - user is vac banned", failed)
	}

	if parsed.Response.Params.PublisherBanned {
		return fmt.Errorf("%v - user is publisher banned", failed)
	}

	if parsed.Response.Params.SteamId != data.Id {
		return fmt.Errorf("%v - response steam id does not match the user provided id (%v vs. %v)", failed, parsed.Response.Params.SteamId, data.Id)
	}

	return nil
}
