package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	lru "go-web-app/thirdparty/golang-lru"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	mod "go-web-app/common/model"
	"go-web-app/common/util/gormutils"
	permissionUtils "go-web-app/common/util/permissionutils"
)

var (
	roleMemCache *lru.CacheWithExpiration = InitCacheWithSizeAndExpiration(100, 10*time.Minute)
	userMemCache *lru.CacheWithExpiration = InitCacheWithSize(100)
)

const (
// errTemplateIDMismatch        = errors.New("templateID mismatch")
//
// errUserNotFound        = errors.New("cannot find specified User")
)

// UserRepository ...
type UserRepository struct {
	modelUser       *mod.User
	modelRole       *mod.Role
	databaseSession *gorm.DB
}

// NewUserRepository ...
func NewUserRepository(databaseSession *gorm.DB) UserRepository {
	return UserRepository{
		modelUser:       &mod.User{},
		modelRole:       &mod.Role{},
		databaseSession: databaseSession,
	}
}

// EnsureIndex ...
func (repo UserRepository) EnsureIndex() error {

	repo.databaseSession.AutoMigrate(repo.modelUser, repo.modelRole)

	repo.EnsureRoles()

	return nil
}

func (repo UserRepository) EnsureRoles() error {
	rolesAndPermissions := []struct {
		roleName    string
		permissions []string
	}{
		{
			roleName: permissionUtils.ROLE_DEFAULT_ADMIN,
			permissions: []string{
				permissionUtils.PERMISSION_READ_ADMIN,
				permissionUtils.PERMISSION_WRITE_ADMIN,
				permissionUtils.PERMISSION_READ_USER,
				permissionUtils.PERMISSION_WRITE_USER,
			},
		},
	}
	for _, roleInfo := range rolesAndPermissions {
		role := &mod.Role{
			Name:        roleInfo.roleName,
			Permissions: strings.Join(roleInfo.permissions, ","),
		}
		roles, _ := repo.RetrieveRolesByRoleNames(role.Name)
		if len(roles) > 0 {
			theRole := roles[0]
			theRole.Permissions = role.Permissions
			repo.UpdateRole(theRole.UUID, theRole)
		} else {
			repo.createRole(role)
		}
	}

	return nil
}

func (repo UserRepository) cacheRoles(roles ...*mod.Role) {
	for _, role := range roles {
		if role == nil {
			continue
		}

		roleMemCache.Add(role.Name, role)
	}
}

func (repo UserRepository) generateKeyByUserName(userName string) string {
	return fmt.Sprintf("user_name_%s", userName)
}
func (repo UserRepository) cacheUsers(users ...*mod.User) {
	for _, user := range users {
		if user == nil {
			continue
		}

		userMemCache.Add(user.ModelID(), user)
		userMemCache.Add(repo.generateKeyByUserName(user.UserName), user)
	}
}

func (repo UserRepository) getCachedUsersByIDs(userIDs ...string) (cachedResult []*mod.User, getCachedSuccess bool) {
	expectedAmount := len(userIDs)

	// nothing to get
	if expectedAmount == 0 {
		return []*mod.User{}, true
	}

	isAllCached := true
	allUsersFound := make([]*mod.User, expectedAmount)
	for index, userID := range userIDs {
		cached, ok := userMemCache.Get(userID)

		if (!ok) || cached == nil {
			isAllCached = false
			break
		} else {
			allUsersFound[index] = cached.(*mod.User)
		}
	}
	if isAllCached {
		return allUsersFound, true
	} else {
		return []*mod.User{}, false
	}
}

func (repo UserRepository) getCachedUsersByUserName(userName string) (cachedResult *mod.User, getCachedSuccess bool) {

	cached, ok := userMemCache.Get(repo.generateKeyByUserName(userName))
	if (!ok) || cached == nil {
		return nil, false
	} else {
		return cached.(*mod.User), true
	}
}

