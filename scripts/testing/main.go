// Package main is used to bootstrap a DB for
// work in a development environment
package main

import (
	"crypto/sha1"
	"fmt"
	"log"

	"github.com/johncoleman83/cerebrum/pkg/utl/config"
	"github.com/johncoleman83/cerebrum/pkg/utl/datastore"
	"github.com/johncoleman83/cerebrum/pkg/utl/models"
	"github.com/johncoleman83/cerebrum/pkg/utl/secure"
	"github.com/johncoleman83/cerebrum/pkg/utl/support"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // for use with gorm
)

const (
	adminUsername = "rocinante"
	adminPassword = "zvuEFGa84598705027345SDfhlasdfasjzqGRFs"
)

func createUser(cfg *config.Configuration, sec *secure.Service, r uint, e, f, l, u, p string) models.User {
	user := models.User{
		Email:     e,
		FirstName: f,
		LastName:  l,
		Username:  u,
		RoleID:    r,
		AccountID: 1,
		TeamID:    1,
		Password:  p,
	}
	if ok := sec.Password(user.Password, user.FirstName, user.LastName, user.Username, user.Email); !ok {
		log.Fatal(fmt.Sprintf("Password %v is not strong enough", user.Password))
	}
	user.Password = sec.Hash(user.Password)
	return user
}

// buildQueries creates some SQL queries into a string slice
func buildQueries() []string {
	return []string{
		"INSERT INTO accounts VALUES (1, now(), now(), NULL, 'admin_account', true);",
		"INSERT INTO teams VALUES (1, now(), now(), NULL, 'admin_team', true, 'admin_description', 1);",
		"INSERT INTO roles VALUES (1, 100, 'SUPER_ADMIN');",
		"INSERT INTO roles VALUES (2, 110, 'ADMIN');",
		"INSERT INTO roles VALUES (3, 120, 'ACCOUNT_ADMIN');",
		"INSERT INTO roles VALUES (4, 130, 'TEAM_ADMIN');",
		"INSERT INTO roles VALUES (5, 200, 'USER');",
	}
}

// main bootstrap a db
func main() {
	cfgPath, err := support.ExtractPathFromFlags()
	if err != nil {
		panic(err.Error())
	}
	cfg, err := config.LoadConfigFrom(cfgPath)
	if err != nil {
		log.Fatal(err)
	}
	if cfg == nil {
		log.Fatal("unknown error loading yaml file")
	}
	db, err := datastore.NewMySQLGormDb(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	queries := buildQueries()
	createSchema(db, &models.Account{}, &models.Team{}, models.Role{}, &models.User{})
	for _, v := range queries[0:len(queries)] {
		db.Exec(v)
	}
	sec := secure.New(cfg.App.MinPasswordStr, sha1.New())
	adminUser := createUser(
		cfg,
		sec,
		1,
		"rocinante@mail.com",
		"Rocinante",
		"DeLaMancha",
		adminUsername,
		adminPassword)
	if err := db.Create(&adminUser).Error; err != nil {
		log.Fatal(err)
	}
	var checkUser = new(models.User)
	if err := db.Set("gorm:auto_preload", true).Where("id = ?", adminUser.ID).First(&checkUser).Error; err != nil {
		log.Fatal(err)
	}
	if ok := sec.HashMatchesPassword(checkUser.Password, adminPassword); !ok {
		log.Println("ADMIN PASSWORD DOES NOT MATCH")
	}
	log.Println("ADMIN PASSWORD DOES MATCH!!")
	userUser := createUser(
		cfg,
		sec,
		1,
		"user1@mail.com",
		"user1_first",
		"user1_last",
		"user1",
		adminPassword)
	if err := db.Create(&userUser).Error; err != nil {
		log.Fatal(err)
	}
	checkUser = new(models.User)
	if err := db.Set("gorm:auto_preload", true).Where("id = ?", userUser.ID).First(&checkUser).Error; err != nil {
		log.Fatal(err)
	}
	if ok := sec.HashMatchesPassword(checkUser.Password, adminPassword); !ok {
		log.Println("USER PASSWORD DOES NOT MATCH")
	}
	log.Println("USER PASSWORD DOES MATCH!!")
}

func createSchema(db *gorm.DB, modelsList ...interface{}) {
	for _, model := range modelsList {
		if db.HasTable(model) {
			log.Printf("dropping table for ")
			db.DropTable(model)
		}
		if err := db.CreateTable(model).Error; err != nil {
			log.Fatal(err)
		}
	}
}
