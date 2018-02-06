package test

import (
	"encoding/json"
	"fmt"
	"github.com/Schtolc/mooncore/dependencies"
	"github.com/gavv/httpexpect"
	"math/rand"
	"net/http"
	"net/url"
	"testing"
)

var (
	conf       = dependencies.ConfigInstance()
	localhost  = url.URL{Scheme: "http", Host: conf.Server.Hostbase.Host + ":" + conf.Server.Hostbase.Port}
	graphqlURL = "/graphql"
)

type graphQLQuery struct {
	Query string `json:"query"`
}

func randString() string {
	var result string
	l := rand.Int() % 50
	for i := 0; i < l; i++ {
		result += string(rand.Int()%('z'-'a') + 'a')
	}
	return result
}

func graphQLBody(query string, a ...interface{}) []byte {
	body, _ := json.Marshal(graphQLQuery{
		fmt.Sprintf(query, a...),
	})
	return body
}

func expect(t *testing.T) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  localhost.String() + conf.Server.APIPrefix,
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: nil,
	})
}

func TestApiCall(t *testing.T) {
	e := expect(t)

	e.OPTIONS(graphqlURL).Expect().Status(http.StatusOK)

	// test empty request
	e.POST(graphqlURL).Expect().Status(http.StatusBadRequest)

	// test request without query
	e.POST(graphqlURL).WithText("{}").
		Expect().Status(http.StatusBadRequest).JSON().Object().Value("data").Equal("No query in request")

	// test bad query format
	e.POST(graphqlURL).WithText("{\"query\":[1,2,3]}").
		Expect().Status(http.StatusBadRequest).JSON().Object().Value("data").Equal("Bad query format")

	// test bad variables format
	e.POST(graphqlURL).WithText("{\"query\":\"123\",\"variables\":[1,2,3]}").
		Expect().Status(http.StatusBadRequest).JSON().Object().Value("data").Equal("Bad variables format")

}
