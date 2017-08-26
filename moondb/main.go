package main

import(
	 _ "github.com/go-sql-driver/mysql"
	"log"
	"moondb/models"
	"github.com/jinzhu/gorm"
	"fmt"
)

func main() {
	database := "moondb";
	create_db(database)

	db, err := gorm.Open("mysql", "root@/"+database)
	check_error(err)

	err = db.DB().Ping()
	check_error(err)

	db.AutoMigrate(
		&models.Metric{},
	)
}



func create_db(name string) {
	db, err := gorm.Open("mysql", "root@/")
	check_error(err)

	defer db.Close()

	db.Exec("CREATE DATABASE IF NOT EXISTS " + name)
}

func check_error(err error){
	if err != nil {
		log.Fatal(err)
	}
	return;
}