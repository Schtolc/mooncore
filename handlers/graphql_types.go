package handlers

import (
	"github.com/Schtolc/mooncore/dependencies"
	"github.com/Schtolc/mooncore/models"
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
)

// AddressObject is a graphql object for address
var AddressObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "Address",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"lat": &graphql.Field{
			Type: graphql.Float,
		},
		"lon": &graphql.Field{
			Type: graphql.Float,
		},
	},
})

// PhotoObject is a graphql object for photo
var PhotoObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "Photo",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"path": &graphql.Field{
			Type: graphql.String,
		},
	},
})

// UserObject is a graphql object for user
var UserObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"email": &graphql.Field{
			Type: graphql.String,
		},
		"address": &graphql.Field{
			Type: AddressObject,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				address := models.Address{}
				if dbc := dependencies.DBInstance().First(&address, p.Source.(models.User).AddressID); dbc.Error != nil {
					logrus.Println(dbc.Error)
					return nil, dbc.Error
				}

				return address, nil
			},
		},
		"photo": &graphql.Field{
			Type: PhotoObject,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				photo := models.Photo{}
				if dbc := dependencies.DBInstance().First(&photo, p.Source.(models.User).PhotoID); dbc.Error != nil {
					logrus.Println(dbc.Error)
					return nil, dbc.Error
				}
				return photo, nil
			},
		},
	},
})
