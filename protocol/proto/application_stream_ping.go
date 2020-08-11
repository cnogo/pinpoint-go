package proto

import (
	"encoding/binary"
	"io"
	"log"
	"net"
)

type ApplicationStreamPing struct {
	Type      int16
	ChannelID int
	RequestID int
}

// NewApplicationStreamPing ...
func NewApplicationStreamPing() *ApplicationStreamPing {
	return &ApplicationStreamPing{
		Type: APPLICATION_STREAM_PING,
	}
}

// Decode ...
func (a *ApplicationStreamPing) Decode(conn net.Conn, reader io.Reader) error {
	buf := make([]byte, 8)
	if _, err := io.ReadFull(reader, buf); err != nil {
		log.Println("application stream error ", err)
		return err
	}
	a.ChannelID = int(binary.BigEndian.Uint32(buf[:4]))
	a.RequestID = int(binary.BigEndian.Uint32(buf[4:8]))

	return nil
}

// Encode ...
func (a *ApplicationStreamPing) Encode() ([]byte, error) {
	buf := make([]byte, 10)
	binary.BigEndian.PutUint16(buf[0:2], uint16(a.Type))
	binary.BigEndian.PutUint32(buf[2:6], uint32(a.ChannelID))
	binary.BigEndian.PutUint32(buf[6:10], uint32(a.RequestID))

	return buf, nil
}

// GetPacketType ...
func (a *ApplicationStreamPing) GetPacketType() int16 {
	return a.Type
}

// GetPayload ...
func (a *ApplicationStreamPing) GetPayload() []byte {
	return nil
}

// GetRequestID ...
func (a *ApplicationStreamPing) GetRequestID() int {
	return a.RequestID
}
