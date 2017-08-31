package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"mooncore/cfg"
	"mooncore/models"
)

func ConnectDB() (db *gorm.DB) {
	args := cfg.GetBdConfig("db.conf")

	db, err := gorm.Open(args.Dialect, args.User+"@/"+args.Dbname)
	db.AutoMigrate(
		&models.User{},
		&models.Metric{},
	)
	check_error(err)

	err = db.DB().Ping()
	check_error(err)

	db.DB().SetMaxOpenConns(10)

	return db
}

func check_error(err error) {
	if err != nil {
		log.Fatal(err)
	}
	return
}
