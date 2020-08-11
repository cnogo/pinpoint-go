package proto

import (
	"encoding/binary"
	"io"
	"net"
)

type ControlPong struct {
	Type int16
}

func NewControlPong() *ControlPong {
	return &ControlPong{
		Type: CONTROL_PING,
	}
}

func (c *ControlPong) Decode(conn net.Conn, reader io.Reader) error {

	return nil
}

// Encode ...
func (c *ControlPong) Encode() ([]byte, error) {
	body := make([]byte, 2)
	binary.BigEndian.PutUint16(body[0:2], uint16(c.Type))
	return body, nil
}

// GetPacketType ...
func (c *ControlPong) GetPacketType() int16 {
	return c.Type
}

// GetPayload ...
func (c *ControlPong) GetPayload() []byte {
	return nil
}

// GetRequestID ...
func (c *ControlPong) GetRequestID() int {
	return 0
}