func (repo UserRepository) getCachedRolesByRoleNames(roleNames ...string) (cachedResult []*mod.Role, getCachedSuccess bool) {
	expectedAmount := len(roleNames)

	// nothing to get
	if expectedAmount == 0 {
		return []*mod.Role{}, true
	}

	isAllCached := true
	allRolesFound := make([]*mod.Role, expectedAmount)
	for index, roleName := range roleNames {
		cached, ok := roleMemCache.Get(roleName)

		if (!ok) || cached == nil {
			isAllCached = false
			break
		} else {
			allRolesFound[index] = cached.(*mod.Role)
		}
	}
	if isAllCached {
		return allRolesFound, true
	} else {
		return []*mod.Role{}, false
	}
}

func (repo UserRepository) clearUserPassword(users ...*mod.User) {
	for _, user := range users {
		if user == nil {
			continue
		}

		user.Password = ""
	}
}

// Public

// Get ...
func (repo UserRepository) Get(withPassword bool, userIDs ...string) ([]*mod.User, error) {
	var users []*mod.User

	noPassword := !withPassword
	if noPassword {
		cachedResult, getCachedSuccess := repo.getCachedUsersByIDs(userIDs...)
		if getCachedSuccess && repo.checkPasswordStatus(withPassword, cachedResult...) {
			return cachedResult, nil
		}
	}

	// nothing to get
	if len(userIDs) == 0 {
		return []*mod.User{}, nil
	}

	ids := userIDs

	queryResult := repo.databaseSession.Preload(clause.Associations).Find(&users, gormutils.GetQueryWhereIn(repo.databaseSession, UUIDFieldName, ids))

	if noPassword {
		repo.clearUserPassword(users...)
	}
	if queryResult.Error != nil {
		// return nil, errors.WithMessage(err, "db error")
		log.WithFields(log.Fields{
			"err": queryResult.Error,
		}).Debugf("repository.Get => Get error: %v", queryResult.Error)
		return users, queryResult.Error
	}

	if noPassword {
		repo.cacheUsers(users...)
	}

	return users, nil
}

// RetrieveRolesByRoleNames ...
func (repo UserRepository) RetrieveRolesByRoleNames(roleNames ...string) ([]*mod.Role, error) {
	var err error
	var roles []*mod.Role

	cachedResult, getCachedSuccess := repo.getCachedRolesByRoleNames(roleNames...)
	if getCachedSuccess {
		return cachedResult, nil
	}

	// get all Role in RolePool
	roles, err = repo.retrieveRolesByRoleNames(roleNames...)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Debugf("repository.RetrieveRolesByRoleNames => retrieveRolesByRoleNames error: %v", err)
		return nil, err
	}

	return roles, err
}

// RetrieveUserByUserName ...
func (repo UserRepository) RetrieveUserByUserName(withPassword bool, userName string) (*mod.User, error) {
	var user *mod.User

	if !withPassword {
		cachedResult, getCachedSuccess := repo.getCachedUsersByUserName(userName)
		if getCachedSuccess && repo.checkPasswordStatus(withPassword, cachedResult) {
			return cachedResult, nil
		}
	}

	filter := map[string]interface{}{}
	userNameExisting := userName != ""
	if userNameExisting {
		filter[userAttributeUserNameFieldName] = userName
	}
	if !userNameExisting {
		return nil, errors.New("Leak userName")
	}

	queryResult := repo.databaseSession.Preload(clause.Associations).Find(&user, filter)
	if !withPassword {
		repo.clearUserPassword(user)
	}
	if queryResult.Error != nil {
		log.WithFields(log.Fields{
			"err": queryResult.Error,
		}).Debugf("repository.RetrieveUserByUserName => retrieveUserByUserName error: %v", queryResult.Error)
		return nil, queryResult.Error
	}
	if user.UUID == "" {
		return nil, nil
	}

	if !withPassword {
		repo.cacheUsers(user)
	}

	return user, nil
}

// CreateUser ...
func (repo UserRepository) CreateUser(user *mod.User) error {
	var err error

	if user.UUID == "" {
		theUUID, _ := uuid.NewRandom()
		user.UUID = theUUID.String()
	}

	now := time.Now().UTC()
	user.CreatedAt = now
	user.UpdatedAt = now

	err = repo.createUser(user)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Debugf("repository.CreateUser => createUser error: %v", err)
		return err
	}

	return err
}

