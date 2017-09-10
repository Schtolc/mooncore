package database

import (
	"github.com/Schtolc/mooncore/config"
	"github.com/Schtolc/mooncore/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

// Init database: open connection + migrate schema + set connection params
func Init(config config.Config) (db *gorm.DB) {
	db, err := gorm.Open(config.Database.Dialect, config.Database.User+"@tcp("+config.Database.Host+")/"+config.Database.Dbname)
	if err != nil {

		log.Fatal(err)
	}

	err = db.DB().Ping()
	if err != nil {
		log.Fatal(err)
	}

	db.DB().SetMaxOpenConns(config.Database.MaxOpenConns)
	db.AutoMigrate(
		&models.Metric{},
	)
	return db
}
