package proto

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"
)

type ApplicationSend struct {
	Type    int16
	Length  int
	Payload []byte
}

func NewApplicationSend() *ApplicationSend {
	return &ApplicationSend{
		Type: APPLICATION_SEND,
	}
}

func (p *ApplicationSend) Encode() ([]byte, error) {
	body := make([]byte, 6)
	binary.BigEndian.PutUint16(body[0:2], uint16(p.Type))
	binary.BigEndian.PutUint32(body[2:6], uint32(len(p.Payload)))
	bys := bytes.NewBuffer(body)
	bys.Write(p.Payload)

	return bys.Bytes(), nil
}

func (p *ApplicationSend) Decode(conn net.Conn, reader io.Reader) error {
	buf := make([]byte, 4)

	if _, err := io.ReadFull(reader, buf); err != nil {
		log.Println("decode application request error")
		return err
	}

	p.Length = int(binary.BigEndian.Uint32(buf[:4]))

	if p.Length <= 0 {
		return nil
	}

	p.Payload = make([]byte, p.Length)

	if _, err := io.ReadFull(reader, p.Payload); err != nil {
		return err
	}

	return nil
}

func (p *ApplicationSend) GetPacketType() int16 {
	return p.Type
}

func (p *ApplicationSend) GetPayload() []byte {
	return p.Payload
}

func (p *ApplicationSend) GetRequestID() int {
	return 0
}
