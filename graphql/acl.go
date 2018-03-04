package graphql

import (
	"errors"
	"github.com/Schtolc/mooncore/models"
	"github.com/Schtolc/mooncore/utils"
	"github.com/graphql-go/graphql"
)

// CheckRights function to check access rights
func CheckRights(right int, p graphql.ResolveParams) (*models.User, error) {
	user := p.Context.Value(utils.GraphQLContextUserKey)
	if user == nil {
		return nil, errors.New("AccessDeny")
	}
	userModel := user.(*models.User)
	if userModel.Role&right == userModel.Role {
		return userModel, nil
	}
	return nil, errors.New("AccessDeny")
}
