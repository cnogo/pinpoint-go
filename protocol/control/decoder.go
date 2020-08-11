package control

import (
	"encoding/binary"
	"io"
	"math"
)

type Decoder struct {
	buf    []byte
	rpos   int
	length int
}

func NewDecoder(payload []byte) *Decoder {
	return &Decoder{buf: payload, rpos: 0, length: len(payload)}
}

func (p *Decoder) readByte() (byte, error) {
	if p.rpos+1 > p.length {
		return 0, io.EOF
	}
	b := p.buf[p.rpos]
	p.rpos++
	return b, nil
}

func (p *Decoder) readFull(b []byte) (int, error) {
	preLen := len(b)
	if p.rpos+preLen > p.length {
		return 0, io.EOF
	}

	copy(b, p.buf[p.rpos:])
	p.rpos += preLen

	return preLen, nil
}

func (p *Decoder) lookupByte() (byte, error) {
	if p.rpos >= p.length {
		return 0, io.EOF
	}

	return p.buf[p.rpos], nil
}

func (p *Decoder) Decode() interface{} {
	Type, err := p.readByte()

	if err != nil {
		return nil
	}

	switch int32(Type) {
	case TYPE_NIL:
		return nil
	case TYPE_BOOL_TRUE:
		return true
	case TYPE_BOOL_FALSE:
		return false
	case TYPE_INT:
		if v, err := p.decodeInt32(); err != nil {
			return nil
		} else {
			return v
		}
	case TYPE_LONG:
		if v, err := p.decodeInt64(); err != nil {
			return nil
		} else {
			return v
		}
	case TYPE_DOUBLE:
		if v, err := p.decodeDouble(); err != nil {
			return nil
		} else {
			return v
		}
	case TYPE_STRING:
		if v, err := p.decodeString(); err != nil {
			return nil
		} else {
			return v
		}
	case TYPE_LIST_START:
		var s []interface{}
		for {
			tp, err := p.lookupByte()
			if err != nil {
				return nil
			}

			if int32(tp) == TYPE_LIST_END {
				break
			}
			v := p.Decode()
			s = append(s, v)
		}
		//上面for循环检测到下个标识为TYPE_LIST_END,所以得把它得读了。
		_, err := p.readByte()

		if err != nil {
			return nil
		}

		return s

	case TYPE_MAP_START:
		m := make(map[string]interface{})
		for {
			tp, err := p.lookupByte()
			if err != nil {
				return nil
			}
			if int32(tp) == TYPE_MAP_END {
				break
			}
			key := p.Decode()
			if key == nil {
				return nil
			}
			if _, ok := key.(string); !ok {
				return nil
			}
			value := p.Decode()
			m[key.(string)] = value
		}

		//上面for循环检测到下个标识为TYPE_MAP_END,所以得把它得读了。
		_, err := p.readByte()

		if err != nil {
			return nil
		}

		return m
	}

	return nil
}

func (p *Decoder) decodeInt32() (int32, error) {
	buf := make([]byte, 4)

	if _, err := p.readFull(buf); err != nil {
		return 0, err
	}
	return int32(binary.BigEndian.Uint32(buf)), nil
}

func (p *Decoder) decodeInt64() (int64, error) {
	buf := make([]byte, 8)

	if _, err := p.readFull(buf); err != nil {
		return 0, err
	}
	return int64(binary.BigEndian.Uint64(buf)), nil
}

func (p *Decoder) decodeDouble() (float64, error) {
	buf := make([]byte, 8)

	if _, err := p.readFull(buf); err != nil {
		return 0, err
	}

	return math.Float64frombits(binary.BigEndian.Uint64(buf)), nil
}

func (p *Decoder) decodeString() (string, error) {
	length, err := p.decodeLen()
	if err != nil {
		return "", err
	}
	buf := make([]byte, length)

	if _, err = p.readFull(buf); err != nil {
		return "", err
	}

	return string(buf), nil
}

func (p *Decoder) decodeLen() (int, error) {
	var result, shift uint
	for {
		b, err := p.readByte()
		if err != nil {
			return 0, err
		}

		result |= uint(b&0x7f) << shift

		if (b & 0x80) != 128 {
			break
		}
		shift += 7
	}

	return int(result), nil
}
