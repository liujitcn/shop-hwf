package version

import (
	"gitee.com/liujit/shop/server/cmd/migrate/migration"
	"gitee.com/liujit/shop/server/cmd/migrate/models"
	"gorm.io/gorm"
	"runtime"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _initialTables)
}

func _initialTables(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		tb := make([]interface{}, 0)
		tb = append(tb, models.TableCreate()...)
		err := tx.Debug().Migrator().AutoMigrate(
			tb...,
		)
		if err != nil {
			return err
		}
		err = models.InitDb(tx)
		if err != nil {
			return err
		}

		return tx.Create(&models.Migration{
			Version: version,
		}).Error
	})
}
