package handlers

import (
	"encoding/json"
	"errors"
	"github.com/Schtolc/mooncore/dependencies"
	"github.com/Schtolc/mooncore/models"
	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
)

func getUserDetails(email string, db *gorm.DB) (*models.UserDetails, error) {
	user := &models.User{}
	userDetails := &models.UserDetails{}
	if dbc := db.Where("email in (?)", email).First(user); dbc.Error != nil {
		logrus.Error(dbc.Error)
		return nil, dbc.Error
	}
	if dbc := db.Where("user_id in (?)", user.ID).First(userDetails); dbc.Error != nil {
		logrus.Error(dbc.Error)
		return nil, dbc.Error
	}
	return userDetails, nil
}

func getRootMutation(db *gorm.DB) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
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
						Email:    params.Args["email"].(string),
						Password: params.Args["password"].(string),
						Role:     params.Args["role"].(int),
					}
					if err := tx.Create(user).Error; err != nil {
						tx.Rollback()
						return nil, err
					}
					userDetails := &models.UserDetails{
						UserID:    user.ID,
						AddressID: models.DefaultAddress.ID,
						PhotoID:   models.DefaultAvatar.ID,
					}
					if err := tx.Create(userDetails).Error; err != nil {
						tx.Rollback()
						return nil, err
					}
					tx.Commit()
					return user, nil
				},
			},
			"editUserProfile": &graphql.Field{
				Type:        UserDetailsObject,
				Description: "Edit user fields",
				Args: graphql.FieldConfigArgument{
					"email": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"address_id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"avatar_id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					userDetails, err := getUserDetails(params.Args["email"].(string), db)
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
			},
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
			},
			"createAvatar": &graphql.Field{
				Type:        PhotoObject,
				Description: "Create new address",
				Args: graphql.FieldConfigArgument{
					"path": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"tags": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.Int),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
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
					return *photo, nil
				},
			},
			"addPhoto": &graphql.Field{
				Type:        PhotoObject,
				Description: "Add new photo - master",
				Args: graphql.FieldConfigArgument{
					"path": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"tags": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.Int),
					},
					"email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					userDetails, err := getUserDetails(params.Args["email"].(string), db)
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
			},
			"addSigns": &graphql.Field{
				Type:        UserDetailsObject,
				Description: "Add signs -  admin",
				Args: graphql.FieldConfigArgument{
					"signs": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.NewList(graphql.Int)),
					},
					"email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					userDetails, err := getUserDetails(params.Args["email"].(string), db)
					if err != nil {
						logrus.Warn(err)
						return nil, err
					}
					signs := []models.Sign{}
					if dbc := db.Where("id in (?)", params.Args["signs"]).Find(&signs); dbc.Error != nil {
						logrus.Error(dbc.Error)
						return nil, dbc.Error
					}
					if len(signs) == 0 {
						return nil, errors.New("No such signs")
					}
					tx := db.Begin()
					for _, sign := range signs {
						if err := tx.Model(&userDetails).Association("signs").Append([]models.Sign{sign}).Error; err != nil {
							tx.Rollback()
							return nil, err
						}
					}
					tx.Commit()
					return *userDetails, nil
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
			"addressListInArea": &graphql.Field{
				Type:        graphql.NewList(AddressObject),
				Description: "Get all addresses in this area",
				Args: graphql.FieldConfigArgument{
					"lat1": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Float),
					},
					"lon1": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Float),
					},
					"lat2": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Float),
					},
					"lon2": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Float),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					lat1 := params.Args["lat1"]	// первая точка сверху слева
					lon1 := params.Args["lon1"]	// вторая снизу и справа
					lat2 := params.Args["lat2"]
					lon2 := params.Args["lon2"]
					var address []models.Address
					query := "lat > ? AND lat < ? AND lon < ? AND lon > ?"
					if dbc := db.Where(query, lat1, lat2, lon1, lon2).Find(&address); dbc.Error != nil {
						logrus.Error(dbc.Error)
						return nil, dbc.Error
					}
					return address, nil
				},
			},
			"getPhoto": &graphql.Field{
				Type:        PhotoObject,
				Description: "Get single photo", // done
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					photo := &models.Photo{}

					if dbc := db.First(&photo, params.Args["id"].(int)); dbc.Error != nil {
						logrus.Error(dbc.Error)
						return nil, dbc.Error
					}
					if dbc := db.Model(&photo).Related(&photo.Tags, "Tags"); dbc.Error != nil {
						logrus.Error(dbc.Error)
						return nil, dbc.Error
					}
					return *photo, nil
				},
			},
			"getUserPhotos": &graphql.Field{
				Type:        graphql.NewList(PhotoObject),
				Description: "Get single photo", //done
				Args: graphql.FieldConfigArgument{
					"email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					userDetails, err := getUserDetails(params.Args["email"].(string), db)
					if err != nil {
						logrus.Warn(err)
						return nil, err
					}
					db.Model(&userDetails).Association("photos").Find(&userDetails.Photos)
					return userDetails.Photos, nil
				},
			},
			// done
			"getUser": &graphql.Field{
				Type:        UserDetailsObject,
				Description: "Get single user",
				Args: graphql.FieldConfigArgument{
					"email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					userDetails, err := getUserDetails(params.Args["email"].(string), db)
					if err != nil {
						logrus.Error(err)
						return nil, err
					}
					return *userDetails, nil
				},
			},
			// done
			"listUsers": &graphql.Field{
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
			"getSigns": &graphql.Field{
				Type:        graphql.NewList(SignObject),
				Description: "Get user signs", //done
				Args: graphql.FieldConfigArgument{
					"email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					userDetails, err := getUserDetails(params.Args["email"].(string), db)
					if err != nil {
						logrus.Warn(err)
						return nil, err
					}
					signs := []models.Sign{}
					db.Model(&userDetails).Association("signs").Find(&signs)
					return signs, nil
				},
			},
			"getFeed": &graphql.Field{
				Type:        graphql.NewList(UserDetailsObject),
				Description: "feed",
				Args: graphql.FieldConfigArgument{
					"limit": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var offset int
					db.Table("user_details").Count(&offset)
					offset = rand.Intn(offset)
					users := []models.UserDetails{}

					if dbc := db.Limit(p.Args["limit"].(int)).Offset(0).Find(&users); dbc.Error != nil {
						logrus.Error(dbc.Error)
						return nil, dbc.Error
					}

					for i, user := range users {
						if dbc := db.Model(&user).Association("Photos").Find(&user.Photos); dbc.Error != nil {
							logrus.Error(dbc.Error)
							return nil, dbc.Error
						}
						for index, photo := range user.Photos {
							if dbc := db.Model(&photo).Related(&photo.Tags, "Tags"); dbc.Error != nil {
								logrus.Error(dbc.Error)
								return nil, dbc.Error
							}
							user.Photos[index] = photo
						}
						if dbc := db.Model(&user).Association("Signs").Find(&user.Signs); dbc.Error != nil {
							logrus.Error(dbc.Error)
							return nil, dbc.Error
						}
						users[i] = user
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
