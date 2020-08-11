package thrift

import (
	"fmt"
	"log"
	"os"
	"github.com/cnogo/pinpoint-go/protocol/thrift/pinpoint"
	"github.com/cnogo/pinpoint-go/protocol/thrift/trace"
	"testing"
	"time"
)

func Test_AgentInfo(t *testing.T) {
	agentInfo := pinpoint.NewTAgentInfo()
	hostName, err := os.Hostname()
	if err != nil {
		log.Println("get hostname error ", err)
		hostName = "Unknown"
	}

	agentInfo.Hostname = hostName
	agentInfo.IP = "0.0.0.0"
	agentInfo.Ports = "8080"
	agentInfo.AgentId = "test_golang"
	agentInfo.ApplicationName = "golang"
	agentInfo.ServiceType = 1000
	agentInfo.Pid = int32(os.Getpid())
	agentInfo.AgentVersion = "1.8.2"
	agentInfo.VmVersion = ""
	agentInfo.StartTimestamp = time.Now().UnixNano() / 1000

	payload := Serialize(agentInfo)

	fmt.Println(agentInfo)
	tAgentInfo := Deserialize(payload)
	fmt.Println(tAgentInfo)
}

func Test_ApiMetaInfo(t *testing.T) {
	apiMetaData := trace.NewTApiMetaData()
	apiMetaData.AgentId = "dfljlkdf"
	apiMetaData.AgentStartTime = time.Now().UnixNano() / 1000
	apiMetaData.ApiId = 20
	apiMetaData.ApiInfo = "/hello"

	payload := Serialize(apiMetaData)
	tAgentInfo := Deserialize(payload)
	fmt.Println(tAgentInfo)
}
