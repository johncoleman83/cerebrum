// Package datastore creates an ORM connection to the
// databse and persistence layer
package datastore

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // for use with gorm

	"github.com/johncoleman83/cerebrum/pkg/utl/config"
)

// FormatDSN creates the datastore name string for database connections
func FormatDSN(dbConfig *config.Database) string {
	return fmt.Sprintf(
		"%s:%s@%s(%s:%s)/%s?%s",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Protocol,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
		dbConfig.Settings,
	)
}

// NewMySQLGormDb creates new database connection to a mysql database
func NewMySQLGormDb(dbConfig *config.Database) (*gorm.DB, error) {
	dsn := FormatDSN(dbConfig)

	db, err := gorm.Open(dbConfig.Dialect, dsn)
	if err != nil {
		return db, err
	}

	db.LogMode(true)
	if err = db.Exec("SELECT 1").Error; err != nil {
		return db, err
	}

	return db, nil
}
