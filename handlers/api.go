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
	"net/http"
	//"math/rand"
	"math/rand"
	"strconv"
)



func getRootMutation(db *gorm.DB) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			//"addUserPhoto": &graphql.Field{
			//	Type:        PhotoObject,
			//	Description: "Create new photo",
			//	Args: graphql.FieldConfigArgument{
			//		"path": &graphql.ArgumentConfig{
			//			Type: graphql.NewNonNull(graphql.String),
			//		},
			//		"tags": &graphql.ArgumentConfig{
			//			Type: graphql.NewList(graphql.Int),
			//		},
			//		"email": &graphql.ArgumentConfig{
			//			Type: graphql.NewNonNull(graphql.String),
			//		},
			//	},
			//	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			//		userDetails, err := getUserDetails(params.Args["email"].(string), db)
			//		if err != nil {
			//			logrus.Warn(err)
			//			return nil, err
			//		}
			//
			//	},
			//},
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
					"lat": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"lon": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"description": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"path": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"tags": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.Int),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					userDetails, err := getUserDetails(params.Args["email"].(string), db)
					if err != nil {
						logrus.Warn(err)
						return nil, err
					}
					tx := db.Begin()
					if params.Args["name"] != nil && params.Args["name"].(string) != userDetails.Name {
						if err := tx.Model(userDetails).Where("id = ?", userDetails.ID).Update("name", params.Args["name"].(string)).Error; err != nil {
							tx.Rollback()
							return nil, err
						}
					}
					if params.Args["lat"] != nil || params.Args["lon"] != nil || params.Args["description"] != nil {
						logrus.Warn(params.Args["lat"].(string))
						lat, err := strconv.ParseFloat( params.Args["lat"].(string), 64)
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
						if err := tx.Create(address).Error; err != nil {
							tx.Rollback()
							return nil, err
						}
						userDetails.AddressID = address.ID
					}
					if params.Args["path"] != nil  {
						tags := []models.Tag{}
						if dbc := db.Where("id in (?)", params.Args["tags"]).Find(&tags); dbc.Error != nil {
							logrus.Error(dbc.Error)
							return nil, dbc.Error
						}
						if len(tags) == 0 && params.Args["tags"] != nil {
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
						userDetails.PhotoID = photo.ID
					}
					logrus.Warn(userDetails)
					if err := tx.Save(userDetails).Error; err != nil {
						tx.Rollback()
						return nil, err
					}
					tx.Commit()
					return *userDetails, nil
				},
			},
			//"addSigns": &graphql.Field{
			//	Type:        UserDetailsObject,
			//	Description: "Add signs to user",
			//	Args: graphql.FieldConfigArgument{
			//		"signs": &graphql.ArgumentConfig{
			//			Type: graphql.NewNonNull(graphql.NewList(graphql.Int)),
			//		},
			//		"email": &graphql.ArgumentConfig{
			//			Type: graphql.NewNonNull(graphql.String),
			//		},
			//	},
			//	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			//		error =
			//		return nil, error
			//	},
			//},
		},
	})
}
func getUserDetails(email string, db *gorm.DB) (*models.UserDetails, error) {
	user := &models.User{}
	userDetails := &models.UserDetails{}
	if dbc := db.Where("email in (?)",email).First(user); dbc.Error != nil {
		logrus.Error(dbc.Error)
		return nil, dbc.Error
	}
	if dbc := db.Where("user_id in (?)",user.ID).First(userDetails); dbc.Error != nil {
		logrus.Error(dbc.Error)
		return nil, dbc.Error
	}
	return userDetails, nil
}

