package httputils

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"

	"go-web-app/common/model"
	authService "go-web-app/common/service/auth"
	"go-web-app/errtrace"
)

func GetHttpRequestInfo(r *fasthttp.RequestCtx) *model.HTTPRequestInfo {
	return &model.HTTPRequestInfo{
		Remote:        r.RemoteAddr().String(),
		Method:        string(r.Method()),
		RequestURI:    string(r.RequestURI()),
		Protocol:      string(r.Request.Header.Protocol()),
		Host:          string(r.Host()),
		ContentLength: int64(r.Request.Header.ContentLength()),
		Referer:       string(r.Referer()),
		UserAgent:     string(r.UserAgent()),
	}
}

// NewHttpLogger
// field names in httpRequest map fields should follow Stackdriver Logging API v2 HttpRequest
// https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#HttpRequest
func NewHttpLogger(logger *logrus.Entry, r *fasthttp.RequestCtx) *logrus.Entry {
	return logger.WithField("httpRequest", map[string]interface{}{
		"requestMethod": string(r.Method()),
		"requestUrl":    string(r.RequestURI()),
		"requestSize":   r.Request.Header.ContentLength(),
		//"status": number,
		//"responseSize": string,
		"userAgent": string(r.UserAgent()),
		"remoteIp":  r.RemoteAddr().String(),
		//"serverIp":  string,
		"referer": string(r.Referer()),
		//"latency": string,
		//"cacheLookup": boolean,
		//"cacheHit": boolean,
		//"cacheValidatedWithOriginServer": boolean,
		//"cacheFillBytes": string,
		//"protocol": r.Request.Proto,
		"protocol": string(r.Request.Header.Protocol()),
	})
}

// TryGetCookie ...
func TryGetCookie(c *fasthttp.RequestCtx, name string) (string, error) {
	targetCookie := string(c.Request.Header.Cookie(name))
	//targetCookie, cookieTokenErr := c.Request.Header.Cookie(name)
	//if cookieTokenErr != nil {
	//	return "", cookieTokenErr
	//}

	var target string
	var urlQueryUnescapeErr error
	if targetCookie != "" {
		target, urlQueryUnescapeErr = url.QueryUnescape(targetCookie)
	}

	return target, urlQueryUnescapeErr
}

// CheckLoginStatus ...
func CheckLoginStatus(c *fasthttp.RequestCtx) (*model.AuthToken, *errtrace.Error) {
	authTokenJWTString := GetBearerAuthTokenString(c)
	if authService.DebugMode {
		fmt.Println("Auth:", authTokenJWTString)
	}

	// authTokenJWTString, cookieTokenErr := TryGetCookie(c, authService.KeyCookieToken)
	// if cookieTokenErr != nil || authTokenJWTString == "" {

	if authTokenJWTString == "" {
		response := ResponseAuthError(
			c,
			authService.ErrorAuthTokenInvalid,
			"AuthToken is not found or accessing cookie error",
			nil, true)

		return nil, &response
	}
	authToken, checkErr := authService.CheckAuthTokenValidation(authTokenJWTString)
	if authToken == nil || checkErr != nil {
		if checkErr.ErrCode == authService.ErrorAuthTokenInvalidEmpty {
			response := ResponseAuthError(
				c,
				checkErr.ErrCode,
				"AuthToken is empty",
				checkErr, true)

			return nil, &response
		}
		if checkErr.ErrCode == authService.ErrorAuthTokenInvalidTtlTimeout {
			response := ResponseAuthError(
				c,
				checkErr.ErrCode,
				"AuthToken is out of date",
				checkErr, true)

			return nil, &response
		}
		if checkErr.ErrCode == authService.ErrorAuthTokenInvalidDecodeError {
			response := ResponseAuthError(
				c,
				checkErr.ErrCode,
				"AuthToken is out of date",
				checkErr, true)

			return nil, &response
		}

		// Unexpected errors
		response := ResponseAuthError(
			c,
			authService.ErrorAuthTokenInvalidParseError,
			"AuthToken couldn't be parsed",
			checkErr, true)
		return nil, &response
	}

	return authToken, nil
}

func GetBearerAuthTokenString(c *fasthttp.RequestCtx) string {
	str := string(c.Request.Header.Peek(authService.HeaderAuthBearerKey))
	if idx := strings.Index(str, authService.HeaderAuthBearerPrefix); idx != -1 {
		str = str[idx+authService.HeaderAuthBearerPrefixLen:]
	}
	return str
}

func CheckLoginStatusOrAbort(c *fasthttp.RequestCtx, funcLogger *logrus.Entry, messageOnError string) (*model.AuthToken, *errtrace.Error) {
	authToken, checkLoginErr := CheckLoginStatus(c)

	if checkLoginErr != nil || authToken == nil {
		var checkLoginErrStr string
		if checkLoginErr != nil {
			checkLoginErrStr = checkLoginErr.Error()
		}

		fields := logrus.Fields{
			"ip":        ReadUserIP(c),
			"authToken": authToken,
			"error":     checkLoginErrStr,

			"code":   401,
			"status": 401,
		}
		funcLogger.WithFields(fields).Infof(messageOnError)
		DoJSONWrite(c, http.StatusUnauthorized, fields)
	}

	return authToken, checkLoginErr
}

// ReadUserIP ...
func ReadUserIP(r *fasthttp.RequestCtx) string {
	IPAddress := string(r.Request.Header.Peek("X-Real-Ip"))
	if IPAddress == "" {
		IPAddress = string(r.Request.Header.Peek("X-Forwarded-For"))
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr().String()
	}
	return IPAddress
}
