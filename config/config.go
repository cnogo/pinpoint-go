package config

type Config struct {
	AgentID         string
	ApplicationName string
	Pinpoint        struct {
		InfoAddr string //tcp agent info
		StatAddr string //udp agent stat
		SpanAddr string //udp span
	}
}

var Conf *Config

func InitConfig(conf *Config) {
	Conf = conf
}
