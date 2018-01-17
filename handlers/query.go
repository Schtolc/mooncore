package handlers

import (
	"github.com/Schtolc/mooncore/models"
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
)

var address = &graphql.Field{
	Type:        AddressObject,
	Description: "Get single address",
	Args: graphql.FieldConfigArgument{
		"id": notNull(graphql.Int),
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		address := models.Address{}
		if dbc := db.First(&address, params.Args["id"].(int)); dbc.Error != nil {
			logrus.Error(dbc.Error)
			return nil, dbc.Error
		}
		return address, nil
	},
}

var addressList = &graphql.Field{
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
}

var addressListInArea = &graphql.Field{
	Type:        graphql.NewList(AddressObject),
	Description: "Get all addresses in this area",
	Args: graphql.FieldConfigArgument{
		"lat1": notNull(graphql.String),
		"lon1": notNull(graphql.String),
		"lat2": notNull(graphql.String),
		"lon2": notNull(graphql.String),
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		lat1, err := strconv.ParseFloat(params.Args["lat1"].(string), 64)
		if err != nil {
			logrus.Error(err) // первая точка сверху слева
			return nil, err   // вторая снизу и справа
		}
		lon1, err := strconv.ParseFloat(params.Args["lon1"].(string), 64)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		lat2, err := strconv.ParseFloat(params.Args["lat2"].(string), 64)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		lon2, err := strconv.ParseFloat(params.Args["lon2"].(string), 64)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		var addresses []models.Address
		query := "lat > ? AND lat < ? AND lon < ? AND lon > ?"
		if dbc := db.Where(query, lat1, lat2, lon1, lon2).Find(&addresses); dbc.Error != nil {
			logrus.Error(dbc.Error)
			return nil, dbc.Error
		}
		return addresses, nil
	},
}

var getUserPhotos = &graphql.Field{
	Type:        graphql.NewList(PhotoObject),
	Description: "Get single photo", //done
	Args: graphql.FieldConfigArgument{
		"email": notNull(graphql.String),
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		userDetails, err := getUserDetails("email", params.Args["email"].(string))
		if err != nil {
			logrus.Warn(err)
			return nil, err
		}
		db.Model(&userDetails).Association("photos").Find(&userDetails.Photos)
		return userDetails.Photos, nil
	},
}

var getUser = &graphql.Field{
	Type:        UserDetailsObject,
	Description: "Get single user",
	Args: graphql.FieldConfigArgument{
		"id": notNull(graphql.Int),
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		userDetails, err := getUserDetails("id", params.Args["id"].(int))
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		return *userDetails, nil
	},
}

var listUsers = &graphql.Field{
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
}

var getSigns = &graphql.Field{
	Type:        graphql.NewList(SignObject),
	Description: "Get user signs", //done
	Args: graphql.FieldConfigArgument{
		"email": notNull(graphql.String),
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		userDetails, err := getUserDetails("email", params.Args["email"].(string))
		if err != nil {
			logrus.Warn(err)
			return nil, err
		}
		var signs []models.Sign
		db.Model(&userDetails).Association("signs").Find(&signs)
		return signs, nil
	},
}

var feed = &graphql.Field{
	Type:        graphql.NewList(UserDetailsObject),
	Description: "feed",
	Args: graphql.FieldConfigArgument{
		"limit": notNull(graphql.Int),
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		var userCount int
		db.Table("user_details").Count(&userCount)
		offset := 1
		if userCount-p.Args["limit"].(int) > 0 {
			offset = userCount - p.Args["limit"].(int)
		}
		offset = rand.Intn(offset)
		var users []models.UserDetails
		if dbc := db.Limit(p.Args["limit"].(int)).Offset(offset).Find(&users); dbc.Error != nil {
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
}

var getPhoto = &graphql.Field{
	Type:        PhotoObject,
	Description: "Get single photo", // done
	Args: graphql.FieldConfigArgument{
		"id": notNull(graphql.Int),
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
}
