package migration

import (
	"gorm.io/gorm"
)

// CreateTables creates database tables for the given models using GORM's AutoMigrate function.
func CreateTables(db *gorm.DB, models ...interface{}) error {
	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return err
		}
	}
	return nil
}

// DropTables drops database tables for the given models using GORM's DropTable function.
func DropTables(db *gorm.DB, models ...interface{}) error {
	for _, model := range models {
		if err := db.Migrator().DropTable(model); err != nil {
			return err
		}
	}
	return nil
}
