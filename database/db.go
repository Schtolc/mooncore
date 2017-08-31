package db



import(
	 _ "github.com/go-sql-driver/mysql"
	"mooncore/models"
	"github.com/jinzhu/gorm"
	"log"
)


func InitDB() (db *gorm.DB) {
	config := cfg.GetBdConfig("db.conf")

	db, err := gorm.Open(config.Dialect, cconfig.User + "@/" + config.Dbname)
	check_error(err)

	err = db.DB().Ping()
	check_error(err)

	db.DB().SetMaxOpenConns(10)
	db.AutoMigrate(
		&models.User{},
		&models.Metric{},
	)
	return db
}

//
//func create_db(name string) {
//	db, err := gorm.Open("mysql", "root@/")
//	check_error(err)
//
//	defer db.Close()
//
//	db.Exec("CREATE DATABASE IF NOT EXISTS " + name)
//}

func check_error(err error){
	if err != nil {
		log.Fatal(err)
	}
	return
}
