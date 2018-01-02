package handlers

import (
	"errors"
	"github.com/Schtolc/mooncore/models"
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
	"strconv"
)

var createUser = &graphql.Field{
	Type:        UserObject,
	Description: "Create new user",
	Args: graphql.FieldConfigArgument{
		"email":    notNull(graphql.String),
		"password": notNull(graphql.String),
		"role":     notNull(graphql.Int),
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		tx := db.Begin()
		user := &models.User{
			Email:    params.Args["email"].(string),
			Password: params.Args["password"].(string),
			Role:     params.Args["role"].(int),
		}
		if err := tx.Create(user).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		userDetails := &models.UserDetails{
			UserID:    user.ID,
			AddressID: models.DefaultAddress.ID,
			PhotoID:   models.DefaultAvatar.ID,
		}
		if err := tx.Create(userDetails).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		tx.Commit()
		return user, nil
	},
}

var editUserProfile = &graphql.Field{
	Type:        UserDetailsObject,
	Description: "Edit user fields",
	Args: graphql.FieldConfigArgument{
		"email":      just(graphql.String),
		"name":       just(graphql.String),
		"address_id": just(graphql.Int),
		"avatar_id":  just(graphql.Int),
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		userDetails, err := getUserDetails(params.Args["email"].(string))
		if err != nil {
			logrus.Warn(err)
			return nil, err
		}
		if params.Args["name"] != nil {
			userDetails.Name = params.Args["name"].(string)
		}
		if params.Args["address_id"] != nil {
			address := &models.Address{}
			if dbc := db.Where("id = ?", params.Args["address_id"].(int)).First(address); dbc.Error != nil {
				logrus.Error(dbc.Error)
				return nil, dbc.Error
			}
			userDetails.AddressID = address.ID
		}
		if params.Args["photo"] != nil {
			photo := &models.Photo{}
			if dbc := db.Where("id = ?", params.Args["avatar_id"].(int)).First(photo); dbc.Error != nil {
				logrus.Error(dbc.Error)
				return nil, dbc.Error
			}
			userDetails.PhotoID = photo.ID
		}
		if dbc := db.Save(&userDetails); dbc.Error != nil {
			logrus.Error(dbc.Error)
			return nil, dbc.Error
		}
		return *userDetails, nil
	},
}

var createAddress = &graphql.Field{
	Type:        AddressObject,
	Description: "Create new address",
	Args: graphql.FieldConfigArgument{
		"lat":         notNull(graphql.String),
		"lon":         notNull(graphql.String),
		"description": notNull(graphql.String),
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		lat, err := strconv.ParseFloat(params.Args["lat"].(string), 64)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		lon, err := strconv.ParseFloat(params.Args["lon"].(string), 64)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		address := &models.Address{
			Lat:         lat,
			Lon:         lon,
			Description: params.Args["description"].(string),
		}
		if dbc := db.Create(address); dbc.Error != nil {
			logrus.Error(dbc.Error)
			return nil, dbc.Error
		}
		return *address, nil
	},
}

var createAvatar = &graphql.Field{
	Type:        PhotoObject,
	Description: "Create new address",
	Args: graphql.FieldConfigArgument{
		"path": notNull(graphql.String),
		"tags": listof(graphql.Int),
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		photo := &models.Photo{
			Path: params.Args["path"].(string),
			Tags: []models.Tag{},
		}
		if dbc := db.Where("id in (?)", params.Args["tags"]).Find(&photo.Tags); dbc.Error != nil {
			logrus.Error(dbc.Error)
			return nil, dbc.Error
		}
		if len(photo.Tags) == 0 && params.Args["tags"] != nil {
			logrus.Error("No such tags")
			return nil, errors.New("No such tags")
		}

		if dbc := db.Create(photo); dbc.Error != nil {
			logrus.Error(dbc.Error)
			return nil, dbc.Error
		}
		return *photo, nil
	},
}

var addPhoto = &graphql.Field{
	Type:        PhotoObject,
	Description: "Add new photo - master",
	Args: graphql.FieldConfigArgument{
		"path":  notNull(graphql.String),
		"tags":  listof(graphql.Int),
		"email": notNull(graphql.String),
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		userDetails, err := getUserDetails(params.Args["email"].(string))
		if err != nil {
			logrus.Warn(err)
			return nil, err
		}
		photo := &models.Photo{
			Path: params.Args["path"].(string),
			Tags: []models.Tag{},
		}
		if dbc := db.Where("id in (?)", params.Args["tags"]).Find(&photo.Tags); dbc.Error != nil {
			logrus.Error(dbc.Error)
			return nil, dbc.Error
		}
		if len(photo.Tags) == 0 && params.Args["tags"] != nil {
			logrus.Error("No such tags")
			return nil, errors.New("No such tags")
		}
		if dbc := db.Create(photo); dbc.Error != nil {
			logrus.Error(dbc.Error)
			return nil, dbc.Error
		}
		if dbc := db.Model(&userDetails).Association("photos").Append([]models.Photo{*photo}); dbc.Error != nil {
			logrus.Error(dbc.Error)
			return nil, dbc.Error
		}
		return *photo, nil

	},
}

var addSigns = &graphql.Field{
	Type:        UserDetailsObject,
	Description: "Add signs -  admin",
	Args: graphql.FieldConfigArgument{
		"signs": notNull(graphql.NewList(graphql.Int)),
		"email": notNull(graphql.String),
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		userDetails, err := getUserDetails(params.Args["email"].(string))
		if err != nil {
			logrus.Warn(err)
			return nil, err
		}
		signs := []models.Sign{}
		if dbc := db.Where("id in (?)", params.Args["signs"]).Find(&signs); dbc.Error != nil {
			logrus.Error(dbc.Error)
			return nil, dbc.Error
		}
		if len(signs) == 0 {
			return nil, errors.New("No such signs")
		}
		tx := db.Begin()
		for _, sign := range signs {
			if err := tx.Model(&userDetails).Association("signs").Append([]models.Sign{sign}).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
		tx.Commit()
		return *userDetails, nil
	},
}
