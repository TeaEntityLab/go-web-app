package middleware

import (
	"fmt"

	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"

	"go-web-app/common/util/httputils"
)

// AuthWithExceptionRoutes general Auth checking
func AuthWithExceptionRoutes(logger *logrus.Entry, next fasthttp.RequestHandler, exceptionalRoutes ...ExceptionalRoute) fasthttp.RequestHandler {
	return func(c *fasthttp.RequestCtx) {

		matchedPath := c.UserValue("matchedPath")

		matchedException := false
		path := string(c.Request.URI().Path())
		method := string(c.Request.Header.Method())
		for _, route := range exceptionalRoutes {
			if (route.Path == path || route.Path == matchedPath) && route.Method == method {
				matchedException = true
				break
			}
		}
		if matchedException {
			next(c)
			return
		}

		Auth(logger, next)
	}
}

// Auth general Auth checking
func Auth(logger *logrus.Entry, next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(c *fasthttp.RequestCtx) {

		requestXID := xid.New().String()
		funcLogger := httputils.NewHttpLogger(logger, c).
			WithField("requestXID", requestXID).
			WithField("func", "go-web-app.AuthMiddleware")

		authToken, checkLoginErr := httputils.CheckLoginStatusOrAbort(c, funcLogger, fmt.Sprintf("string(c.Request.URI().Path()): %v", string(c.Request.URI().Path())))
		if checkLoginErr != nil || authToken == nil {
			return
		}

		c.SetUserValue("authToken", authToken)

		next(c)
	}
}
