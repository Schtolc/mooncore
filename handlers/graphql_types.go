package handlers

import (
	"github.com/Schtolc/mooncore/dependencies"
	"github.com/Schtolc/mooncore/models"
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
)

// AddressObject is a graphql object for address
var AddressObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "AddressObject",
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
		"description": &graphql.Field{
			Type: graphql.String,
		},
	},
})
// TagObject is a graphql object for photos tags
var TagObject =  graphql.NewObject(graphql.ObjectConfig{
	Name: "TagObject",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
	},
})

// ManicureType is a graphql object for photos tags
var ManicureTypeObject =  graphql.NewObject(graphql.ObjectConfig{
	Name: "ManicureTypeObject",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
	},
})

// ServiceObject is a graphql object for photos tags
var ServiceObject =  graphql.NewObject(graphql.ObjectConfig{
	Name: "ServiceObject",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"description": &graphql.Field{
			Type: graphql.String,
		},
		"price": &graphql.Field{
			Type: graphql.Float,
		},
		"master": &graphql.Field{
			Type: UserDetailsObject,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				master := models.UserDetails{}
				if dbc := dependencies.DBInstance().First(&master, p.Source.(models.Service).MasterID); dbc.Error != nil {
					logrus.Println(dbc.Error)
					return nil, dbc.Error
				}
				return master, nil
			},
		},
		"manicure": &graphql.Field{
			Type: ManicureTypeObject,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				manicure := models.ManicureType{}
				if dbc := dependencies.DBInstance().First(&manicure, p.Source.(models.Service).ManicureTypeID); dbc.Error != nil {
					logrus.Println(dbc.Error)
					return nil, dbc.Error
				}
				return manicure, nil
			},
		},
	},
})


// PhotoObject is a graphql object for photo
var PhotoObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "PhotoObject",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"path": &graphql.Field{
			Type: graphql.String,
		},
		"tags": &graphql.Field{
			Type: graphql.NewList(TagObject),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				logrus.Warn(p.Source.(*models.Photo))
				logrus.Warn(p.Source.(*models.Photo).Tags)
				return p.Source.(*models.Photo).Tags, nil
			},
		},
	},
})
// SignObject is a graphql object for photos tags
var SignObject =  graphql.NewObject(graphql.ObjectConfig{
	Name: "SignObject",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"description":&graphql.Field{
			Type: graphql.String,
		},
		"photo": &graphql.Field{
			Type: PhotoObject,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				photo := models.Photo{}
				if dbc := dependencies.DBInstance().First(&photo, p.Source.(models.Sign).PhotoID); dbc.Error != nil {
					logrus.Error(dbc.Error)
					return nil, dbc.Error
				}
				return photo, nil
			},
		},
	},
})

// UserAuth is a graphql object for user
var UserObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "UserObject",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"email": &graphql.Field{
			Type: graphql.String,
		},
		"role": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

// UserDetails is a graphql object for user
var  UserDetailsObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "UserDetailsObject",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"address": &graphql.Field{
			Type: AddressObject,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				address := &models.Address{}
				if dbc := dependencies.DBInstance().First(address, p.Source.(*models.UserDetails).AddressID); dbc.Error != nil {
					logrus.Error(dbc.Error)
					return nil, dbc.Error
				}
				return address, nil
			},
		},
		"avatar": &graphql.Field{
			Type: PhotoObject,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				photo := &models.Photo{}
				if dbc := dependencies.DBInstance().First(photo, p.Source.(*models.UserDetails).PhotoID); dbc.Error != nil {
					logrus.Error(dbc.Error)
					return nil, dbc.Error
				}
				return photo, nil
			},
		},
		"photos": &graphql.Field{
			Type: graphql.NewList(PhotoObject),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(*models.UserDetails).Photos, nil
			},
		},
		"user_id":&graphql.Field{
			Type: graphql.Int,
		},
		"stars": &graphql.Field{
			Type: graphql.Float,
		},
		"signs": &graphql.Field{
			Type: graphql.NewList( SignObject ),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(*models.UserDetails).Signs, nil
			},
		},
		"services": &graphql.Field{
			Type: graphql.NewList( SignObject ),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				services := []models.Sign{}
				if dbc := dependencies.DBInstance().Find(&services); dbc.Error != nil {
					logrus.Error(dbc.Error)
					return nil, dbc.Error
				}
				return services , nil
			},
		},
	},
})
