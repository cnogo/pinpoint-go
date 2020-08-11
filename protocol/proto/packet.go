package proto

import (
	"io"
	"net"
)

//PackageType
var (
	APPLICATION_SEND           int16 = 1
	APPLICATION_TRACE_SEND     int16 = 2
	APPLICATION_TRACE_SEND_ACK int16 = 3

	APPLICATION_REQUEST  int16 = 5
	APPLICATION_RESPONSE int16 = 6

	APPLICATION_STREAM_CREATE         int16 = 10
	APPLICATION_STREAM_CREATE_SUCCESS int16 = 12
	APPLICATION_STREAM_CREATE_FAIL    int16 = 14

	APPLICATION_STREAM_CLOSE int16 = 15

	APPLICATION_STREAM_PING int16 = 17
	APPLICATION_STREAM_PONG int16 = 18

	APPLICATION_STREAM_RESPONSE int16 = 20

	CONTROL_CLIENT_CLOSE int16 = 100
	CONTROL_SERVER_CLOSE int16 = 110

	// control packet
	CONTROL_HANDSHAKE          int16 = 150
	CONTROL_HANDSHAKE_RESPONSE int16 = 151

	// keep stay because of performance in case of ping and pong. others removed.
	// CONTROL_PING will be deprecated. caused : Two payload types are used in one control packet.
	// since 1.7.0, use CONTROL_PING_SIMPLE, CONTROL_PING_PAYLOAD
	//@Deprecated
	CONTROL_PING int16 = 200
	CONTROL_PONG int16 = 201

	CONTROL_PING_SIMPLE  int16 = 210
	CONTROL_PING_PAYLOAD int16 = 211

	UNKNOWN int16 = 500

	PACKET_TYPE_SIZE int16 = 2
)

type Packet interface {
	Decode(conn net.Conn, reader io.Reader) error
	Encode() ([]byte, error)
	GetPacketType() int16
	GetPayload() []byte
	GetRequestID() int
}
