package graphql

import (
	"github.com/Schtolc/mooncore/dao"
	"github.com/Schtolc/mooncore/models"
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
)

var feed = &graphql.Field{
	Type:        graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(MasterObject))),
	Description: "feed",
	Args: graphql.FieldConfigArgument{
		"offset": notNull(graphql.Int),
		"limit":  notNull(graphql.Int),
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

var viewer = &graphql.Field{
	Type:        UserObject,
	Description: "current user",
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		return p.Context.Value(UserKey), nil
	},
}
