package proto

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"
)

type ControlHandShake struct {
	Type      int16
	RequestID int
	Length    int
	Payload   []byte
}

func NewControlHandShake() *ControlHandShake {
	return &ControlHandShake{
		Type: CONTROL_HANDSHAKE,
	}
}

// Decode ...
func (c *ControlHandShake) Decode(conn net.Conn, reader io.Reader) error {
	buf := make([]byte, 8)
	if _, err := io.ReadFull(reader, buf); err != nil {
		log.Println("control hanshake decode error ", err)
		return err
	}
	c.RequestID = int(binary.BigEndian.Uint32(buf[:4]))
	c.Length = int(binary.BigEndian.Uint32(buf[4:8]))

	if c.Length <= 0 {
		return nil
	}

	c.Payload = make([]byte, c.Length)

	if _, err := io.ReadFull(reader, c.Payload); err != nil {
		return err
	}

	//log.Println(string(c.Payload))
	return nil
}

// Encode ...
func (c *ControlHandShake) Encode() ([]byte, error) {
	buf := make([]byte, 10)
	binary.BigEndian.PutUint16(buf[0:2], uint16(c.Type))
	binary.BigEndian.PutUint32(buf[2:6], uint32(c.RequestID))
	binary.BigEndian.PutUint32(buf[6:10], uint32(len(c.Payload)))
	bys := bytes.NewBuffer(buf)
	bys.Write(c.Payload)

	return bys.Bytes(), nil
}

// GetPacketType ...
func (c *ControlHandShake) GetPacketType() int16 {
	return c.Type
}

// GetPayload ...
func (c *ControlHandShake) GetPayload() []byte {
	return c.Payload
}

func (c *ControlHandShake) GetRequestID() int {
	return c.RequestID
}
