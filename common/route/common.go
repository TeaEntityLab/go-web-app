package route

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
	// log "github.com/sirupsen/logrus"
	lru "go-web-app/thirdparty/golang-lru"
	// "go-web-app/common/util/httputils"
)

var (
	organizationInfoForPermissionCache = initCacheWithSizeAndExpiration(1000, 5*time.Minute)
)

type CachedOrganizationInfoForPermission struct {
	TenantID string
}

func generateTenantKeyForCache(TenantID string) string {
	return fmt.Sprintf("Tenant_%s", TenantID)
}

func getCachedOrganizationInfoForPermissionOrUpdateIt(key string, getInfo func() *CachedOrganizationInfoForPermission) *CachedOrganizationInfoForPermission {
	value := getCachedOrganizationInfoForPermission(key)
	if value == nil {
		evalResult := getInfo()
		cacheOrganizationInfoForPermission(key, evalResult)
		return evalResult
	}

	return value
}
func getCachedOrganizationInfoForPermission(key string) *CachedOrganizationInfoForPermission {
	value, ok := organizationInfoForPermissionCache.Get(key)
	if (!ok) || value == nil {
		return nil
	}

	return value.(*CachedOrganizationInfoForPermission)
}

func cacheOrganizationInfoForPermission(key string, info *CachedOrganizationInfoForPermission) bool {
	return organizationInfoForPermissionCache.Add(key, info)
}

func initCacheWithSizeAndExpiration(size int, expiration time.Duration) *lru.CacheWithExpiration {
	cache, createLRUErr := lru.NewCacheWithExpiration(size, expiration)
	if createLRUErr != nil || cache == nil {
		panic([]interface{}{
			createLRUErr,
			cache,
		})
	}

	return cache
}

var (
	AccessControlAllowMethods = strings.Join([]string{
		http.MethodOptions,
		http.MethodGet,
		//http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		//http.MethodDelete,
		//http.MethodTrace,
		//http.MethodConnect,
		http.MethodPatch,
	}, ", ")

	AccessControlAllowHeaders = strings.Join([]string{
		"Origin",
		"Content-Type",
		"Content-Length",
		"Accept-Encoding",
		"X-CSRF-Token",
		"Authorization",
	}, ", ")
)

func SimpleCORSMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(c *fasthttp.RequestCtx) {
		// If a request may contain a `Access-Control-Allow-Origin` with different values, then the host should always respond with `Vary: Origin`,
		// even for responses without an `Access-Control-Allow-Origin` header.
		// If the header isn't always present, it would be possible to fill the cache with incorrect values.
		c.Response.Header.Add("Vary", "Origin")

		// cors setting
		origin := string(c.Request.Header.Peek("Origin"))
		if origin == "" {
			origin = "*"
		}
		c.Response.Header.Set("Access-Control-Allow-Origin", origin)

		c.Response.Header.Set("Access-Control-Max-Age", "86400")
		c.Response.Header.Set("Access-Control-Allow-Methods", AccessControlAllowMethods)
		c.Response.Header.Set("Access-Control-Allow-Headers", AccessControlAllowHeaders)
		c.Response.Header.Set("Access-Control-Expose-Headers", "Content-Length")
		c.Response.Header.Set("Access-Control-Allow-Credentials", "true")

		if string(c.Request.Header.Method()) == http.MethodOptions {
			//c.AbortWithStatus(200)
			c.SetStatusCode(200)
		} else {
			next(c)
		}
	}
}
