package thrift

import (
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/cnogo/pinpoint-go/protocol/thrift/pinpoint"
	"github.com/cnogo/pinpoint-go/protocol/thrift/trace"
)

const (
	UKNOWN           int16 = -1
	NETWORK_CHECK    int16 = 10
	SPAN             int16 = 40
	AGENT_INFO       int16 = 50
	AGENT_STAT       int16 = 55
	AGENT_STAT_BATCH int16 = 56
	SPANCHUNK        int16 = 70
	SPANEVENT        int16 = 80

	SQL_META_DATA    int16 = 300
	API_META_DATA    int16 = 310
	RESULT           int16 = 320
	STRING_META_DATA int16 = 330

	//#################

	HEADER_SIGNATURE byte = 0xef
	HEADER_VERSION   byte = 0x10
)

type Header struct {
	Signature byte
	Version   byte
	Type      int16
}

func NewHeader(Type int16) *Header {
	return &Header{
		Signature: HEADER_SIGNATURE,
		Version:   HEADER_VERSION,
		Type:      Type,
	}
}

func HeaderLookup(tStrcut thrift.TStruct) *Header {
	switch tStrcut.(type) {
	case *trace.TSpan:
		return NewHeader(SPAN)
	case *trace.TSpanChunk:
		return NewHeader(SPANCHUNK)
	case *trace.TSpanEvent:
		return NewHeader(SPANEVENT)
	case *pinpoint.TAgentInfo:
		return NewHeader(AGENT_INFO)
	case *pinpoint.TAgentStat:
		return NewHeader(AGENT_STAT)
	case *pinpoint.TAgentStatBatch:
		return NewHeader(AGENT_STAT_BATCH)
	case *trace.TSqlMetaData:
		return NewHeader(SQL_META_DATA)
	case *trace.TApiMetaData:
		return NewHeader(API_META_DATA)
	case *trace.TStringMetaData:
		return NewHeader(STRING_META_DATA)
	case *trace.TResult_:
		return NewHeader(RESULT)
	}

	return nil
}

func TBaseLookup(Type int16) thrift.TStruct {
	switch Type {
	case SPAN:
		return trace.NewTSpan()
	case AGENT_INFO:
		return pinpoint.NewTAgentInfo()
	case AGENT_STAT:
		return pinpoint.NewTAgentStat()
	case AGENT_STAT_BATCH:
		return pinpoint.NewTAgentStatBatch()
	case SPANCHUNK:
		return trace.NewTSpanChunk()
	case SPANEVENT:
		return trace.NewTSpanEvent()
	case SQL_META_DATA:
		return trace.NewTSqlMetaData()
	case API_META_DATA:
		return trace.NewTApiMetaData()
	case STRING_META_DATA:
		return trace.NewTStringMetaData()
	case RESULT:
		return trace.NewTResult_()
	}

	return nil
}