func createAddress(db *gorm.DB, latStr string, lonStr string, description string) (*models.Address, error){
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	address := &models.Address{
		Lat:         lat,
		Lon:         lon,
		Description: description,
	}
	if dbc := db.Create(address); dbc.Error != nil {
		logrus.Error(dbc.Error)
		return nil, dbc.Error
	}
	return address, nil
}
func createPhoto (path string, tags []models.Tag, db *gorm.DB) (*models.Photo, error) {
	photo := &models.Photo{
		Path: path,
		Tags: []models.Tag{},
	}
	if dbc := db.Where("id in (?)", tags).Find(&photo.Tags); dbc.Error != nil {
		logrus.Error(dbc.Error)
		return nil, dbc.Error
	}
	if len(photo.Tags) == 0 && tags != nil {
		logrus.Error("No such tags")
		return nil, errors.New("No such tags")
	}

	if dbc := db.Create(photo); dbc.Error != nil {
		logrus.Error(dbc.Error)
		return nil, dbc.Error
	}
	return photo, nil
}



//func getRootMutation(db *gorm.DB) *graphql.Object {
//	return graphql.NewObject(graphql.ObjectConfig{
//		Name: "RootMutation",
//		Fields: graphql.Fields{
//			"addUserPhoto": &graphql.Field{
//				Type:        PhotoObject,
//				Description: "Create new photo",
//				Args: graphql.FieldConfigArgument{
//					"path": &graphql.ArgumentConfig{
//						Type: graphql.NewNonNull(graphql.String),
//					},
//					"tags": &graphql.ArgumentConfig{
//						Type: graphql.NewList(graphql.Int),
//					},
//					"email": &graphql.ArgumentConfig{
//						Type: graphql.NewNonNull(graphql.String),
//					},
//				},
//				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
//					userDetails, err := getUserDetails(params.Args["email"].(string), db)
//					if err != nil {
//						logrus.Warn(err)
//						return nil, err
//					}
//					photo := &models.Photo{
//						Path: params.Args["path"].(string),
//						Tags: []models.Tag{},
//					}
//					if dbc := db.Where("id in (?)", params.Args["tags"]).Find(&photo.Tags); dbc.Error != nil {
//						logrus.Error(dbc.Error)
//						return nil, dbc.Error
//					}
//					if len(photo.Tags) == 0 && params.Args["tags"] != nil {
//						logrus.Error("No such tags")
//						return nil, errors.New("No such tags")
//					}
//
//					if dbc := db.Create(photo); dbc.Error != nil {
//						logrus.Error(dbc.Error)
//						return nil, dbc.Error
//					}
//					if dbc := db.Model(&userDetails).Association("photos").Append([]models.Photo{*photo}); dbc.Error != nil {
//						logrus.Error(dbc.Error)
//						return nil, dbc.Error
//					}
//					return photo, nil
//				},
//			},
//			"createUser": &graphql.Field{
//				Type:        UserObject,
//				Description: "Create new user",
//				Args: graphql.FieldConfigArgument{
//					"email": &graphql.ArgumentConfig{
//						Type: graphql.NewNonNull(graphql.String),
//					},
//					"password": &graphql.ArgumentConfig{
//						Type: graphql.NewNonNull(graphql.String),
//					},
//					"role": &graphql.ArgumentConfig{
//						Type: graphql.NewNonNull(graphql.Int),
//					},
//				},
//				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
//					tx := db.Begin()
//					user := &models.User{
//						Email:    params.Args["email"].(string),
//						Password: params.Args["password"].(string),
//						Role:     params.Args["role"].(int),
//					}
//					if err := tx.Create(user).Error; err != nil {
//						tx.Rollback()
//						return nil, err
//					}
//
//					userDetails := &models.UserDetails{
//						UserID:    user.ID,
//						AddressID: models.DefaultAddress.ID,
//						PhotoID:   models.DefaultAvatar.ID,
//					}
//
//					if err := tx.Create(userDetails).Error; err != nil {
//						tx.Rollback()
//						return nil, err
//					}
//					tx.Commit()
//					return user, nil
//				},
//			},
//			"editUserProfile": &graphql.Field{
//				Type:        UserDetailsObject,
//				Description: "Edit user fields",
//				Args: graphql.FieldConfigArgument{
//					"email": &graphql.ArgumentConfig{
//						Type: graphql.String,
//					},
//					"name": &graphql.ArgumentConfig{
//						Type: graphql.String,
//					},
//					"address": &graphql.ArgumentConfig{
//						Type: AddressObject,
//					},
//					"photo": &graphql.ArgumentConfig{
//						Type: PhotoObject,
//					},
//					"photos": &graphql.ArgumentConfig{
//						Type: graphql.NewList(graphql.Int),
//					},
//				},
//				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
//					userDetails, err := getUserDetails(params.Args["email"].(string), db)
//					if err != nil {
//						logrus.Warn(err)
//						return nil, err
//					}
//					tx := db.Begin()
//					if params.Args["name"] != nil {
//						if dbc := db.Model(userDetails).Where("id = ?", userDetails.ID).Update("name", params.Args["name"].(string)); dbc.Error != nil {
//							logrus.Error(dbc.Error)
//							return nil, dbc.Error
//						}
//					}
//					//if params.Args["address"] != (models.Address{}) {
//					//	address := &models.Address{}
//					//	if dbc := db.Where("id = ?", params.Args["address_id"].(int)).First(address); dbc.Error != nil {
//					//		logrus.Error(dbc.Error)
//					//		return nil, dbc.Error
//					//	}
//					//	if dbc := db.Model(user).Where("id = ?", userDetails.ID).Update("address_id", params.Args["address_id"].(int)); dbc.Error != nil {
//					//		logrus.Error(dbc.Error)
//					//		return nil, dbc.Error
//					//	}
//					//}
//					if params.Args["photo"] != nil {
//						photo := createPhoto(params.Args[photo])
//						photo := &models.Photo{}
//						if dbc := db.Where("id = ?", params.Args["avatar_id"].(int)).First(photo); dbc.Error != nil {
//							logrus.Error(dbc.Error)
//							return nil, dbc.Error
//						}
//						if dbc := db.Model(user).Where("id = ?", userAuth.ID).Update("avatar_id", params.Args["avatar_id"].(int)); dbc.Error != nil {
//							logrus.Error(dbc.Error)
//							return nil, dbc.Error
//						}
//					}
//
//					if params.Args["photos"] != nil {
//						if dbc := db.Where("id in (?)", params.Args["photos"]).Find(&user.Photos); dbc.Error != nil {
//							logrus.Error(dbc.Error)
//							return nil, dbc.Error
//						}
//					}
//
//					tx.Commit()
//					return user, nil
//				},
//			},
//			"addSigns": &graphql.Field{
//				Type:        UserDetailsObject,
//				Description: "Add signs to user",
//				Args: graphql.FieldConfigArgument{
//					"signs": &graphql.ArgumentConfig{
//						Type: graphql.NewNonNull(graphql.NewList(graphql.Int)),
//					},
//					"email": &graphql.ArgumentConfig{
//						Type: graphql.NewNonNull(graphql.String),
//					},
//				},
//				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
//					// return if not admin
//					userAuth := &models.User{}
//					user := &models.UserDetails{}
//					signs := []models.Sign{}
//
//					if dbc := db.Where("email = ?", params.Args["email"].(string)).First(userAuth); dbc.Error != nil {
//						logrus.Error(dbc.Error)
//						return nil, dbc.Error
//					}
//
//					if dbc := db.Where("user_id = ?", userAuth.ID).First(user); dbc.Error != nil {
//						logrus.Error(dbc.Error)
//						return nil, dbc.Error
//					}
//
//					if dbc := db.Where("id in (?)", params.Args["signs"]).Find(&signs); dbc.Error != nil {
//						logrus.Error(dbc.Error)
//						return nil, dbc.Error
//					}
//
//					if len(signs) == 0 {
//						return nil, errors.New("No such signs")
//					}
//					if dbc := db.Model(user).Where("id = ?", user.ID).Update("signs", signs); dbc.Error != nil {
//						logrus.Error(dbc.Error)
//						return nil, dbc.Error
//					}
//
//					return user, nil
//				},
//			},
//		},
//	})
//}

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
						logrus.Warn(err)
						return nil, err
					}
					return *userDetails, nil
				},
			},

			"listUser": &graphql.Field{
				Type:        graphql.NewList(UserDetailsObject),
				Description: "List of users",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var users []models.UserDetails
					if dbc := db.Find(&users); dbc.Error != nil {
						logrus.Error(dbc.Error)
						return nil, dbc.Error
					}
					logrus.Warn(users)
					return users, nil
				},
			},
			"getFeed": &graphql.Field{
				Type:	graphql.NewList(FeedElemObject),
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

					var users []models.UserDetails
					if dbc := db.Limit(p.Args["limit"].(int)).Offset(offset).Find(&users); dbc.Error != nil {
						logrus.Error(dbc.Error)
						return nil, dbc.Error
					}

					for _, user := range users {
						if dbc := db.Where("id = ?", user.AddressID).First(&user.Address); dbc.Error != nil {
							logrus.Error(dbc.Error)
							return nil, dbc.Error
						}
						if dbc := db.Where("id = ?", user.PhotoID).First(&user.Photo); dbc.Error != nil {
							logrus.Error(dbc.Error)
							return nil, dbc.Error
						}
						if dbc := db.Model(&user).Related(&user.Photos, "user_photos");dbc.Error != nil {
							logrus.Error(dbc.Error)
							return nil, dbc.Error
						}
						logrus.Warn(user.Photos)
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
	logrus.Warn("executeQuery")
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	logrus.Warn(result)
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
	logrus.Warn(query)
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

//"addUserPhoto": &graphql.Field{
//	Type:        PhotoObject,
//	Description: "Create new photo",
//	Args: graphql.FieldConfigArgument{
//		"path": &graphql.ArgumentConfig{
//			Type: graphql.NewNonNull(graphql.String),
//		},
//		"tags": &graphql.ArgumentConfig{
//			Type: graphql.NewList(graphql.Int),
//		},
//	},
//	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
//		tags := []models.Tag{}
//		if dbc := db.Where("id in (?)", params.Args["tags"]).Find(&tags); dbc.Error != nil {
//			logrus.Error(dbc.Error)
//			return nil, dbc.Error
//		}
//		if len(tags) == 0 && params.Args["tags"] != nil {
//			return nil, errors.New("No such tags")
//		}
//
//		photo := &models.Photo{
//			Path: params.Args["path"].(string),
//			Tags: tags,
//		}
//		if dbc := db.Create(photo); dbc.Error != nil {
//			logrus.Error(dbc.Error)
//			return nil, dbc.Error
//		}
//		return photo, nil
//	},
//},

//"createUserAddress": &graphql.Field{
//		Type:        AddressObject,
//		Description: "Create new address",
//		Args: graphql.FieldConfigArgument{
//			"lat": &graphql.ArgumentConfig{
//				Type: graphql.NewNonNull(graphql.String),
//			},
//			"lon": &graphql.ArgumentConfig{
//				Type: graphql.NewNonNull(graphql.String),
//			},
//			"description": &graphql.ArgumentConfig{
//				Type: graphql.NewNonNull(graphql.String),
//			},
//		},
//		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
//			lat, err := strconv.ParseFloat(params.Args["lat"].(string), 64)
//			if err != nil {
//				logrus.Error(err)
//				return nil, err
//			}
//			lon, _ := strconv.ParseFloat(params.Args["lon"].(string), 64)
//			if err != nil {
//				logrus.Error(err)
//				return nil, err
//			}
//			address := &models.Address{
//				Lat:         lat,
//				Lon:         lon,
//				Description: params.Args["description"].(string),
//			}
//			if dbc := db.Create(address); dbc.Error != nil {
//				logrus.Error(dbc.Error)
//				return nil, dbc.Error
//			}
//			return address, nil
//		},
//},