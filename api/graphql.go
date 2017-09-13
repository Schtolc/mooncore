package api

import (
	"fmt"
	"github.com/Schtolc/mooncore/models"
	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	"log"
	"strconv"
)

func getRootMutation(db *gorm.DB) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"createAddress": &graphql.Field{
				Type:        AddressType,
				Description: "Create new address",
				Args: graphql.FieldConfigArgument{
					"lat": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"lon": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {

					slat, _ := params.Args["lat"].(string)
					slon, _ := params.Args["lon"].(string)

					lat, _ := strconv.ParseFloat(slat, 32)
					lon, _ := strconv.ParseFloat(slon, 32)

					println(slat, slon)

					address := &models.Address{
						Lat: float32(lat),
						Lon: float32(lon),
					}
					if dbc := db.Create(address); dbc.Error != nil {
						log.Println(dbc.Error)
						return models.Address{}, dbc.Error
					}
					return address, nil
				},
			},

			"createPhoto": &graphql.Field{
				Type:        PhotoType,
				Description: "Create new photo",
				Args: graphql.FieldConfigArgument{
					"path": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {

					path, _ := params.Args["path"].(string)

					photo := &models.Photo{
						Path: path,
					}
					if dbc := db.Create(photo); dbc.Error != nil {
						log.Println(dbc.Error)
						return models.Address{}, dbc.Error
					}
					return photo, nil
				},
			},

			"createUser": &graphql.Field{
				Type:        UserType,
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
						Type: graphql.NewNonNull(graphql.String),
					},
					"photo_id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {

					name, _ := params.Args["name"].(string)
					password, _ := params.Args["password"].(string)
					email, _ := params.Args["email"].(string)
					saddress_id, _ := params.Args["address_id"].(string)
					sphoto_id, _ := params.Args["photo_id"].(string)

					address_id, _ := strconv.ParseInt(saddress_id, 10, 32)
					photo_id, _ := strconv.ParseInt(sphoto_id, 10, 32)

					user := &models.User{
						Name:      name,
						Email:     email,
						Password:  password,
						AddressID: int(address_id),
						PhotoID:   int(photo_id),
					}

					if dbc := db.Create(user); dbc.Error != nil {
						log.Println(dbc.Error)
						return models.Address{}, dbc.Error
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
				Type:        AddressType,
				Description: "Get single address",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {

					idQuery, _ := params.Args["id"].(int)

					address := models.Address{}
					db.First(&address, idQuery)

					return address, nil
				},
			},

			"addressList": &graphql.Field{
				Type:        graphql.NewList(AddressType),
				Description: "List of address",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var addresses []models.Address
					db.Find(&addresses)

					fmt.Println(addresses)
					return addresses, nil
				},
			},

			"photo": &graphql.Field{
				Type:        PhotoType,
				Description: "Get single photo",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {

					idQuery, _ := params.Args["id"].(int)

					photo := models.Photo{}
					db.First(&photo, idQuery)

					return photo, nil
				},
			},

			"user": &graphql.Field{
				Type:        UserType,
				Description: "Get single user",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {

					idQuery, _ := params.Args["id"].(int)

					user := models.User{}
					db.First(&user, idQuery)

					return user, nil
				},
			},
		},
	})
}

func CreateSchema(db *gorm.DB) graphql.Schema {
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query:    getRootQuery(db),
		Mutation: getRootMutation(db),
	})
	return schema
}

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
