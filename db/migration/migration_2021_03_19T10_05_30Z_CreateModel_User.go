package migration

import (
	"gorm.io/gorm"

	"github.com/go-gormigrate/gormigrate/v2"
	mod "go-web-app/common/model"
)

func migration_2021_03_19T10_05_30Z_CreateModel_User() {
	migrations = append(migrations,
		&gormigrate.Migration{
			ID: "2021_03_19T10_05_30Z",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&mod.User{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(&mod.User{})
			},
		})
}
