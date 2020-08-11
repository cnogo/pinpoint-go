package proto

import (
	"errors"
	"io"
	"net"
)

func DecodePacketFactory(packetType int16, conn net.Conn, reader io.Reader) (Packet, error) {

	var pack Packet
	var err error
	switch packetType {
	case APPLICATION_SEND:
		pack = NewApplicationSend()
	case APPLICATION_TRACE_SEND:
		pack = NewApplicationTraceSend()
	case APPLICATION_TRACE_SEND_ACK:
		pack = NewApplicationTraceSendAck()
	case APPLICATION_REQUEST:
		pack = NewApplicationRequest()
	case APPLICATION_RESPONSE:
		pack = NewApplicationResponse()
	case APPLICATION_STREAM_CREATE:
		pack = NewApplicationStreamCreate()
	case APPLICATION_STREAM_CREATE_SUCCESS:
		pack = NewApplicationStreamCreateSuccess()
	case APPLICATION_STREAM_CREATE_FAIL:
		pack = NewApplicationStreamCreateFail()
	case APPLICATION_STREAM_CLOSE:
		pack = NewApplicationStreamClose()
	case APPLICATION_STREAM_PING:
		pack = NewApplicationStreamPing()
	case APPLICATION_STREAM_PONG:
		pack = NewApplicationStreamPong()
	case APPLICATION_STREAM_RESPONSE:
		pack = NewApplicationStreamResponse()
	case CONTROL_CLIENT_CLOSE:
		pack = NewControlClientClose()
	case CONTROL_SERVER_CLOSE:
		pack = NewControlServerClose()
	case CONTROL_HANDSHAKE:
		pack = NewControlHandShake()
	case CONTROL_HANDSHAKE_RESPONSE:
		pack = NewControlHandShakeResponse()
	case CONTROL_PONG:
		pack = NewControlPong()
	case CONTROL_PING_SIMPLE:
		pack = NewControlPingSimple()
	case CONTROL_PING_PAYLOAD:
		pack = NewControlPingPayload()

	default:
		err = errors.New("unknown packet type")
	}

	if pack != nil {
		err = pack.Decode(conn, reader)
	}

	return pack, err
}
