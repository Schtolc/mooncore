package graphql

import (
	"context"
	"encoding/json"
	"github.com/Schtolc/mooncore/rest"
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

var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		"createAddress": createAddress, // tested
		"createMaster":  createMaster,  // tested
		"createClient":  createClient,  // tested
		"signIn":        signIn,        // tested
		//"editMaster":    editMaster,
		//"editClient":    editClient,
		//"addService":    addService,
		//"removeService": removeService,
	},
})

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"address":         address, // tested
		"addressesInArea": addressesInArea,
		"master":          master, //tested
		"client":          client, // tested
		"feed":            feed,   // tested
		"viewer":          viewer, // tested
	},
})

var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})

func executeQuery(query string, variables map[string]interface{}, c echo.Context) *graphql.Result {
	return graphql.Do(graphql.Params{
		Schema:         schema,
		RequestString:  query,
		VariableValues: variables,
		Context:        context.WithValue(c.Request().Context(), utils.GraphQLContextUserKey(rest.UserKey), c.Get(rest.UserKey)),
	})
}

// API GraphQL handler
func API(context echo.Context) error {
	body, err := ioutil.ReadAll(context.Request().Body)
	if err != nil {
		logrus.Error(err)
		return utils.SendResponse(context, http.StatusBadRequest, err.Error())
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		logrus.Error(err)
		return utils.SendResponse(context, http.StatusBadRequest, err.Error())
	}
	rawQuery, ok := data["query"]
	if !ok {
		strErr := "No query in request"
		logrus.Error(strErr)
		return utils.SendResponse(context, http.StatusBadRequest, strErr)
	}

	query, ok := rawQuery.(string)
	if !ok {
		strErr := "Bad query format"
		logrus.Error(strErr)
		return utils.SendResponse(context, http.StatusBadRequest, strErr)
	}

	rawVariables, ok := data["variables"]
	if !ok {
		rawVariables = make(map[string]interface{})
	}

	variables, ok := rawVariables.(map[string]interface{})
	if !ok {
		strErr := "Bad variables format"
		logrus.Error(strErr)
		return utils.SendResponse(context, http.StatusBadRequest, strErr)
	}

	result := executeQuery(query, variables, context)

	if len(result.Errors) > 0 {
		return utils.SendResponse(context, http.StatusNotFound, result.Errors)
	}

	return utils.SendResponse(context, http.StatusOK, result.Data)
}
