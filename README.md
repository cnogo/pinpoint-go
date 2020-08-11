### github.com/cnogo/pinpoint-go 项目包说明
+ agent包：处理agent连接、数据收发
+ config包：github.com/cnogo/pinpoint-go初始化配置
+ protocol: pinpoint相关协议处理
+ sqlx：mysql trace工具
+ tool: http-client、echo web框架 trace工具

---

### 接入方式

go get github.com/cnogo/pinpoint-go

--- 

### tool工具包说明

+ echo_tool说明：

实现了对echo web服务框架的请求的采集，只需一键接入即可，不入侵业务 e.Use(tool.Echo_PPTrace).
并把信息注入到echo.Context里面，以便后续请求拿到agent.Context记录信息 

+ http_tool: 

使用tool.NewRequest和tool.NewPPHttpClient可以把链路信息传给下一个服务。

### Sample

```
//初始化Config配置

address := "192.168.99.100"
conf := &config.Config{
	ApplicationName: "test-go",
	Pinpoint: struct {
		InfoAddr string
		StatAddr string
		SpanAddr string
	}{InfoAddr: address + ":9994", StatAddr: address + ":9995", SpanAddr: address + ":9996"},
}

config.InitConfig 

//创建一个Agent全局单例
agent.NewAgent()

//在合适位置启动
agent.GAgent.Start()

```

+ Echo Web框架中间件的使用

导入`tool.EchoPinpointTrace`中间件，即可启动对echo的trace，中间件会在echo的`context`中注入`*agent.TraceContext`,后续处理相关请求时，可以通过`context`获取`*agent.TraceContext`来展开对其他链路的追踪。
```

e := echo.New()

e.Use(tool.EchoPinpointTrace)

```

+ Http Client的使用

通过`tool.NewPPHttpClient(*agent.TraceContext)`调用获取一个client进行请求，内部自动追踪了重定向请求
```
client := tool.NewPPHttpClient(ctx)
resq, err := client.Get("http://localhost:6789/test")
```

### Sample
```

func TestAgent(t *testing.T) {
	address := "192.168.99.100"
	conf := &config.Config{
		AgentID:         "golang_pinpoint",
		ApplicationName: "test-go",
		Pinpoint: struct {
			InfoAddr string
			StatAddr string
			SpanAddr string
		}{InfoAddr: address + ":9994", StatAddr: address + ":9995", SpanAddr: address + ":9996"},
	}
        //初始化配置
	config.InitConfig(conf)
	//创建全局Agent
	agent.NewAgent()
    //启动Agent
	go agent.GAgent.Start()

	e := echo.New()
	
	//注册Echo中间件
	e.Use(tool.EchoPinpointTrace)



	e.GET("/user/:name", func(c echo.Context) error {
		return c.String(http.StatusOK, c.Param("name"))
	})

	e.GET("/age", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/info")
	})

	e.GET("/info", func(c echo.Context) error {
		return c.String(http.StatusOK,"the info age")
	})

	e.GET("/more", func(c echo.Context) error {


		client := tool.NewPPHttpClient(c.Get(agent.PP_CTX))
		resq, err := client.Get("http://localhost:6789/test")

		defer resq.Body.Close()

		if err != nil {
			fmt.Println(err)
		} else {
			p := make([]byte, 1204)
			io.ReadFull(resq.Body, p)
			fmt.Println(string(p))
		}


		return c.String(http.StatusOK, "its more")
	})

	fmt.Println(e.Start(":6789"))
}
```