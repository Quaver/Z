package packets

import (
	"encoding/json"
	"log"
)

type ClientPong struct {
	Packet
	ProcessList string `json:"p"`
}

type Processes struct {
	Processes []Process `json:"Processes"`
}

type Process struct {
	Name        string `json:"Name"`
	WindowTitle string `json:"WindowTitle"`
	FileName    string `json:"FileName"`
}

func (p *ClientPong) ParseProcessList() []Process {
	var data Processes

	err := json.Unmarshal([]byte(p.ProcessList), &data)

	if err != nil {
		log.Println(err)
		return nil
	}

	return data.Processes
}
