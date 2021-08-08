package ginutils

import (
	"github.com/gin-gonic/gin"
	"github.com/valyala/fasthttp"

	authService "go-web-app/common/service/auth"
	"go-web-app/errtrace"
)

// ResponseAuthError ...
func ResponseAuthError(c *fasthttp.RequestCtx, errCode int, detail string, raw error, skipReponse bool) errtrace.Error {
	response := authService.GenerateError(
		errCode,
		detail,
		raw)

	if skipReponse {
		return response
	}

	return responseError(c, response)
}

// responseError ...
func responseError(c *fasthttp.RequestCtx, response errtrace.Error) errtrace.Error {
	statusCode := response.StatusCode
	c.AbortWithStatusJSON(statusCode, gin.H{
		"title":  response.ErrTitle,
		"path":   c.Request.URL.Path,
		"method": string(c.Request.Header.Method()),

		"code":   response.ErrCode,
		"status": statusCode,
	})

	return response
}
