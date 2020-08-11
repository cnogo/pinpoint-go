package proto

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"
)

type ApplicationStreamResponse struct {
	Type      int16
	ChannelID int
	Length    int
	Payload   []byte
}

func NewApplicationStreamResponse() *ApplicationStreamResponse {
	return &ApplicationStreamResponse{
		Type: APPLICATION_STREAM_RESPONSE,
	}
}

// Decode ...
func (a *ApplicationStreamResponse) Decode(conn net.Conn, reader io.Reader) error {
	buf := make([]byte, 8)
	if _, err := io.ReadFull(reader, buf); err != nil {
		log.Println("application stream response error ", err)
		return err
	}
	a.ChannelID = int(binary.BigEndian.Uint32(buf[:4]))
	a.Length = int(binary.BigEndian.Uint32(buf[4:8]))

	if a.Length <= 0 {
		return nil
	}

	a.Payload = make([]byte, a.Length)
	if _, err := io.ReadFull(reader, a.Payload); err != nil {
		return err
	}
	//log.Println(string(a.Payload))
	return nil
}

// Encode ...
func (a *ApplicationStreamResponse) Encode() ([]byte, error) {
	body := make([]byte, 10)
	binary.BigEndian.PutUint16(body[0:2], uint16(a.Type))
	binary.BigEndian.PutUint32(body[2:6], uint32(a.ChannelID))
	binary.BigEndian.PutUint32(body[6:10], uint32(a.Length))

	bys := bytes.NewBuffer(body)
	bys.Write(a.Payload)

	return bys.Bytes(), nil
}

// GetPacketType ...
func (a *ApplicationStreamResponse) GetPacketType() int16 {
	return a.Type
}

// GetPayload ...
func (a *ApplicationStreamResponse) GetPayload() []byte {
	return a.Payload
}

// GetRequestID ...
func (a *ApplicationStreamResponse) GetRequestID() int {
	return 0
}
