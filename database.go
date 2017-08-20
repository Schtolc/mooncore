package main


import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
	"fmt"
)

type Product struct {
	gorm.Model
	Code string
	Time time.Time
}


func ConnectDB() ( db *gorm.DB) {
	db, err := gorm.Open("mysql", "patrik:qweqwe@/qwa")
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}
	db.AutoMigrate(&Product{})
	return db
}
