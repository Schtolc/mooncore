package graphql

import (
	"github.com/Schtolc/mooncore/dao"
	"github.com/Schtolc/mooncore/utils"
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
	"strconv"
)

var master = &graphql.Field{
	Type:        MasterObject, // == nil if not found
	Description: "Get master by id",
	Args: graphql.FieldConfigArgument{
		"id": notNull(graphql.ID),
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		id, err := strconv.ParseInt(params.Args["id"].(string), 10, 64)
		if err != nil {
			return nil, err
		}
		return dao.GetMasterByID(id)
	},
}

var client = &graphql.Field{
	Type:        ClientObject, // == nil if not found
	Description: "Get client by id",
	Args: graphql.FieldConfigArgument{
		"id": notNull(graphql.ID),
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		id, err := strconv.ParseInt(params.Args["id"].(string), 10, 64)
		if err != nil {
			return nil, err
		}
		return dao.GetClientByID(id)
	},
}

var address = &graphql.Field{
	Type:        AddressObject,
	Description: "Get address by id",
	Args: graphql.FieldConfigArgument{
		"id": notNull(graphql.ID),
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		id, err := strconv.ParseInt(params.Args["id"].(string), 10, 64)
		if err != nil {
			return nil, err
		}
		return dao.GetAddressByID(id)
	},
}

var addressesInArea = &graphql.Field{
	Type:        graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(AddressObject))),
	Description: "Get addresses in the area",
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
		return dao.GetAddressesInArea(lat1, lon1, lat2, lon2)
	},
}

var feed = &graphql.Field{
	Type:        graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(MasterObject))),
	Description: "feed",
	Args: graphql.FieldConfigArgument{
		"offset": notNull(graphql.Int),
		"limit":  notNull(graphql.Int),
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		return dao.Feed(p.Args["offset"].(int), p.Args["limit"].(int))
	},
}

var viewer = &graphql.Field{
	Type:        UserObject,
	Description: "current logged user",
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		return p.Context.Value(utils.GraphQLContextUserKey), nil
	},
}
