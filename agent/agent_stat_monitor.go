package agent

import (
	"github.com/cnogo/pinpoint-go/protocol/thrift"
	"github.com/cnogo/pinpoint-go/protocol/thrift/pinpoint"
	"time"
)

type AgentStatMonitor struct {
	client    *UDPClient
	agentStat []*pinpoint.TAgentStat
}

func NewAgentStatMonitor(addr string) *AgentStatMonitor {
	return &AgentStatMonitor{
		client:    NewUDPClient(addr),
		agentStat: make([]*pinpoint.TAgentStat, 0, AGENT_STAT_BATCH_NUM),
	}
}

func (p *AgentStatMonitor) Start() {
	quitC := make(chan struct{})
	defer func() {
		close(quitC)
	}()

	go p.collectAgentStat(quitC)
	p.client.Start()
}

func (p *AgentStatMonitor) collectAgentStat(quitC <-chan struct{}) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			agentStat := pinpoint.NewTAgentStat()
			startTime := GAgent.startTime

			agentStat.AgentId = &GAgent.agentID
			agentStat.StartTimestamp = &startTime

			p.agentStat = append(p.agentStat, agentStat)

			if len(p.agentStat) == AGENT_STAT_BATCH_NUM {
				agentStatBatch := pinpoint.NewTAgentStatBatch()
				agentStatBatch.AgentId = GAgent.agentID
				agentStatBatch.StartTimestamp = GAgent.startTime
				agentStatBatch.AgentStats = p.agentStat[0:AGENT_STAT_BATCH_NUM]

				payload := thrift.Serialize(agentStatBatch)
				p.client.SendPacket(payload)
				p.agentStat = p.agentStat[:0]
			}

		case <-quitC:
			return
		}
	}
}
