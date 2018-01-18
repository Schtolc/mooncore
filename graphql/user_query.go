package graphql

import (
	"github.com/Schtolc/mooncore/dao"
	"github.com/graphql-go/graphql"
)

var master = &graphql.Field{
	Type:        MasterObject, // == nil if not found
	Description: "Get single master by id",
	Args: graphql.FieldConfigArgument{
		"id": notNull(graphql.Int),
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		return dao.GetMasterById(params.Args["id"].(int64))
	},
}

var client = &graphql.Field{
	Type:        ClientObject, // == nil if not found
	Description: "Get single client by id",
	Args: graphql.FieldConfigArgument{
		"id": notNull(graphql.Int),
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		return dao.GetClientById(params.Args["id"].(int64))
	},
}
