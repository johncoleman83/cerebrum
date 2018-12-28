package main

import (
	"fmt"
	"log"

	"github.com/johncoleman83/cerebrum/pkg/utl/config"
	"github.com/johncoleman83/cerebrum/pkg/utl/datastore"
	"github.com/johncoleman83/cerebrum/pkg/utl/model"
	"github.com/johncoleman83/cerebrum/pkg/utl/secure"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func buildQueries() ([7]string) {
	return [7]string{
		"INSERT INTO companies VALUES (1, now(), now(), NULL, 'admin_company', true);",
		"INSERT INTO locations VALUES (1, now(), now(), NULL, 'admin_location', true, 'admin_address', 1);",
		"INSERT INTO roles VALUES (100, 100, 'SUPER_ADMIN');",
		"INSERT INTO roles VALUES (110, 110, 'ADMIN');",
		"INSERT INTO roles VALUES (120, 120, 'COMPANY_ADMIN');",
		"INSERT INTO roles VALUES (130, 130, 'LOCATION_ADMIN');",
		"INSERT INTO roles VALUES (200, 200, 'USER');",
	}
}

func main() {
	queries := buildQueries()
	cfg, err := config.LoadConfigFromFlags()
	checkErr(err)
	if cfg == nil {
		log.Fatal("unknown error loading yaml file")
	}
	db, err := datastore.NewMySQLGormDb(cfg.DB)
	checkErr(err)

	createSchema(db, &cerebrum.Company{}, &cerebrum.Location{}, &cerebrum.Role{}, &cerebrum.User{})

	for _, v := range queries[0 : len(queries)-1] {
		db.Exec(v)
	}

	sec := secure.New(1, nil)

	userInsert := `INSERT INTO users (id, created_at, updated_at, first_name, last_name, username, password, email, active, role_id, company_id, location_id) VALUES (1, now(),now(),'Admin', 'Admin', 'admin', '%s', 'johndoe@mail.com', true, 100, 1, 1);`
	db.Exec(fmt.Sprintf(userInsert, sec.Hash("admin")))
	fmt.Println(fmt.Sprintf("bootstrap finished with %d errors", len(db.GetErrors())))
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func createSchema(db *gorm.DB, models ...interface{}) {
	for _, model := range models {
		if db.HasTable(model) {
			log.Printf("dropping table for ")
			db.DropTable(model)
		}
		checkErr(db.CreateTable(model).Error)
	}
}