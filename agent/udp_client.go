package agent

import (
	"log"
	"net"
	"time"
)

type UDPClient struct {
	conn     net.Conn
	addr     string
	sendChan chan []byte
	isRestart bool

}

func NewUDPClient(addr string) *UDPClient {
	return &UDPClient{addr: addr, isRestart: true}
}

func (p *UDPClient) Start() {
	var err error

	p.sendChan = make(chan []byte, UDP_PACKET_COUNT)

	defer func() {
		if err := recover(); err != nil {
			log.Println("udp client error: ", err)
		}

		if p.conn != nil {
			p.conn.Close()
		}

		if p.isRestart {
			time.AfterFunc(5 * time.Second, func() {
				go p.Start()
			})
		}
	}()

	for {
		p.conn, err = net.Dial("udp", p.addr)
		if err != nil {
			log.Println("dial ", p.addr, " error：", err)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	for {
		select {
		case payload, ok := <-p.sendChan:
			if !ok {
				break
			}

			_, err = p.conn.Write(payload)
			if err != nil {
				log.Println("error: ", err)
				return
			}
		}
	}
}

func (p *UDPClient) SendPacket(pack []byte) {
	if p.sendChan == nil {
		return
	}
	p.sendChan <- pack
}

func (p *UDPClient) Close(){
	p.isRestart = false

	if p.conn != nil {
		p.conn.Close()
	}
}