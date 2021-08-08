package servicePermission

import (
	"strings"

	mod "go-web-app/common/model"
	repo "go-web-app/common/repository"
	"go-web-app/errtrace"
)

type UserPermissionRequestInfo struct {
	AuthToken *mod.AuthToken

	UserID string

	UserRepo *repo.UserRepository
}

type UserPermissionGrantedInfo struct {
	Info *UserPermissionRequestInfo

	Global *mod.CachedPermission
}

// CheckUserPermissionsByRequestInfo ...
func CheckUserPermissionsByRequestInfo(requestInfo *UserPermissionRequestInfo, permissions ...string) (grantedInfo *UserPermissionGrantedInfo, err error) {
	grantedInfo = &UserPermissionGrantedInfo{
		Info: requestInfo,
	}

	var permissionMap *mod.CachedPermission

	var userID string
	if requestInfo.AuthToken == nil {
		return nil, GenerateError(ErrorPermissionInsufficient, "AuthToken is empty", nil)
	}
	userID = requestInfo.AuthToken.UserID

	// Global
	if requestInfo.UserRepo != nil {
		permissionMap, err = GetGlobalPermissionsByUserID(requestInfo.UserRepo, userID)
	}
	if err != nil {
		return grantedInfo, err
	}
	grantedInfo.Global = permissionMap

	//grantedInfoStr, _ := jsoniter.Marshal(grantedInfo)
	//fmt.Println("grantedInfoStr: ", string(grantedInfoStr))

	for _, permission := range permissions {
		anyMatched := false
		if grantedInfo.Global != nil {
			matched := (*grantedInfo.Global)[permission]
			anyMatched = anyMatched || matched
		}

		if !anyMatched {
			return grantedInfo, GenerateError(ErrorPermissionInsufficient, "No Permission/Forbidden", nil)
		}
	}

	return grantedInfo, err
}

/**
Block: Global
*/

// GetGlobalPermissionsByAuthToken ...
func GetGlobalPermissionsByAuthToken(userRepo *repo.UserRepository, authToken *mod.AuthToken) (cachedPermission *mod.CachedPermission, err error) {
	return GetGlobalPermissionsByUserID(userRepo, authToken.UserID)
}

// GetGlobalPermissionsByUserID ...
func GetGlobalPermissionsByUserID(userRepo *repo.UserRepository, userID string) (cachedPermission *mod.CachedPermission, err error) {
	// NOTE Check CachedObject
	cachedPermissions := getPermissionsByUserIDFromMemCache(userID)
	if cachedPermissions != nil && cachedPermissions.Global != nil {
		return cachedPermissions.Global, nil
	}

	users, retrieveUserByUserIDErr := userRepo.Get(false, userID)
	if retrieveUserByUserIDErr != nil || users == nil || len(users) <= 0 {
		// fields := log.Fields{
		// 	"error": retrieveUserByUserIDErr,
		//
		// 	"code":   500,
		// 	"status": 500,
		// }
		// log.WithFields(fields).Errorf("string(c.Request.URI().Path()): %v", string(c.Request.URI().Path()))
		// httputils.DoJSONWrite(c, http.StatusInternalServerError, fields)
		// return
		return nil, retrieveUserByUserIDErr
	}
	return GetGlobalPermissionsByUser(userRepo, users[0])
}

// GetGlobalPermissionsByUser ...
func GetGlobalPermissionsByUser(userRepo *repo.UserRepository, user *mod.User) (cachedPermission *mod.CachedPermission, err error) {
	var permissionMap *mod.CachedPermission
	permissionMap, err = GetPermissionsByRoleNames(userRepo, strings.Split(user.RoleNames, ","))
	if err != nil {
		return permissionMap, err
	}
	if user.Freezed {
		permissionMap = &mod.CachedPermission{}
	}
	if permissionMap != nil {
		updateGlobalPermissionsByUserIDToMemCache(user.ModelID(), permissionMap)
	}
	return permissionMap, err
}

// Block: Common

// CheckPermissionsByRoleNames ...
func CheckPermissionsByRoleNames(userRepo *repo.UserRepository, roleNames []string, permissions ...string) (cachedPermission *mod.CachedPermission, err error) {

	var permissionMap *mod.CachedPermission
	permissionMap, err = GetPermissionsByRoleNames(userRepo, roleNames)
	if err != nil {
		return permissionMap, err
	}
	return CheckPermissions(permissionMap, permissions...)
}

// GetPermissionsByRoleNames ...
func GetPermissionsByRoleNames(userRepo *repo.UserRepository, roleNames []string) (cachedPermission *mod.CachedPermission, err error) {
	//fmt.Println("roleNames: ", roleNames)
	roles, rolesError := userRepo.RetrieveRolesByRoleNames(roleNames...)
	if rolesError != nil {
		return nil, &errtrace.Error{
			ErrCode: ErrorRolesFetchError,
			ErrRef:  rolesError,
		}
	}

	var permissionMap *mod.CachedPermission
	permissionMap = GetPermissionsByRoles(roles...)
	return permissionMap, err
}

// CheckPermissionsByRoles ...
func CheckPermissionsByRoles(roles []*mod.Role, permissions ...string) (permissionMap *mod.CachedPermission, err error) {
	permissionMap = GetPermissionsByRoles(roles...)

	return CheckPermissions(permissionMap, permissions...)
}

// CheckPermissions ...
func CheckPermissions(permissionMap *mod.CachedPermission, permissions ...string) (cachedPermission *mod.CachedPermission, err error) {
	if permissionMap == nil {
		return nil, GenerateError(ErrorPermissionInsufficient, "", nil)
	}

	var isAnyDismatch = false
	for _, permission := range permissions {
		if (*permissionMap)[permission] != true {
			isAnyDismatch = true
			break
		}
	}

	if isAnyDismatch {
		return nil, GenerateError(ErrorPermissionInsufficient, "", nil)
	}

	return permissionMap, nil
}

// GetPermissionsByRoles ...
func GetPermissionsByRoles(roles ...*mod.Role) (permissionMap *mod.CachedPermission) {
	permissions := mod.CachedPermission{}

	for _, role := range roles {
		if role == nil {
			continue
		}

		for _, permission := range strings.Split((*role).Permissions, ",") {
			permissions[permission] = true
		}
	}

	return &permissions
}
