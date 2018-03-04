package graphql

import (
	"github.com/Schtolc/mooncore/models"
	"github.com/Schtolc/mooncore/utils"
	"github.com/graphql-go/graphql"
)

// CheckRights function to check access rights
func CheckRights(right int, p graphql.ResolveParams) bool {
	user := p.Context.Value(utils.GraphQLContextUserKey)
	if user == nil {
		return false
	}
	userModel := user.(*models.User)
	if userModel.Role&right == userModel.Role {
		return true
	}
	return false
}
