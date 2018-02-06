package dependencies

import (
	"github.com/Schtolc/mooncore/models"
	_ "github.com/go-sql-driver/mysql" // mysql driver for gorm.Open
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"sync"
)

func initDatabase(config *Config) *gorm.DB {
	db, err := gorm.Open(config.Database.Dialect, config.Database.User+"@/"+config.Database.Dbname)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("connected to database")
	if err = db.DB().Ping(); err != nil {
		logrus.Fatal(err)
	}
	db.DB().SetMaxOpenConns(config.Database.MaxOpenConns)
	db.AutoMigrate(
		&models.Mock{},
		&models.Address{},
		&models.Tag{},
		&models.Photo{},
		&models.Sign{},
		&models.User{},
		&models.Salon{},
		&models.Material{},
		&models.Master{},
		&models.ServiceType{},
		&models.Client{},
		&models.Service{},
	)

	db.Table("photo_tags").AddForeignKey("photo_id", "photos(id)", "CASCADE", "CASCADE")
	db.Table("photo_tags").AddForeignKey("tag_id", "tags(id)", "CASCADE", "CASCADE")

	db.Table("service_photos").AddForeignKey("service_id", "services(id)", "CASCADE", "CASCADE")
	db.Table("service_photos").AddForeignKey("photo_id", "photos(id)", "CASCADE", "CASCADE")

	db.Table("service_materials").AddForeignKey("service_id", "services(id)", "CASCADE", "CASCADE")
	db.Table("service_materials").AddForeignKey("material_id", "materials(id)", "CASCADE", "CASCADE")

	db.Table("master_photos").AddForeignKey("master_id", "masters(id)", "CASCADE", "CASCADE")
	db.Table("master_photos").AddForeignKey("photo_id", "photos(id)", "CASCADE", "CASCADE")
	db.Table("master_signs").AddForeignKey("master_id", "masters(id)", "CASCADE", "CASCADE")
	db.Table("master_signs").AddForeignKey("sign_id", "signs(id)", "CASCADE", "CASCADE")

	db.Table("client_favorites").AddForeignKey("client_id", "clients(id)", "CASCADE", "CASCADE")
	db.Table("client_favorites").AddForeignKey("master_id", "masters(id)", "CASCADE", "CASCADE")

	logrus.Info("models migrated")

	models.InsertDefaultValues(db)
	models.InsertConstValues(db)

	logrus.Info("create default and constant values in database")
	return db
}

var dbInstance *gorm.DB
var dbMutex = &sync.Mutex{}

// DBInstance returns database instance
func DBInstance() *gorm.DB {
	if dbInstance != nil {
		return dbInstance
	}
	dbMutex.Lock()
	defer dbMutex.Unlock()
	if dbInstance == nil {
		dbInstance = initDatabase(ConfigInstance())
	}
	return dbInstance
}
