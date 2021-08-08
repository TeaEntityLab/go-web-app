package route

import (
	"fmt"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	// "github.com/pkg/errors"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"go-web-app/common/model"
	repo "go-web-app/common/repository"
	authService "go-web-app/common/service/auth"
	"go-web-app/common/util/httputils"
)

type CommonTokenResponse struct {
	CommonResponse

	Token string `json:"token"`
	Title string `json:"title"`
}

// CheckUsernamePasswordHTTPAPIHandler godoc
// @Summary Check the username & password correctness for login
// @Description Check username & password by json
// @ID check-username-password-by-json
// @Accept  json
// @Produce  json
// @Param loginForm body model.AuthLogin true "Login Form"
// @Success 200 {object} CommonTokenResponse{data=string}
// @Failure 400,401,403 {object} CommonErrorResponse
// @Failure 500 {object} CommonErrorResponse
// @Failure default {object} CommonErrorResponse
// @Router /api/v1/auth/login [post]
func CheckUsernamePasswordHTTPAPIHandler(c *fasthttp.RequestCtx) {
	requestXID := xid.New().String()
	funcLogger := httputils.NewHttpLogger(Logger, c).
		WithField("requestXID", requestXID).
		WithField("func", "dashboard-backend.CheckUsernamePasswordHTTPAPIHandler")

	dbClient := c.UserValue("dbClient").(*gorm.DB)

	userRepo := repo.NewUserRepository(dbClient)

	// appKey := c.Params.ByName("appKey")

	auth := model.AuthLogin{}
	body := c.PostBody()
	//formErr := c.BindJSON(&auth)
	formErr := jsoniter.Unmarshal(body, &auth)
	if formErr != nil || (!authService.CheckAuthInfoValidation(&auth)) {

		response := httputils.ResponseAuthError(
			c,
			authService.ErrorAuthLoginInfoInsufficient,
			"The auth infos from the form is not enough",
			formErr, false)

		funcLogger.WithFields(logrus.Fields{
			"ip":    httputils.ReadUserIP(c),
			"error": response.Error(),
		}).Infof("JSON binding error")

		return
	}

	user, getUserErr := userRepo.RetrieveUserByUserName(true, auth.UserName)
	if getUserErr != nil || user == nil {
		response := httputils.ResponseAuthError(
			c,
			authService.ErrorAuthLoginUsernameNotFound,
			"User not found",
			getUserErr, false)

		funcLogger.WithFields(logrus.Fields{
			"ip":    httputils.ReadUserIP(c),
			"error": response.Error(),
		}).Infof("db error")

		return
	}

	loginErr := authService.CheckLoginUserNamePassword(auth.Password, user.Password)
	if loginErr != nil {
		response := httputils.ResponseAuthError(
			c,
			authService.ErrorAuthLoginPasswordNotMatching,
			"Password is not matching the one in the database",
			loginErr, false)

		funcLogger.WithFields(logrus.Fields{
			"ip":    httputils.ReadUserIP(c),
			"error": response.Error(),
		}).Infof("login error")

		return
	}

	token, tokenErr := authService.GenerateJWTTokenForUser(user)
	if tokenErr != nil {
		response := httputils.ResponseAuthError(
			c,
			authService.ErrorAuthLoginInfoInsufficient,
			fmt.Sprintf("Token generation errors: %v", tokenErr.Error()),
			tokenErr, false)

		funcLogger.WithFields(logrus.Fields{
			"ip":    httputils.ReadUserIP(c),
			"error": response.Error(),
		}).Infof("jwt token error")

		return
	}

	cookie := fasthttp.Cookie{}
	cookie.SetKey(authService.KeyCookieToken)
	cookie.SetValue(token)
	cookie.SetMaxAge(3600)
	cookie.SetPath("/")
	cookie.SetDomain(string(c.Request.URI().Host()))
	cookie.SetSecure(true)
	cookie.SetHTTPOnly(true)
	c.Response.Header.SetCookie(&cookie)
	//cookie.SetExpire(time.Now().Add(3600 * time.Second))
	//cookie.SetSameSite(fasthttp.CookieSameSiteDefaultMode)
	//c.Response.Header.SetCookie(authService.KeyCookieToken, token, 3600, "/", c.Request.URI().Host(), true, true)

	httputils.DoJSONWrite(c, http.StatusOK, CommonTokenResponse{
		Token: token,
		Title: "success",

		CommonResponse: CommonResponse{
			Code:   COMMON_RESPONSE_CODE_SUCCESS,
			Status: COMMON_RESPONSE_STATUS_SUCCESS_STR,
		},
	})
}

