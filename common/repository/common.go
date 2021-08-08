package repository

import (
	"context"
	"sync"
	"time"
	// "go-web-app/common/model"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	lru "go-web-app/thirdparty/golang-lru"
	"gopkg.in/redis.v5"
	"gorm.io/gorm"

	"go-web-app/common/model"
	"go-web-app/common/model/modelchecker"
	"go-web-app/db"
)

type Repository interface {
	EnsureIndex() (err error)
}
type TextSearchableRepository interface {
	TextSearchIndexes() []string
}

type SortOrder int32

const (
	SortOrder_DESC SortOrder = 0
	SortOrder_ASC  SortOrder = 1
)

var SortOrder_name = map[int32]string{
	0: "DESC",
	1: "ASC",
}
var SortOrder_value = map[string]int32{
	"DESC": 0,
	"ASC":  1,
}

func (x SortOrder) Enum() *SortOrder {
	p := new(SortOrder)
	*p = x
	return p
}
func (x SortOrder) String() string {
	return model.EnumName(SortOrder_name, int32(x))
}
func (x *SortOrder) UnmarshalJSON(data []byte) error {
	value, err := model.UnmarshalJSONEnum(SortOrder_value, data, "SortOrder")
	if err != nil {
		return err
	}
	*x = SortOrder(value)
	return nil
}

// func (SortOrder) EnumDescriptor() ([]byte, []int) {
// 	return fileDescriptor2, []int{0}
// }

func NewSortOrderByString(status string) SortOrder {
	val, ok := SortOrder_value[status]
	if !ok {
		return SortOrder_DESC
	}

	return SortOrder(val)
}

type SortField int32

const (
	SortField_created_at SortField = 0
	SortField_updated_at SortField = 1
)

var SortField_name = map[int32]string{
	0: CreatedAtFieldName,
	1: UpdatedAtFieldName,
}
var SortField_value = map[string]int32{
	CreatedAtFieldName: 0,
	UpdatedAtFieldName: 1,
}

func (x SortField) Enum() *SortField {
	p := new(SortField)
	*p = x
	return p
}
func (x SortField) String() string {
	return model.EnumName(SortField_name, int32(x))
}
func (x *SortField) UnmarshalJSON(data []byte) error {
	value, err := model.UnmarshalJSONEnum(SortField_value, data, "SortField")
	if err != nil {
		return err
	}
	*x = SortField(value)
	return nil
}

// func (SortField) EnumDescriptor() ([]byte, []int) {
// 	return fileDescriptor2, []int{0}
// }

func NewSortFieldByString(status string) SortField {
	val, ok := SortField_value[status]
	if !ok {
		return SortField_created_at
	}

	return SortField(val)
}

var (
	ConnectionDuration time.Duration

	defaultRedisClient *redis.Cmdable
	defaultRedisLock   sync.Mutex

	defaultDBLock sync.Mutex

	// ErrNotFound ...
	//ErrNotFound = mgo.ErrNotFound
	// ErrBadRequest ...
	ErrBadRequest = errors.New("bad request")

	DefaultPageLimit        = 20
	DefaultPageLimitMaximum = 1000
)

const (
	InvalidPageIndex = -1

	//MongoDBObjectID = "_id"
	UUIDFieldName = "uuid"

	userAttributeRoleNamesFieldName = "role_names"
	descriptionFieldName            = "description"
	CreatedAtFieldName              = "created_at"
	CreatedByFieldName              = "created_by"
	UpdatedAtFieldName              = "updated_at"
	UpdatedByFieldName              = "updated_by"

	/** Pool Names **/

	/** Field Names **/

	// User
	userAttributeUserNameFieldName  = "user_name"
	userAttributePasswordFieldName  = "password"
	userAttributeEmailFieldName     = "email"
	userAttributeFreezedFieldName   = "freezed"
	userAttributeUserExtraFieldName = "user_extra"

	// Role
	roleNameFieldName        = "name"
	rolePermissionsFieldName = "permissions"
)

// GetRedisClientDefault ...
func GetRedisClientDefault(redisEndpoints string) *redis.Cmdable {
	defaultRedisLock.Lock()
	defer defaultRedisLock.Unlock()

	if defaultRedisClient == nil {
		if redisEndpoints == "" {
			log.WithField("redisEndpoints", redisEndpoints).Warn("Failed to parse REDIS_ENDPOINTS config")
			// redisEndpoints = redisEndpointsDefault
			// log.WithField("redisEndpoints", redisEndpoints).Warn("Use default REDIS_ENDPOINTS config")
		}

		var c redis.Cmdable
		var err error
		c, err = db.NewRedisClient(redisEndpoints, 100, true)
		if err != nil {
			log.WithField("Error", err).Fatalf("Failed to connect Redis")
		}
		defaultRedisClient = &c
	} else {
		//
	}

	return defaultRedisClient
}

func GetContextWithTimeoutFromConnectionDurationCommon() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), ConnectionDuration)
}

func NewBackgroundContextWithTimeout(timeout time.Duration) (ctx context.Context) {
	ctx, _ = context.WithTimeout(context.Background(), timeout)
	return ctx
}

func NewBackgroundContextWithDefaultTimeout() (ctx context.Context) {
	return NewBackgroundContextWithTimeout(ConnectionDuration)
}

func Init(databaseSession *gorm.DB) {

	repos := []Repository{
		NewUserRepository(databaseSession),
	}
	for _, repo := range repos {
		err := repo.EnsureIndex()

		if err != nil {
			log.WithError(err).Errorln("Error!")
		}
	}
}

func InitCacheWithSize(size int) *lru.CacheWithExpiration {
	return InitCacheWithSizeAndExpiration(size, 30*time.Second)
}

func InitCacheWithSizeAndExpiration(size int, expiration time.Duration) *lru.CacheWithExpiration {
	cache, createLRUErr := lru.NewCacheWithExpiration(size, expiration)
	if createLRUErr != nil || cache == nil {
		panic([]interface{}{
			createLRUErr,
			cache,
		})
	}

	return cache
}

func SetDebug(isOn bool) {
	//var logger MongoLog
	//if isOn {
	//	logger = MongoLog{}
	//}
	//
	//mgo.SetDebug(isOn)
	//mgo.SetLogger(logger)
}

//// IsNotFound ...
//func IsNotFound(err error) bool {
//	return err == ErrNotFound
//}

// IsBadRequest ...
func IsBadRequest(err error) bool {
	return errors.Cause(err) == ErrBadRequest
}

// IsTypeCheckInvalidFields ...
func IsTypeCheckInvalidFields(err error) bool {
	return errors.Cause(err) == modelchecker.ErrTypeCheckInvalidFields
}

// IsTypeCheckNonExistentObject ...
func IsTypeCheckNonExistentObject(err error) bool {
	return errors.Cause(err) == modelchecker.ErrTypeCheckNonExistentObject
}
