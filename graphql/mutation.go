package graphql

import (
	"github.com/Schtolc/mooncore/dao"
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
	"strconv"
)

var createMaster = &graphql.Field{
	Type:        MasterObject,
	Description: "Create new master",
	Args: graphql.FieldConfigArgument{
		"username":    just(graphql.String),
		"email":       notNull(graphql.String),
		"password":    notNull(graphql.String),
		"name":        notNull(graphql.String),
		"lat":         notNull(graphql.String),
		"lon":         notNull(graphql.String),
		"description": notNull(graphql.String),
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		username, ok := params.Args["username"].(string)
		if !ok {
			username = ""
		}

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

		address, err := dao.CreateAddress(lat, lon, params.Args["description"].(string))

		if err != nil {
			return nil, err
		}

		return dao.CreateMaster(
			username,
			params.Args["email"].(string),
			params.Args["password"].(string),
			params.Args["name"].(string),
			address.ID)
	},
}

var createClient = &graphql.Field{
	Type:        ClientObject,
	Description: "Create new client",
	Args: graphql.FieldConfigArgument{
		"username": just(graphql.String),
		"email":    notNull(graphql.String),
		"password": notNull(graphql.String),
		"name":     notNull(graphql.String),
		"photo_id": notNull(graphql.ID),
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		username, ok := params.Args["username"].(string)
		if !ok {
			username = ""
		}
		photoID, err := strconv.ParseInt(params.Args["photo_id"].(string), 10, 64)
		if err != nil {
			return nil, err
		}
		return dao.CreateClient(
			username,
			params.Args["email"].(string),
			params.Args["password"].(string),
			params.Args["name"].(string),
			photoID)
	},
}

var signIn = &graphql.Field{
	Type:        TokenObject, // nil if user not found
	Description: "Sign in",
	Args: graphql.FieldConfigArgument{
		//"username": notNull(graphql.String),
		"email":    notNull(graphql.String),
		"password": notNull(graphql.String),
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		return dao.SignIn(
			params.Args["email"].(string),
			params.Args["password"].(string))
	},
}