// RenewAuthTokenHTTPAPIHandler godoc
// @Summary Renew authToken to avoid expirations
// @Description Renew authToken to avoid expirations by old authToken
// @ID renew-auth-token-by-auth-token
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} CommonTokenResponse{data=string}
// @Failure 400,401,403 {object} CommonErrorResponse
// @Failure 500 {object} CommonErrorResponse
// @Failure default {object} CommonErrorResponse
// @Router /api/v1/auth/renew [post]
func RenewAuthTokenHTTPAPIHandler(c *fasthttp.RequestCtx) {
	requestXID := xid.New().String()
	funcLogger := httputils.NewHttpLogger(Logger, c).
		WithField("requestXID", requestXID).
		WithField("func", "dashboard-backend.RenewAuthTokenHTTPAPIHandler")

	dbClient := c.UserValue("dbClient").(*gorm.DB)

	userRepo := repo.NewUserRepository(dbClient)

	authToken, checkLoginErr := httputils.CheckLoginStatusOrAbort(c, funcLogger, "CheckLoginStatusOrAbort error")
	if checkLoginErr != nil || authToken == nil {
		return
	}

	users, getUserErr := userRepo.Get(false, authToken.UserID)
	var user *model.User
	if len(users) > 0 {
		user = users[0]
	}
	if getUserErr != nil || user == nil {
		response := httputils.ResponseAuthError(
			c,
			authService.ErrorAuthLoginUsernameNotFound,
			"User not found",
			getUserErr, false)

		funcLogger.WithFields(logrus.Fields{
			"ip":        httputils.ReadUserIP(c),
			"authToken": authToken,
			"error":     response.Error(),
		}).Errorf("db error")

		return
	}

	token, tokenErr := authService.GenerateJWTTokenForUser(user)
	if tokenErr != nil {
		response := httputils.ResponseAuthError(
			c,
			authService.ErrorAuthLoginInfoInsufficient,
			fmt.Sprintf("Token generation errors: %v", tokenErr.Error()),
			tokenErr, false)

		funcLogger.WithFields(logrus.Fields{
			"ip":        httputils.ReadUserIP(c),
			"authToken": authToken,
			"error":     response.Error(),
		}).Errorf("jwt token error")

		return
	}

	cookie := fasthttp.Cookie{}
	cookie.SetKey(authService.KeyCookieToken)
	cookie.SetValue(token)
	cookie.SetMaxAge(3600)
	cookie.SetPath("/")
	cookie.SetDomain(string(c.Request.URI().Host()))
	cookie.SetSecure(true)
	cookie.SetHTTPOnly(true)
	c.Response.Header.SetCookie(&cookie)
	//cookie.SetExpire(time.Now().Add(3600 * time.Second))
	//cookie.SetSameSite(fasthttp.CookieSameSiteDefaultMode)
	//c.Response.Header.SetCookie(authService.KeyCookieToken, token, 3600, "/", c.Request.URI().Host, true, true)

	httputils.DoJSONWrite(c, http.StatusOK, CommonTokenResponse{
		Token: token,
		Title: "success",

		CommonResponse: CommonResponse{
			Code:   COMMON_RESPONSE_CODE_SUCCESS,
			Status: COMMON_RESPONSE_STATUS_SUCCESS_STR,
		},
	})
}
