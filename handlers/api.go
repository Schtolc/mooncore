package handlers

import (
	"encoding/json"
	"github.com/Schtolc/mooncore/dependencies"
	"github.com/Schtolc/mooncore/models"
	"github.com/graphql-go/graphql"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func just(param graphql.Type) *graphql.ArgumentConfig {
	return &graphql.ArgumentConfig{
		Type: param,
	}
}
func notNull(param graphql.Type) *graphql.ArgumentConfig {
	return just(graphql.NewNonNull(param))
}
func listOf(param graphql.Type) *graphql.ArgumentConfig {
	return just(graphql.NewList(param))
}

var db = dependencies.DBInstance()

func getUserDetails(email string) (*models.UserDetails, error) {
	user := &models.User{}
	userDetails := &models.UserDetails{}
	if dbc := db.Where("email in (?)", email).First(user); dbc.Error != nil {
		logrus.Error(dbc.Error)
		return nil, dbc.Error
	}
	if dbc := db.Where("user_id in (?)", user.ID).First(userDetails); dbc.Error != nil {
		logrus.Error(dbc.Error)
		return nil, dbc.Error
	}
	return userDetails, nil
}

func getRootMutation() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"createUser":      createUser,
			"editUserProfile": editUserProfile,
			"createAddress":   createAddress,
			"createAvatar":    createAvatar,
			"addPhoto":        addPhoto,
			"addSigns":        addSigns,
		},
	})
}

func getRootQuery() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"address":           address,
			"addressList":       addressList,
			"addressListInArea": addressListInArea,
			"getPhoto":          getPhoto,
			"getUserPhotos":     getUserPhotos,
			"getUser":           getUser,
			"listUsers":         listUsers,
			"getSigns":          getSigns,
			"feed":              feed,
			"viewer":            viewer,
		},
	})
}

var schema *graphql.Schema

func createSchema() graphql.Schema {
	if schema != nil {
		return *schema
	}
	_schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query:    getRootQuery(),
		Mutation: getRootMutation(),
	})
	schema = &_schema
	return *schema
}

func executeQuery(query string, schema graphql.Schema, c echo.Context) *graphql.Result {
	params := graphql.Params{
		Schema:        schema,
		RequestString: query,
		Context:       c.Request().Context(),
	}

	result := graphql.Do(params)
	return result
}

// API GraphQL handler
func API(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		logrus.Error(err)
		return sendResponse(c, http.StatusBadRequest, err.Error())
	}
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		logrus.Error(err)
		return sendResponse(c, http.StatusBadRequest, err.Error())
	}
	query, ok := data["query"]
	if !ok {
		strErr := "No query in request"
		logrus.Error(strErr)
		return sendResponse(c, http.StatusBadRequest, strErr)
	}
	result := executeQuery(query.(string), createSchema(), c)

	if len(result.Errors) > 0 {
		return sendResponse(c, http.StatusNotFound, result.Errors)
	}

	return sendResponse(c, http.StatusOK, result.Data)
}
