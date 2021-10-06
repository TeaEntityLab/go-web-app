package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	// "strconv"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/fasthttp/router"
	"github.com/golang-jwt/jwt"
	gormLogrus "github.com/onrik/gorm-logrus"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	gormlogger "gorm.io/gorm/logger"

	"go-web-app/common/middleware"
	"go-web-app/common/repository"
	routeCommon "go-web-app/common/route"
	authService "go-web-app/common/service/auth"
	serviceLog "go-web-app/common/service/log"
	"go-web-app/common/util/httputils"
	// servicePermission "go-web-app/common/service/permission"
	_ "go-web-app/app/docs"
	"go-web-app/app/route"
	"go-web-app/db"
	_ "go-web-app/lib/logrus_env"
)

const (
	defaultServicePort = "8080"
	defaultDBMode      = "release"
	// defaultRedisEndpoints         = "localhost:6379"
)

var (
	isDebugMode = false

	ServiceName    = "go-web-app"
	ServiceVersion = "PLACEHOLDER"

	cfg = new(Config)
)

type Config struct {
	db.CommonConfig

	RedisEndpoints string `json:"REDIS_ENDPOINTS" env:"REDIS_ENDPOINTS,required"`
	HostName       string `json:"HOST_NAME" env:"HOST_NAME" envDefault:"0.0.0.0"`
	ServicePort    string `json:"SERVICE_PORT" env:"SERVICE_PORT" envDefault:"8080"`

	JwtTokenPrivateKeyPath string `json:"JWT_TOKEN_PRIVATEKEY_PATH" env:"JWT_TOKEN_PRIVATEKEY_PATH,required"`
	JwtTokenPublicKeyPath  string `json:"JWT_TOKEN_PUBLICKEY_PATH" env:"JWT_TOKEN_PUBLICKEY_PATH,required"`

	SlackHookURL string `json:"SLACK_HOOK_URL" env:"SLACK_HOOK_URL" envDefault:""`
	SlackChannel string `json:"SLACK_CHANNEL" env:"SLACK_CHANNEL" envDefault:""`

	DBMode string `json:"DB_MODE" env:"DB_MODE" envDefault:"info"`

	EnableDocument     string `json:"ENABLE_DOCUMENT" env:"ENABLE_DOCUMENT" envDefault:"false"`
	DocumentRouteGroup string `json:"DOCUMENT_ROUTE_GROUP" env:"DOCUMENT_ROUTE_GROUP" envDefault:""`
}

func main() {
	funcLogger := logrus.WithField("func", "main")
	funcLogger.WithField("ServiceName", ServiceName).WithField("ServiceVersion", ServiceVersion).Debugf("service start")

	if !parseEnv(funcLogger) {
		return
	}
	if dbErr := db.InitDefaultDatabase(cfg.DBType, cfg.DBEndpoints); dbErr != nil {
		funcLogger.WithError(dbErr).Fatalf("db.InitDefaultDatabase() error")
	}
	if cfg.SlackHookURL != "" && cfg.SlackChannel != "" {
		logrus.AddHook(&serviceLog.SlackrusHook{
			HookURL:        cfg.SlackHookURL,
			AcceptedLevels: serviceLog.LevelThreshold(logrus.WarnLevel),
			Channel:        cfg.SlackChannel,
			IconEmoji:      ":ghost:",
			Username:       "Oops",
		})
	}

	route.Logger = logrus.WithField("object", "route")
	r := router.New()

	//r.Use(static.ServeRoot("/static", "./web/static"))
	utilsMiddleware := func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(c *fasthttp.RequestCtx) {
			c.SetUserValue("dbClient", db.GetDefaultDatabase())
			c.SetUserValue("cacheStore", repository.NewCacheStoreWithRedisClient(repository.GetRedisClientDefault(cfg.RedisEndpoints)))
			next(c)
		}
	}
	commonMiddlewares := func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return middleware.Pipe(next, utilsMiddleware)
	}

	// GeneralAPI router group

	//api := r.Group("/api")

	r.POST("/operation/ensureIndex", func(context *fasthttp.RequestCtx) {
		repository.Init(db.GetDefaultDatabase())
	})

	// Docker/Kubernetes health check
	r.GET("/healthz", func(c *fasthttp.RequestCtx) {
		httputils.DoJSONWrite(c, http.StatusOK, http.StatusText(http.StatusOK))
	})

	// Default Kubernetes L7 Loadbalancing health check
	r.GET("/", func(c *fasthttp.RequestCtx) {
		httputils.DoJSONWrite(c, http.StatusOK, http.StatusText(http.StatusOK))
	})

	apiV1 := r.Group("/api/v1")
	apiV1.POST("/auth/login", commonMiddlewares(route.CheckUsernamePasswordHTTPAPIHandler))

	authMiddleware := func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return middleware.Auth(
			route.Logger.WithField("middleware", "auth"),
			//middleware.ExceptionalRoute{
			//	Path:   "/auth/login",
			//	Method: "POST",
			//}
			next,
		)
	}

	apiV1.POST("/auth/renew",
		middleware.Pipe(
			route.RenewAuthTokenHTTPAPIHandler,
			commonMiddlewares, authMiddleware,
		),
	)

	if cfg.EnableDocument == "true" {
		//r.GET(cfg.DocumentRouteGroup+"/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		r.GET(cfg.DocumentRouteGroup+"/swagger/{filepath:*}", fasthttpadaptor.NewFastHTTPHandlerFunc(httpSwagger.WrapHandler))
	}

	//r.GET("/proxy/*url", func(c *fasthttp.RequestCtx) {
	//	target := strings.TrimPrefix(c.Params.ByName("url"), "/")
	//	httputils.ReverseProxy(target)(c)
	//})

	// run the server

	//serverErr := r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	serviceErr := fasthttp.ListenAndServe(
		fmt.Sprintf(
			"%s:%s", cfg.HostName, cfg.ServicePort,
		),
		routeCommon.SimpleCORSMiddleware(r.Handler),
	)

	if serviceErr != nil {
		funcLogger.WithError(serviceErr).Fatalf("fasthttp.ListenAndServe() error")
	}
}

