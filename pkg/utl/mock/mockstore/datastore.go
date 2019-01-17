package mockstore

import (
	"errors"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // for use with gorm

	"github.com/johncoleman83/cerebrum/pkg/utl/config"
	"github.com/johncoleman83/cerebrum/pkg/utl/datastore"
	cerebrum "github.com/johncoleman83/cerebrum/pkg/utl/model"
	"github.com/johncoleman83/cerebrum/pkg/utl/support"
)

// NewDBConn instantiates new mysql database connection via docker container
func NewDBConn(t *testing.T, cfg *config.Configuration, models ...interface{}) *gorm.DB {
	db, err := datastore.NewMySQLGormDb(cfg.DB)
	if err != nil {
		t.Fatal(err)
	}
	for _, model := range models {
		if err := db.CreateTable(model).Error; err != nil {
			t.Fatal(err)
		}
	}
	return db
}

// NewDataBaseConnection creates and returns a new GORM connection to the test DB
func NewDataBaseConnection() (*gorm.DB, error) {
	cfgPath := support.TestingConfigPath()
	cfg, err := config.LoadConfigFrom(cfgPath)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, errors.New("unknown error loading testing yaml file")
	}
	db, err := datastore.NewMySQLGormDb(cfg.DB)
	if err != nil {
		return nil, err
	}
	if err := DropAndCreateAllTablesFor(db); err != nil {
		return nil, err
	}
	return db, nil
}

// DropAndCreateAllTablesFor drops all tables in input db and recreates the ones listed in the function
func DropAndCreateAllTablesFor(db *gorm.DB) error {
	models := []interface{}{
		&cerebrum.Company{},
		&cerebrum.Location{},
		&cerebrum.Role{},
		&cerebrum.User{},
	}
	for _, model := range models {
		if db.HasTable(model) {
			if err := db.DropTable(model).Error; err != nil {
				return err
			}
		}
		if err := db.CreateTable(model).Error; err != nil {
			return err
		}
	}
	return nil
}

// InsertRowsFor inserts multiple values into database
func InsertRowsFor(db *gorm.DB, models ...interface{}) error {
	for _, v := range models {
		if err := db.Create(v).Error; err != nil {
			return err
		}
	}
	return nil
}
