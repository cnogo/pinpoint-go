package pinpoint_tracing

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/cnogo/pinpoint-go/agent"
	"github.com/cnogo/pinpoint-go/config"
	"github.com/cnogo/pinpoint-go/tool"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestAgent(t *testing.T) {
	address := "192.168.99.100"
	conf := &config.Config{
		//AgentID:         "192.168.99.10",
		ApplicationName: "test-go",
		Pinpoint: struct {
			InfoAddr string
			StatAddr string
			SpanAddr string
		}{InfoAddr: address + ":9994", StatAddr: address + ":9995", SpanAddr: address + ":9996"},
	}

	config.InitConfig(conf)
	agent.NewAgent()

	ghttp := tool.NewPPHttpClient(tool.WithHttpTimeOut(10*time.Second))

	go agent.GAgent.Start()

	e := echo.New()
	e.Use(tool.EchoPinpointTrace)


	e.GET("/user/:name", func(c echo.Context) error {
		return c.String(http.StatusOK, c.Param("name"))
	})

	e.GET("/age", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/info")
	})

	e.GET("/info", func(c echo.Context) error {
		return c.String(http.StatusOK, "the info age")
	})

	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "its test ok")
	})

	e.GET("/more", func(c echo.Context) error {

		url := "http://localhost:6789/test"
		request, _ := tool.NewRequest(c, http.MethodGet, url, nil)
		resq, err := ghttp.Do(request)
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
