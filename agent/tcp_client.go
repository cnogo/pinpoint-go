package agent

import (
	"bufio"
	"encoding/binary"
	"io"
	"log"
	"net"
	"os"
	"github.com/cnogo/pinpoint-go/config"
	"github.com/cnogo/pinpoint-go/protocol/proto"
	"github.com/cnogo/pinpoint-go/protocol/thrift"
	"github.com/cnogo/pinpoint-go/protocol/thrift/pinpoint"
	"sync/atomic"
	"time"
)

type TCPClient struct {
	conn       net.Conn
	sockState  *SocketState
	sendChan   chan proto.Packet
	pingID     int32
	socketID   int32
	requestID  int
	handShake  *HandShake
	isRestart bool
}

func NewTCPClient() *TCPClient {
	return &TCPClient{isRestart: true}
}

func (p *TCPClient) Init() {
	var err error
	p.handShake = NewHandShake(p)
	p.sockState = NewSocketState()
	p.sendChan = make(chan proto.Packet, TCP_PACKET_COUNT)
	atomic.AddInt32(&p.socketID, 1)

	quitC := make(chan struct{})

	defer func() {
		if err := recover(); err != nil {
			//先用系统的日志
			log.Println("tcpclient init recover: ", err)
		}

		if p.isRestart {
			time.AfterFunc(5 * time.Second, func() {
				go p.Init()
			})
		}
	}()

	defer func() {
		log.Println("its closed???")
		close(quitC)
		if p.conn != nil {
			_ = p.conn.Close()
		}
	}()

	for {
		p.conn, err = net.Dial("tcp", config.Conf.Pinpoint.InfoAddr)
		if err != nil {
			log.Println("dial ", config.Conf.Pinpoint.InfoAddr, " error：", err)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	tcpConn := p.conn.(*net.TCPConn)
	err = tcpConn.SetReadBuffer(TCP_READ_BUFFER_LEN)
	if err != nil {
		log.Println("set read buffer error: ", err)
		return
	}

	err = tcpConn.SetWriteBuffer(TCP_WRITE_BUFFER_LEN)
	if err != nil {
		log.Println("set write buffer error: ", err)
		return
	}

	p.sockState.toConnected()
	log.Println("connect success ", config.Conf.Pinpoint.InfoAddr)
	p.sockState.toRunWithoutHandShake()

	go p.writePacket(quitC)
	go p.doPing(quitC)
	go p.handShake.Start(quitC)
	go p.sendAgentInfo(quitC)

	//待优化
	GAgent.sendAllMetaData()

	p.handleRead()
}

func (p *TCPClient) GetSocketID() int32 {
	return atomic.LoadInt32(&p.socketID)
}

func (p *TCPClient) SendPacket(pack proto.Packet) {
	if p.sendChan == nil {
		return
	}
	p.sendChan <- pack
}

func (p *TCPClient) sendAgentInfo(quitC <-chan struct{}) {
	agentInfo := pinpoint.NewTAgentInfo()
	hostName, err := os.Hostname()
	if err != nil {
		log.Println("get hostname error ", err)
		hostName = "Unknown"
	}

	agentInfo.Hostname = hostName
	agentInfo.IP = "0.0.0.0"
	agentInfo.Ports = ""
	agentInfo.AgentId = GAgent.agentID
	agentInfo.ApplicationName = GAgent.applicationName
	agentInfo.ServiceType = int16(GAgent.serverType)
	agentInfo.Pid = int32(os.Getpid())
	agentInfo.AgentVersion = PINPOINT_AGENT_VERSION
	agentInfo.VmVersion = ""
	agentInfo.StartTimestamp = GAgent.startTime

	payload := thrift.Serialize(agentInfo)

	request := proto.NewApplicationRequest()
	request.Payload = payload

	ticker := time.NewTicker(AGENT_INFO_INTERVAL)
	defer ticker.Stop()

	firstTicker := time.After(1 * time.Second)

	for {
		select {
		case <-firstTicker:
			p.SendPacket(request)
		case <-ticker.C:
			p.SendPacket(request)
		case <-quitC:
			return
		}
	}

}

func (p *TCPClient) writePacket(quitC <-chan struct{}) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("write packet error: ", err)
		}
	}()

	for {
		select {
		case pack, ok := <-p.sendChan:
			if !ok {
				break
			}

			if pack.GetPacketType() == proto.APPLICATION_REQUEST {
				if v, ok := pack.(*proto.ApplicationRequest); ok {
					v.RequestID = p.requestID
					p.requestID++
				}
			}

			payload, err := pack.Encode()
			if err != nil {
				break
			}
			_, err = p.conn.Write(payload)
			if err != nil {
				log.Println("error: ", err)
			}

		case <-quitC:
			return
		}
	}
}

func (p *TCPClient) doPing(quitC <-chan struct{}) {
	firstTick := time.After(300 * time.Millisecond)
	tick := time.NewTicker(PING_INTERVAL)
	defer tick.Stop()

	ping := func() {
		pingPacket := proto.NewControlPingPayload()
		pingPacket.StateCode = p.sockState.getCurrentState()
		pingPacket.StateVersion = 0
		pingPacket.PingID = int(atomic.AddInt32(&p.pingID, 1))
		p.SendPacket(pingPacket)
	}

	for {
		select {
		case <-firstTick:
			ping()
		case <-tick.C:
			ping()

		case <-quitC:
			return
		}
	}
}

func (p *TCPClient) handleRead() {
	var err error
	reader := bufio.NewReaderSize(p.conn, TCP_READ_BUFFER_LEN)
	typeBuf := make([]byte, 2)
	//
	for {
		if _, err = io.ReadFull(reader, typeBuf); err != nil {
			log.Println("io.ReadFull", err)
			return
		}
		packetType := int16(binary.BigEndian.Uint16(typeBuf[0:2]))

		pack, err := proto.DecodePacketFactory(packetType, p.conn, reader)

		if err != nil {
			return
		}

		switch packetType {
		case proto.CONTROL_HANDSHAKE_RESPONSE:
			p.handShake.handlerReponse(pack.(*proto.ControlHandShakeResponse))

		case proto.CONTROL_PING_PAYLOAD:
			fallthrough
		case proto.CONTROL_PING_SIMPLE:
			fallthrough
		case proto.CONTROL_PING:
			pack := proto.NewControlPong()
			p.SendPacket(pack)
		}
	}
}

func (p *TCPClient) Close() {
	p.isRestart = false

	if p.conn != nil {
		p.conn.Close()
	}

	if p.sockState != nil {
		p.sockState.toClosedByClient()
	}
}
