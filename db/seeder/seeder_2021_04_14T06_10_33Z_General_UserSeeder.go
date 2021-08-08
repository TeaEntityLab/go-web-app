package seeder

import (
	// "github.com/pkg/errors"

	"go-web-app/common/model"
	"go-web-app/common/repository"
	serviceAuth "go-web-app/common/service/auth"
	permissionUtils "go-web-app/common/util/permissionutils"
	"go-web-app/db"
)

const (
	USER_NAME_ADMIN = "admin"
)

type UserSeeder struct{}

func (seeder UserSeeder) Run() error {

	databaseSession := db.GetDefaultDatabase()
	userRepository := repository.NewUserRepository(databaseSession)
	userRepository.EnsureIndex()

	userNameList := []struct {
		UserName string
	}{
		{USER_NAME_ADMIN},
	}
	for _, item := range userNameList {
		userName := item.UserName

		adminUser, err := userRepository.RetrieveUserByUserName(true, userName)
		if err != nil {
			return err
		}

		var defaultPassword string
		defaultPassword, err = serviceAuth.GenerateLoginHashedPassword("admin")

		newUser := &model.User{
			UserName: userName,
			Password: defaultPassword,

			RoleNames: permissionUtils.ROLE_DEFAULT_ADMIN,
		}
		if adminUser == nil {
			if err != nil {
				return err
			}

			err = userRepository.CreateUser(newUser)
			if err != nil {
				return err
			}
		} else {
			//adminUser.Password = defaultPassword

			adminUser.RoleNames = newUser.RoleNames
			err = userRepository.UpdateUser(adminUser.UUID, adminUser)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (seeder UserSeeder) SeederName() string {
	return "UserSeeder"
}

func seeder_2021_04_14T06_10_33Z_General_UserSeeder() {
	seeder := UserSeeder{}
	seederMap["UserSeeder"] = seeder
	seeders = append(seeders, seeder)
}
