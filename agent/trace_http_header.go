package agent

import (
	"net/http"
	"strconv"
)

type TraceHttpHeader struct {
	HttpType int
	TransactionID string
	SpanID int64
	PSpanID int64
	Sampled bool
	Flag int64
	PAppName string
	PAppType int64
	PHostName string
}

func NewTraceHttpHeader(header *http.Header) *TraceHttpHeader{
	var traceHeader = new(TraceHttpHeader)
	traceHeader.Sampled = true
	traceHeader.HttpType = VALID_HTTP_HEADER
	if header == nil {
		traceHeader.HttpType = NULL_HTTP_HEADER
		return traceHeader
	}

	sampled := header.Get(HTTP_SAMPLED)
	if sampled == SAMPLING_RATE_FALSE {
		traceHeader.Sampled = false
		return traceHeader
	}

	traceHeader.TransactionID = header.Get(HTTP_TRACE_ID)
	if traceHeader.TransactionID == "" {
		traceHeader.HttpType = NULL_HTTP_HEADER
		return traceHeader
	}

	spanID, err := strconv.ParseInt(header.Get(HTTP_SPAN_ID), 10, 64)
	traceHeader.SpanID = spanID
	if err != nil {
		traceHeader.HttpType = INVALID_HTTP_HEADER
		return traceHeader
	}

	traceHeader.PSpanID, err = strconv.ParseInt(header.Get(HTTP_PARENT_SPAN_ID), 10, 64)
	if err != nil {
		traceHeader.HttpType = INVALID_HTTP_HEADER
		return traceHeader
	}

	traceHeader.Flag, err = strconv.ParseInt(header.Get(HTTP_FLAGS), 10, 64)
	traceHeader.PAppName = header.Get(HTTP_PARENT_APPLICATION_NAME)
	traceHeader.PAppType, err = strconv.ParseInt(header.Get(HTTP_PARENT_APPLICATION_TYPE), 10, 64)
	traceHeader.PHostName = header.Get(HTTP_HOST)
	return traceHeader
}

//param httpHeader is nil means its is root trace
func StartTrace(httpHeader *http.Header) *TraceContext {
	//如果agent没有准备好
	if GAgent == nil {
		return nil
	}

	traceHeader := NewTraceHttpHeader(httpHeader)

	if traceHeader.HttpType == INVALID_HTTP_HEADER || !traceHeader.Sampled {
		return nil
	}

	var traceID *TraceID

	if traceHeader.HttpType == VALID_HTTP_HEADER {
		//transID := ParseTransID([]byte(traceHeader.TransactionID))
		//traceID = NewTraceID(transID, traceHeader.PSpanID, traceHeader.SpanID, traceHeader.Flag)
		traceID = NewTraceID(traceHeader.TransactionID, traceHeader.PSpanID, traceHeader.SpanID, traceHeader.Flag)
		if traceID == nil {
			return nil
		}
	} else {
		transSeq := NewTransIDSeq()
		transID := FmtTransIDString(GAgent.agentID, GAgent.startTime, transSeq)
		traceID = &TraceID{
			AgentID: GAgent.agentID,
			AgentStartTime: GAgent.startTime,
			PSpanID: -1,
			SpanID: NewSpanID(),
			TransSeq: transSeq,
			Flags: 0,
			TransID: transID,
		}
	}

	span := NewSpan(traceID, traceHeader)

	if !traceID.IsRoot() {
		span.setAcceptorHost(traceHeader.PHostName)
	}
	span.MarkStartTime()
	span.SetEndPoint(GAgent.endPoint)

	return NewTraceContext(traceHeader, traceID, span)
}

func FinishTrace(traceContext *TraceContext) {
	if traceContext == nil {
		return
	}
	traceContext.finish()
	GAgent.SendSpan(traceContext.Span.TSpan)
}


