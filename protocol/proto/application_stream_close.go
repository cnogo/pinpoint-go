package proto

import (
	"encoding/binary"
	"io"
	"log"
	"net"
)

type ApplicationStreamClose struct {
	Type      int16
	ChannelID int
	Code      int16
}

func NewApplicationStreamClose() *ApplicationStreamClose {
	return &ApplicationStreamClose{
		Type: APPLICATION_STREAM_CLOSE,
	}
}

func (p *ApplicationStreamClose) Encode() ([]byte, error) {
	body := make([]byte, 8)
	binary.BigEndian.PutUint16(body[0:2], uint16(p.Type))
	binary.BigEndian.PutUint32(body[2:6], uint32(p.ChannelID))
	binary.BigEndian.PutUint16(body[6:8], uint16(p.Code))

	return body, nil
}

func (p *ApplicationStreamClose) Decode(conn net.Conn, reader io.Reader) error {
	buf := make([]byte, 6)

	if _, err := io.ReadFull(reader, buf); err != nil {
		log.Println("decode application request error")
		return err
	}

	p.ChannelID = int(binary.BigEndian.Uint32(buf[:4]))
	p.Code = int16(binary.BigEndian.Uint16(buf[4:6]))

	return nil
}

func (p *ApplicationStreamClose) GetPacketType() int16 {
	return p.Type
}

func (p *ApplicationStreamClose) GetPayload() []byte {
	return nil
}

func (p *ApplicationStreamClose) GetRequestID() int {
	return 0
}
