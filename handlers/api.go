package handlers

import (
	"encoding/json"
	"github.com/Schtolc/mooncore/dependencies"
	"github.com/Schtolc/mooncore/models"
	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
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
						logrus.Println(dbc.Error)
						return nil, dbc.Error
					}
					return address, nil
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
						logrus.Println(dbc.Error)
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
					if dbc := db.First(&address, params.Args["id"].(int)); dbc.Error != nil {
						logrus.Error(dbc.Error)
						return nil, dbc.Error
					}
					return address, nil
				},
			},

			"addressList": &graphql.Field{
				Type:        graphql.NewList(AddressObject),
				Description: "List of addresses",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var addresses []models.Address
					if dbc := db.Find(&addresses); dbc.Error != nil {
						logrus.Error(dbc.Error)
						return nil, dbc.Error
					}
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
					if dbc := db.First(&photo, params.Args["id"].(int)); dbc.Error != nil {
						logrus.Error(dbc.Error)
						return nil, dbc.Error
					}
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
					if dbc := db.First(&user, params.Args["id"].(int)); dbc.Error != nil {
						logrus.Error(dbc.Error)
						return nil, dbc.Error
					}
					return user, nil
				},
			},

			"usersList": &graphql.Field{
				Type:        graphql.NewList(UserObject),
				Description: "List of users",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var users []models.User
					if dbc := db.Find(&users); dbc.Error != nil {
						logrus.Error(dbc.Error)
						return nil, dbc.Error
					}
					return users, nil
				},
			},
		},
	})
}

var schema *graphql.Schema

func createSchema() graphql.Schema {
	if schema != nil {
		return *schema
	}
	_schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query:    getRootQuery(dependencies.DBInstance()),
		Mutation: getRootMutation(dependencies.DBInstance()),
	})
	schema = &_schema
	return *schema
}

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	return result
}

// API GraphQL handler
func API(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		logrus.Error(err)
		return sendResponse(c, http.StatusBadRequest, err.Error())
	}
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		logrus.Error(err)
		return sendResponse(c, http.StatusBadRequest, err.Error())
	}
	query, ok := data["query"]
	if !ok {
		strErr := "No query in request"
		logrus.Error(strErr)
		return sendResponse(c, http.StatusBadRequest, strErr)
	}
	result := executeQuery(query.(string), createSchema())

	if len(result.Errors) > 0 {
		return sendResponse(c, http.StatusNotFound, result.Errors)
	}

	return sendResponse(c, http.StatusOK, result.Data)
}
