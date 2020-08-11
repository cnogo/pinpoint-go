package agent

import (
	"testing"
)

func TestAgent(t *testing.T) {
	//address := "192.168.99.100"
	//conf := &config.Config{
	//	AgentID:         "golang_pinpoint",
	//	ApplicationName: "test-go",
	//	Pinpoint: struct {
	//		InfoAddr string
	//		StatAddr string
	//		SpanAddr string
	//	}{InfoAddr: address + ":9994", StatAddr: address + ":9995", SpanAddr: address + ":9996"},
	//}
	//
	//config.InitConfig(conf)
	//
	//agent := NewAgent()
	//
	////time.AfterFunc(2 * time.Second, func() {
	////	agent.AddAPIMeta("http.Get", 1,0)
	////	agent.AddAPIMeta("http.Post", 1, 0)
	////})
	////
	////time.AfterFunc(3 * time.Second, func() {
	////	simutorTrace()
	////})
	//
	//go agent.Start()
	//e := echo.New()
	//e.Use(Echo_PPTrace)
	//e.GET("/user/:name", func(c echo.Context) error {
	//	return c.String(http.StatusOK, c.Param("name"))
	//})
	//
	//e.GET("/age", func(c echo.Context) error {
	//	return c.Redirect(http.StatusTemporaryRedirect, "/info")
	//})
	//
	//e.GET("/info", func(c echo.Context) error {
	//
	//	//tc := c.Get(PP_CTX).(*TraceContext)
	//	////
	//	////
	//	//
	//	//sqlx.MustOpen()
	//	//sqlx.ConnectContext()
	//	//sqlx.MustConnect()
	//	//sqlx.Connect()
	//	////db, _ := NewPPSqlX(tc, "root:123456@tcp(192.168.99.100:3306)/demo")
	//	////db.Get()
	//
	//	return c.String(http.StatusOK,"the info age")
	//})
	//
	//e.GET("/more", func(c echo.Context) error {
	//
	//	client := NewPPHttpClient(c.Get(PP_CTX))
	//	resq, err := client.Get("http://localhost:6789/age")
	//
	//	defer resq.Body.Close()
	//
	//	if err != nil {
	//		fmt.Println(err)
	//	} else {
	//		p := make([]byte, 1204)
	//		io.ReadFull(resq.Body, p)
	//		fmt.Println(string(p))
	//	}
	//
	//
	//	return c.String(http.StatusOK, "its more")
	//})

	//fmt.Println(e.Start(":6789"))
}

