package models

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

const (
	// ClientRole for user model
	ClientRole = 1 << iota // 00001 = 1
	// MasterRole for user model
	MasterRole = 1 << iota // 00010 = 2
	// SalonRole for user model
	SalonRole = 1 << iota // 00100 = 4
	// AdminRole for user model
	AdminRole = 1 << iota // 01000 = 8
	// AnonRole for user model
	AnonRole = (1 << iota) - 1
)

var (
	// DefaultAddress for unknown address
	DefaultAddress = &Address{
		Lat:         0,
		Lon:         0,
		Description: "default",
	}
	// DefaultAvatar  for unknown photos
	DefaultAvatar = &Photo{
		Path: "default",
		Tags: []Tag{},
	}
	// Tags array
	Tags []Tag
	// Signs array
	Signs []Sign
	// AnonUser model
	AnonUser = &User{
		ID:           1,
		Email:        "asd",
		PasswordHash: "asd",
		Role:         AnonRole,
		Ctime:        1,
		Atime:        1,
	}
)

func createConstValues() {
	tagsName := []string{
		"Shellac", "Paraffin", "Mirror", "Vinylux", "3D_Nails",
	}
	for _, tagName := range tagsName {
		Tags = append(Tags, Tag{Name: tagName})
	}

	Signs = append(Signs, Sign{
		Name:        "accuracy",
		Icon:        "default",
		Description: "means accuracy",
	}, Sign{
		Name:        "politeness",
		Icon:        "default",
		Description: "means politeness",
	}, Sign{
		Name:        "varnish resistance",
		Icon:        "default",
		Description: "means varnish resistance",
	})
}

// InsertDefaultValues of address and photos
func InsertDefaultValues(db *gorm.DB) {
	tx := db.Begin()
	if err := tx.FirstOrCreate(&DefaultAddress).Error; err != nil {
		tx.Rollback()
		logrus.Fatal(err)
	}
	if err := tx.FirstOrCreate(&DefaultAvatar).Error; err != nil {
		tx.Rollback()
		logrus.Fatal(err)
	}
	tx.Commit()
}

// InsertConstValues for tags and signs
func InsertConstValues(db *gorm.DB) {
	createConstValues()
	tx := db.Begin()
	for _, sign := range Signs {
		if err := tx.FirstOrCreate(&Sign{}, &sign).Error; err != nil {
			tx.Rollback()
			logrus.Fatal(err)
		}
	}
	for _, tag := range Tags {
		if err := tx.FirstOrCreate(&Tag{}, &tag).Error; err != nil {
			tx.Rollback()
			logrus.Fatal(err)
		}
	}
	tx.Commit()
}
