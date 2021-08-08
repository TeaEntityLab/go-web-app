package middleware

import (
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"
)

// ExceptionalRoute route descriptions for exceptions
type ExceptionalRoute struct {
	Path   string
	Method string
}

// Common middleware set common http header to http response
func Common(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(c *fasthttp.RequestCtx) {

		matchedPath := c.Request.URI().String()
		//for _, p := range c.Params {
		//	matchedPath = strings.Replace(matchedPath, p.Value, ":"+p.Key, 1)
		//}
		c.VisitUserValues(func(key []byte, val interface{}) {
			matchedPath = strings.Replace(matchedPath, fmt.Sprintf("%v", val), ":"+string(key), 1)
		})
		c.SetUserValue("matchedPath", matchedPath)

		// cors setting
		requestOriginDomain := string(c.Request.Header.Peek("Origin"))
		if len(requestOriginDomain) > 0 {
			c.Response.Header.Set("Access-Control-Allow-Origin", requestOriginDomain)
		} else {
			c.Response.Header.Set("Access-Control-Allow-Origin", "*")
		}

		c.Response.Header.Set("Access-Control-Allow-Methods", "POST, GET, PUT, PATCH, OPTIONS, DELETE")
		c.Response.Header.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Cache-Control, X-Requested-With")
		c.Response.Header.Set("Access-Control-Expose-Headers", "Content-Length")
		c.Response.Header.Set("Access-Control-Allow-Credentials", "true")
		c.Response.Header.Set("Cache-Control", "no-cache")

		if string(c.Request.Header.Method()) == "OPTIONS" {
			//c.AbortWithStatus(200)
			c.SetStatusCode(200)
		} else {
			next(c)
		}
	}
}

// Pipe middleware set common http header to http response
func Pipe(next fasthttp.RequestHandler, handlers ...func(fasthttp.RequestHandler) fasthttp.RequestHandler) fasthttp.RequestHandler {
	for i := len(handlers) - 1; i <= 0; i-- {
		next = handlers[i](next)
	}
	return next
}
