package models

import (
	"github.com/jinzhu/gorm"
	"github.com/google/uuid"
)

type User struct {
	Id   string
	Name string 	`sql:"size:255;index"`
	Age  int
}



func (m *User) BeforeCreate(scope *gorm.Scope) {
	if m.Id == "" {
		uuid := uuid.New().String()
		scope.SetColumn("id", uuid)
	}
}