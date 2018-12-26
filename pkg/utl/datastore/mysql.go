package datastore

import (
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/johncoleman83/cerebrum/pkg/utl/config"
)

type Result struct {
	Date  time.Time
	Total int64
}

func FormatDSN(dbConfig *config.Database) (string) {
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

func ExtractHostAndPortFrom(addr string) (string, string) {
	addrSplit := strings.Split(addr, ":")
	return addrSplit[0], addrSplit[1]
}

// New creates new database connection to a mysql database
func NewMySQLGormDb(dbConfig *config.Database) (*gorm.DB, error) {
	dsn := FormatDSN(dbConfig)

	fmt.Println(dsn)
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
