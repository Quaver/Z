package utils

import (
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"io"
	"io/ioutil"
	"net"
	"time"
)

// CloseConnection Closes a connection
func CloseConnection(conn net.Conn) {
	if conn == nil {
		return
	}

	var body = ws.NewCloseFrameBody(1000, "")
	var frame = ws.NewCloseFrame(body)
	if err := ws.WriteHeader(conn, frame.Header); err != nil {
		return
	}
	if _, err := conn.Write(body); err != nil {
		return
	}
	_ = conn.Close()
}

// CloseConnectionDelayed Closes the connection after a specified amount of time
func CloseConnectionDelayed(conn net.Conn) {
	time.AfterFunc(250*time.Millisecond, func() {
		CloseConnection(conn)
	})
}

func ReadData(rw io.ReadWriter, s ws.State, want ws.OpCode) ([]byte, ws.OpCode, error) {
	controlHandler := wsutil.ControlFrameHandler(rw, s)
	rd := wsutil.Reader{
		Source:          rw,
		State:           s,
		CheckUTF8:       true,
		SkipHeaderCheck: false,
		OnIntermediate:  controlHandler,
	}
	for {
		hdr, err := rd.NextFrame()
		if err != nil {
			return nil, 0, err
		}
		if hdr.OpCode.IsControl() {
			if err := controlHandler(hdr, &rd); err != nil {
				return nil, 0, err
			}
			if hdr.OpCode&want != 0 {
				return nil, hdr.OpCode, err
			}
			continue
		}
		if hdr.OpCode&want == 0 {
			if err := rd.Discard(); err != nil {
				return nil, 0, err
			}
			continue
		}

		bts, err := ioutil.ReadAll(&rd)

		return bts, hdr.OpCode, err
	}
}
