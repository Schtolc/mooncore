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
	"regexp"
	"strconv"
	"strings"
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

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
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
	query, ok := data["query"].(string)
	if !ok {
		strErr := "No query in request"
		logrus.Error(strErr)
		return sendResponse(c, http.StatusBadRequest, strErr)
	}

	variables, ok := data["variables"].(map[string]interface{})
	if !ok {
		strErr := "No variables in request"
		logrus.Error(strErr)
		return sendResponse(c, http.StatusBadRequest, strErr)
	}

	braceIdx := strings.Index(query, "{")
	if braceIdx == -1 {
		strErr := "Bad request"
		logrus.Error(strErr)
		return sendResponse(c, http.StatusBadRequest, strErr)
	}
	query = query[braceIdx:] // Отбрасываем все что до первой фигурной скобки

	for key, value := range variables {
		re := regexp.MustCompile(`\$(?i)` + key) // case ignore replacing
		switch value.(type) {
		case float64:
			query = re.ReplaceAllString(query, FloatToString(value.(float64)))
		case string:
			query = re.ReplaceAllString(query, value.(string))
		default:
			strErr := "Unknown type of variable"
			logrus.Error(strErr)
			return sendResponse(c, http.StatusBadRequest, strErr)
		}
	}

	result := executeQuery(query, createSchema())

	if len(result.Errors) > 0 {
		return sendResponse(c, http.StatusNotFound, result.Errors)
	}

	return sendResponse(c, http.StatusOK, result.Data)
}

// FloatToString converts a float number to a string
func FloatToString(inputNum float64) string {
	return strconv.FormatFloat(inputNum, 'f', -1, 64)
}