// UpdateUser ...
func (repo UserRepository) UpdateUser(userID string, user *mod.User) error {
	var err error

	user.UpdatedAt = time.Now().UTC()

	err = repo.updateUser(userID, user)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Debugf("repository.UpdateUser => updateUser error: %v", err)
		return err
	}

	return err
}

// DeleteUser ...
func (repo UserRepository) DeleteUser(userID string) error {
	var err error

	err = repo.deleteUser(userID)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Debugf("repository.DeleteUser => deleteUser error: %v", err)
		return err
	}

	return err
}

// ChangeUserFreezedStatus ...
func (repo UserRepository) ChangeUserFreezedStatus(freezed bool, userIDs ...string) (int64, error) {

	// nothing to patch
	if len(userIDs) == 0 {
		return 0, nil
	}

	return repo.setField(userAttributeFreezedFieldName, freezed, userIDs...)
}

// getAllUsersAndCountWithFilter ...
func (repo UserRepository) getAllUsersAndCountWithFilter(sortField SortField, sortOrder SortOrder, pageLimit, pageIndex int, filter []clause.Expression) ([]*mod.User, int64, error) {
	var users []*mod.User
	var count int64
	var err error

	if pageLimit <= 0 {
		pageLimit = DefaultPageLimit
	} else if pageLimit > DefaultPageLimitMaximum {
		pageLimit = DefaultPageLimitMaximum
	}

	var queryResult *gorm.DB

	targetModel := repo.modelUser
	commonStatement := gormutils.GetModelCommonQuery(repo.databaseSession, targetModel, filter)

	//// Sort
	//filterAggregation = appendFilterSortByTextSearchOrCommon(filter, &filterAggregation)

	//query := repo.userPool.Find(filter).Sort("-" + createdAtFieldName)
	queryResult = gormutils.GetQueryWithPagination(commonStatement, sortField.String(), sortOrder.String(), pageLimit, pageIndex).Preload(clause.Associations).Find(&users)
	if queryResult.Error != nil {
		return nil, 0, queryResult.Error
	}
	// Get Count only
	queryResult = gormutils.GetModelCommonQuery(repo.databaseSession, targetModel, filter).Count(&count)
	if queryResult.Error != nil {
		return nil, 0, queryResult.Error
	}
	if err != nil {
		return nil, 0, err
	}

	return users, count, err
}

// RetrieveUsersByFilters ...
func (repo UserRepository) RetrieveUsersByFilters(userIDs []string, roleName string, keyword string, userStatus *bool, userStatusNotEqual *bool, startTime *time.Time, endTime *time.Time, sortField SortField, sortOrder SortOrder, pageLimit, pageIndex int) ([]*mod.User, int64, error) {
	var users []*mod.User
	var count int64
	var err error

	// get all User in userPool
	var filter []clause.Expression

	//appendKeywordFilter(&filter, keyword, repo)
	if userIDs != nil && len(userIDs) > 0 {
		filter = append(filter, clause.Eq{Column: UUIDFieldName, Value: gormutils.GetQueryWhereIn(repo.databaseSession, UUIDFieldName, userIDs)})
	}
	if roleName != "" {
		filter = append(filter, clause.Eq{Column: userAttributeRoleNamesFieldName, Value: roleName})
	}
	// UserStatus: EQ & NE
	if userStatus != nil {
		filter = append(filter, clause.Eq{Column: userAttributeFreezedFieldName, Value: *userStatus})
	} else if userStatusNotEqual != nil {
		filter = append(filter, clause.Expr{SQL: userAttributeFreezedFieldName + " != ?", Vars: []interface{}{*userStatusNotEqual}})
	}
	// StartTime/EndTime
	if startTime != nil && endTime != nil {
		filter = append(filter, clause.Expr{SQL: CreatedAtFieldName + " BETWEEN ? AND ?", Vars: []interface{}{startTime, endTime}})
	}

	if strings.TrimSpace(keyword) != "" {
		filter = append(filter, clause.Expr{SQL: userAttributeUserNameFieldName + " LIKE ?", Vars: []interface{}{keyword}})
	}

	users, count, err = repo.getAllUsersAndCountWithFilter(sortField, sortOrder, pageLimit, pageIndex, filter)

	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Debugf("repository.RetrieveUsersByFilters => error: %v", err)
		return nil, 0, err
	}
	// else if users == nil {
	// 	log.WithFields(log.Fields{}).Debugf("repository.RetrieveUsersByFilters => users == nil")
	// 	err = ErrNotFound
	// 	return
	// }

	return users, count, nil
}

