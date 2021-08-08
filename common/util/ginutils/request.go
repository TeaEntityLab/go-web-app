package ginutils

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"go-web-app/common/model"
	authService "go-web-app/common/service/auth"
	"go-web-app/errtrace"
)

func GetHttpRequestInfo(r *http.Request) *model.HTTPRequestInfo {
	return &model.HTTPRequestInfo{
		Remote:        r.RemoteAddr,
		Method:        r.Method,
		RequestURI:    r.RequestURI,
		Protocol:      r.Proto,
		Host:          r.Host,
		ContentLength: r.ContentLength,
		Referer:       r.Referer(),
		UserAgent:     r.UserAgent(),
	}
}

// NewHttpLogger
// field names in httpRequest map fields should follow Stackdriver Logging API v2 HttpRequest
// https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#HttpRequest
func NewHttpLogger(logger *logrus.Entry, r *http.Request) *logrus.Entry {
	return logger.WithField("httpRequest", map[string]interface{}{
		"requestMethod": r.Method,
		"requestUrl":    r.RequestURI,
		"requestSize":   r.ContentLength,
		//"status": number,
		//"responseSize": string,
		"userAgent": r.UserAgent(),
		"remoteIp":  r.RemoteAddr,
		//"serverIp":  string,
		"referer": r.Referer(),
		//"latency": string,
		//"cacheLookup": boolean,
		//"cacheHit": boolean,
		//"cacheValidatedWithOriginServer": boolean,
		//"cacheFillBytes": string,
		"protocol": r.Proto,
	})
}

// TryGetCookie ...
func TryGetCookie(c *gin.Context, name string) (string, error) {
	targetCookie, cookieTokenErr := c.Request.Cookie(name)
	if cookieTokenErr != nil {
		return "", cookieTokenErr
	}

	var target string
	var urlQueryUnescapeErr error
	if targetCookie != nil {
		target, urlQueryUnescapeErr = url.QueryUnescape(targetCookie.Value)
	}

	return target, urlQueryUnescapeErr
}

// CheckLoginStatus ...
func CheckLoginStatus(c *gin.Context) (*model.AuthToken, *errtrace.Error) {
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
		if checkErr != nil && checkErr.ErrCode == authService.ErrorAuthTokenInvalidEmpty {
			response := ResponseAuthError(
				c,
				checkErr.ErrCode,
				"AuthToken is empty",
				checkErr, true)

			return nil, &response
		}
		if checkErr != nil && checkErr.ErrCode == authService.ErrorAuthTokenInvalidTtlTimeout {
			response := ResponseAuthError(
				c,
				checkErr.ErrCode,
				"AuthToken is out of date",
				checkErr, true)

			return nil, &response
		}
		if checkErr != nil && checkErr.ErrCode == authService.ErrorAuthTokenInvalidDecodeError {
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

func GetBearerAuthTokenString(c *gin.Context) string {
	str := c.Request.Header.Get(authService.HeaderAuthBearerKey)
	if idx := strings.Index(str, authService.HeaderAuthBearerPrefix); idx != -1 {
		str = str[idx+authService.HeaderAuthBearerPrefixLen:]
	}
	return str
}

func CheckLoginStatusOrAbort(c *gin.Context, funcLogger *logrus.Entry, messageOnError string) (*model.AuthToken, *errtrace.Error) {
	authToken, checkLoginErr := CheckLoginStatus(c)

	if checkLoginErr != nil || authToken == nil {
		var checkLoginErrStr string
		if checkLoginErr != nil {
			checkLoginErrStr = checkLoginErr.Error()
		}

		fields := logrus.Fields{
			"ip":        ReadUserIP(c.Request),
			"authToken": authToken,
			"error":     checkLoginErrStr,

			"code":   401,
			"status": 401,
		}
		funcLogger.WithFields(fields).Infof(messageOnError)
		c.AbortWithStatusJSON(http.StatusUnauthorized, fields)
	}

	return authToken, checkLoginErr
}

// ReadUserIP ...
func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}
