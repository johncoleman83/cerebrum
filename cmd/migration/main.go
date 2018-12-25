package main

import (
	"fmt"
	"log"

	"github.com/johncoleman83/cerebrum/pkg/utl/config"
	"github.com/johncoleman83/cerebrum/pkg/utl/model"
	"github.com/johncoleman83/cerebrum/pkg/utl/secure"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	queries := [7]string{
		"INSERT INTO companies VALUES (1, now(), now(), NULL, 'admin_company', true);",
		"INSERT INTO locations VALUES (1, now(), now(), NULL, 'admin_location', true, 'admin_address', 1);",
		"INSERT INTO roles VALUES (100, 100, 'SUPER_ADMIN');",
		"INSERT INTO roles VALUES (110, 110, 'ADMIN');",
		"INSERT INTO roles VALUES (120, 120, 'COMPANY_ADMIN');",
		"INSERT INTO roles VALUES (130, 130, 'LOCATION_ADMIN');",
		"INSERT INTO roles VALUES (200, 200, 'USER');",
	}
	cfg, err := config.LoadConfig()
	checkErr(err)
	if cfg == nil {
		log.Fatal("unknown error loading yaml file")
	}
	var args = fmt.Sprintf(
		"%s:%s@%s(%s:%s)/%s?%s",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Protocol,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
		cfg.DB.Params,
	)
	
	db, err := gorm.Open(cfg.DB.Dialect, args)
	checkErr(err)

	db.LogMode(true)
	db.Exec("SELECT 1")

	createSchema(db, &cerebrum.Company{}, &cerebrum.Location{}, &cerebrum.Role{}, &cerebrum.User{})

	for _, v := range queries[0 : len(queries)-1] {
		db.Exec(v)
	}

	sec := secure.New(1, nil)

	userInsert := `INSERT INTO users (id, created_at, updated_at, first_name, last_name, username, password, email, active, role_id, company_id, location_id) VALUES (1, now(),now(),'Admin', 'Admin', 'admin', '%s', 'johndoe@mail.com', true, 100, 1, 1);`
	db.Exec(fmt.Sprintf(userInsert, sec.Hash("admin")))
	fmt.Println(fmt.Sprintf("migration finished with %d errors", len(db.GetErrors())))
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func createSchema(db *gorm.DB, models ...interface{}) {
	for _, model := range models {
		if db.HasTable(model) {
			db.DropTable(model)
		}
		db.CreateTable(model)
	}
}
