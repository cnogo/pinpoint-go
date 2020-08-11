package agent

var ServiceType map[int]string

const (
	STAND_ALONE int16 = 1000

	ACTIVEMQ_CLIENT          int16 = 8310
	ACTIVEMQ_CLIENT_INTERNAL int16 = 8311

	AKKA_HTTP_SERVER int16 = 1310

	STR1311 int16 = 9998

	ARCUS                    int16 = 8100
	ARCUS_FUTURE_GET         int16 = 8101
	ARCUS_EHCACHE_FUTURE_GET int16 = 8102
	ARCUS_INTERNAL           int16 = 8103

	MEMCACHED            int16 = 8050
	MEMCACHED_FUTURE_GET int16 = 8051

	CASSANDRA               int16 = 2600
	CASSANDRA_EXECUTE_QUERY int16 = 2601

	CUBRID               int16 = 2400
	CUBRID_EXECUTE_QUERY int16 = 2401

	CXF_CLIENT int16 = 9080

	DBCP  int16 = 6050
	DBCP2 int16 = 6052

	DUBBO_PROVIDER int16 = 1110
	DUBBO_CONSUMER int16 = 9110
	DUBBO          int16 = 9111

	GOOGLE_HTTP_CLIENT_INTERNAL int16 = 9054
	GSON                        int16 = 5010
	HIKARICP                    int16 = 6060

	HTTP_CLIENT_3          int16 = 9050
	HTTP_CLIENT_4          int16 = 9052
	HTTP_CLIENT_4_INTERNAL int16 = 9053

	HYSTRIX_COMMAND          int16 = 9120
	HYSTRIX_COMMAND_INTERNAL int16 = 9121

	IBATIS        int16 = 5500
	IBATIS_SPRING int16 = 5501

	JACKSON int16 = 5011

	JBOSS        int16 = 1040
	JBOSS_METHOD int16 = 1041

	JDK_HTTPURLCONNECTOR int16 = 9055

	JETTY        int16 = 1030
	JETTY_METHOD int16 = 1031

	JSON_LIB int16 = 5012
	JSP      int16 = 5005

	MSSQLSERVER         int16 = 2200
	MSSQL_EXECUTE_QUERY int16 = 2201

	KAFKA_CLIENT          int16 = 8660
	KAFKA_CLIENT_INTERNAL int16 = 8661

	MARIADB               int16 = 2150
	MARIADB_EXECUTE_QUERY int16 = 2151

	MYBATIS int16 = 5510

	MYSQL               int16 = 2100
	MYSQL_EXECUTE_QUERY int16 = 2101

	NETTY          int16 = 9150
	NETTY_INTERNAL int16 = 9151
	NETTY_HTTP     int16 = 9152

	ASYNC_HTTP_CLIENT          int16 = 9056
	ASYNC_HTTP_CLIENT_INTERNAL int16 = 9057

	OK_HTTP_CLIENT          int16 = 9058
	OK_HTTP_CLIENT_INTERNAL int16 = 9059

	ORACLE               int16 = 2300
	ORACLE_EXECUTE_QUERY int16 = 2301

	PHP               int16 = 1500
	PHP_METHOD        int16 = 1501
	PHP_REMOTE_METHOD int16 = 9700

	POSTGRESQL               int16 = 2500
	POSTGRESQL_EXECUTE_QUERY int16 = 2501

	RABBITMQ_CLIENT          int16 = 8300
	RABBITMQ_CLIENT_INTERNAL int16 = 8301

	REDIS int16 = 8200
	RESIN int16 = 1200

	RESIN_METHOD  int16 = 1201
	REST_TEMPLATE int16 = 9140

	RX_JAVA          int16 = 6500
	RX_JAVA_INTERNAL int16 = 6501

	SPRING       int16 = 5050
	SPRING_BEAN  int16 = 5071
	SPRING_ASYNC int16 = 5052
	SPRING_MVC   int16 = 5051

	NAME int16 = 1210

	THRIFT_SERVER          int16 = 1100
	THRIFT_CLIENT          int16 = 9100
	THRIFT_SERVER_INTERNAL int16 = 1101
	THRIFT_CLIENT_INTERNAL int16 = 9101

	TOMCAT        int16 = 1010
	TOMCAT_METHOD int16 = 1011

	UNDERTOW        int16 = 1120
	UNDERTOW_METHOD int16 = 1121

	VERTX                      int16 = 1050
	VERTX_INTERNAL             int16 = 1051
	VERTX_HTTP_SERVER          int16 = 1052
	VERTX_HTTP_SERVER_INTERNAL int16 = 1053
	VERTX_HTTP_CLIENT          int16 = 9130
	VERTX_HTTP_CLIENT_INTERNAL int16 = 9131

	WEBLOGIC        int16 = 1070
	WEBLOGIC_METHOD int16 = 1071

	WEBSPHERE        int16 = 1060
	WEBSPHERE_METHOD int16 = 1061
)

