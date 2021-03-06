package test

import (
	"fmt"
	"github.com/Schtolc/mooncore/dao"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net/http"
	"testing"
)

func TestSignIn(t *testing.T) {
	e := expect(t)

	address, err := dao.CreateAddress(rand.Float64()+55, rand.Float64()+37)

	if err != nil {
		t.Error("cannot create address")
	}

	defer dao.DeleteAddress(address.ID)

	photo, err := dao.CreatePhoto(randString(), nil)

	if err != nil {
		t.Error("cannot create photo")
	}

	defer dao.DeletePhoto(photo.ID)

	email := randString()
	password := randString()
	name := randString()

	master, err := dao.CreateMaster("", email, password, name, address.ID)
	if err != nil {
		t.Error("cannot create master")
	}
	defer dao.DeleteMaster(master.ID)

	reqParams := fmt.Sprintf("email:\"%s\", password:\"%s\"", email, password)
	respParams := "token"

	query := graphQLBody("mutation {signIn(%s){%s}}", reqParams, respParams)
	token := e.POST(graphqlURL).
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object().Value("data").
		Object().Value("signIn").Object().ContainsKey("token").Value("token").String().Raw()

	assert.NotEmpty(t, token, "empty token")

	query = graphQLBody("query { viewer{... on Master {name}, ... on Client {name }}}")

	root := e.POST(graphqlURL).WithBytes(query).WithHeader("Authorization", "Bearer "+token).
		Expect().Status(http.StatusOK).JSON().Object().Value("data").Object().Value("viewer").Object()
	root.ContainsKey("name").Value("name").String().Equal(master.Name)
}
