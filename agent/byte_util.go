package agent

import (
	"bytes"
	"encoding/binary"
	"log"
)

type ByteUtil struct {
	buf *bytes.Buffer
}

func NewByteUtil(parseStr string) *ByteUtil {
	if parseStr == "" {
		return &ByteUtil{buf: new(bytes.Buffer)}
	}
	return &ByteUtil{buf: bytes.NewBufferString(parseStr)}
}

func (p *ByteUtil) GetBytes() []byte {
	return p.buf.Bytes()
}

func (p *ByteUtil) PutPrefixedString(str string) {
	if str == "" {
		p.PutVar32(-1)
	} else {
		p.PutSVar(int32(len(str)))
		p.PutString(str)
	}
}

func (p *ByteUtil) ReadPrefixedString() string {
	size := p.ReadSVarLen()
	if size == -1 || size == 0{
		return ""
	}

	return p.ReadString(size)
}

func (p *ByteUtil) PutString(str string) {
	p.buf.WriteString(str)
}

func (p ByteUtil) ReadString(size int32) string {
	bys := make([]byte, size)
	n, err := p.buf.Read(bys)
	if err != nil || n != int(size) {
		return ""
	}

	return string(bys)
}

func (p *ByteUtil) PutSVar(v int32) {
	value := p.intToZigZag(v)
	p.PutVar32(value)
}

func (p *ByteUtil) ReadSVarLen() int32 {
	return p.zigZagToInt(p.ReadVarInt32())
}

func (p *ByteUtil) PutVar32(v int32) {
	value := v
	for {
		if (value & ^0x7f) == 0 {
			p.buf.WriteByte(byte(value))
			return
		}

		bys := byte((byte(value) & byte(0x7f)) | byte(0x80))
		p.buf.WriteByte(bys)
		value = int32(uint32(value) >> 7)
	}
}

func (p *ByteUtil) ReadVarInt32() int32 {
	var result int32
	for i := 0; i < 6; i++ {
		value := p.ReadByte()
		result = result | int32(uint32(byte(value) & byte(0x7f)) << uint(7 * i))
		if int8(value) >= 0 {
			return result
		}
	}

	return result
}

func (p *ByteUtil) PutInt64(v int64) {
	bys := make([]byte, 8)
	binary.BigEndian.PutUint64(bys, uint64(v))
	p.buf.Write(bys)
}

func (p *ByteUtil) ReadInt64() int64 {
	bys := make([]byte, 8)
	n, err := p.buf.Read(bys)
	if err != nil || n != 8{
		return 0
	}
	result :=  binary.BigEndian.Uint64(bys)
	return int64(result)
}

func (p *ByteUtil) PutByte(v byte) {
	p.buf.WriteByte(v)
}

func (p *ByteUtil) ReadByte() byte {
	bys, err := p.buf.ReadByte()
	if err != nil {
		log.Println("read error")
	}
	return bys
}

func (p *ByteUtil) PutVar64(v int64) {
	value := v
	for {
		if (value & ^0x7f) == 0 {
			p.buf.WriteByte(byte(value))
			return
		}

		p.buf.WriteByte(byte((value & 0x7f) | 0x80))
		value = int64(uint64(value) >> 7)
	}
}

func (p *ByteUtil) ReadVar64() int64 {
	var result int64
	for i := 0; i < 10; i++ {
		value := p.ReadByte()
		result = int64(uint64(value & 0x7f) << uint(7 * i))  | result
		if int8(value) >= 0 {
			return result
		}
	}

	return result
}

func (p *ByteUtil) intToZigZag(n int32) int32 {
	return (n << 1) ^ (n >> 31)
}

func (p *ByteUtil) zigZagToInt(n int32) int32 {
	return int32(uint32(n)  >> 1) ^ -(n & 1)
}

