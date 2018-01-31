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

	address, err := dao.CreateAddress(rand.Float64(), rand.Float64(), randString())

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

	master, err := dao.CreateMaster("", email, password, name, address.ID, photo.ID)

	if err != nil {
		t.Error("cannot create master")
	}

	defer dao.DeleteMaster(master.ID)

	reqParams := fmt.Sprintf("email:\"%s\", password:\"%s\"", email, password)
	respParams := "token"

	query := graphQLBody("mutation {signIn(%s){%s}}", reqParams, respParams)

	token := e.POST(graphqlUrl).
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object().Value("data").
		Object().Value("signIn").Object().ContainsKey("token").Value("token").String().Raw()

	assert.NotEmpty(t, token, "empty token")

	query = graphQLBody("query { viewer { id, username, email, role } }")

	root := e.POST(graphqlUrl).WithBytes(query).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Object()

	root.ContainsKey("id")
	root.ContainsKey("username")
	root.ContainsKey("email").Value("email").String().Equal(email)
	root.ContainsKey("role")
}
