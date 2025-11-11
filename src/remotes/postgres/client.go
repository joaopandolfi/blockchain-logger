package postgres

import (
	"fmt"
	"logger/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var driver *gorm.DB

// Init postgres database
func Init(conf config.Config) error {
	/*
		postgres.Config{
			DSN: "gorm:gorm@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True&loc=Local", // data source name
			DefaultStringSize: 256, // default size for string fields
			DisableDatetimePrecision: true, // disable datetime precision, which not supported before postgres 5.6
			DontSupportRenameIndex: true, // drop & create when rename index, rename index not supported before postgres 5.7, MariaDB
			DontSupportRenameColumn: true, // `change` when rename column, rename column not supported before postgres 8, MariaDB
			SkipInitializeWithVersion: false, // auto configure based on currently postgres version
		}
	*/

	db, err := gorm.Open(postgres.Open(conf.PostgreSQL), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("connecting on postgres database: %w", err)
	}
	driver = db
	return nil
}

// Close pool connection
func Close() {
	// This driver close automatically
}

// Driver Return postgres driver
func Driver() *gorm.DB {
	if driver == nil {
		Init(config.Get())
	}
	return driver
}
