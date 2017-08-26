package main

import (
	"log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)



func ConnectDB() (db *gorm.DB) {
	db, err := gorm.Open("mysql", "root@/moondb")
	check_error(err)

	err = db.DB().Ping()
	check_error(err)

	db.DB().SetMaxOpenConns(10)

	return db
}


func check_error(err error){
	if err != nil {
		log.Fatal(err)
	}
	return;
}
