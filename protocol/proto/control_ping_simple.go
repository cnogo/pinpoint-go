package proto

import (
	"encoding/binary"
	"io"
	"net"
)

type ControlPingSimple struct {
	Type int16
}

func NewControlPingSimple() *ControlPingSimple {
	return &ControlPingSimple{
		Type: CONTROL_PING_SIMPLE,
	}
}

// Decode ...
func (c *ControlPingSimple) Decode(conn net.Conn, reader io.Reader) error {
	return nil
}

// Encode ...
func (c *ControlPingSimple) Encode() ([]byte, error) {
	body := make([]byte, 2)
	binary.BigEndian.PutUint16(body[0:2], uint16(c.Type))
	return body, nil
}

// GetPacketType ...
func (c *ControlPingSimple) GetPacketType() int16 {
	return c.Type
}

// GetPayload ...
func (c *ControlPingSimple) GetPayload() []byte {
	return nil
}

func (c *ControlPingSimple) GetRequestID() int {
	return 0
}
