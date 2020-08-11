package agent

import (
	"github.com/cnogo/pinpoint-go/protocol/thrift/trace"
)

type Span struct {
	TSpan *trace.TSpan
}

func NewSpan(traceID *TraceID, header *TraceHttpHeader) *Span {
	span := &Span{
		TSpan: trace.NewTSpan(),
	}

	if header != nil && header.HttpType == VALID_HTTP_HEADER{
		pAppName := header.PAppName
		pAppType := int16(header.PAppType)
		span.TSpan.ParentApplicationName = &pAppName
		span.TSpan.ParentApplicationType = &pAppType
	}

	span.TSpan.TransactionId = FmtTransIDByte(traceID.AgentID, traceID.AgentStartTime, traceID.TransSeq)
	span.TSpan.AgentId = GAgent.agentID
	span.TSpan.ApplicationName = GAgent.applicationName
	span.TSpan.AgentStartTime = GAgent.startTime
	span.TSpan.ServiceType = GAgent.serverType
	span.TSpan.SpanId = traceID.SpanID
	span.TSpan.ParentSpanId = traceID.PSpanID
	span.TSpan.Flag = int16(traceID.Flags)

	return span
}

func (p *Span) MarkStartTime() {
	p.TSpan.StartTime = GetNowMSec()
}

func (p *Span) MarkEndTime() {
	p.TSpan.Elapsed = int32(GetNowMSec() - p.TSpan.StartTime)
}

func (p *Span) SetServiceType(serviceType int16) {
	p.TSpan.ServiceType = serviceType
}

func (p *Span) GetStartTime() int64 {
	return p.TSpan.StartTime
}

func (p *Span)  GetEndTime() int64 {
	return p.TSpan.StartTime + int64(p.TSpan.Elapsed)
}

func (p *Span) SetAPIID(apiID int32) {
	p.TSpan.ApiId = &apiID
}

func (p *Span) SetEndPoint(endPoint string) {
	p.TSpan.EndPoint = &endPoint
}

func (p *Span) SetRemoteAddr(remote string) {
	p.TSpan.RemoteAddr = &remote
}

func (p *Span) setAcceptorHost(host string) {
	p.TSpan.AcceptorHost = &host
}

func (p *Span) SetRpc(rpc string) {
	p.TSpan.RPC = &rpc
}

func (p *Span) AddSpanEvent(spanEvent *trace.TSpanEvent) {
	p.TSpan.SpanEventList = append(p.TSpan.SpanEventList, spanEvent)
}

func (p *Span) AddAnnotation(annotation *trace.TAnnotation) {
	p.TSpan.Annotations = append(p.TSpan.Annotations, annotation)
}

func (p *Span) SetExceptionInfo(errID int32, errMsg string) {
	value := trace.NewTIntStringValue()
	value.IntValue = errID
	value.StringValue = &errMsg
	p.TSpan.ExceptionInfo = value
}

type SpanEvent struct {
	TSpanEvent *trace.TSpanEvent
	Span *Span
}

func NewSpanEvent(span *Span, traceID *TraceID) *SpanEvent {
	if traceID == nil || span == nil {
		return nil
	}

	spanEvent := &SpanEvent{
		TSpanEvent: trace.NewTSpanEvent(),
		Span: span,
	}
	spanEvent.Span = span
	spanID := traceID.SpanID
	spanEvent.TSpanEvent.SpanId = &spanID

	return spanEvent
}

func (p *SpanEvent) MarkStartTime() {
	elapsed := int32(GetNowMSec() - p.Span.GetStartTime())
	p.TSpanEvent.StartElapsed = elapsed
}

func (p *SpanEvent) MarkEndTime() {
	elapsed := int32(GetNowMSec() - p.Span.GetStartTime())
	p.TSpanEvent.EndElapsed = elapsed
}

func (p *SpanEvent) SetApiID(apiID int32) {
	p.TSpanEvent.ApiId = &apiID
}

func (p *SpanEvent) SetServiceType(serviceType int16) {
	p.TSpanEvent.ServiceType = serviceType
}

func (p *SpanEvent) AddAnnotation(annotaion *trace.TAnnotation) {
	p.TSpanEvent.Annotations = append(p.TSpanEvent.Annotations, annotaion)
}

func (p *SpanEvent) SetEndPoint(endPoint string ) {
	p.TSpanEvent.EndPoint = &endPoint
}

func (p *SpanEvent) SetRPC(rpc string) {
	p.TSpanEvent.RPC = &rpc
}

func (p *SpanEvent) SetExceptionInfo(errID int32, errMsg string) {
	value := trace.NewTIntStringValue()
	value.IntValue = errID
	value.StringValue = &errMsg
	p.TSpanEvent.ExceptionInfo = value
}

func (p *SpanEvent) SetNextSpanID(nextSpanID int64) {
	p.TSpanEvent.NextSpanId = nextSpanID
}

func (p *SpanEvent) SetDestinationID(desID string) {
	p.TSpanEvent.DestinationId = &desID
}

func (p *SpanEvent) Finish() {
	p.MarkEndTime()
}

func (p *SpanEvent) setSequence(seq int16) {
	p.TSpanEvent.Sequence = seq
}

func (p *SpanEvent) setDepth(depth int32) {
	p.TSpanEvent.Depth = depth
}




