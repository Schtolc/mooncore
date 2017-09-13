package api

import (
	"github.com/graphql-go/graphql"
)

var AddressType = graphql.NewObject(graphql.ObjectConfig{
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

var PhotoType = graphql.NewObject(graphql.ObjectConfig{
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

var UserType = graphql.NewObject(graphql.ObjectConfig{
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
			Type: AddressType,
		},
		"photo": &graphql.Field{
			Type: PhotoType,
		},
	},
})
