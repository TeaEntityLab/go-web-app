package route

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	mod "go-web-app/common/model"
	repo "go-web-app/common/repository"
	"go-web-app/common/util/ginutils"
)

const (
	KeyCacheLinkWithUniqueKey = "LINK_%s"

	ErrorHTTPFailureDataCenterServerUpdateModelByIDs = 200100

	COMMON_RESPONSE_CODE_SUCCESS                   = 0
	COMMON_RESPONSE_CODE_SUCCESS_STR               = "0"
	COMMON_RESPONSE_CODE_BAD_REQUEST               = 400
	COMMON_RESPONSE_CODE_BAD_REQUEST_STR           = "400"
	COMMON_RESPONSE_CODE_FORBIDDEN                 = 403
	COMMON_RESPONSE_CODE_FORBIDDEN_STR             = "403"
	COMMON_RESPONSE_CODE_INTERNAL_SERVER_ERROR     = 500
	COMMON_RESPONSE_CODE_INTERNAL_SERVER_ERROR_STR = "500"

	COMMON_RESPONSE_STATUS_SUCCESS                   = 200
	COMMON_RESPONSE_STATUS_SUCCESS_STR               = "200"
	COMMON_RESPONSE_STATUS_BAD_REQUEST               = 400
	COMMON_RESPONSE_STATUS_BAD_REQUEST_STR           = "400"
	COMMON_RESPONSE_STATUS_FORBIDDEN                 = 403
	COMMON_RESPONSE_STATUS_FORBIDDEN_STR             = "403"
	COMMON_RESPONSE_STATUS_INTERNAL_SERVER_ERROR     = 500
	COMMON_RESPONSE_STATUS_INTERNAL_SERVER_ERROR_STR = "500"
)

var (
	Logger *logrus.Entry

	CustomDomainBindingEndpointDomainName string

	DataCenterServerInternalAccessToken string
)

type commonRequest struct {
	mod.HTTPRequestInfo

	RequestId string `json:"request_id"`
}
type CommonResponse struct {
	RequestId     string `json:"request_id"`
	Status        string `json:"status"`
	StatusMessage string `json:"status_message,omitempty"`
	Code          int    `json:"code"`

	Data  interface{} `json:"data,omitempty"`
	Count *int64      `json:"count,omitempty"`
}
type CommonErrorResponse struct {
	CommonResponse

	AuthToken *mod.AuthToken
	Ip        string `json:"ip"`

	Details interface{} `json:"details"`

	Error string `json:"error"`
}

func checkModelErrorOrAbort(c *gin.Context, funcLogger *logrus.Entry, requestId string, authToken *mod.AuthToken, modelError error) bool {
	if modelError != nil {
		httpStatus := http.StatusInternalServerError
		code := 500
		status := 500
		if repo.IsTypeCheckInvalidFields(modelError) ||
			repo.IsTypeCheckNonExistentObject(modelError) ||
			repo.IsBadRequest(modelError) {
			httpStatus = http.StatusBadRequest
			code = 400
			status = 400
		}

		fields := logrus.Fields{
			"ip":        ginutils.ReadUserIP(c.Request),
			"authToken": authToken,
			"error":     modelError.Error(),
			"requestId": requestId,

			"code":   code,
			"status": status,
		}
		funcLogger.WithFields(fields).Errorf("db error")

		c.AbortWithStatusJSON(httpStatus, CommonErrorResponse{
			Ip:        ginutils.ReadUserIP(c.Request),
			AuthToken: authToken,
			Error:     modelError.Error(),
			CommonResponse: CommonResponse{
				RequestId: requestId,
				Status:    strconv.Itoa(status),
				Code:      code,
			},
		})
		return true
	}

	return false
}
