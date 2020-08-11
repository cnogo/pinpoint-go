package proto

import (
	"encoding/binary"
	"io"
	"log"
	"net"
)

type ApplicationStreamCreateFail struct {
	Type      int16
	ChannelID int
	Code      int16
}

func NewApplicationStreamCreateFail() *ApplicationStreamCreateFail {
	return &ApplicationStreamCreateFail{
		Type: APPLICATION_STREAM_CREATE_FAIL,
	}
}

func (p *ApplicationStreamCreateFail) Encode() ([]byte, error) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint16(buf[0:2], uint16(p.Type))
	binary.BigEndian.PutUint32(buf[2:6], uint32(p.ChannelID))
	binary.BigEndian.PutUint16(buf[6:8], uint16(p.Code))

	return buf, nil
}

func (p *ApplicationStreamCreateFail) Decode(conn net.Conn, reader io.Reader) error {
	buf := make([]byte, 6)

	if _, err := io.ReadFull(reader, buf); err != nil {
		log.Println("decode application request error")
		return err
	}

	p.ChannelID = int(binary.BigEndian.Uint32(buf[:4]))
	p.Code = int16(binary.BigEndian.Uint16(buf[4:6]))

	return nil
}

func (p *ApplicationStreamCreateFail) GetPacketType() int16 {
	return p.Type
}

func (p *ApplicationStreamCreateFail) GetPayload() []byte {
	return nil
}

func (p *ApplicationStreamCreateFail) GetRequestID() int {
	return 0
}
