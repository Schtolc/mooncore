package graphql

import (
	// "github.com/Schtolc/mooncore/dao"
	"github.com/Schtolc/mooncore/dao"
	"github.com/Schtolc/mooncore/models"
	"github.com/graphql-go/graphql"
)

var AddressMetroObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "AddressMetroObject",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"line": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"color": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"distance": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Float),
		},
	},
})

// AddressObject is a graphql object for address
var AddressObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "Address",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.ID),
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
		"stations": &graphql.Field{
			Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(AddressMetroObject))),
		},
	},
})

// PhotoObject is a graphql object for photo
var PhotoObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "Photo",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.ID),
		},
		"path": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"tags": &graphql.Field{
			Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(TagObject))),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return dao.PhotoTags(p.Source.(*models.Photo))
			},
		},
	},
})

// TagObject is a graphql object for photo tag
var TagObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "Tag",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.ID),
		},
		"name": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
})

// SignObject is a graphql object for master sign
var SignObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "Sign",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.ID),
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

// UserObject is a graphql object for user
var UserObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.ID),
		},
		"username": &graphql.Field{
			Type: graphql.String,
		},
		"email": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"role": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
	},
})

// MasterObject is a graphql object for master
var MasterObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "Master",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.ID),
		},
		"user": &graphql.Field{
			Type: graphql.NewNonNull(UserObject),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return dao.GetUserByID(p.Source.(*models.Master).UserID)
			},
		},
		"name": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"address": &graphql.Field{
			Type: graphql.NewNonNull(AddressObject),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return dao.GetAddressByID(p.Source.(*models.Master).AddressID)
			},
		},
		"avatar": &graphql.Field{
			Type: PhotoObject,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if !p.Source.(*models.Master).PhotoID.Valid {
					return nil, nil
				}
				return dao.GetPhotoByID(p.Source.(*models.Master).PhotoID.Int64)
			},
		},
		"photos": &graphql.Field{
			Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(PhotoObject))),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return dao.MasterPhotos(p.Source.(*models.Master))
			},
		},
		"stars": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Float),
		},
		"signs": &graphql.Field{
			Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(SignObject))),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return dao.MasterSigns(p.Source.(*models.Master))
			},
		},
		"services": &graphql.Field{
			Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(ServiceObject))),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return dao.MasterServices(p.Source.(*models.Master))
			},
		},
	},
})

// ClientObject is a graphql object for client
var ClientObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "Client",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.ID),
		},
		"user": &graphql.Field{
			Type: graphql.NewNonNull(UserObject),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return dao.GetUserByID(p.Source.(*models.Client).UserID)
			},
		},
		"name": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"avatar": &graphql.Field{
			Type: graphql.NewNonNull(PhotoObject),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return dao.GetPhotoByID(p.Source.(*models.Client).PhotoID)
			},
		},
		"favorites": &graphql.Field{
			Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(MasterObject))),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return dao.ClientFavorites(p.Source.(*models.Client))
			},
		},
	},
})

// ServiceObject is a graphql object for service
var ServiceObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "Service",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.ID),
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
	},
})

// TokenObject is a graphql object for jwt token
var TokenObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "Token",
	Fields: graphql.Fields{
		"token": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
})
