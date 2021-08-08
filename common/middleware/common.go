package middleware

import (
	"github.com/gin-gonic/gin"
	"strings"

	"github.com/valyala/fasthttp"
)

// ExceptionalRoute route descriptions for exceptions
type ExceptionalRoute struct {
	Path   string
	Method string
}

// Common middleware set common http header to http response
func Common() gin.HandlerFunc {
	return func(c *fasthttp.RequestCtx) {

		matchedPath := c.Request.URL.String()
		for _, p := range c.Params {
			matchedPath = strings.Replace(matchedPath, p.Value, ":"+p.Key, 1)
		}
		c.Set("matchedPath", matchedPath)

		// cors setting
		requestOriginDomain := string(c.Request.Header.Peek("Origin"))
		if len(requestOriginDomain) > 0 {
			c.Writer.Header().Set("Access-Control-Allow-Origin", requestOriginDomain)
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, PATCH, OPTIONS, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Cache-Control", "no-cache")

		if string(c.Request.Header.Method()) == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}
