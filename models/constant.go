package models

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

var (
	//DefaultAddress
	DefaultAddress Address
	//DefaultAvatar
	DefaultAvatar Photo
	//Tags array
	Tags []Tag
	//Signs array
	Signs []Sign
)

func createConstValues() {
	tags_name := []string{
		"Shellac", "Paraffin", "Mirror", "Vinylux", "3D_Nails",
	}
	for _, tag_name := range tags_name {
		Tags = append(Tags, Tag{Name: tag_name})
	}

	Signs = append(Signs, Sign{
		Name:        "accuracy",
		PhotoID:     DefaultAvatar.ID,
		Description: "means accuracy",
	}, Sign{
		Name:        "politeness",
		PhotoID:     DefaultAvatar.ID,
		Description: "means politeness",
	}, Sign{
		Name:        "varnish resistance",
		PhotoID:     DefaultAvatar.ID,
		Description: "means varnish resistance",
	})
}
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
