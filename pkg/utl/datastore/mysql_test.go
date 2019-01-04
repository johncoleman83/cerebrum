package datastore_test

// TODO: NEED TO UPDATE THIS TO ACCOUNT FOR NEW CONTAINER SOLUTION
import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johncoleman83/cerebrum/pkg/utl/config"
	"github.com/johncoleman83/cerebrum/pkg/utl/datastore"
	"github.com/johncoleman83/cerebrum/pkg/utl/mock/mockdb"
)

func TestNew(t *testing.T) {
	container := mockdb.NewMySQLDockerTestContainer(t)
	db, cfg, pool, resource := container.DB, container.Configuration, container.Pool, container.Resource

	dsn := datastore.FormatDSN(cfg.DB)
	expectedDsn := "mysql_test_user:mysql_test_password" +
		fmt.Sprintf("@tcp(localhost:%s)/cerebrum_mysql_test_db", cfg.DB.Port) +
		"?tls=skip-verify&charset=utf8&parseTime=True&loc=Local&autocommit=true&timeout=20s"
	assert.Equal(t, expectedDsn, dsn, "dsn should be properly formated")

	corruptedDBcfg := &config.Database{
		Dialect:  cfg.DB.Dialect,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		Name:     cfg.DB.Name,
		Protocol: cfg.DB.Protocol,
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Settings: cfg.DB.Settings,
	}
	corruptedDBcfg.Host, corruptedDBcfg.Port = "pluto", "53456345634563"
	_, err := datastore.NewMySQLGormDb(corruptedDBcfg)
	assert.EqualError(t, err, err.Error(), "there should be an error connecting to mysql with bad config")

	corruptedDBcfg.Host, corruptedDBcfg.Port = cfg.DB.Host, cfg.DB.Port
	corruptedDBcfg.Password = "root"
	_, err = datastore.NewMySQLGormDb(corruptedDBcfg)
	assert.EqualError(t, err, err.Error(), "there should be an error connecting to mysql with bad config")

	corruptedDBcfg.Password = "admin"
	_, err = datastore.NewMySQLGormDb(corruptedDBcfg)
	assert.EqualError(t, err, err.Error(), "there should be an error connecting to mysql with bad config")

	db, err = datastore.NewMySQLGormDb(cfg.DB)
	if err != nil {
		t.Fatalf("Error establishing connection %v", err)
	}

	assert.Nil(t, db.Close(), "there should not be an error closing the DB")
	if err := pool.Purge(resource); err != nil {
		t.Fatal(fmt.Sprintf("Could not purge resource: %v", err))
	}
}