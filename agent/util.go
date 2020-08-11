package agent

import (
	"math/rand"
	"net"
	"os"
	"time"
)

func GetHostName() string {
	hostName, err := os.Hostname()
	if err != nil {
		hostName = "Unknow-host"
	}

	return hostName
}

func GetHostIP() string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return "0.0.0.0"
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String()
					}
				}
			}
		}
	}

	return "0.0.0.0"
}

func GetNowMSec() int64 {
	return time.Now().UnixNano()/1000/1000
}

func GetRandomChar(num int) string {
	//0-9a-zA-Z
	count := 10 + 26 + 26

	str := ""
	for i := 0; i < num; i++ {
		index := rand.Intn(count)
		if index < 10 {
			str += string('0' + index)
		} else if index < 36 {
			str += string('a' + index - 10)
		} else {
			str += string('A' + index - 36)
		}
	}
	return str
}