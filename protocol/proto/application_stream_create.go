package proto

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"
)

type ApplicationStreamCreate struct {
	Type      int16
	ChannelID int
	Length    int
	Payload   []byte
}

func NewApplicationStreamCreate() *ApplicationStreamCreate {
	return &ApplicationStreamCreate{
		Type: APPLICATION_STREAM_CREATE,
	}
}

func (p *ApplicationStreamCreate) Encode() ([]byte, error) {
	body := make([]byte, 10)
	binary.BigEndian.PutUint16(body[0:2], uint16(p.Type))
	binary.BigEndian.PutUint32(body[2:6], uint32(p.ChannelID))
	binary.BigEndian.PutUint32(body[6:10], uint32(p.Length))

	bys := bytes.NewBuffer(body)
	bys.Write(p.Payload)

	return bys.Bytes(), nil
}

func (p *ApplicationStreamCreate) Decode(conn net.Conn, reader io.Reader) error {
	buf := make([]byte, 8)

	if _, err := io.ReadFull(reader, buf); err != nil {
		log.Println("decode application request error")
		return err
	}

	p.ChannelID = int(binary.BigEndian.Uint32(buf[:4]))
	p.Length = int(binary.BigEndian.Uint32(buf[4:8]))

	if p.Length <= 0 {
		return nil
	}

	p.Payload = make([]byte, p.Length)

	if _, err := io.ReadFull(reader, p.Payload); err != nil {
		return err
	}

	return nil
}

func (p *ApplicationStreamCreate) GetPacketType() int16 {
	return p.Type
}

func (p *ApplicationStreamCreate) GetPayload() []byte {
	return p.Payload
}

func (p *ApplicationStreamCreate) GetRequestID() int {
	return 0
}
