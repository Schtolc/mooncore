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
	"database/sql"
	"strconv"
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


var editMaster = &graphql.Field{
	Type:        MasterObject, // nil if user not found
	Description: "edit Master",
	Args: graphql.FieldConfigArgument{
		"name":  notNull(graphql.String),
		"photo": notNull(graphql.String),
		"lat": notNull(graphql.String),
		"lon": notNull(graphql.String),
	},
	Resolve: resolveMiddleware(models.MasterRole, func(params graphql.ResolveParams) (interface{}, error) {
		master, err := dao.GetMasterFromContext(params)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		tx := db.Begin()
		master.Name = params.Args["name"].(string)
		if master.PhotoID.Valid {
			if err := dao.UpdatePhoto(master.PhotoID.Int64, params.Args["photo"].(string), []int64{}, tx); err != nil {
				tx.Rollback()
				logrus.Error(err)
				return nil, err
			}
		} else {
			photo, err := dao.CreatePhoto(params.Args["photo"].(string), []int64{}, tx);
			if err != nil {
				tx.Rollback()
				logrus.Error(err)
				return nil, err
			}
			master.PhotoID = sql.NullInt64{
				Int64: photo.ID,
				Valid: true,
			}
		}
		lat, err := strconv.ParseFloat(params.Args["lat"].(string), 64)
		if err != nil {
			tx.Rollback()
			logrus.Error(err)
			return nil, err
		}

		lon, err := strconv.ParseFloat(params.Args["lon"].(string), 64)
		if err != nil {
			tx.Rollback()
			logrus.Error(err)
			return nil, err
		}
		if master.AddressID.Valid {
			if err := dao.UpdateAddress(master.AddressID.Int64, lat, lon, tx); err != nil {
				tx.Rollback()
				logrus.Error(err)
				return nil, err
			}
		} else {
			address, err := dao.CreateAddress(lat, lon, tx)
			if err != nil {
				tx.Rollback()
				logrus.Error(err)
				return nil, err
			}
			master.AddressID = sql.NullInt64{
				Int64: address.ID,
				Valid: true,
			}
		}
		if err := dao.UpdateMaster(master, tx); err != nil {
			tx.Rollback()
			logrus.Error(err)
			return nil, err
		}
		tx.Commit()
		return master, nil
	}),
}

