package api

import (
	"fmt"
	"github.com/Schtolc/mooncore/models"
	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	"log"
)

func getRootMutation(db *gorm.DB) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"createAddress": &graphql.Field{
				Type:        AddressObject,
				Description: "Create new address",
				Args: graphql.FieldConfigArgument{
					"lat": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Float),
					},
					"lon": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Float),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					address := &models.Address{
						Lat: params.Args["lat"].(float64),
						Lon: params.Args["lon"].(float64),
					}
					if dbc := db.Create(address); dbc.Error != nil {
						log.Println(dbc.Error)
						return nil, dbc.Error
					}
					return address, nil
				},
			},

			"createPhoto": &graphql.Field{
				Type:        PhotoObject,
				Description: "Create new photo",
				Args: graphql.FieldConfigArgument{
					"path": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					photo := &models.Photo{
						Path: params.Args["path"].(string),
					}
					if dbc := db.Create(photo); dbc.Error != nil {
						log.Println(dbc.Error)
						return nil, dbc.Error
					}
					return photo, nil
				},
			},

			"createUser": &graphql.Field{
				Type:        UserObject,
				Description: "Create new user",
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"address_id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"photo_id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					user := &models.User{
						Name:      params.Args["name"].(string),
						Email:     params.Args["email"].(string),
						Password:  params.Args["password"].(string),
						AddressID: params.Args["address_id"].(int),
						PhotoID:   params.Args["photo_id"].(int),
					}

					if dbc := db.Create(user); dbc.Error != nil {
						log.Println(dbc.Error)
						return nil, dbc.Error
					}
					return user, nil
				},
			},
		},
	})
}

func getRootQuery(db *gorm.DB) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"address": &graphql.Field{
				Type:        AddressObject,
				Description: "Get single address",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					address := models.Address{}
					db.First(&address, params.Args["id"].(int))

					return address, nil
				},
			},

			"addressList": &graphql.Field{
				Type:        graphql.NewList(AddressObject),
				Description: "List of address",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var addresses []models.Address
					db.Find(&addresses)
					return addresses, nil
				},
			},

			"photo": &graphql.Field{
				Type:        PhotoObject,
				Description: "Get single photo",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					photo := models.Photo{}
					db.First(&photo, params.Args["id"].(int))

					return photo, nil
				},
			},

			"user": &graphql.Field{
				Type:        UserObject,
				Description: "Get single user",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					user := models.User{}
					db.First(&user, params.Args["id"].(int))

					return user, nil
				},
			},

			"usersList": &graphql.Field{
				Type:        graphql.NewList(UserObject),
				Description: "List of users",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var users []models.User
					db.Find(&users)
					return users, nil
				},
			},
		},
	})
}

//CreateSchema returns graphql schema
func CreateSchema(db *gorm.DB) graphql.Schema {
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query:    getRootQuery(db),
		Mutation: getRootMutation(db),
	})
	return schema
}

// ExecuteQuery executes graphql query
func ExecuteQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}
