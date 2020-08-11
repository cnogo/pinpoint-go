package agent

import "time"

const AGENT_STAT_BATCH_NUM = 6

const PING_INTERVAL = 5 * time.Minute
const AGENT_INFO_INTERVAL = 24 * time.Hour

const PINPOINT_AGENT_VERSION = "1.8.2"

const PP_CTX = "PinpointContext"

const (
	MYSQL_EXEC    = "mysql.Exec"
	MYSQL_CONNECT = "mysql.Connect"
)

const (
	STR_ERROR     = "Error"
	STR_EXCEPTION = "Exception"
	STR_WARN      = "Warn"
)

//tcp udp相关定义
const (
	TCP_READ_BUFFER_LEN  int = 1024 * 64
	TCP_WRITE_BUFFER_LEN int = 1024 * 64
	TCP_PACKET_COUNT         = 1024 * 10

	UDP_PACKET_COUNT = 1024
)

//HTTP请求本服务的baggage信息
const (
	HTTP_TRACE_ID                = "Pinpoint-TraceID"
	HTTP_SPAN_ID                 = "Pinpoint-SpanID"
	HTTP_PARENT_SPAN_ID          = "Pinpoint-pSpanID"
	HTTP_SAMPLED                 = "Pinpoint-Sampled"
	HTTP_FLAGS                   = "Pinpoint-Flags"
	HTTP_PARENT_APPLICATION_NAME = "Pinpoint-pAppName"
	HTTP_PARENT_APPLICATION_TYPE = "Pinpoint-pAppType"
	HTTP_HOST                    = "Pinpoint-Host"
)

const (
	SAMPLING_RATE_FALSE = "s0"
	SAMPLING_RATE_TRUE  = "s1"

	NULL_HTTP_HEADER    = 0
	INVALID_HTTP_HEADER = 1
	VALID_HTTP_HEADER   = 2
)

const (
	HANDSHAKE_INIT      byte = 0
	HANDSHAKE_STARTED   byte = 1
	HANDSHAKE_FINISHED  byte = 2
	HANDSHAKE_MAX_COUNT      = 10
	HANDSHAKE_INTERVAL       = 60 * time.Second
)

const (
	METHOD_TYPE_DEFAULT     int32 = 0
	METHOD_TYPE_EXCEPTION   int32 = 1
	METHOD_TYPE_ANNOTATION  int32 = 2
	METHOD_TYPE_PARAMETER   int32 = 3
	METHOD_TYPE_WEB_REQUEST int32 = 100
	METHOD_TYPE_INVOCATION  int32 = 200
	METHOD_TYPE_CORRUPTED   int32 = 900
)

//Annotation Key

const (
	AK_API          int32 = 12
	AK_API_METADATA int32 = 13
	AK_RETURN_DATA  int32 = 14
	AK_API_TAG      int32 = 10015

	AK_ERR_API_METADATA_ERROR                  int32 = 10000010
	AK_ERR_API_METADATA_AGENT_INFO_NOT_FOUND   int32 = 10000011
	AK_ERR_API_METADATA_IDENTIFIER_CHECK_ERROR int32 = 10000012
	AK_ERR_API_METADATA_NOT_FOUND              int32 = 10000013
	AK_ERR_API_METADATA_DID_COLLSION           int32 = 10000014

	AK_SQL_ID        int32 = 20
	AK_SQL           int32 = 21
	AK_SQL_METADATA  int32 = 22
	AK_SQL_PARAM     int32 = 23
	AK_SQL_BINDVALUE int32 = 24

	AK_STRING_ID int32 = 30

	AK_HTTP_URL              int32 = 40
	AK_HTTP_PARAM            int32 = 41
	AK_HTTP_PARAM_ENTITY     int32 = 42
	AK_HTTP_COOKIE           int32 = 45
	AK_HTTP_STATUS_CODE      int32 = 46
	AK_HTTP_INTERNAL_DISPLAY int32 = 48
	AK_HTTP_IO               int32 = 49

	AK_MSG_QUEUE_URI int32 = 100

	AK_ARGS0 int32 = -1
	AK_ARGS1 int32 = -2
	AK_ARGS2 int32 = -3
	AK_ARGS3 int32 = -4
	AK_ARGS4 int32 = -5
	AK_ARGS5 int32 = -6
	AK_ARGS6 int32 = -7
	AK_ARGS7 int32 = -8
	AK_ARGS8 int32 = -9
	AK_ARGS9 int32 = -10
	AK_ARGSN int32 = -11

	AK_CACHE_ARGS0 int32 = -30
	AK_CACHE_ARGS1 int32 = -31
	AK_CACHE_ARGS2 int32 = -32
	AK_CACHE_ARGS3 int32 = -33
	AK_CACHE_ARGS4 int32 = -34
	AK_CACHE_ARGS5 int32 = -35
	AK_CACHE_ARGS6 int32 = -36
	AK_CACHE_ARGS7 int32 = -37
	AK_CACHE_ARGS8 int32 = -38
	AK_CACHE_ARGS9 int32 = -39
	AK_CACHE_ARGSN int32 = -40
)
