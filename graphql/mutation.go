package graphql

import (
	"errors"
	"github.com/Schtolc/mooncore/dao"
	"github.com/Schtolc/mooncore/models"
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
	"strconv"
	"unicode"
)

var createMaster = &graphql.Field{
	Type:        MasterObject,
	Description: "Create new master",
	Args: graphql.FieldConfigArgument{
		"username":   notNull(graphql.String),
		"email":      notNull(graphql.String),
		"password":   notNull(graphql.String),
		"name":       notNull(graphql.String),
		"address_id": notNull(graphql.Int),
		"photo_id":   notNull(graphql.Int),
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		return dao.CreateMaster(
			params.Args["username"].(string),
			params.Args["email"].(string),
			params.Args["password"].(string),
			params.Args["name"].(string),
			params.Args["address_id"].(int64),
			params.Args["photo_id"].(int64))
	},
}

var createClient = &graphql.Field{
	Type:        ClientObject,
	Description: "Create new master",
	Args: graphql.FieldConfigArgument{
		"username": notNull(graphql.String),
		"email":    notNull(graphql.String),
		"password": notNull(graphql.String),
		"name":     notNull(graphql.String),
		"photo_id": notNull(graphql.Int),
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		return dao.CreateClient(
			params.Args["username"].(string),
			params.Args["email"].(string),
			params.Args["password"].(string),
			params.Args["name"].(string),
			params.Args["photo_id"].(int64))
	},
}

var signIn = &graphql.Field{
	Type:        TokenObject, // nil if user not found
	Description: "Sign in",
	Args: graphql.FieldConfigArgument{
		"username": notNull(graphql.String),
		"email":    notNull(graphql.String),
		"password": notNull(graphql.String),
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		return dao.SignIn(
			params.Args["username"].(string),
			params.Args["email"].(string),
			params.Args["password"].(string))
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

var addPhoto = &graphql.Field{
	Type:        PhotoObject,
	Description: "Add new photo - master",
	Args: graphql.FieldConfigArgument{
		"path":  notNull(graphql.String),
		"tags":  listOf(graphql.Int),
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

