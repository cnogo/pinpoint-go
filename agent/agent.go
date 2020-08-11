package agent

import (
	"math/rand"
	"net/http"
	"github.com/cnogo/pinpoint-go/config"
	"github.com/cnogo/pinpoint-go/protocol/proto"
	"github.com/cnogo/pinpoint-go/protocol/thrift"
	"github.com/cnogo/pinpoint-go/protocol/thrift/trace"
	"time"
)

func Init() {
	rand.Seed(time.Now().UnixNano())
}

func InitCommonAPI() {
	GAgent.AddAPIMeta("http." + http.MethodHead, -1, -1)
	GAgent.AddAPIMeta("http." + http.MethodGet, -1, -1)
	GAgent.AddAPIMeta("http." + http.MethodDelete, -1, -1)
	GAgent.AddAPIMeta("http." + http.MethodConnect, -1, -1)
	GAgent.AddAPIMeta("http." + http.MethodOptions, -1, -1)
	GAgent.AddAPIMeta("http." + http.MethodPatch, -1, -1)
	GAgent.AddAPIMeta("http." + http.MethodPost, -1, -1)
	GAgent.AddAPIMeta("http." + http.MethodTrace, -1, -1)
	GAgent.AddAPIMeta("http." + http.MethodPut, -1, -1)


	GAgent.AddAPIMeta(MYSQL_EXEC, -1, -1)
	GAgent.AddAPIMeta(MYSQL_CONNECT, -1, -1)

	GAgent.AddStrMeta(STR_ERROR)
	GAgent.AddStrMeta(STR_EXCEPTION)
	GAgent.AddStrMeta(STR_WARN)

}

type Agent struct {
	agentID         string
	applicationName string
	client          *TCPClient
	statMonitor     *AgentStatMonitor
	spanClient      *AgentSpanClient
	apiMetaManager  *IDManager
	sqlMetaManager  *IDManager
	strMetaManager  *IDManager
	startTime       int64
	serverType      int16
	hostName        string
	endPoint        string
}

var GAgent *Agent

func NewAgent() *Agent {

	agentId := GetHostIP() + GetRandomChar(4)
	config.Conf.AgentID = agentId

	GAgent = &Agent{
		client:      NewTCPClient(),
		statMonitor: NewAgentStatMonitor(config.Conf.Pinpoint.StatAddr),
		spanClient:  NewAgentSpanClient(config.Conf.Pinpoint.SpanAddr),
		startTime:   GetNowMSec(),
		serverType:  STAND_ALONE,
		apiMetaManager: NewIDManager(),
		sqlMetaManager: NewIDManager(),
		strMetaManager: NewIDManager(),
		hostName: GetHostName(),
		endPoint: GetHostIP(),
	}

	return GAgent
}

func (p *Agent) Start() {
	p.agentID = config.Conf.AgentID
	p.applicationName = config.Conf.ApplicationName
	go p.client.Init()
	//go p.statMonitor.Start()
	go p.spanClient.Start()

	InitCommonAPI()
}

func (p *Agent) AddAPIMeta(api string, line, apiType int32) int32 {
	apiID, isNew := p.apiMetaManager.GetOrCreateID(api)

	//非新的说明之前已经发过了
	if !isNew {
		return apiID
	}

	apiMetaObject := &APIMeta{
		Api: api,
		ApiID: apiID,
		Line: line,
		Type: apiType,
	}

	p.apiMetaManager.SetMetaObject(apiID, apiMetaObject)

	p.sendApiMeta(api, apiID, line, apiType)

	return apiID
}

func (p *Agent) sendApiMeta(api string, apiID, line, apiType int32) {
	apiMetaData := trace.NewTApiMetaData()
	apiMetaData.AgentId = GAgent.agentID
	apiMetaData.AgentStartTime = GAgent.startTime
	apiMetaData.ApiId = apiID
	apiMetaData.ApiInfo = api

	if line != -1 {
		apiMetaData.Line = &line
	}
	if apiType != -1 {
		apiMetaData.Type = &apiType
	}

	payload := thrift.Serialize(apiMetaData)
	request := proto.NewApplicationRequest()
	request.Payload = payload
	p.client.SendPacket(request)
}


func (p *Agent) GetApiID(metaInfo string) int32 {
	return p.apiMetaManager.GetID(metaInfo)
}

func (p *Agent) AddSQLMeta(sql string) int32 {
	sqlID, isNew := p.sqlMetaManager.GetOrCreateID(sql)

	if !isNew {
		return sqlID
	}

	sqlMetaObject := &SQLMeta{
		SqlID: sqlID,
		Sql: sql,
	}

	p.sqlMetaManager.SetMetaObject(sqlID, sqlMetaObject)

	p.sendSqlMeta(sqlID, sql)
	return sqlID
}

func (p *Agent) sendSqlMeta(sqlID int32, sql string) {
	sqlMetaData := trace.NewTSqlMetaData()
	sqlMetaData.AgentId = GAgent.agentID
	sqlMetaData.AgentStartTime = GAgent.startTime
	sqlMetaData.Sql = sql
	sqlMetaData.SqlId = sqlID

	payload := thrift.Serialize(sqlMetaData)
	request := proto.NewApplicationRequest()
	request.Payload = payload
	p.client.SendPacket(request)
}

func (p *Agent) AddStrMeta(metaInfo string) int32 {
	strID, isNew := p.strMetaManager.GetOrCreateID(metaInfo)
	if !isNew {
		return strID
	}

	strMetaObject := &StrMeta{
		StrID: strID,
		Str: metaInfo,
	}

	p.strMetaManager.SetMetaObject(strID, strMetaObject)
	p.sendStrMeta(strID, metaInfo)
	return strID
}

func (p *Agent) sendStrMeta(strID int32, str string) {
	strMetaData := trace.NewTStringMetaData()
	strMetaData.AgentId = GAgent.agentID
	strMetaData.AgentStartTime = GAgent.startTime
	strMetaData.StringId = strID
	strMetaData.StringValue = str

	payload := thrift.Serialize(strMetaData)
	request := proto.NewApplicationRequest()
	request.Payload = payload
	p.client.SendPacket(request)
}

func (p *Agent) GetStrID(metaInfo string) int32 {
	return p.strMetaManager.GetID(metaInfo)
}

func (p *Agent) sendAllMetaData() {
	sqlMetaMap := p.sqlMetaManager.ID2MetaObject

	sqlMetaMap.Range(func(key, metaObject interface{}) bool {
		sqlMeta, ok := metaObject.(*SQLMeta)
		if !ok {
			return true
		}
		p.sendSqlMeta(sqlMeta.SqlID, sqlMeta.Sql)
		return true
	})

	apiMetaMap := p.apiMetaManager.ID2MetaObject

	apiMetaMap.Range(func(key, metaObject interface{}) bool {
		apiMeta, ok := metaObject.(*APIMeta)
		if !ok {
			return true
		}
		p.sendApiMeta(apiMeta.Api, apiMeta.ApiID, apiMeta.Line, apiMeta.Type)
		return true
	})

	strMetaMap := p.strMetaManager.ID2MetaObject
	strMetaMap.Range(func(key, metaObject interface{}) bool {
		strMeta, ok := metaObject.(*StrMeta)
		if !ok {
			return true
		}
		p.sendStrMeta(strMeta.StrID, strMeta.Str)
		return true
	})
}

func (p *Agent) GetSQLID(metaInfo string) int32 {
	return p.sqlMetaManager.GetID(metaInfo)
}


func (p *Agent) SendSpan(tspan *trace.TSpan) {
	p.spanClient.SendTrace(tspan)
}