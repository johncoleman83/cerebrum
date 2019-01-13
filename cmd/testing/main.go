// Package testing is used only to execute golang commands
// for test uses, it is like a go playground
package testing

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	//"time"
	"database/sql"

	"github.com/jinzhu/gorm"
	"github.com/johncoleman83/cerebrum/pkg/utl/config"
	"github.com/johncoleman83/cerebrum/pkg/utl/datastore"
	cerebrum "github.com/johncoleman83/cerebrum/pkg/utl/model"
	"github.com/johncoleman83/cerebrum/pkg/utl/support"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

// main testing go playground
func main() {
	fmt.Println(basepath)
	fmt.Println("^^ basepath")
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
	dsn := datastore.FormatDSN(cfg.DB)
	dbSQL, err := sql.Open(cfg.DB.Dialect, dsn)
	if err != nil {
		fmt.Println("************************")
		fmt.Println("ERROR!!!!")
		fmt.Println(err)
	} else {
		fmt.Println("************************")
		fmt.Println("SUCCESSSSSFUL Going to Ping!!!")
		fmt.Println(dbSQL.Ping())
	}
	dbSQL.Close()
	db, err := datastore.NewMySQLGormDb(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	user := &cerebrum.User{}
	if err := db.Where("id = ?", 4).First(&user).Error; gorm.IsRecordNotFoundError(err) {
		fmt.Println("Record not found error")
		fmt.Println(gorm.IsRecordNotFoundError(err))
		fmt.Println(err)
		err = db.First(&user, 4).Error
		err2 := db.Error
		fmt.Println(err)
		fmt.Println(err2)
	}

	res := db.First(&user, 4).RecordNotFound()

	userTwo := &cerebrum.User{}
	db.Raw("SELECT * FROM `users`  WHERE `users`.`deleted_at` IS NULL AND ((`users`.`id` = 1)) ORDER BY `users`.`id` ASC LIMIT 1").Scan(&userTwo)
	fmt.Println("FINISHED NOW INSPECTING")
	fmt.Println(res)
	fmt.Println(userTwo)
	db.Close()
}
