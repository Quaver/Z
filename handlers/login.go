package handlers

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"example.com/Quaver/Z/chat"
	"example.com/Quaver/Z/common"
	"example.com/Quaver/Z/config"
	"example.com/Quaver/Z/db"
	"example.com/Quaver/Z/packets"
	"example.com/Quaver/Z/sessions"
	"example.com/Quaver/Z/utils"
	"example.com/Quaver/Z/webhooks"
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
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

	// If the user is reconnecting to the server
	IsReconnection bool `json:"is_reconnection"`
}

// HandleLogin Handles the login of a client
func HandleLogin(conn net.Conn, r *http.Request) error {
	data, err := parseLoginData(r)

	if err != nil {
		return logFailedLogin(conn, err)
	}

	err = authenticateSteamTicket(data)

	if err != nil {
		return logFailedLogin(conn, err)
	}

	err = checkSteamAppOwnership(data.Id)

	if err != nil {
		return logFailedLogin(conn, err)
	}

	user, err := db.GetUserBySteamId(data.Id)

	if err != nil {
		// User does not exist yet, so prompt them to select a username for their account.
		if err == sql.ErrNoRows {
			_ = sessions.SendPacketToConnection(packets.NewServerChooseUsername(), conn)
			utils.CloseConnectionDelayed(conn)
			log.Printf("[%v] %v logged in but does not have an account yet.\n", conn.RemoteAddr(), data.Id)
			return nil
		}

		return err
	}

	if !user.Allowed {
		_ = sessions.SendPacketToConnection(packets.NewServerNotificationError("You are banned. You can appeal your ban at: discord.gg/quaver"), conn)
		utils.CloseConnectionDelayed(conn)
		log.Printf("[%v - #%v] Attempted to login, but they are banned\n", user.Username, user.Id)
		return nil
	}

	err = verifyGameBuild(data)

	if err != nil {
		if err == sql.ErrNoRows {
			if !handleCustomGameBuildUsage(conn, user, data.Client) {
				return nil
			}
		} else {
			return err
		}
	}

	ip := conn.RemoteAddr().String()
	ipHeader := r.Header.Get("X-Forwarded-For")

	if ipHeader != "" {
		ip = ipHeader
	}

	err = db.InsertLoginIpAddress(user.Id, ip)

	if err != nil {
		return err
	}

	err = db.UpdateUserLatestActivity(user.Id)

	if err != nil {
		return err
	}

	err = updateUserAvatar(user)

	if err != nil {
		log.Println("Failed to update steam avatar: ", err)
	}

	sessionUser, err := getOrCreateSession(conn, user, data.IsReconnection)

	if err != nil {
		return err
	}

	err = sessionUser.SetStats()

	if err != nil {
		return err
	}

	err = sessions.AddUser(sessionUser)

	if err != nil {
		return err
	}

	err = sendLoginPackets(sessionUser)

	if err != nil {
		return err
	}

	log.Printf("[%v #%v] Logged in (%v users online).\n", user.Username, user.Id, sessions.GetOnlineUserCount())
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
	if config.Instance.BypassSteamLogin {
		return nil
	}

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

	const failed string = "failed to authenticate steam ticket"

	if resp.IsError() {
		return fmt.Errorf("%v %v - %v", failed, resp.StatusCode(), string(resp.Body()))
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
		return fmt.Errorf("%v - json unmarshal - %v - %v", failed, err, string(resp.Body()))
	}

	if parsed.Response.Error != nil || parsed.Response.Params.Result != "OK" {
		return fmt.Errorf("%v - invalid response result - %v", failed, string(resp.Body()))
	}

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

// Checks if the user actually owns the game on Steam
func checkSteamAppOwnership(steamId string) error {
	if config.Instance.BypassSteamLogin {
		return nil
	}

	resp, err := resty.New().R().
		SetQueryParams(map[string]string{
			"key":     config.Instance.Steam.PublisherKey,
			"appid":   strconv.Itoa(config.Instance.Steam.AppId),
			"steamid": steamId,
		}).
		Get("https://partner.steam-api.com/ISteamUser/CheckAppOwnership/v2/")

	if err != nil {
		return err
	}

	const failed string = "failed to check steam app ownership"

	if resp.IsError() {
		return fmt.Errorf("%v %v - %v", resp.StatusCode(), failed, string(resp.Body()))
	}

	type checkSteamAppOwnershipResponse struct {
		AppOwnership struct {
			OwnsApp bool `json:"ownsapp"`
		} `json:"appownership,omitempty"`

		Error interface{} `json:"error,omitempty"`
	}

	var parsed checkSteamAppOwnershipResponse
	err = json.Unmarshal(resp.Body(), &parsed)

	if err != nil {
		return fmt.Errorf("%v - json unmarshal - %v - %v", failed, err, string(resp.Body()))
	}

	if parsed.Error != nil {
		return fmt.Errorf("%v - %v", failed, string(resp.Body()))
	}

	if !parsed.AppOwnership.OwnsApp {
		return fmt.Errorf("%v - user does not own Quaver on Steam", failed)
	}

	return nil
}

// Checks the client signatures to see if the build they are using is valid
func verifyGameBuild(data *LoginData) error {
	split := strings.Split(data.Client, "|")

	if len(split) != 5 {
		return fmt.Errorf("user provided an incorrect amount of client signatures - %v", data.Client)
	}

	err := db.VerifyGameBuild(db.GameBuild{
		QuaverAPIDll:          split[1],
		QuaverServerClientDll: split[2],
		QuaverServerCommonDll: split[3],
		QuaverSharedDll:       split[4],
	})

	if err != nil {
		return err
	}

	return nil
}

// Formats an invalid client build into a readable json
func formatCustomGameBuild(client string) string {
	split := strings.Split(client, "|")

	if len(split) != 5 {
		return client
	}

	data, err := json.MarshalIndent(struct {
		QuaverAPIDll          string `json:"quaver_api_dll"`
		QuaverServerClientDll string `json:"quaver_server_client_dll"`
		QuaverServerCommonDll string `json:"quaver_server_common_dll"`
		QuaverSharedDll       string `json:"quaver_shared_dll"`
	}{
		QuaverAPIDll:          split[1],
		QuaverServerClientDll: split[2],
		QuaverServerCommonDll: split[3],
		QuaverSharedDll:       split[4],
	}, "", "    ")

	if err != nil {
		return client
	}

	return string(data)
}

// Returns if a user is eligible to use a custom client
func canUserUseCustomGameBuild(user *db.User) bool {
	return common.HasUserGroup(user.UserGroups, common.UserGroupSwan) ||
		common.HasUserGroup(user.UserGroups, common.UserGroupDeveloper) ||
		common.HasUserGroup(user.UserGroups, common.UserGroupAdmin) ||
		common.HasUserGroup(user.UserGroups, common.UserGroupContributor)
}

// Sends webhook and disconnects a user for invalid client usage. Returns if the user is allowed to login
func handleCustomGameBuildUsage(conn net.Conn, user *db.User, client string) bool {
	clientStr := fmt.Sprintf("```json\n%v```", formatCustomGameBuild(client))
	webhooks.SendAntiCheat(user.Username, user.Id, user.GetProfileUrl(), user.AvatarUrl.String, "Invalid Game Build", clientStr)

	if !canUserUseCustomGameBuild(user) {
		_ = sessions.SendPacketToConnection(packets.NewServerNotificationError("Please update your client before attempting to login."), conn)
		utils.CloseConnectionDelayed(conn)
		return false
	}

	log.Printf("[%v - #%v] Attempted to login, but is using an game build: %v\n", user.Username, user.Id, clientStr)
	return true
}

// Updates the avatar for the user and sets the new one.
func updateUserAvatar(user *db.User) error {
	if config.Instance.BypassSteamLogin {
		return nil
	}

	avatar, err := db.UpdateUserSteamAvatar(user.SteamId)

	if err != nil {
		return err
	}

	// Make sure the avatar is the most up to date version
	user.AvatarUrl.String = avatar
	return nil
}

// Checks to see if a user had a previous session and returns that, otherwise
func getOrCreateSession(conn net.Conn, user *db.User, reconnection bool) (sessionUser *sessions.User, err error) {
	sessionUser = sessions.GetUserById(user.Id)

	// User is logging in for the first time.
	if sessionUser == nil {
		return sessions.NewUser(conn, user), nil
	}

	// User is reconnecting to the server, so update their existing session connection with the incoming one.
	if reconnection && sessionUser.GetLastTemporaryDisconnectionTimestamp() != -1 {
		sessionUser.SetLastTemporaryDisconnectionTimestamp(-1)
		sessionUser.Conn = conn

		log.Println("User", user.Username, "reconnected after temporary disconnection")
		return sessionUser, nil
	}

	// User is logging in from another location, so disconnect the old and connect the new.
	return sessions.NewUser(conn, user), removePreviousLoginSession(sessionUser)
}

// Checks to see if the user is already logged in and removes the previous session
func removePreviousLoginSession(u *sessions.User) error {
	sessions.SendPacketToUser(packets.NewServerNotificationError("You are being logged out due to logging in from a different location"), u)

	err := sessions.RemoveUser(u)

	if err != nil {
		return err
	}

	utils.CloseConnectionDelayed(u.Conn)
	return nil
}

// Sends initial packets to log the user in
func sendLoginPackets(user *sessions.User) error {
	sessions.SendPacketToUser(packets.NewServerLoginReply(user.SerializeForPacket(), user.GetStatsSlice(), user.GetToken()), user)
	sessions.SendPacketToUser(packets.NewServerUsersOnline(sessions.GetOnlineUserIds()), user)
	sessions.SendPacketToUser(packets.NewServerUserInfo(sessions.GetSerializedOnlineUsers()), user)
	sessions.SendPacketToUser(packets.NewServerTwitchConnection(user.Info.TwitchUsername.String), user)
	sessions.SendPacketToAllUsers(packets.NewServerUserConnected(user.SerializeForPacket()))
	joinChatChannels(user)

	friends, err := db.GetUserFriendsList(user.Info.Id)

	if err != nil {
		return err
	}

	sessions.SendPacketToUser(packets.NewServerFriendsList(friends), user)
	return nil
}

// Joins an available chat channel
func joinChatChannels(user *sessions.User) {
	channels := chat.GetAvailableChannels(user.Info.UserGroups)

	for _, channel := range channels {
		sessions.SendPacketToUser(packets.NewServerAvailableChatChannel(channel.Name, channel.Description), user)

		if channel.AutoJoin {
			channel.AddUser(user)
		}
	}
}

// Logs a generic login failure
func logFailedLogin(conn net.Conn, err error) error {
	return fmt.Errorf("[%v] login failed - %v", conn.RemoteAddr(), err)
}
