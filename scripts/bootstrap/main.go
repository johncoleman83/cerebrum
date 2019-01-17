// Package main is used to bootstrap a DB for
// work in a development environment
package main

import (
	"crypto/sha1"
	"fmt"
	"log"

	"github.com/johncoleman83/cerebrum/pkg/utl/config"
	"github.com/johncoleman83/cerebrum/pkg/utl/datastore"
	cerebrum "github.com/johncoleman83/cerebrum/pkg/utl/model"
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
		"INSERT INTO roles VALUES (100, 100, 'SUPER_ADMIN');",
		"INSERT INTO roles VALUES (110, 110, 'ADMIN');",
		"INSERT INTO roles VALUES (120, 120, 'COMPANY_ADMIN');",
		"INSERT INTO roles VALUES (130, 130, 'LOCATION_ADMIN');",
		"INSERT INTO roles VALUES (200, 200, 'USER');",
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
	user := cerebrum.User{
		Base: cerebrum.Base{
			Model: gorm.Model{
				ID: 1,
			},
		},
		Email:      "rocinante@mail.com",
		FirstName:  "Rocinante",
		LastName:   "DeLaMancha",
		Username:   adminUsername,
		RoleID:     cerebrum.AccessRole(100),
		CompanyID:  1,
		LocationID: 1,
		Password:   adminPassword,
	}
	company := cerebrum.Company{
		Base: cerebrum.Base{
			Model: gorm.Model{
				ID: 1,
			},
		},
		Name:    "admin_company",
		OwnerID: user.ID,
	}
	location := cerebrum.Location{
		Base: cerebrum.Base{
			Model: gorm.Model{
				ID: 1,
			},
		},
		Name:      "admin_location",
		Address:   "admin_address",
		CompanyID: company.ID,
	}

	if !sec.Password(user.Password, user.FirstName, user.LastName, user.Username, user.Email) {
		log.Fatal(fmt.Sprintf("Password %v is not strong enough", user.Password))
	}
	user.Password = sec.Hash(user.Password)
	if err := db.Create(&user).Error; err != nil {
		log.Fatal(err)
	}
	if err := db.Create(&company).Error; err != nil {
		log.Fatal(err)
	}
	if err := db.Create(&location).Error; err != nil {
		log.Fatal(err)
	}
	fmt.Println(fmt.Sprintf("bootstrap finished with %d db errors", len(db.GetErrors())))
}

func createSchema(db *gorm.DB) {
	models := []interface{}{
		&cerebrum.Company{},
		&cerebrum.Location{},
		&cerebrum.Role{},
		&cerebrum.User{},
	}
	for _, model := range models {
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
