package main

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	//"time"
	"database/sql"

	"github.com/johncoleman83/cerebrum/pkg/utl/config"
	"github.com/johncoleman83/cerebrum/pkg/utl/datastore"
	//"github.com/johncoleman83/cerebrum/pkg/utl/mock/mockdb"
	//"github.com/johncoleman83/cerebrum/pkg/utl/mock/docker"
	"github.com/johncoleman83/cerebrum/pkg/utl/model"	
)

var (
	_, b, _, _ 		 = runtime.Caller(0)
	basepath   		 = filepath.Dir(b)
)

func WaitFunc(addr string) error {
	return nil
}

func main() {
	fmt.Println(basepath)
	fmt.Println("^^ basepath")
	cfg, errConfig := config.LoadConfigFromFlags()
	checkErr(errConfig)
	if cfg == nil {
		log.Fatal("unknown error loading yaml file")
	}
	dsn := datastore.FormatDSN(cfg.DB)
	dbSQL, errSQL := sql.Open(cfg.DB.Dialect, dsn)
	if errSQL != nil {
		fmt.Println("************************")
		fmt.Println("ERROR!!!!")
		fmt.Println(errSQL)
	} else {
		fmt.Println("************************")
		fmt.Println("SUCCESSSSSFUL!!!")
		fmt.Println(dbSQL.Ping())
	}
	dbSQL.Close()
	db, errDB := datastore.NewMySQLGormDb(cfg.DB)
	checkErr(errDB)
	
	user := &cerebrum.User{}
	res := db.First(&user, 4).RecordNotFound()

	userTwo := &cerebrum.User{}
	db.Raw("SELECT * FROM `users`  WHERE `users`.`deleted_at` IS NULL AND ((`users`.`id` = 1)) ORDER BY `users`.`id` ASC LIMIT 1").Scan(&userTwo)
	fmt.Println("FINISHED NOW INSPECTING")
	fmt.Println(res)
	fmt.Println(userTwo)
	db.Close()
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
