package graphql

import (
	"errors"
	"github.com/Schtolc/mooncore/models"
	"github.com/Schtolc/mooncore/utils"
	"github.com/graphql-go/graphql"
	//"fmt"
)

var (
	// Permissions
	Perm = map[string]int{
		"Anon":   15,
		"Master": 2, // 10
		"Salon":  4, // 100
		"Admin":  8, // 1000
	}
	// MethodAccess
	MethodAccess = map[string]int{
		"feed":      Perm["Anon"],
		"master":    Perm["Anon"],
		"addMaster": Perm["Salon"],
		"addPhoto":  Perm["Salon"] + Perm["Master"],
		"addSign":   Perm["Admin"],
	}
)

// CheckRights
func CheckRights(p graphql.ResolveParams) (*models.User, error) {
	user := p.Context.Value(utils.UserKey)
	right := MethodAccess[p.Info.FieldName]
	if user == nil {
		return nil, errors.New("AccessDeny")
	} else {
		userModel := user.(*models.User)
		if userModel.Role&right == userModel.Role {
			return userModel, nil
		} else {
			return nil, errors.New("AccessDeny")
		}
	}
}
