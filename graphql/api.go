package graphql

import (
	"context"
	"encoding/json"
	"errors"
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
		"signUp":     signUp,     // tested
		"signIn":     signIn,     // tested
		"editMaster": editMaster, // tested
	},
})

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"master": master, //tested
		"feed":   feed,   // tested
		"viewer": viewer, // tested
	},
})

var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})

// resolveMiddleware check access rights before resolving function
func resolveMiddleware(right int, next graphql.FieldResolveFn) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		if result := CheckRights(right, p); !result {
			err := errors.New("AccessDenied")
			return func(params graphql.ResolveParams) (interface{}, error) {
				return nil, err
			}(p)
		}
		return next(p)
	}
}

func executeQuery(query string, variables map[string]interface{}, c echo.Context) *graphql.Result {
	return graphql.Do(graphql.Params{
		Schema:         schema,
		RequestString:  query,
		VariableValues: variables,
		Context:        context.WithValue(c.Request().Context(), utils.GraphQLContextUserKey, c.Get(utils.UserKey)),
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