func init() {
	ServiceType = make(map[int]string)

	ServiceType[1000] = "STAND_ALONE"

	// activemq.client
	ServiceType[8310] = "ACTIVEMQ_CLIENT"
	ServiceType[8311] = "ACTIVEMQ_CLIENT_INTERNAL"

	// akka.http

	ServiceType[1310] = "AKKA_HTTP_SERVER"
	ServiceType[9998] = "1311"

	// arcus
	ServiceType[8100] = "ARCUS"
	ServiceType[8101] = "ARCUS_FUTURE_GET"
	ServiceType[8102] = "ARCUS_EHCACHE_FUTURE_GET"

	ServiceType[8103] = "ARCUS_INTERNAL"
	ServiceType[8050] = "MEMCACHED"
	ServiceType[8051] = "MEMCACHED_FUTURE_GET"

	// cassandra
	ServiceType[2600] = "CASSANDRA"
	ServiceType[2601] = "CASSANDRA_EXECUTE_QUERY"

	// cubrid
	ServiceType[2400] = "CUBRID"
	ServiceType[2401] = "CUBRID_EXECUTE_QUERY"

	// cxf
	ServiceType[9080] = "CXF_CLIENT"

	// dbcp
	ServiceType[6050] = "DBCP"
	// dbcp2
	ServiceType[6052] = "DBCP2"

	// dubbo
	ServiceType[1110] = "DUBBO_PROVIDER"
	ServiceType[9110] = "DUBBO_CONSUMER"
	ServiceType[9111] = "DUBBO"

	// httpclient
	ServiceType[9054] = "GOOGLE_HTTP_CLIENT_INTERNAL"

	// gson
	ServiceType[5010] = "GSON"

	// hikaricp
	ServiceType[6060] = "HIKARICP"

	// httpclient3
	ServiceType[9050] = "HTTP_CLIENT_3"

	// httpclient4
	ServiceType[9052] = "HTTP_CLIENT_4"
	ServiceType[9053] = "HTTP_CLIENT_4_INTERNAL"

	// hystrix
	ServiceType[9120] = "HYSTRIX_COMMAND"
	ServiceType[9121] = "HYSTRIX_COMMAND_INTERNAL"

	// ibatis
	ServiceType[5500] = "IBATIS"
	ServiceType[5501] = "IBATIS_SPRING"

	// jackson
	ServiceType[5011] = "JACKSON"

	// jboss
	ServiceType[1040] = "JBOSS"
	ServiceType[1041] = "JBOSS_METHOD"

	// jdk.http
	ServiceType[9055] = "JDK_HTTPURLCONNECTOR"

	// jetty
	ServiceType[1030] = "JETTY"
	ServiceType[1031] = "JETTY_METHOD"

	// json_lib
	ServiceType[5012] = "JSON-LIB"

	// jsp
	ServiceType[5005] = "JSP"

	// jtds
	ServiceType[2200] = "MSSQLSERVER"
	ServiceType[2201] = "MSSQL_EXECUTE_QUERY"
	// kafka
	ServiceType[8660] = "KAFKA_CLIENT"
	ServiceType[8661] = "KAFKA_CLIENT_INTERNAL"
	// mariadb
	ServiceType[2150] = "MARIADB"
	ServiceType[2151] = "MARIADB_EXECUTE_QUERY"
	// mybatis
	ServiceType[5510] = "MYBATIS"

	// mysql
	ServiceType[2100] = "MYSQL"
	ServiceType[2101] = "MYSQL_EXECUTE_QUERY"

	// netty
	ServiceType[9150] = "NETTY"
	ServiceType[9151] = "NETTY_INTERNAL"
	ServiceType[9152] = "NETTY_HTTP"

	// asynchttpclient
	ServiceType[9056] = "ASYNC_HTTP_CLIENT"
	ServiceType[9057] = "ASYNC_HTTP_CLIENT_INTERNAL"

	// okhttp
	ServiceType[9058] = "OK_HTTP_CLIENT"
	ServiceType[9059] = "OK_HTTP_CLIENT_INTERNAL"

	// oracle
	ServiceType[2300] = "ORACLE"
	ServiceType[2301] = "ORACLE_EXECUTE_QUERY"

	// php
	ServiceType[1500] = "PHP"
	ServiceType[1501] = "PHP_METHOD"
	ServiceType[9700] = "PHP_REMOTE_METHOD"

	// postgresql
	ServiceType[1500] = "PHP"
	ServiceType[1501] = "PHP_METHOD"
	ServiceType[9700] = "PHP_REMOTE_METHOD"

	// php
	ServiceType[1500] = "PHP"
	ServiceType[1501] = "PHP_METHOD"
	ServiceType[9700] = "PHP_REMOTE_METHOD"

	// postgresql
	ServiceType[2500] = "POSTGRESQL"
	ServiceType[2501] = "POSTGRESQL_EXECUTE_QUERY"

	// rabbitmq.client
	ServiceType[8300] = "RABBITMQ_CLIENT"
	ServiceType[8301] = "RABBITMQ_CLIENT_INTERNAL"

	// redis
	ServiceType[8200] = "REDIS"

	// resin
	ServiceType[1200] = "RESIN"
	ServiceType[1201] = "RESIN_METHOD"

	// resttemplate
	ServiceType[9140] = "REST_TEMPLATE"

	// rxjava
	ServiceType[6500] = "RX_JAVA"
	ServiceType[6501] = "RX_JAVA_INTERNAL"

	// spring.beans
	ServiceType[5071] = "SPRING_BEAN"
	// spring.async
	ServiceType[5052] = "SPRING_ASYNC"
	// spring.web
	ServiceType[5051] = "SPRING_MVC"

	// spring.boot
	ServiceType[1210] = "NAME"

	ServiceType[1100] = "THRIFT_SERVER"
	ServiceType[9100] = "THRIFT_CLIENT"
	ServiceType[1101] = "THRIFT_SERVER_INTERNAL"
	ServiceType[9101] = "THRIFT_CLIENT_INTERNAL"

	// tomcat
	ServiceType[1010] = "TOMCAT"
	ServiceType[1011] = "TOMCAT_METHOD"

	// tomcat
	ServiceType[1010] = "TOMCAT"
	ServiceType[1011] = "TOMCAT_METHOD"

	// undertow
	ServiceType[1120] = "UNDERTOW"
	ServiceType[1121] = "UNDERTOW_METHOD"

	// vertx
	ServiceType[1050] = "VERTX"
	ServiceType[1051] = "VERTX_INTERNAL"
	ServiceType[1052] = "VERTX_HTTP_SERVER"
	ServiceType[1053] = "VERTX_HTTP_SERVER_INTERNAL"
	ServiceType[9130] = "VERTX_HTTP_CLIENT"
	ServiceType[9131] = "VERTX_HTTP_CLIENT_INTERNAL"

	// weblogic
	ServiceType[1070] = "WEBLOGIC"
	ServiceType[1071] = "WEBLOGIC_METHOD"

	// weblogic
	ServiceType[1060] = "WEBSPHERE"
	ServiceType[1061] = "WEBSPHERE_METHOD"
}
