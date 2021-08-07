package serviceAuth

import (
	"crypto/rsa"
	"net/http"
	"time"

	errtrace "go-web-app/errtrace"
)

// Constant values
const (
	KeyCookieToken = "app_cookie_token"

	ErrorAuthLoginInfoInsufficient     = 100100
	ErrorAuthLoginUsernameNotFound     = 100200
	ErrorAuthLoginPasswordNotMatching  = 100300
	ErrorAuthLoginTokenGenerationError = 100400

	ErrorAuthTokenInvalid            = 100500
	ErrorAuthTokenInvalidEmpty       = 100510
	ErrorAuthTokenInvalidTtlTimeout  = 100520
	ErrorAuthTokenInvalidParseError  = 100530
	ErrorAuthTokenInvalidDecodeError = 100540

	IntervalTtlToken = int64(7 * 24 * time.Hour)

	HeaderAuthBearerPrefix    = "Bearer "
	HeaderAuthBearerPrefixLen = len(HeaderAuthBearerPrefix)
	HeaderAuthBearerKey       = "Authorization"
)

var errorTitle = map[int]string{
	ErrorAuthLoginInfoInsufficient:    "ErrorAuthLoginInfoInsufficient",
	ErrorAuthLoginUsernameNotFound:    "ErrorAuthLoginInfoInsufficient",
	ErrorAuthLoginPasswordNotMatching: "ErrorAuthLoginInfoInsufficient",

	ErrorAuthLoginTokenGenerationError: "500 Internal Server Error",

	ErrorAuthTokenInvalid:            "401 Unauthorized",
	ErrorAuthTokenInvalidEmpty:       "401 Unauthorized",
	ErrorAuthTokenInvalidTtlTimeout:  "401 Unauthorized",
	ErrorAuthTokenInvalidParseError:  "401 Unauthorized",
	ErrorAuthTokenInvalidDecodeError: "401 Unauthorized",
}
var errorStatus = map[int]int{
	ErrorAuthLoginInfoInsufficient:    http.StatusBadRequest,
	ErrorAuthLoginUsernameNotFound:    http.StatusBadRequest,
	ErrorAuthLoginPasswordNotMatching: http.StatusBadRequest,

	ErrorAuthLoginTokenGenerationError: http.StatusInternalServerError,

	ErrorAuthTokenInvalid:            http.StatusUnauthorized,
	ErrorAuthTokenInvalidEmpty:       http.StatusUnauthorized,
	ErrorAuthTokenInvalidTtlTimeout:  http.StatusUnauthorized,
	ErrorAuthTokenInvalidParseError:  http.StatusUnauthorized,
	ErrorAuthTokenInvalidDecodeError: http.StatusUnauthorized,
}

var (
	DebugMode = false
)

// GenerateError Get an error object
func GenerateError(code int, detail string, raw error) errtrace.Error {
	return errtrace.Error{
		ErrCode:    code,
		ErrTitle:   errorTitle[code],
		StatusCode: errorStatus[code],

		ErrDetail: detail,
		ErrRef:    raw,
	}
}

// EnvJWTToken Settings
var (
	//EnvJWTTokenPrivateKeyPath string
	//EnvJWTTokenPublicKeyPath  string

	EnvJWTTokenPrivateKey *rsa.PrivateKey
	EnvJWTTokenPublicKey  *rsa.PublicKey
)
