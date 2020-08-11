package control

import (
	"bytes"
	"encoding/binary"
	"math"
)

type Encoder struct {
	buf *bytes.Buffer
}

//New Encoder ...
func NewEncoder() *Encoder {
	return &Encoder{buf: &bytes.Buffer{}}
}

//Encode ...
func (p *Encoder) Encode(data map[string]interface{}) []byte {
	p.buf.Reset()
	p.encode(data)
	return p.buf.Bytes()
}

func (p *Encoder) encode(data interface{}) {
	if data == nil {
		p.encodeNil()
	}

	switch data.(type) {
	case int16:
		p.encodeInt32(int32(data.(int16)))
	case int:
		p.encodeInt32(int32(data.(int)))
	case int32:
		p.encodeInt32(data.(int32))
	case int64:
		p.encodeInt64(data.(int64))
	case string:
		p.encodeString(data.(string))
	case bool:
		p.encodeBool(data.(bool))
	case float32:
		p.encodeDouble(float64(data.(float32)))
	case float64:
		p.encodeDouble(data.(float64))
	case []interface{}:
		p.encodeSlice(data.([]interface{}))
	case map[string]interface{}:
		p.encodeMap(data.(map[string]interface{}))
	}
}

func (p *Encoder) encodeSlice(s []interface{}) {
	p.buf.WriteByte(byte(TYPE_LIST_START))
	for value := range s {
		p.encode(value)
	}
	p.buf.WriteByte(byte(TYPE_LIST_END))
}

func (p *Encoder) encodeMap(m map[string]interface{}) {
	p.buf.WriteByte(byte(TYPE_MAP_START))

	for key, value := range m {
		p.encodeString(key)
		p.encode(value)
	}

	p.buf.WriteByte(byte(TYPE_MAP_END))
}

func (p *Encoder) encodeInt32(i int32) {
	p.buf.WriteByte(byte(TYPE_INT))
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(i))
	p.buf.Write(buf)
}

func (p *Encoder) encodeString(s string) {
	p.buf.WriteByte(byte(TYPE_STRING))
	p.encodeLen(uint32(len(s)))
	p.buf.WriteString(s)
}

func (p *Encoder) encodeNil() {
	p.buf.WriteByte(byte(TYPE_NIL))
}

func (p *Encoder) encodeBool(v bool) {
	if v {
		p.buf.WriteByte(byte(TYPE_BOOL_TRUE))
	} else {
		p.buf.WriteByte(byte(TYPE_BOOL_FALSE))
	}
}

func (p *Encoder) encodeInt64(v int64) {
	p.buf.WriteByte(byte(TYPE_LONG))
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(v))
	p.buf.Write(buf)
}

func (p *Encoder) encodeDouble(value float64) {
	p.buf.WriteByte(byte(TYPE_DOUBLE))
	longValue := math.Float64bits(value)
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, longValue)
	p.buf.Write(buf)
}

func (p *Encoder) encodeLen(len uint32) {
	lenBuf := make([]byte, 5)
	idx := 0
	for {
		if (len & 0xffffff80) == 0 {
			lenBuf[idx] = byte(len)
			idx++
			break

		}
		lenBuf[idx] = byte(len&0x7f | 0x80)
		idx++
		len >>= 7
	}

	p.buf.Write(lenBuf[:idx])
}
