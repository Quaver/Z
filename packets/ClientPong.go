package packets

import (
	"encoding/json"
	"log"
)

type ClientPong struct {
	Packet
	Data string `json:"p"`
}

type PongPacketData struct {
	Processes []Process `json:"Processes"`
	Libraries []string  `json:"Libraries"`
}

type Process struct {
	Name        string `json:"Name"`
	WindowTitle string `json:"WindowTitle"`
	FileName    string `json:"FileName"`
}

func (p *ClientPong) Parse() *PongPacketData {
	var data PongPacketData

	err := json.Unmarshal([]byte(p.Data), &data)

	if err != nil {
		log.Println(err)
		return nil
	}

	return &data
}
