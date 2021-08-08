package servicePermission

import (
	"net/http"
	"sync"
	"time"

	"go-web-app/common/model"
	"go-web-app/errtrace"
	lru "go-web-app/thirdparty/golang-lru"
)

// Constant values
const (
	ErrorPermissionInsufficient = 100100

	ErrorRolesFetchError = 200100

	DefaultCacheSize = 100
)

var errorTitle = map[int]string{
	ErrorPermissionInsufficient: "ErrorPermissionInsufficient",
	ErrorRolesFetchError:        "ErrorRolesFetchError",
}
var errorStatus = map[int]int{
	ErrorPermissionInsufficient: http.StatusBadRequest,
	ErrorRolesFetchError:        http.StatusBadRequest,
}

var (
	memCache *lru.CacheWithExpiration
	mutex    *sync.RWMutex

	DebugMode = false
)

func init() {
	InitCacheWithSize(DefaultCacheSize)
	mutex = &sync.RWMutex{}
}

func InitCacheWithSize(size int) {
	cache, createLRUErr := lru.NewCacheWithExpiration(size, 30*time.Second)
	if createLRUErr != nil || cache == nil {
		panic([]interface{}{
			createLRUErr,
			cache,
		})
	}
	memCache = cache
}

func getPermissionsByUserIDFromMemCache(userID string) *model.CachedPermissionEntry {
	mutex.RLock()
	defer mutex.RUnlock()

	return _getPermissionEntryByUserIDFromMemCache(userID)
}

func updateGlobalPermissionsByUserIDToMemCache(userID string, permissionMap *model.CachedPermission) {
	mutex.Lock()

	entry := _checkOrNewCachedPermissionEntry(userID, _getPermissionEntryByUserIDFromMemCache(userID))
	entry.Global = permissionMap

	_putPermissionEntryByUserIDToMemCache(userID, entry)
	mutex.Unlock()
}

func _checkOrNewCachedPermissionEntry(userID string, entry *model.CachedPermissionEntry) *model.CachedPermissionEntry {
	if entry == nil {
		memCache.Remove(userID)

		entry = &model.CachedPermissionEntry{}
		entry.CreatedTime = time.Now().UTC()
	}

	return entry
}

func _getPermissionEntryByUserIDFromMemCache(userID string) *model.CachedPermissionEntry {
	object, ok := memCache.Get(userID)
	if ok && object != nil {
		return object.(*model.CachedPermissionEntry)
	}

	return nil
}

func _putPermissionEntryByUserIDToMemCache(userID string, entry *model.CachedPermissionEntry) {
	memCache.Add(userID, entry)
}

func CleanPermissionsByUserIDToMemCache(userID string) {
	mutex.Lock()

	memCache.Remove(userID)

	mutex.Unlock()
}

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
