package database

import (
	"github.com/Schtolc/mooncore/config"
	"github.com/Schtolc/mooncore/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
)

func Init(config config.Config) (db *gorm.DB) {
	db, err := gorm.Open(config.Database.Dialect, config.Database.User+"@/"+config.Database.Dbname)
	check_error(err)

	err = db.DB().Ping()
	check_error(err)

	db.DB().SetMaxOpenConns(config.Database.MaxOpenConns)
	db.AutoMigrate(
		&models.Metric{},
	)
	return db
}

func check_error(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
