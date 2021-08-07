package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"

	"go-web-app/common/util/ginutils"
)

// Auth general Auth checking
func Auth(logger *logrus.Entry, exceptionalRoutes ...ExceptionalRoute) gin.HandlerFunc {
	return func(c *gin.Context) {

		matchedPath, _ := c.Get("matchedPath")

		matchedException := false
		path := c.Request.URL.Path
		method := c.Request.Method
		for _, route := range exceptionalRoutes {
			if (route.Path == path || route.Path == matchedPath) && route.Method == method {
				matchedException = true
				break
			}
		}
		if matchedException {
			c.Next()
			return
		}

		requestXID := xid.New().String()
		funcLogger := ginutils.NewHttpLogger(logger, c.Request).
			WithField("requestXID", requestXID).
			WithField("func", "go-web-app.AuthMiddleware")

		authToken, checkLoginErr := ginutils.CheckLoginStatusOrAbort(c, funcLogger, fmt.Sprintf("c.Request.URL.Path: %v", c.Request.URL.Path))
		if checkLoginErr != nil || authToken == nil {
			return
		}

		c.Set("authToken", authToken)

		c.Next()
	}
}
