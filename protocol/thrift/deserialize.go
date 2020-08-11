package thrift

import (
	"encoding/binary"
	"github.com/apache/thrift/lib/go/thrift"
	"sync"
)

type Deserializer struct {
	transportOut *thrift.TMemoryBuffer
	protocolOut  thrift.TProtocol
}

func NewDeserializer() *Deserializer {
	transportOut := thrift.NewTMemoryBuffer()
	protocolOut := thrift.NewTCompactProtocolFactory().GetProtocol(transportOut)
	return &Deserializer{transportOut: transportOut, protocolOut: protocolOut}
}

var deserializePool sync.Pool

func Deserialize(payload []byte) thrift.TStruct {
	var deserializer *Deserializer
	v := deserializePool.Get()
	if v == nil {
		deserializer = NewDeserializer()
	} else {
		deserializer = v.(*Deserializer)
	}

	deserializer.transportOut.Reset()
	_, err := deserializer.transportOut.Write(payload)
	if err != nil {
		return nil
	}
	header := readHeader(deserializer.protocolOut)
	if validae(header) {
		tStruct := TBaseLookup(header.Type)
		tStruct.Read(deserializer.protocolOut)
		deserializePool.Put(deserializer)
		return tStruct
	}

	deserializePool.Put(deserializer)
	return nil
}

func validae(header *Header) bool {
	if header.Signature == HEADER_SIGNATURE {
		return true
	}
	return false
}

func readHeader(protocol thrift.TProtocol) *Header {
	Signature, err0 := readByte(protocol)
	Version, err1 := readByte(protocol)
	byte1, err2 := readByte(protocol)
	byte2, err3 := readByte(protocol)
	if err0 != nil || err1 != nil || err2 != nil || err3 != nil {
		return nil
	}

	Type := bytesToInt16(byte1, byte2)
	return &Header{
		Signature: Signature,
		Version:   Version,
		Type:      Type,
	}
}

func bytesToInt16(byte1, byte2 byte) int16 {
	buf := make([]byte, 2)
	buf[0] = byte1
	buf[1] = byte2
	return int16(binary.BigEndian.Uint16(buf))
}

func readByte(protocol thrift.TProtocol) (byte, error) {
	buf, err := protocol.ReadByte()
	if err != nil {
		return byte(buf), err
	}

	return byte(buf), nil
}
