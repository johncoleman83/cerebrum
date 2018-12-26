package datastore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johncoleman83/cerebrum/pkg/utl/config"
	"github.com/johncoleman83/cerebrum/pkg/utl/datastore"
	"github.com/johncoleman83/cerebrum/pkg/utl/mock/mockdb"
	"github.com/johncoleman83/cerebrum/pkg/utl/mock/docker"
	"github.com/johncoleman83/cerebrum/pkg/utl/model"
	"github.com/johncoleman83/cerebrum/pkg/utl/support"
)

var (
	testConfigPath = support.TestingConfigPath()
)

func TestNew(t *testing.T) {
	cfg, err := config.LoadConfigFrom(testConfigPath)
	if cfg == nil || err != nil {
		t.Fatalf("Error loading test config %v", err)
	}
	args := mockdb.BuildDockerArgs(cfg.DB)
	container, err := docker.RunContainer("mysql:latest", "3306", mockdb.WaitFunc, args...)
	if err != nil {
		t.Fatalf("Error starting container %v", err)
	}

	defer container.Shutdown()

	corruptedDBcfg := cfg.DB
	corruptedDBcfg.Host, corruptedDBcfg.Port  = "pluto", "53456345634563"
	_, err = datastore.NewMySQLGormDb(corruptedDBcfg)
	if err == nil {
		t.Error("Expected error due to improper host")
	}

	corruptedDBcfg = cfg.DB
	corruptedDBcfg.Password = "root"
	_, err = datastore.NewMySQLGormDb(corruptedDBcfg)
	if err == nil {
		t.Error("Expected error due to incorrect password")
	}

	corruptedDBcfg.Password = "admin"
	_, err = datastore.NewMySQLGormDb(corruptedDBcfg)
	if err == nil {
		t.Error("Expected error due to incorrect password")
	}

	db, err := datastore.NewMySQLGormDb(cfg.DB)
	if err != nil {
		t.Fatalf("Error establishing connection %v", err)
	}

	var user cerebrum.User
	found := db.First(&user).RecordNotFound()

	assert.Nil(t, db.Error, "there should not be an error in querying the DB")
	assert.True(t, found, "there should be a proper response from RecordNotFound when the DB is empty")

	assert.Nil(t, db.Close().Error, "there should not be an error closing the DB")
}
