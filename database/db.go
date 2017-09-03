package database

import (
	"github.com/Schtolc/mooncore/config"
	"github.com/Schtolc/mooncore/models"
	_ "github.com/go-sql-driver/mysql" // mysql driver for gorm.Open
	"github.com/jinzhu/gorm"
	"log"
)

// Init connects to database specified on config.yml. After successful connection models are migrated.
// If any error occurs program is terminated.
func Init(config config.Config) (db *gorm.DB) {
	db, err := gorm.Open(config.Database.Dialect, config.Database.User+"@/"+config.Database.Dbname)
	checkError(err)

	err = db.DB().Ping()
	checkError(err)

	db.DB().SetMaxOpenConns(config.Database.MaxOpenConns)
	db.AutoMigrate(
		&models.Metric{},
	)
	return db
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
