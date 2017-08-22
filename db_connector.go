package main

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Metric struct {
	Id   string
	Path string
	Time time.Time
}

func (m *Metric) BeforeCreate(scope *gorm.Scope) {
	if m.Id == "" {
		uuid := uuid.New().String()
		scope.SetColumn("id", uuid)
	}
}

func ConnectDB() (db *gorm.DB) {
	db, err := gorm.Open("mysql", "root@/datadata")
	if err != nil {
		log.Fatal(err)
		panic("failed to connect database")
	}
	return db
}
