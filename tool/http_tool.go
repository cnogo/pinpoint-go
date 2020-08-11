package tool

import (
	"context"
	"crypto/tls"
	"github.com/labstack/echo/v4"
	"github.com/cnogo/pinpoint-go/agent"
	"io"
	"net"
	"net/http"
	"time"
)

type ppHttpConfig struct {
	TimeOut               time.Duration
	MaxIdleConnsPerHost   int
	IdleConnTimeout       time.Duration
	TLSHandshakeTimeout   time.Duration
	ExpectContinueTimeout time.Duration
	DialTimeout           time.Duration
	DialKeepAlive         time.Duration
	TlsConfig             *tls.Config
}

//设置默认配置
func (p *ppHttpConfig) setDefault() {
	p.TimeOut = 0
	p.MaxIdleConnsPerHost = 500
	p.IdleConnTimeout = 90 * time.Second
	p.TLSHandshakeTimeout = 10 * time.Second
	p.ExpectContinueTimeout = 1 * time.Second
	p.DialTimeout = 30 * time.Second
	p.DialKeepAlive = 30 * time.Second
}

type httpConfigurer func(configer *ppHttpConfig)

//请求超时设置
func WithHttpTimeOut(timeOut time.Duration) httpConfigurer {
	return func(configer *ppHttpConfig) {
		configer.TimeOut = timeOut
	}
}

//每个host最大空闲链接
func WithHttpMaxIdleConnsPerHost(count int) httpConfigurer {
	return func(configer *ppHttpConfig) {
		configer.MaxIdleConnsPerHost = count
	}
}

//空闲链接多少释放
func WithHttpIdleConnTimeout(timeout time.Duration) httpConfigurer {
	return func(configer *ppHttpConfig) {
		configer.IdleConnTimeout = timeout
	}
}

//TLS握手超时
func WithHttpTLSHandshakeTimeout(timeout time.Duration) httpConfigurer {
	return func(configer *ppHttpConfig) {
		configer.TLSHandshakeTimeout = timeout
	}
}

func WithHttpExpectContinueTimeout(timeout time.Duration) httpConfigurer {
	return func(configer *ppHttpConfig) {
		configer.ExpectContinueTimeout = timeout
	}
}

func WithHttpDialTimeout(timeout time.Duration) httpConfigurer {
	return func(configer *ppHttpConfig) {
		configer.DialTimeout = timeout
	}
}

func WithHttpDialKeepAlive(timeout time.Duration) httpConfigurer {
	return func(configer *ppHttpConfig) {
		configer.DialKeepAlive = timeout
	}
}

func WithHttpTlsClientConfig(tlsConfig *tls.Config) httpConfigurer {
	return func(configer *ppHttpConfig) {
		configer.TlsConfig = tlsConfig
	}
}

var DefaultTransport http.RoundTripper = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	MaxIdleConnsPerHost:   200,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

type Transport struct {
	http.RoundTripper
}

func (p *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	rt := p.RoundTripper
	if rt == nil {
		rt = DefaultTransport
	}

	//获取trace上下文
	traceContext, _ := req.Context().Value(agent.PP_CTX).(*agent.TraceContext)

	//说明是个根请求,尝试创建一个trace
	if traceContext == nil {
		tc := agent.StartTrace(nil)
		defer agent.FinishTrace(tc)

		//说明不采集
		if tc == nil {
			return rt.RoundTrip(req)
		}

		traceContext = tc

		span := tc.Span
		span.SetAPIID(agent.GAgent.GetApiID("http." + req.Method))
		span.SetRpc(req.URL.Path)

	}

	spanEvent := traceContext.StartTraceSpanEvent()
	spanEvent.SetApiID(agent.GAgent.GetApiID("http." + req.Method))
	spanEvent.SetServiceType(agent.PHP_METHOD)
	spanEvent.AddAnnotation(agent.NewStringAnnotation(agent.AK_ARGS0, req.URL.Path))

	traceHead := traceContext.GetNextSpanInfo()
	for k, v := range traceHead {
		req.Header.Set(k, v)
	}
	resp, err := rt.RoundTrip(req)

	if resp != nil {
		spanEvent.AddAnnotation(agent.NewInt32Annotation(agent.AK_HTTP_STATUS_CODE, int32(resp.StatusCode)))
	}

	if err != nil {
		spanEvent.SetExceptionInfo(agent.GAgent.GetStrID(agent.STR_ERROR), err.Error())
	}

	spanEvent.Finish()

	return resp, err
}


func NewRequest(ctx interface{}, method, url string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	cx, _ := ctx.(*agent.TraceContext)
	echoCtx, _ := ctx.(echo.Context)
	if cx == nil && echoCtx != nil {
		cx, _ = echoCtx.Get(agent.PP_CTX).(*agent.TraceContext)
	}

	if cx == nil {
		return request, nil
	}

	parentCtx := request.Context()
	request = request.WithContext(context.WithValue(parentCtx, agent.PP_CTX, cx))
	return request, err
}

func NewPPHttpClient(configurers ...httpConfigurer) *http.Client {
	httpConfig := new(ppHttpConfig)
	httpConfig.setDefault()

	//更新配置
	for _, configurer := range configurers {
		configurer(httpConfig)
	}

	transport :=  &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   httpConfig.DialTimeout,
			KeepAlive: httpConfig.DialKeepAlive,
		}).DialContext,
		MaxIdleConnsPerHost:   httpConfig.MaxIdleConnsPerHost,
		IdleConnTimeout:       httpConfig.IdleConnTimeout,
		TLSHandshakeTimeout:   httpConfig.TLSHandshakeTimeout,
		ExpectContinueTimeout: httpConfig.ExpectContinueTimeout,
		TLSClientConfig: httpConfig.TlsConfig,
	}

	return &http.Client{Transport: &Transport{RoundTripper: transport}, Timeout: httpConfig.TimeOut}
}
