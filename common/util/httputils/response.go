package httputils

import (
	"github.com/sirupsen/logrus"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"

	authService "go-web-app/common/service/auth"
	"go-web-app/errtrace"
)

var (
	StrContentType     = []byte("Content-Type")
	StrApplicationJSON = []byte("application/json")
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
	DoJSONWrite(c, statusCode, H{
		"title":  response.ErrTitle,
		"path":   string(c.Request.URI().Path()),
		"method": string(c.Request.Header.Method()),

		"code":   response.ErrCode,
		"status": statusCode,
	})

	return response
}

func DoJSONWrite(ctx *fasthttp.RequestCtx, code int, obj interface{}) {
	ctx.Response.Header.SetCanonical(StrContentType, StrApplicationJSON)
	ctx.Response.SetStatusCode(code)
	start := time.Now()
	if err := jsoniter.NewEncoder(ctx).Encode(obj); err != nil {
		elapsed := time.Since(start)
		logrus.Error("Error: ", elapsed, err.Error(), obj)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}
