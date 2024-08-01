package tests

import (
	"encoding/base64"
	"encoding/json"
	"example.com/Quaver/Z/config"
	"example.com/Quaver/Z/db"
	"example.com/Quaver/Z/handlers"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
)

type Client struct {
	User          *db.User
	Conn          *websocket.Conn
	enableLogging bool
}

func newClient(user *db.User) *Client {
	return &Client{User: user}
}

func (client *Client) host() string {
	return fmt.Sprintf("localhost:%v", config.Instance.Server.Port)
}

func (client *Client) loginData() handlers.LoginData {
	return handlers.LoginData{
		Id:        client.User.SteamId,
		PTicket:   "aGVsbG8=",
		PcbTicket: 0,
		Client:    "1|2|3|4|5",
	}
}

func (client *Client) loginDataEncoded() string {
	data := client.loginData()
	dataJson, _ := json.Marshal(data)

	var encoded = make([]byte, base64.StdEncoding.EncodedLen(len(dataJson)))
	base64.StdEncoding.Encode(encoded, dataJson)

	return string(encoded)
}
func (client *Client) login() {
	u := url.URL{
		Scheme:   "ws",
		Host:     client.host(),
		Path:     "/",
		RawQuery: fmt.Sprintf("login=%v", client.loginDataEncoded()),
	}

	var err error

	client.Conn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		log.Fatal("dial:", err)
	}

	defer client.Conn.Close()
	client.handleInterrupt(client.readMessages())
}

func (client *Client) readMessages() chan struct{} {
	done := make(chan struct{})

	go func() {
		defer close(done)

		for {
			_, message, err := client.Conn.ReadMessage()

			if err != nil {
				log.Println("read:", err)
				return
			}

			if client.enableLogging {
				log.Printf("recv: %s", message)
			}
		}
	}()

	return done
}

func (client *Client) handleInterrupt(doneReadingChan chan struct{}) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for {
		select {
		case <-doneReadingChan:
			return
		case <-interrupt:
			err := client.Conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

			if err != nil {
				log.Println("write close:", err)
				return
			}
		}
	}
}
