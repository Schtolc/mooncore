package test

import (
	"fmt"
	"github.com/Schtolc/mooncore/dao"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strconv"
	"testing"
)

func TestCreateClient(t *testing.T) {
	e := expect(t)

	photo, err := dao.CreatePhoto(randString(), nil)

	if err != nil {
		t.Error("cannot create photo")
	}

	defer dao.DeletePhoto(photo.ID)

	email := randString()
	password := randString()
	name := randString()

	reqParams := fmt.Sprintf("email:\"%s\", password:\"%s\", name:\"%s\", photo_id:\"%d\"", email, password, name, photo.ID)
	respParams := "id, user {id, username, email, role}, name, avatar {id, path, tags { id, name } }, favorites { id, name }"

	query := graphQLBody("mutation{createClient(%s){%s}}", reqParams, respParams)

	root := e.POST(graphqlURL).
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object().Value("data").Object().Value("createClient").Object()

	id, err := strconv.ParseInt(root.ContainsKey("id").Value("id").String().Raw(), 10, 64)

	if err != nil {
		t.Error("cannot parse id")
	}

	root.ContainsKey("name").Value("name").String().Equal(name)

	user := root.ContainsKey("user").Value("user").Object()
	user.ContainsKey("email").Value("email").String().Equal(email)

	responsePhoto := root.ContainsKey("avatar").Value("avatar").Object()
	photoID, err := strconv.ParseInt(responsePhoto.ContainsKey("id").Value("id").String().Raw(), 10, 64)
	if err != nil {
		t.Error("cannot parse avatar.id")
	}

	assert.Equal(t, photo.ID, photoID, "photo id in response differs")
	responsePhoto.ContainsKey("path").Value("path").Equal(photo.Path)

	client, err := dao.GetClientByID(id)
	if err != nil {
		t.Error("cannot load client from database")
	}

	assert.Equal(t, client.PhotoID, photo.ID, "photo id in database differs")
	assert.Equal(t, client.Name, name, "name in database differs")

	dao.DeleteClient(id)
}

func TestGetClient(t *testing.T) {
	e := expect(t)

	photo, err := dao.CreatePhoto(randString(), nil)

	if err != nil {
		t.Error("cannot create photo")
	}

	defer dao.DeletePhoto(photo.ID)

	email := randString()
	password := randString()
	name := randString()

	client, err := dao.CreateClient("", email, password, name, photo.ID)

	if err != nil {
		t.Error("cannot create client")
	}

	defer dao.DeleteClient(client.ID)

	reqParams := fmt.Sprintf("email:\"%s\", password:\"%s\"", email, password)
	respParams := "token"

	query := graphQLBody("mutation {signIn(%s){%s}}", reqParams, respParams)
	token := e.POST(graphqlURL).
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object().Value("data").
		Object().Value("signIn").Object().ContainsKey("token").Value("token").String().Raw()

	assert.NotEmpty(t, token, "empty token")

	query = graphQLBody("query { viewer{... on Master {name}, ... on Client { user {id, username, email, role}, name, avatar {id, path, tags { id, name } }, favorites { id, name }}}}")

	root := e.POST(graphqlURL).WithBytes(query).WithHeader("Authorization", "Bearer "+token).
		Expect().Status(http.StatusOK).JSON().Object().Value("data").Object().Value("viewer").Object()

	root.ContainsKey("name").Value("name").String().Equal(name)

	user := root.ContainsKey("user").Value("user").Object()
	user.ContainsKey("email").Value("email").String().Equal(email)
	user.ContainsKey("id").Value("id").String().Equal(strconv.FormatInt(client.UserID, 10))

	responsePhoto := root.ContainsKey("avatar").Value("avatar").Object()
	photoID, err := strconv.ParseInt(responsePhoto.ContainsKey("id").Value("id").String().Raw(), 10, 64)
	if err != nil {
		t.Error("cannot parse avatar.id")
	}
	assert.Equal(t, photo.ID, photoID, "photo id in response differs")
	responsePhoto.ContainsKey("path").Value("path").Equal(photo.Path)
}
