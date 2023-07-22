package packets

type ServerSongRequest struct {
	Packet
	SongRequest
}

type SongRequestGame int

const (
	SongRequestGameQuaver SongRequestGame = iota
	SongRequestGameOsu
)

type SongRequest struct {
	TwitchUsername   string          `json:"tw"`
	UserId           int             `json:"u"`
	Game             SongRequestGame `json:"g"`
	MapId            int             `json:"mid"`
	MapsetId         int             `json:"msid"`
	MapMd5           string          `json:"md5"`
	Artist           string          `json:"a"`
	Title            string          `json:"t"`
	DifficultyName   string          `json:"d"`
	Creator          string          `json:"c"`
	DifficultyRating float64         `json:"r"`
}

func NewServerSongRequest(request SongRequest) *ServerSongRequest {
	return &ServerSongRequest{
		Packet:      Packet{Id: PacketIdServerSongRequest},
		SongRequest: request,
	}
}
