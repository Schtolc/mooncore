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
	"errors"
	"strconv"
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
						Type: graphql.NewNonNull(graphql.String),
					},
					"lon": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"description": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					lat, err := strconv.ParseFloat(params.Args["lat"].(string), 64)
					if err != nil {
						logrus.Error(err)
						return nil, errors.New("InvalidParam: lat")
					}
					lon, _ := strconv.ParseFloat(params.Args["lon"].(string), 64)
					if err != nil {
						logrus.Error(err)
						return nil, errors.New("InvalidParam: lon")
					}
					address := &models.Address{
						Lat: lat,
						Lon: lon,
						Description: params.Args["description"].(string),
					}
					if dbc := db.Create(address); dbc.Error != nil {
						logrus.Error(dbc.Error)
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
					"tags": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.NewList( graphql.Int )),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					tags := []models.Tag{}
					if dbc := db.Where("id in (?)", params.Args["tags"]).Find(&tags); dbc.Error != nil {
						logrus.Error(dbc.Error)
						return nil, dbc.Error
					}
					if len(tags) == 0 {
						return nil, errors.New("No such tags")
					}

					photo := &models.Photo{
						Path: params.Args["path"].(string),
						Tags: tags,
					}
					if dbc := db.Create(photo); dbc.Error != nil {
						logrus.Error(dbc.Error)
						return nil, dbc.Error
					}
					return photo, nil
				},
			},
			"createUser": &graphql.Field{
				Type:        UserObject,
				Description: "Create new user",
				Args: graphql.FieldConfigArgument{
					"email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"role": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					tx := db.Begin()
						user := &models.User{
							Email:  params.Args["email"].(string),
							Password:  params.Args["password"].(string),
							Role:  params.Args["role"].(int),
						}
						if err := tx.Create(user).Error; err != nil {
							tx.Rollback()
							return nil, err
						}

						userDetails := &models.UserDetails{
							UserID : user.ID,
							AddressID: models.DefaultAddress.ID,
							PhotoID: models.DefaultAvatar.ID,
						}

						if err := tx.Create(userDetails).Error; err != nil {
							tx.Rollback()
							return nil, err
						}
					tx.Commit()
					return user, nil
				},
			},
			"createUserProfile": &graphql.Field{
				Type:        UserDetailsObject,
				Description: "Edit user fields",
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"address_id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"avatar_id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"email": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"photos": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.NewList( graphql.Int )),
					},
					"signs": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.NewList( graphql.Int )),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					logrus.Warn(params)
					userAuth := &models.User{}
					user     := &models.UserDetails{}


					if dbc := db.Where("email = ?", params.Args["email"].(string)).First(userAuth); dbc.Error != nil {
						logrus.Error(dbc.Error)
						return nil, dbc.Error
					}
					if dbc := db.Where("user_id = ?",userAuth.ID).First(user); dbc.Error != nil {
						logrus.Error(dbc.Error)
						return nil, dbc.Error
					}
					tx := db.Begin()
						if params.Args["name"] != nil {
							if dbc := db.Model(user).Where("id = ?", userAuth.ID).Update("name", params.Args["name"].(string) ); dbc.Error != nil {
								logrus.Error(dbc.Error)
								return nil, dbc.Error
							}
						}
						if params.Args["address_id"] != nil {
							address := &models.Address{}
							if dbc := db.Where("id = ?", params.Args["address_id"].(int)).First(address); dbc.Error != nil {
								logrus.Error(dbc.Error)
								return nil, dbc.Error
							}
							if dbc := db.Model(user).Where("id = ?", userAuth.ID).Update("address_id", params.Args["address_id"].(int) ); dbc.Error != nil {
								logrus.Error(dbc.Error)
								return nil, dbc.Error
							}
						}
						if params.Args["avatar_id"] != nil {
							photo := &models.Photo{}
							if dbc := db.Where("id = ?", params.Args["avatar_id"].(int)).First(photo); dbc.Error != nil {
								logrus.Error(dbc.Error)
								return nil, dbc.Error
							}
							if dbc := db.Model(user).Where("id = ?", userAuth.ID).Update("avatar_id", params.Args["avatar_id"].(int) ); dbc.Error != nil {
								logrus.Error(dbc.Error)
								return nil, dbc.Error
							}
						}
						logrus.Warn(params.Args["photos"])
						if params.Args["photos"] != nil {
							photos := []models.Photo{}
							if dbc := db.Where("id in (?)", params.Args["photos"]).Find(&photos);dbc.Error != nil {
								logrus.Error(dbc.Error)
								return nil, dbc.Error
							}
							user.Photos = photos
						}
						logrus.Warn(params.Args["signs"])
						if params.Args["signs"] != nil {
							signs := []models.Sign{}
							if dbc := db.Where("id in (?)", params.Args["signs"]).Find(&signs);dbc.Error != nil {
								logrus.Error(dbc.Error)
								return nil, dbc.Error
							}
							user.Signs = signs
						}
					logrus.Warn(user)
					tx.Commit()
					return user, nil
				},
			},
			"addSigns": &graphql.Field{
				Type:        UserDetailsObject,
				Description: "Edit user fields",
				Args: graphql.FieldConfigArgument{
					"signs": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.NewList( graphql.Int )),
					},
					"email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// return if not admin
					userAuth := &models.User{}
					user     := &models.UserDetails{}
					signs    := []models.Sign{}
					logrus.Warn(params.Args["email"].(string))
					if dbc := db.Where("email = ?",params.Args["email"].(string)).First(userAuth); dbc.Error != nil {
						logrus.Error(dbc.Error)
						return nil, dbc.Error
					}
					logrus.Warn(userAuth.ID)
					if dbc := db.Where("user_id = ?",userAuth.ID).First(user); dbc.Error != nil {
						logrus.Error(dbc.Error)
						return nil, dbc.Error
					}
					logrus.Warn(user.ID)
					if dbc := db.Where("id in (?)", params.Args["signs"]).Find(&signs);dbc.Error != nil {
						logrus.Error(dbc.Error)
						return nil, dbc.Error
					}
					logrus.Warn(signs)
					if len(signs) == 0 {
						return nil, errors.New("No such signs")
					}
					if dbc := db.Model(user).Where("id = ?", user.ID).Update("signs", signs); dbc.Error != nil {
						logrus.Error(dbc.Error)
						return nil, dbc.Error
					}
					logrus.Warn(user)
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
					if dbc := db.Model(&photo).Related(&photo.Tags, "Tags"); dbc.Error != nil {
						logrus.Error(dbc.Error)
						return nil, dbc.Error
					}
					return &photo, nil
				},
			},

			"user": &graphql.Field{
				Type:        UserDetailsObject,
				Description: "Get single user",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					user := models.UserDetails{}
					if dbc := db.First(&user, params.Args["id"].(int)); dbc.Error != nil {
						logrus.Error(dbc.Error)
						return nil, dbc.Error
					}
					return user, nil
				},
			},

			"usersList": &graphql.Field{
				Type:        graphql.NewList(UserDetailsObject),
				Description: "List of users",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var users []models.UserDetails
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
<<<<<<< Updated upstream
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

=======
	logrus.Warn(c.QueryParams().Get("query"))
	result := executeQuery(c.QueryParams().Get("query"), createSchema())
	response := Response{}
>>>>>>> Stashed changes
	if len(result.Errors) > 0 {
		return sendResponse(c, http.StatusNotFound, result.Errors)
	}

	return sendResponse(c, http.StatusOK, result.Data)
}
