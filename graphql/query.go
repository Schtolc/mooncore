package graphql

import (
	"github.com/Schtolc/mooncore/dao"

	"github.com/graphql-go/graphql"
	"strconv"

	"errors"
	"github.com/Schtolc/mooncore/dependencies"
	"github.com/Schtolc/mooncore/models"
	"github.com/Schtolc/mooncore/utils"
)

var master = &graphql.Field{
	Type:        MasterObject, // == nil if not found
	Description: "Get master by id",
	Args: graphql.FieldConfigArgument{
		"id": notNull(graphql.ID),
	},
	Resolve: resolveMiddleware(func(params graphql.ResolveParams) (interface{}, error) {
		id, err := strconv.ParseInt(params.Args["id"].(string), 10, 64)
		if err != nil {
			return nil, err
		}
		return dao.GetMasterByID(id)
	}),
}

var feed = &graphql.Field{
	Type:        graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(MasterObject))),
	Description: "feed",
	Args: graphql.FieldConfigArgument{
		"offset": notNull(graphql.Int),
		"limit":  notNull(graphql.Int),
	},
	Resolve: resolveMiddleware(func(p graphql.ResolveParams) (interface{}, error) {
		return dao.Feed(p.Args["offset"].(int), p.Args["limit"].(int))
	}),
}

var viewer = &graphql.Field{
	Type: UserType,
	Resolve: resolveMiddleware(func(params graphql.ResolveParams) (interface{}, error) {
		user := params.Context.Value(utils.GraphQLContextUserKey).(*models.User)
		if user.Role == models.MasterRole {
			master := models.Master{}
			db.Where("user_id = ?", user.ID).First(&master)
			return &master, nil
		} else if user.Role == models.ClientRole {
			client := models.Client{}
			db.First(&client).Where(client.UserID == user.ID)
			return &client, nil
		} else {
			err := errors.New("this Role is not available to viewer")
			return nil, err
		}
	}),
}

var db = dependencies.DBInstance()