func parseEnv(funcLogger *logrus.Entry) bool {
	var envParseErr error
	envParseErr = env.Parse(cfg)
	if envParseErr != nil {
		funcLogger.WithError(envParseErr).Fatalf("env.Parse error")
		return false
	}

	var errLoadKey error
	var errParseKey error
	var keyBytes []byte
	keyBytes, errLoadKey = ioutil.ReadFile(cfg.JwtTokenPrivateKeyPath)
	authService.EnvJWTTokenPrivateKey, errParseKey = jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
	if errLoadKey != nil || errParseKey != nil {
		funcLogger.Fatalf("Failed to get EnvJWTTokenPrivateKey at %v. Abort.\nerrLoadKey: %v \nerrParseKey: %v\n", cfg.JwtTokenPrivateKeyPath, errLoadKey, errParseKey)
	}
	keyBytes, errLoadKey = ioutil.ReadFile(cfg.JwtTokenPublicKeyPath)
	authService.EnvJWTTokenPublicKey, errParseKey = jwt.ParseRSAPublicKeyFromPEM(keyBytes)
	if errLoadKey != nil || errParseKey != nil {
		funcLogger.Infoln(string(keyBytes))
		funcLogger.Fatalf("Failed to get EnvJWTTokenPublicKey at %v. Abort.\nerrLoadKey: %v \nerrParseKey: %v\n", cfg.JwtTokenPublicKeyPath, errLoadKey, errParseKey)
	}

	// DB Mode
	switch cfg.DBMode {
	case "silent":
		logger := gormLogrus.New()
		db.Logger = logger
		logger.LogMode(gormlogger.Silent)
	case "error":
		logger := gormLogrus.New()
		db.Logger = logger
		logger.LogMode(gormlogger.Error)
	case "warn":
		logger := gormLogrus.New()
		db.Logger = logger
		logger.LogMode(gormlogger.Warn)
	case "info":
		logger := gormLogrus.New()
		db.Logger = logger
		logger.LogMode(gormlogger.Info)
	case "debug":
		db.SetLogger(gormlogger.New(
			log.New(os.Stdout, "[DB] \r\n", log.LstdFlags), // io writer
			gormlogger.Config{
				SlowThreshold: time.Second,       // Slow SQL threshold
				LogLevel:      gormlogger.Silent, // Log level
				Colorful:      false,             // Disable color
			},
		))
	default:
		funcLogger.Infof("Failed to get env GIN_MODE. Use default: `%s`", defaultDBMode)
	}
	return true
}
