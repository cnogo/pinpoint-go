package thrift

import (
	"encoding/binary"
	"github.com/apache/thrift/lib/go/thrift"
	"sync"
)

type Serializer struct {
	transportOut *thrift.TMemoryBuffer
	protocolOut  thrift.TProtocol
}

func NewSerializer() *Serializer {
	transportOut := thrift.NewTMemoryBufferLen(65507)
	protocolOut := thrift.NewTCompactProtocolFactory().GetProtocol(transportOut)
	return &Serializer{transportOut: transportOut, protocolOut: protocolOut}
}

var serializePool sync.Pool

func Serialize(tStruct thrift.TStruct) []byte {
	var serializer *Serializer
	v := serializePool.Get()
	if v == nil {
		serializer = NewSerializer()
	} else {
		serializer = v.(*Serializer)
	}

	header := HeaderLookup(tStruct)
	serializer.transportOut.Reset()
	writeHeader(serializer.protocolOut, header)
	tStruct.Write(serializer.protocolOut)
	buf := serializer.transportOut.Bytes()
	newBuf := make([]byte, len(buf))
	copy(newBuf, buf)
	serializePool.Put(serializer)
	return newBuf
}

func writeHeader(protocol thrift.TProtocol, header *Header) {
	protocol.WriteByte(int8(header.Signature))
	protocol.WriteByte(int8(header.Version))
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(header.Type))
	protocol.WriteByte(int8(buf[0]))
	protocol.WriteByte(int8(buf[1]))
}
