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
	db.LogMode(true)
	db.DB().SetMaxOpenConns(config.Database.MaxOpenConns)
	db.AutoMigrate(
		&models.Mock{},
		&models.Address{},
		&models.Photo{},
		&models.Tag{},
		&models.ManicureType{},
		&models.Sign{},
		&models.User{},
		&models.UserDetails{},
		&models.Service{},
		&models.UserAuth{},
	)

	db.Table("user_signs").AddForeignKey("user_details_id", "user_details(id)", "CASCADE", "CASCADE")
	db.Table("user_signs").AddForeignKey("sign_id", "signs(id)", "CASCADE", "CASCADE")

	db.Table("service_photos").AddForeignKey("service_id", "services(id)", "CASCADE", "CASCADE")
	db.Table("service_photos").AddForeignKey("photo_id", "photos(id)", "CASCADE", "CASCADE")

	db.Table("user_photos").AddForeignKey("user_details_id", "user_details(id)", "CASCADE", "CASCADE")
	db.Table("user_photos").AddForeignKey("photo_id", "photos(id)", "CASCADE", "CASCADE")

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
