package mockdb

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/johncoleman83/cerebrum/pkg/utl/config"
	"github.com/johncoleman83/cerebrum/pkg/utl/datastore"
	"github.com/johncoleman83/cerebrum/pkg/utl/mock/docker"
	"github.com/johncoleman83/cerebrum/pkg/utl/support"
)

var (
	testConfigPath = support.TestingConfigPath()
)

func WaitFunc(addr string) error {
	cfg, _ := config.LoadConfigFrom(testConfigPath)
	cfg.DB.Host, cfg.DB.Port = datastore.ExtractHostAndPortFrom(addr)
	cfg.DB.User, cfg.DB.Password = "root", ""

	dsn := datastore.FormatDSN(cfg.DB)
	time.Sleep(time.Second * 3)
	db, err := sql.Open(cfg.DB.Dialect, dsn)
	if err == nil {
		fmt.Println(dsn)
		fmt.Println(db)
		fmt.Println("Going to query DB")
		rows, err := db.Query("SELECT * FROM cerebrum_mysql_test_db.users")
		fmt.Println(err)
		fmt.Println(rows.Columns())
		db.Close()
		time.Sleep(time.Second * 3)
	}
	return err
}

func BuildDockerArgs(DB *config.Database) []string {
	randomString := support.NewRandomString(25)
	containerName := fmt.Sprintf("cerebrum_mysql_test_no_%s", randomString)
	return []string{
		"-d", "--name", containerName,
		"-e", "MYSQL_ALLOW_EMPTY_PASSWORD=yes",
		"-e", "MYSQL_DATABASE=" + DB.Name,
		"-e", "MYSQL_USER=" + DB.User,
		"-e", "MYSQL_PASSWORD=" + DB.Password,
	}
}

// docker run --name cerebrum_mysql_dev --detach --env MYSQL_ALLOW_EMPTY_PASSWORD='yes' --env MYSQL_DATABASE='cerebrum_dev' --publish 3306:3306 mysql:latest
// MySQLTestContainer instantiates new PostgreSQL docker container
func MySqlTestContainerConfig(t *testing.T) (*docker.Container, *config.Configuration) {
	cfg, err := config.LoadConfigFrom(testConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	if cfg == nil {
		log.Fatal(errors.New("unknown error loading yaml file"))
	}
	args := BuildDockerArgs(cfg.DB)
	container, err := docker.RunContainer("mysql:latest", "3306", WaitFunc, args...)
	if err != nil {
		t.Fatal(err)
	}
	cfg.DB.Host, cfg.DB.Port = datastore.ExtractHostAndPortFrom(container.Addr)

	return container, cfg
}

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

// InsertMultiple inserts multiple values into database
func InsertMultiple(db *gorm.DB, models ...interface{}) error {
	for _, v := range models {
		if err := db.Create(v).Error; err != nil {
			return err
		}
	}
	return nil
}

func CleanContainers(t *testing.T) {
	name := "cerebrum_mysql_test_no_"
	if err := docker.StopAndRemoveAllContainers(name); err != nil {
		t.Fatal(err)
	}
}
