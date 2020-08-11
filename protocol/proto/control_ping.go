package proto

import (
	"encoding/binary"
	"io"
	"log"
	"net"
)

type ControlPing struct {
	Type         int16
	PingID       int
	StateVersion byte
	StateCode    byte
}

func NewControlPing() *ControlPing {
	return &ControlPing{
		Type: CONTROL_PING,
	}
}

func (c *ControlPing) Decode(conn net.Conn, reader io.Reader) error {
	buf := make([]byte, 6)
	if _, err := io.ReadFull(reader, buf); err != nil {
		log.Println("control ping payload error ", err)
		return err
	}

	c.PingID = int(binary.BigEndian.Uint32(buf[:4]))
	c.StateVersion = buf[4]
	c.StateCode = buf[5]

	return nil
}

// Encode ...
func (c *ControlPing) Encode() ([]byte, error) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint16(buf[0:2], uint16(c.Type))
	binary.BigEndian.PutUint32(buf[2:6], uint32(c.PingID))
	buf[6] = c.StateVersion
	buf[7] = c.StateCode
	return buf, nil
}

// GetPacketType ...
func (c *ControlPing) GetPacketType() int16 {
	return c.Type
}

// GetPayload ...
func (c *ControlPing) GetPayload() []byte {
	return nil
}

// GetRequestID ...
func (c *ControlPing) GetRequestID() int {
	return 0
}
