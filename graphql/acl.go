package graphql

import (
	"errors"
	"github.com/Schtolc/mooncore/models"
	"github.com/Schtolc/mooncore/utils"
	"github.com/graphql-go/graphql"
)

var (
	// MethodAccess for access rights
	MethodAccess = map[string]int{
		"master": models.AnonRole, // access for all users including anon
		"feed":   models.AnonRole,
		"viewer": models.AnonRole,

		"createMaster": models.AnonRole,
		"createClient": models.AnonRole,
		"signIn":       models.AnonRole,
	}
)

// CheckRights function to check access rights
func CheckRights(p graphql.ResolveParams) (*models.User, error) {
	user := p.Context.Value(utils.GraphQLContextUserKey)
	right := MethodAccess[p.Info.FieldName]
	if user == nil {
		return nil, errors.New("AccessDeny")
	}
	userModel := user.(*models.User)
	if userModel.Role&right == userModel.Role {
		return userModel, nil
	}
	return nil, errors.New("AccessDeny")
}
