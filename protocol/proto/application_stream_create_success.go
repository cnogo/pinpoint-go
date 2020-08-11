package proto

import (
	"encoding/binary"
	"io"
	"log"
	"net"
)

type ApplicationStreamCreateSuccess struct {
	Type      int16
	ChannelID int
}

func NewApplicationStreamCreateSuccess() *ApplicationStreamCreateSuccess {
	return &ApplicationStreamCreateSuccess{
		Type: APPLICATION_STREAM_CREATE_SUCCESS,
	}
}

// Decode ...
func (p *ApplicationStreamCreateSuccess) Decode(conn net.Conn, reader io.Reader) error {
	buf := make([]byte, 4)
	if _, err := io.ReadFull(reader, buf); err != nil {
		log.Println("decode application request error")
		return err
	}
	p.ChannelID = int(binary.BigEndian.Uint32(buf[:4]))
	return nil
}

// Encode ...
func (p *ApplicationStreamCreateSuccess) Encode() ([]byte, error) {
	buf := make([]byte, 6)
	binary.BigEndian.PutUint16(buf[0:2], uint16(p.Type))
	binary.BigEndian.PutUint16(buf[2:4], uint16(p.ChannelID))
	return buf, nil
}

// GetPacketType ...
func (p *ApplicationStreamCreateSuccess) GetPacketType() int16 {
	return p.Type
}

// GetPayload ...
func (p *ApplicationStreamCreateSuccess) GetPayload() []byte {
	return nil
}

// GetRequestID ...
func (p *ApplicationStreamCreateSuccess) GetRequestID() int {
	return 0
}
