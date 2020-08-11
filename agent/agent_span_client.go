package agent

import (
	"github.com/cnogo/pinpoint-go/protocol/thrift"
	"github.com/cnogo/pinpoint-go/protocol/thrift/trace"
)

type AgentSpanClient struct {
	client *UDPClient
}

func NewAgentSpanClient(addr string) *AgentSpanClient {
	return &AgentSpanClient{
		client: NewUDPClient(addr),
	}
}

func (p *AgentSpanClient) Start() {
	p.client.Start()
}

func (p *AgentSpanClient) SendTrace(tspan *trace.TSpan) {
	payload := thrift.Serialize(tspan)
	p.client.SendPacket(payload)
}

