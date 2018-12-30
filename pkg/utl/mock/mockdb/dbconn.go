package mockdb

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // for use with gorm
	"github.com/ory/dockertest"

	"github.com/johncoleman83/cerebrum/pkg/utl/config"
	"github.com/johncoleman83/cerebrum/pkg/utl/datastore"
	cerebrum "github.com/johncoleman83/cerebrum/pkg/utl/model"
	"github.com/johncoleman83/cerebrum/pkg/utl/support"
)

var (
	testConfigPath = support.TestingConfigPath()
)

// Container has info helpful for testing with a docker mysql container
type Container struct {
	Configuration *config.Configuration
	Pool          *dockertest.Pool
	Resource      *dockertest.Resource
	DB            *gorm.DB
}

func getTestConfig(t *testing.T) *config.Configuration {
	// load mysql config from config file
	cfg, err := config.LoadConfigFrom(testConfigPath)
	if err != nil {
		t.Fatal(err)
	}
	if cfg == nil {
		t.Fatal(errors.New("unknown error loading yaml file"))
	}
	return cfg
}

func buildDockerOptions(DB *config.Database) *dockertest.RunOptions {
	randomString := support.NewRandomString(25)
	containerName := fmt.Sprintf("cerebrum_mysql_test_db_no_%s", randomString)
	return &dockertest.RunOptions{
		Name:       containerName,
		Repository: "mysql",
		Tag:        "latest",
		Env: []string{
			"MYSQL_ALLOW_EMPTY_PASSWORD=yes",
			"MYSQL_DATABASE=" + DB.Name,
			"MYSQL_USER=" + DB.User,
			"MYSQL_PASSWORD=" + DB.Password,
		},
		Cmd: []string{"mysqld", "--default-authentication-plugin=mysql_native_password"},
	}
}

func loopPingDB(t *testing.T, pool *dockertest.Pool, cfg *config.Configuration) {
	dsn := datastore.FormatDSN(cfg.DB)
	log.Println("pinging db in the docker test container to verify mysql has started up\nDSN: " + dsn)
	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		db, err := sql.Open(cfg.DB.Dialect, dsn)
		if err != nil {
			return err
		}
		if err = db.Ping(); err != nil {
			db.Close()
		}
		return err
	}); err != nil {
		t.Fatal(fmt.Sprintf("Could not connect to docker: %v", err))
	}
	log.Println("end verify mysql has started up")
}

// NewMySQLDockerTestContainer instantiates new mysql docker container
func NewMySQLDockerTestContainer(t *testing.T) *Container {
	// init docker daemon connection
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatal(fmt.Sprintf("Could not connect to docker: %v", err))
	}

	cfg := getTestConfig(t)
	runOptions := buildDockerOptions(cfg.DB)
	// pulls an image, creates a container based on the run options
	resource, err := pool.RunWithOptions(runOptions)
	if err != nil {
		t.Fatal(fmt.Sprintf("Could not start resource: %v", err))
	}

	// update connection host and port based on new docker container
	cfg.DB.Host, cfg.DB.Port = datastore.ExtractHostAndPortFrom(resource.GetHostPort("3306/tcp"))

	loopPingDB(t, pool, cfg)

	return &Container{
		Configuration: cfg,
		Pool:          pool,
		Resource:      resource,
		DB:            NewDBConn(t, cfg, cerebrum.Role{}, &cerebrum.User{}),
	}
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
