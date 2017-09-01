package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"mooncore/config"
	"mooncore/models"
)

func Init(config config.Config) (db *gorm.DB) {
	db, err := gorm.Open(config.Database.Dialect, config.Database.User+"@/"+config.Database.Dbname)
	check_error(err)

	err = db.DB().Ping()
	check_error(err)

	db.DB().SetMaxOpenConns(10)
	db.AutoMigrate(
		&models.Metric{},
	)
	return db
}

func check_error(err error) {
	if err != nil {
		log.Fatal(err)
	}
	return
}
