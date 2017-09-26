package database

import (
	"github.com/Schtolc/mooncore/config"
	"github.com/Schtolc/mooncore/models"
	_ "github.com/go-sql-driver/mysql" // mysql driver for gorm.Open
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"sync"
)

func initDatabase(config config.Config) *gorm.DB {
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
		&models.Mock{},
		&models.Address{},
		&models.Photo{},
		&models.ManicureType{},
		&models.ManicureMaterial{},
		&models.User{},
		&models.Master{},
		&models.Service{},
	)

	db.Table("service_manicure_materials").AddForeignKey("service_id", "services(id)", "CASCADE", "CASCADE")
	db.Table("service_manicure_materials").AddForeignKey("manicure_material_id", "manicure_materials(id)", "CASCADE", "CASCADE")
	db.Table("service_photos").AddForeignKey("service_id", "services(id)", "CASCADE", "CASCADE")
	db.Table("service_photos").AddForeignKey("photo_id", "photos(id)", "CASCADE", "CASCADE")
	db.Table("working_place_photos").AddForeignKey("master_id", "masters(id)", "CASCADE", "CASCADE")
	db.Table("working_place_photos").AddForeignKey("photo_id", "photos(id)", "CASCADE", "CASCADE")

	logrus.Info("models migrated")
	return db
}

var instance *gorm.DB
var mutex = &sync.Mutex{}

// Instance returns database instance
func Instance() *gorm.DB {
	if instance != nil {
		return instance
	}
	mutex.Lock()
	defer mutex.Unlock()
	if instance == nil {
		instance = initDatabase(config.Instance())
	}
	return instance
}
