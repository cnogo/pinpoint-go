package proto

import (
	"encoding/binary"
	"io"
	"log"
	"net"
)

type ApplicationTraceSendAck struct {
	Type    int16
	TraceID int
}

func NewApplicationTraceSendAck() *ApplicationTraceSendAck {
	return &ApplicationTraceSendAck{
		Type: APPLICATION_TRACE_SEND_ACK,
	}
}

func (c *ApplicationTraceSendAck) Decode(conn net.Conn, reader io.Reader) error {
	buf := make([]byte, 4)

	if _, err := io.ReadFull(reader, buf); err != nil {
		log.Println("decode application request error")
		return err
	}

	c.TraceID = int(binary.BigEndian.Uint32(buf[0:4]))
	return nil
}

// Encode ...
func (c *ApplicationTraceSendAck) Encode() ([]byte, error) {
	body := make([]byte, 6)
	binary.BigEndian.PutUint16(body[0:2], uint16(c.Type))
	binary.BigEndian.PutUint32(body[2:6], uint32(c.TraceID))

	return body, nil
}

// GetPacketType ...
func (c *ApplicationTraceSendAck) GetPacketType() int16 {
	return c.Type
}

// GetPayload ...
func (c *ApplicationTraceSendAck) GetPayload() []byte {
	return nil
}

// GetRequestID ...
func (c *ApplicationTraceSendAck) GetRequestID() int {
	return 0
}
