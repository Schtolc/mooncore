package database

import (
	"github.com/Schtolc/mooncore/models"
	"github.com/Schtolc/mooncore/utils"
	_ "github.com/go-sql-driver/mysql" // mysql driver for gorm.Open
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

var instance *gorm.DB

// InitDatabase opens db connection + migrates schema + sets connection params
func InitDatabase(config utils.Config) *gorm.DB {
	if instance != nil {
		return instance
	}
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

	instance = db
	return db
}

// GetInstance returns database instance
func GetInstance() *gorm.DB {
	return instance
}
