package graphql

import (
	"encoding/json"
	"github.com/Schtolc/mooncore/utils"
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

func getRootMutation() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"createAddress": createAddress,
			"createMaster":  createMaster,
			"createClient":  createClient,
			"signIn":        signIn,
			//"editMaster":    editMaster,
			//"editClient":    editClient,
			//"addService":    addService,
			//"removeService": removeService,
		},
	})
}

func getRootQuery() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"address":         address,
			"addressesInArea": addressesInArea,
			"master":          master,
			"client":          client,
			"feed":            feed,
			"viewer":          viewer,
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

func executeQuery(query string, variables map[string]interface{}, schema graphql.Schema, c echo.Context) *graphql.Result {
	params := graphql.Params{
		Schema:         schema,
		RequestString:  query,
		VariableValues: variables,
		Context:        c.Request().Context(),
	}

	result := graphql.Do(params)
	return result
}

// API GraphQL handler
func API(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		logrus.Error(err)
		return utils.SendResponse(c, http.StatusBadRequest, err.Error())
	}
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		logrus.Error(err)
		return utils.SendResponse(c, http.StatusBadRequest, err.Error())
	}
	rawQuery, ok := data["query"]
	if !ok {
		strErr := "No query in request"
		logrus.Error(strErr)
		return utils.SendResponse(c, http.StatusBadRequest, strErr)
	}

	query, ok := rawQuery.(string)
	if !ok {
		strErr := "Bad query format"
		logrus.Error(strErr)
		return utils.SendResponse(c, http.StatusBadRequest, strErr)
	}

	rawVariables, ok := data["variables"]
	if !ok {
		rawVariables = make(map[string]interface{})
	}

	variables, ok := rawVariables.(map[string]interface{})
	if !ok {
		strErr := "Bad variables format"
		logrus.Error(strErr)
		return utils.SendResponse(c, http.StatusBadRequest, strErr)
	}

	result := executeQuery(query, variables, createSchema(), c)

	if len(result.Errors) > 0 {
		return utils.SendResponse(c, http.StatusNotFound, result.Errors)
	}

	return utils.SendResponse(c, http.StatusOK, result.Data)
}