//     /\ /\         /\ /\         /\ /\         /\ /\         /\ /\
//    / / \ \       / / \ \       / / \ \       / / \ \       / / \ \
//   / /   \ \     / /   \ \     / /   \ \     / /   \ \     / /   \ \
//  / /     \ \   / /     \ \   / /     \ \   / /     \ \   / /     \ \
// / /       \ \ / /       \ \ / /       \ \ / /       \ \ / /       \ \
// \/         \/ \/         \/ \/         \/ \/         \/ \/         \/

// multi userPool Operations

// setField ...
func (repo UserRepository) setField(fieldName string, value interface{}, userIDs ...string) (int64, error) {
	ids := userIDs

	commonStatement := gormutils.GetModelCommonQuery(repo.databaseSession, repo.modelUser, gormutils.GetQueryWhereIn(repo.databaseSession, UUIDFieldName, ids))

	return repo.setFieldRaw(commonStatement, fieldName, value)
}

// setFieldByUserName ...
func (repo UserRepository) setFieldByUserName(fieldName string, value interface{}, userName string) (int64, error) {

	commonStatement := gormutils.GetModelCommonQuery(repo.databaseSession, repo.modelUser, map[string]interface{}{
		userAttributeUserNameFieldName: userName,
	})

	return repo.setFieldRaw(commonStatement, fieldName, value)
}

func (repo UserRepository) setFieldRaw(commonStatement *gorm.DB, fieldName string, value interface{}, userIDs ...string) (int64, error) {
	queryResult := commonStatement.Updates(
		map[string]interface{}{
			fieldName: value,
		},
	)
	if queryResult.Error != nil {
		return 0, queryResult.Error
	}

	for _, userID := range userIDs {
		cached, _ := userMemCache.Get(userID)
		if cached != nil {
			user := cached.(*mod.User)
			userMemCache.Remove(repo.generateKeyByUserName(user.UserName))
		}
		userMemCache.Remove(userID)
	}

	return queryResult.RowsAffected, nil
}

//func (repo UserRepository) retrieveUsersByRoleName(roleName string) ([]*mod.User, error) {
//	filter := map[string]interface{}{
//		"role_names": map[string]interface{}{
//			"$elemMatch": map[string]interface{}{
//				"$eq": roleName,
//			},
//		},
//	}
//
//	err = repo.userPool.Find(filter).Select(map[string]interface{}{
//		"password": 0,
//	}).All(&users)
//
//	if err != nil {
//		return nil, err
//	}
//
//	// repo.cacheUsers(users...)
//
//	return users, nil
//}

func (repo UserRepository) createUser(user *mod.User) error {
	// check nil
	if user == nil {
		return nil
	}

	if user.UUID == "" {
		theUUID, err := uuid.NewRandom()
		user.UUID = theUUID.String()
		if err != nil {
			return err
		}
	}

	//// check model
	//err = modelchecker.CheckUser(user)
	//if err != nil {
	//	return err
	//}

	logTime := time.Now().UTC()
	user.CreatedAt = logTime
	user.UpdatedAt = logTime

	// insert to database
	queryResult := repo.databaseSession.Create(user)
	if queryResult.Error != nil {
		return errors.WithMessage(queryResult.Error, "db error")
	}

	repo.cacheUsers(user)

	return nil
}

