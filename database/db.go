package database

import (
	"github.com/Schtolc/mooncore/config"
	"github.com/Schtolc/mooncore/models"
	_ "github.com/go-sql-driver/mysql" // mysql driver for gorm.Open
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

// Init database: open connection + migrate schema + set connection params
func Init(config config.Config) (db *gorm.DB) {
	db, err := gorm.Open(config.Database.Dialect, config.Database.User+"@/"+config.Database.Dbname)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("connected to database")
	err = db.DB().Ping()
	if err != nil {
		logrus.Fatal(err)
	}

	db.DB().SetMaxOpenConns(config.Database.MaxOpenConns)
	db.AutoMigrate(
		&models.Metric{},
	)
	logrus.Info("migrate models")
	return db
}
