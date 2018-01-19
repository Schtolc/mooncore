package graphql

import (
	"github.com/Schtolc/mooncore/dao"
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
			Type: graphql.NewNonNull(graphql.Int),
		},
		"lat": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Float),
		},
		"lon": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Float),
		},
		"description": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
})

// PhotoObject is a graphql object for photo
var PhotoObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "Photo",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"path": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"tags": &graphql.Field{
			Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(TagObject))),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(models.Photo).Tags, nil
			},
		},
	},
})

// TagObject is a graphql object for photos tags
var TagObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "Tag",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"name": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
})

// SignObject is a graphql object for photos tags
var SignObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "Sign",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"name": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"description": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"icon": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
})

//ManicureTypeObject is a graphql object for photos tags
var ManicureTypeObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "ManicureType",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"name": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
})

// UserObject is a graphql object for user
var UserObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"username": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"email": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"role": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
	},
})

// MasterObject is a graphql object for user
var MasterObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "Master",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"user": &graphql.Field{
			Type: graphql.NewNonNull(UserObject),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(models.Master).User, nil
			},
		},
		"name": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"address": &graphql.Field{
			Type: graphql.NewNonNull(AddressObject),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(models.Master).Address, nil
			},
		},
		"avatar": &graphql.Field{
			Type: graphql.NewNonNull(PhotoObject),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(models.Master).Photo, nil
			},
		},
		"photos": &graphql.Field{
			Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(PhotoObject))),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(models.Master).Photos, nil
			},
		},
		"stars": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Float),
		},
		"signs": &graphql.Field{
			Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(SignObject))),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(models.Master).Signs, nil
			},
		},
		"services": &graphql.Field{
			Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(SignObject))),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(models.Master).Services, nil
			},
		},
	},
	//TODO try to remove resolver
})

var ClientObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "Client",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"user": &graphql.Field{
			Type: graphql.NewNonNull(UserObject),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(models.Client).User, nil
			},
		},
		"name": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"avatar": &graphql.Field{
			Type: graphql.NewNonNull(PhotoObject),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(models.Client).Photo, nil
			},
		},
		"favorites": &graphql.Field{
			Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(MasterObject))),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(models.Client).Favorites, nil
			},
		},
	},
})

// ServiceObject is a graphql object for photos tags
var ServiceObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "Service",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"name": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"description": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"price": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Float),
		},
		"master": &graphql.Field{
			Type: MasterObject,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return dao.GetMasterById(p.Source.(models.Service).MasterID)
			},
		},
		"manicure": &graphql.Field{ // TODO delete query
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

var TokenObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "Token",
	Fields: graphql.Fields{
		"token": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
})
