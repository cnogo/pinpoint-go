package proto

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"
)

type ApplicationResponse struct {
	Type      int16
	RequestID int
	Length    int
	Payload   []byte
}

func NewApplicationResponse() *ApplicationResponse {
	return &ApplicationResponse{
		Type: APPLICATION_RESPONSE,
	}
}

func (p *ApplicationResponse) Encode() ([]byte, error) {
	body := make([]byte, 10)
	binary.BigEndian.PutUint16(body[0:2], uint16(p.Type))
	binary.BigEndian.PutUint32(body[2:6], uint32(p.RequestID))
	binary.BigEndian.PutUint32(body[6:10], uint32(len(p.Payload)))
	bys := bytes.NewBuffer(body)
	bys.Write(p.Payload)

	return bys.Bytes(), nil
}

func (p *ApplicationResponse) Decode(conn net.Conn, reader io.Reader) error {
	buf := make([]byte, 8)

	if _, err := io.ReadFull(reader, buf); err != nil {
		log.Println("decode application request error")
		return err
	}

	p.RequestID = int(binary.BigEndian.Uint32(buf[0:4]))
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

func (p *ApplicationResponse) GetPacketType() int16 {
	return p.Type
}

func (p *ApplicationResponse) GetPayload() []byte {
	return p.Payload
}

func (p *ApplicationResponse) GetRequestID() int {
	return p.RequestID
}