func (repo UserRepository) updateUser(userID string, user *mod.User) error {

	// check nil
	if user == nil {
		return nil
	}

	//// check model
	//err = modelchecker.CheckUser(user)
	//if err != nil {
	//	return err
	//}

	logTime := time.Now().UTC()
	user.UpdatedAt = logTime

	// insert to database
	queryResult := repo.databaseSession.Save(user)
	if queryResult.Error != nil {
		return errors.WithMessage(queryResult.Error, "db error")
	}

	repo.cacheUsers(user)

	return nil
}

func (repo UserRepository) deleteUser(userIDs ...string) error {

	// nothing to remove
	if len(userIDs) == 0 {
		return nil
	}

	ids := userIDs

	queryResult := repo.databaseSession.Delete(repo.modelUser, gormutils.GetQueryWhereIn(repo.databaseSession, UUIDFieldName, ids))
	if queryResult.Error != nil {
		return errors.WithMessage(queryResult.Error, "db error")
	}
	//
	return nil

	// userMemCache.Remove(userID)
	userMemCache.Purge()

	return nil
}

//func (repo UserRepository) retrieveRoles() ([]*mod.Role, error) {
//	filter := map[string]interface{}{}
//
//	err = repo.rolePool.Find(filter).Select(map[string]interface{}{
//		mongoDBObjectID: 0,
//	}).All(&roles)
//
//	if err != nil {
//		return nil, err
//	}
//
//	repo.cacheRoles(roles...)
//
//	return roles, nil
//}

func (repo UserRepository) retrieveRolesByRoleNames(roleNames ...string) ([]*mod.Role, error) {
	var roles []*mod.Role
	var err error

	queryResult := repo.databaseSession.Preload(clause.Associations).Find(&roles, gormutils.GetQueryWhereIn(repo.databaseSession, roleNameFieldName, roleNames))
	if queryResult.Error != nil {
		return nil, errors.WithMessage(queryResult.Error, "db error")
	}

	if err != nil {
		return nil, err
	}

	repo.cacheRoles(roles...)

	return roles, nil
}

func (repo UserRepository) createRole(role *mod.Role) error {

	// check nil
	if role == nil {
		return nil
	}

	if role.UUID == "" {
		theUUID, err := uuid.NewRandom()
		role.UUID = theUUID.String()
		if err != nil {
			return err
		}
	}

	//// check model
	//err = modelchecker.CheckRole(role)
	//if err != nil {
	//	return err
	//}

	logTime := time.Now().UTC()
	role.CreatedAt = logTime
	role.UpdatedAt = logTime

	// insert to database
	queryResult := repo.databaseSession.Create(role)
	if queryResult.Error != nil {
		return errors.WithMessage(queryResult.Error, "db error")
	}

	repo.cacheRoles(role)

	return nil
}

// UpdateRole ...
func (repo UserRepository) UpdateRole(roleID string, role *mod.Role) error {

	// check nil
	if role == nil {
		return nil
	}

	if role.UUID == "" {
		theUUID, err := uuid.NewRandom()
		role.UUID = theUUID.String()
		if err != nil {
			return err
		}
	}

	//// check model
	//err = modelchecker.CheckRole(role)
	//if err != nil {
	//	return err
	//}

	logTime := time.Now().UTC()
	role.UpdatedAt = logTime

	// insert to database
	queryResult := repo.databaseSession.Save(role)
	if queryResult.Error != nil {
		return errors.WithMessage(queryResult.Error, "db error")
	}

	repo.cacheRoles(role)

	return nil
}

func (repo UserRepository) deleteRole(roleName string) error {
	queryResult := repo.databaseSession.Delete(repo.modelRole, map[string]interface{}{
		roleNameFieldName: roleName,
	})
	if queryResult.Error != nil {
		return errors.WithMessage(queryResult.Error, "db error")
	}

	roleMemCache.Remove(roleName)

	return nil
}

func (repo UserRepository) checkPasswordStatus(withPassword bool, cachedResult ...*mod.User) bool {
	tempResult := cachedResult
	for _, item := range tempResult {
		if (item.Password == "" && withPassword) ||
			(item.Password != "" && (!withPassword)) {
			cachedResult = nil

			return true
		}
	}

	return false
}
