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

// buildQueries creates some SQL queries into a string slice
func buildQueries() []string {
	return []string{
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
	createSchema(db)

	for _, v := range queries[0:len(queries)] {
		db.Exec(v)
	}

	sec := secure.New(cfg.App.MinPasswordStr, sha1.New())
	user := models.User{
		Base: models.Base{
			Model: gorm.Model{
				ID: 1,
			},
		},
		Email:     "rocinante@mail.com",
		FirstName: "Rocinante",
		LastName:  "DeLaMancha",
		Username:  adminUsername,
		RoleID:    1,
		AccountID: 1,
		TeamID:    1,
		Password:  adminPassword,
	}
	account := models.Account{
		Base: models.Base{
			Model: gorm.Model{
				ID: 1,
			},
		},
		Name:    "admin_account",
		OwnerID: user.ID,
	}
	team := models.Team{
		Base: models.Base{
			Model: gorm.Model{
				ID: 1,
			},
		},
		Name:      "admin_team",
		Description:   "admin_description",
		AccountID: account.ID,
	}

	if !sec.Password(user.Password, user.FirstName, user.LastName, user.Username, user.Email) {
		log.Fatal(fmt.Sprintf("Password %v is not strong enough", user.Password))
	}
	user.Password = sec.Hash(user.Password)
	if err := db.Create(&user).Error; err != nil {
		log.Fatal(err)
	}
	if err := db.Create(&account).Error; err != nil {
		log.Fatal(err)
	}
	if err := db.Create(&team).Error; err != nil {
		log.Fatal(err)
	}
	fmt.Println(fmt.Sprintf("bootstrap finished with %d db errors", len(db.GetErrors())))
}

func createSchema(db *gorm.DB) {
	modelsList := []interface{}{
		&models.Account{},
		&models.Team{},
		&models.Role{},
		&models.User{},
	}
	for _, model := range modelsList {
		if db.HasTable(model) {
			log.Printf("dropping table for ")
			if err := db.DropTable(model).Error; err != nil {
				log.Fatal(err)
			}
		}
		if err := db.CreateTable(model).Error; err != nil {
			log.Fatal(err)
		}
	}
}
