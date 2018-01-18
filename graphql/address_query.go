package graphql

import (
	"github.com/Schtolc/mooncore/dao"
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
	"strconv"
)

var address = &graphql.Field{
	Type:        AddressObject,
	Description: "Get single address by id",
	Args: graphql.FieldConfigArgument{
		"id": notNull(graphql.Int),
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		return dao.GetAddressById(params.Args["id"].(int64))
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
