package tool

import (
	"github.com/labstack/echo/v4"
	"github.com/cnogo/pinpoint-go/agent"
)

// Echo Web框架中间件 具体处理函数可以通过traceContext.StartTraceSpanEven生成一个SpanEvent进行相关记录
func EchoPinpointTrace(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		traceContext := agent.StartTrace(&c.Request().Header)
		defer agent.FinishTrace(traceContext)
		if traceContext == nil {
			return next(c)
		}

		c.Set(agent.PP_CTX, traceContext)

		span := traceContext.Span
		span.SetAPIID(agent.GAgent.GetApiID("http." + c.Request().Method))
		span.SetRemoteAddr(c.RealIP())
		span.SetRpc(c.Request().URL.Path)

		err := next(c)

		if herr, ok := err.(*echo.HTTPError); ok {
			span.AddAnnotation(agent.NewInt32Annotation(agent.AK_HTTP_STATUS_CODE, int32(herr.Code)))
		} else {
			span.AddAnnotation(agent.NewInt32Annotation(agent.AK_HTTP_STATUS_CODE, int32(c.Response().Status)))
		}

		if err != nil {
			span.SetExceptionInfo(agent.GAgent.GetStrID(agent.STR_ERROR), err.Error())
		}

		return err
	}
}
