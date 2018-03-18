package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
	"github.com/Schtolc/mooncore/models"
	"github.com/Schtolc/mooncore/dao"
	"github.com/Schtolc/mooncore/utils"
	"github.com/nbutton23/zxcvbn-go"
	"github.com/badoux/checkmail"
	"errors"
)

var signUp = &graphql.Field{
	Type:        UserType,
	Description: "Sign up",
	Args: graphql.FieldConfigArgument{
		"email":    notNull(graphql.String),
		"password": notNull(graphql.String),
		"role":     notNull(graphql.Int),
	},
	Resolve: resolveMiddleware(models.AnonRole, func(params graphql.ResolveParams) (interface{}, error) {
		password := params.Args["password"].(string)
		email := params.Args["email"].(string)
		role := params.Args["role"].(int)

		if err := checkmail.ValidateFormat(email);err != nil {
			logrus.Error("Wrong format for email")
			return nil, err
		}
		// if err := checkmail.ValidateHost(email); err != nil {
		// 	logrus.Error(err)
		// 	return nil, err
		// }

		result := zxcvbn.PasswordStrength(password, nil)
		if result.CrackTimeDisplay == "instant"{
			return nil, errors.New("Weak password")
		}
		passwordHash, err := utils.HashPassword(password)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		if role == models.MasterRole {
			master, err := dao.CreateMaster(email, passwordHash)
			return master, err
		} else if role == models.ClientRole {
			client, err := dao.CreateClient(email, passwordHash)
			return client, err
		} else if role == models.SalonRole {
			salon, err := dao.CreateSalon(email, passwordHash)
			return salon, err
		} else if role == models.AdminRole {
			admin, err := dao.CreateAdmin(email, passwordHash)
			return admin, err
		}
		return nil, errors.New("unknown role")
	}),
}
var signIn = &graphql.Field{
	Type:        TokenObject, // nil if user not found
	Description: "Sign in",
	Args: graphql.FieldConfigArgument{
		"email":    notNull(graphql.String),
		"password": notNull(graphql.String),
	},
	Resolve: resolveMiddleware(models.AnonRole, func(params graphql.ResolveParams) (interface{}, error) {
		password := params.Args["password"].(string)
		email := params.Args["email"].(string)

		user, err := dao.GetUser(&models.User{ Email: email })
		if err != nil {
			return nil, err
		}
		if !utils.CheckPasswordHash(password, user.PasswordHash) {
			logrus.Info("Wrong password for user: ", email)
			return nil, errors.New("wrong password")
		}

		tokenString, err := utils.CreateJwtToken(user)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		return &models.Token{Token: tokenString}, nil
	}),
}
