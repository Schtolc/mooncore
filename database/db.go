package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"mooncore/cfg"
	"mooncore/models"
)

func InitDB(config cfg.AppConfig) (db *gorm.DB) {
	db, err := gorm.Open(config.Dialect, config.User+"@/"+config.Dbname)
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
