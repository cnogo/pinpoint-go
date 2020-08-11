package agent

import (
	"log"
	"os"
	"github.com/cnogo/pinpoint-go/protocol/control"
	"github.com/cnogo/pinpoint-go/protocol/proto"
	"time"
)

type HandshakeResponseCode struct {
	Code        int32
	SubCode     int32
	CodeMessage string
}

var (
	HS_CODE     = "code"
	HS_SUB_CODE = "subCode"
	HS_CLUSTER  = "cluster"
)

var (
	HANDSHAKE_SUCCESS               = &HandshakeResponseCode{Code: 0, SubCode: 0, CodeMessage: "Success."}
	HANDSHAKE_SIMPLEX_COMMUNICATION = &HandshakeResponseCode{Code: 0, SubCode: 1, CodeMessage: "Simplex Connection successfully established."}
	HANDSHAKE_DUPLEX_COMMUNICATION  = &HandshakeResponseCode{Code: 0, SubCode: 2, CodeMessage: "Duplex Connection successfully established."}

	HANDSHAKE_ALREADY_KNOWN                 = &HandshakeResponseCode{1, 0, "Already Known."}
	HANDSHAKE_ALREADY_SIMPLEX_COMMUNICATION = &HandshakeResponseCode{1, 1, "Already Simplex Connection established."}
	HANDSHAKE_ALREADY_DUPLEX_COMMUNICATION  = &HandshakeResponseCode{1, 2, "Already Duplex Connection established."}

	HANDSHAKE_PROPERTY_ERROR = &HandshakeResponseCode{2, 0, "Property error."}

	HANDSHAKE_PROTOCOL_ERROR = &HandshakeResponseCode{3, 0, "Illegal protocol error."}
	HANDSHAKE_UNKNOWN_ERROR  = &HandshakeResponseCode{4, 0, "Unknown Error."}
	HANDSHAKE_UNKNOWN_CODE   = &HandshakeResponseCode{-1, -1, "Unknown Code."}
)

func getHandShakeCode(code int32, subcode int32) *HandshakeResponseCode {
	if code == 0 && subcode == 0 {
		return HANDSHAKE_SUCCESS
	} else if code == 0 && subcode == 1 {
		return HANDSHAKE_SIMPLEX_COMMUNICATION
	} else if code == 0 && subcode == 2 {
		return HANDSHAKE_DUPLEX_COMMUNICATION
	} else if code == 1 && subcode == 0 {
		return HANDSHAKE_ALREADY_KNOWN
	} else if code == 1 && subcode == 1 {
		return HANDSHAKE_ALREADY_SIMPLEX_COMMUNICATION
	} else if code == 1 && subcode == 2 {
		return HANDSHAKE_ALREADY_DUPLEX_COMMUNICATION
	} else if code == 2 && subcode == 0 {
		return HANDSHAKE_PROPERTY_ERROR
	} else if code == 3 && subcode == 0 {
		return HANDSHAKE_PROTOCOL_ERROR
	} else if code == 4 && subcode == 0 {
		return HANDSHAKE_UNKNOWN_CODE
	}

	return HANDSHAKE_UNKNOWN_CODE
}

type HandShake struct {
	client           *TCPClient
	status           byte
	hanshakeCount    int
	handshakeReponse *HandshakeResponseCode
	handshakePacket  proto.Packet
	sendTicker       *time.Ticker
}

func NewHandShake(client *TCPClient) *HandShake {
	return &HandShake{
		client:        client,
		status:        HANDSHAKE_INIT,
		hanshakeCount: 0,
	}
}

func (p *HandShake) createHandShakePacket() proto.Packet {
	hostName, err := os.Hostname()
	if err != nil {
		log.Println("get hostname error ", err)
		hostName = "Unknown"
	}

	var handshakeData = make(map[string]interface{})
	handshakeData["socketId"] = p.client.GetSocketID()
	handshakeData["hostName"] = hostName
	handshakeData["supportServer"] = true
	handshakeData["ip"] = "0.0.0.0"
	handshakeData["agentId"] = GAgent.agentID
	handshakeData["applicationName"] = GAgent.applicationName
	handshakeData["serviceType"] = GAgent.serverType
	handshakeData["pid"] = int32(os.Getegid())
	handshakeData["version"] = PINPOINT_AGENT_VERSION
	handshakeData["startTimestamp"] = GAgent.startTime

	encoder := control.NewEncoder()
	payload := encoder.Encode(handshakeData)

	pack := proto.NewControlHandShake()
	pack.RequestID = 0
	pack.Payload = payload

	return pack
}

func (p *HandShake) Start(quitC <-chan struct{}) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("handshake panic: ", err)
		}
	}()

	p.status = HANDSHAKE_STARTED
	p.handshakePacket = p.createHandShakePacket()

	p.sendTicker = time.NewTicker(HANDSHAKE_INTERVAL)
	defer p.sendTicker.Stop()

	p.client.SendPacket(p.handshakePacket)
	p.hanshakeCount++

	for {
		select {
		case <-p.sendTicker.C:
			if p.hanshakeCount == HANDSHAKE_MAX_COUNT {
				p.status = HANDSHAKE_FINISHED
				return
			}

			if p.status == HANDSHAKE_STARTED {
				p.client.SendPacket(p.handshakePacket)
			}
			p.hanshakeCount++

		case <-quitC:
			return
		}
	}
}

func (p *HandShake) handlerReponse(pack *proto.ControlHandShakeResponse) {
	decoder := control.NewDecoder(pack.GetPayload())
	responseData := decoder.Decode()

	p.status = HANDSHAKE_FINISHED

	var code int32 = -1
	var subCode int32 = -1

	msgData, ok := responseData.(map[string]interface{})
	if !ok {
		p.handshakeReponse = HANDSHAKE_UNKNOWN_CODE
		return
	}

	vcode, ok := msgData[HS_CODE]
	if !ok {
		p.handshakeReponse = HANDSHAKE_UNKNOWN_CODE
		return
	}

	code, ok = vcode.(int32)
	if !ok {
		p.handshakeReponse = HANDSHAKE_UNKNOWN_CODE
		return
	}

	vcode, ok = msgData[HS_SUB_CODE]
	if !ok {
		p.handshakeReponse = HANDSHAKE_UNKNOWN_CODE
		return
	}

	subCode, ok = vcode.(int32)
	if !ok {
		p.handshakeReponse = HANDSHAKE_UNKNOWN_CODE
		return
	}

	p.handshakeReponse = getHandShakeCode(code, subCode)

	if p.handshakeReponse == HANDSHAKE_SUCCESS || p.handshakeReponse == HANDSHAKE_ALREADY_KNOWN {
		p.client.sockState.toRunSimplex()
	} else if p.handshakeReponse == HANDSHAKE_DUPLEX_COMMUNICATION ||
		p.handshakeReponse == HANDSHAKE_ALREADY_DUPLEX_COMMUNICATION {
		p.client.sockState.toRunDuplex()
	} else if p.handshakeReponse == HANDSHAKE_SIMPLEX_COMMUNICATION ||
		p.handshakeReponse == HANDSHAKE_ALREADY_SIMPLEX_COMMUNICATION {
		p.client.sockState.toRunSimplex()
	}

	p.stopSendTicker()
}

func (p *HandShake) stopSendTicker() {
	if p.sendTicker != nil {
		p.sendTicker.Stop()
	}
}
