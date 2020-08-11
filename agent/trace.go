package agent

import (
	"log"
	"math/rand"
	"strconv"
	"strings"
	"sync/atomic"
)

const TVERSION byte = 0
var transIDSeq int64

func NewTransIDSeq() int64 {
	return atomic.AddInt64(&transIDSeq, 1)
}

func FmtTransIDByte(agentID string, agentStartTime int64, transSeq int64) []byte {
	byteUtil := NewByteUtil("")
	byteUtil.PutByte(TVERSION)
	byteUtil.PutPrefixedString(agentID)
	byteUtil.PutVar64(agentStartTime)
	byteUtil.PutVar64(transSeq)
	return byteUtil.GetBytes()
}

func FmtTransIDString(agentID string, agentStartTime int64, transSeq int64) string {
	return agentID + "^" +
		strconv.FormatInt(agentStartTime, 10) + "^" +
		strconv.FormatInt(transSeq, 10)
}

func ParseTransID(bys []byte) string {
	byteUtil := NewByteUtil(string(bys))
	version := byteUtil.ReadByte()
	if version != TVERSION {
		return ""
	}
	agentID := byteUtil.ReadPrefixedString()
	agentStartTime := byteUtil.ReadVar64()
	transSeq := byteUtil.ReadVar64()
	return FmtTransIDString(agentID, agentStartTime, transSeq)
}

func createSpanID() int64 {
	spanID := rand.Int63()
	for spanID == -1 {
		spanID = rand.Int63()
	}
	return spanID
}

func NewSpanID() int64 {
	return createSpanID()
}

func NextSpanID(spanID, pSpanID int64) int64 {
	nextSpanID := createSpanID()
	for spanID == nextSpanID || nextSpanID == pSpanID {
		nextSpanID = createSpanID()
	}
	return nextSpanID
}

type TraceID struct {
	AgentID string
	TransID string
	PSpanID int64
	SpanID  int64
	TransSeq int64
	AgentStartTime int64
	Flags int64
}

func NewTraceID(transID string, pSpanID, spanID, flags int64) *TraceID {
	transInfo := strings.Split(transID, "^")
	if len(transInfo) != 3 {
		log.Println("New TraceID Parse Error")
		return nil
	}

	agentID := transInfo[0]
	agentStartTime, err := strconv.ParseInt(transInfo[1], 10, 64)
	if err != nil {
		log.Println("New TraceID Parse AgentStartTime error")
		return nil
	}

	transSeq, err := strconv.ParseInt(transInfo[2], 10, 64)
	if err != nil {
		log.Println("New TraceID Parse TransSeq error")
		return nil
	}

	return &TraceID{
		AgentID: agentID,
		AgentStartTime: agentStartTime,
		TransSeq: transSeq,
		PSpanID: pSpanID,
		SpanID: spanID,
		Flags: flags,
	}
}

func (p *TraceID) IsRoot() bool {
	return p.PSpanID == -1
}

func (p *TraceID) GetTransID() string {
	return p.TransID
}

type TraceContext struct {
	Span *Span
	TraceHeader *TraceHttpHeader
	SpanEventList []*SpanEvent
	TraceID *TraceID
	Depth int32
	Sequence int32
}

func NewTraceContext(traceHeader *TraceHttpHeader, traceID *TraceID, span *Span) *TraceContext {
	return &TraceContext{
		Span: span,
		TraceHeader: traceHeader,
		TraceID: traceID,
		Sequence: -1,
		Depth: 0,   //从1开始
	}
}

func (p *TraceContext) StartTraceSpanEvent() *SpanEvent {
	spanEvent := NewSpanEvent(p.Span, p.TraceID)
	spanEvent.setSequence(int16(atomic.AddInt32(&p.Sequence, 1)))
	spanEvent.setDepth(atomic.AddInt32(&p.Depth, 1))
	spanEvent.MarkStartTime()
	p.SpanEventList = append(p.SpanEventList, spanEvent)
	return spanEvent
}

func (p *TraceContext) GetNextSpanInfo() map[string]string {
	m := make(map[string]string)
	nextSpanID := NextSpanID(p.TraceID.SpanID, p.TraceID.PSpanID)

	m[HTTP_TRACE_ID] = FmtTransIDString(p.TraceID.AgentID, p.TraceID.AgentStartTime, p.TraceID.TransSeq)
	m[HTTP_PARENT_SPAN_ID] = strconv.FormatInt(p.TraceID.SpanID, 10)
	m[HTTP_SPAN_ID] = strconv.FormatInt(nextSpanID, 10)
	m[HTTP_FLAGS] = "0"

	sampled := SAMPLING_RATE_TRUE
	if !p.TraceHeader.Sampled {
		sampled = SAMPLING_RATE_FALSE
	}

	m[HTTP_SAMPLED] = sampled
	m[HTTP_PARENT_APPLICATION_NAME] = GAgent.applicationName
	m[HTTP_PARENT_APPLICATION_TYPE] = strconv.Itoa(int(GAgent.serverType))

	return m
}

func (p *TraceContext) finish() {
	p.Span.MarkEndTime()
	count := len(p.SpanEventList)
	for i := 0; i < count; i++ {
		p.Span.AddSpanEvent(p.SpanEventList[i].TSpanEvent)
	}
}
