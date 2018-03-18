package test

import (
	"github.com/Schtolc/mooncore/dao"
	"github.com/Schtolc/mooncore/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"strconv"
	"fmt"
)
func TestSignUpClient(t *testing.T) {
	e := expect(t)

	email := getEmail()
	password := randString()

	reqParams := fmt.Sprintf("email:\"%s\", password:\"%s\", role: %d", email, password, models.ClientRole)
	respParams := "... on Client{ id }"
	query := graphQLBody("mutation {signUp(%s){%s}}", reqParams, respParams)

	result := e.POST(graphqlURL).
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object().Value("data").
		Object().Value("signUp").Object()

	idString := result.ContainsKey("id").Value("id").String().Raw()
	id, err := strconv.ParseInt(idString, 10, 64)

	_, err = dao.GetClientByID(id)
	if err != nil {
		t.Error("cannot get client after signUp")
	}
	// assert(client.User != nil, "user for client doesnot exist")
	defer dao.DeleteClient(id)
}

func TestSignUpMaster(t *testing.T) {
	e := expect(t)

	email := getEmail()
	password := randString()

	reqParams := fmt.Sprintf("email:\"%s\", password:\"%s\", role: %d", email, password, models.MasterRole)
	respParams := "... on Master{ id }"
	query := graphQLBody("mutation {signUp(%s){%s}}", reqParams, respParams)

	result := e.POST(graphqlURL).
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object().Value("data").
		Object().Value("signUp").Object()

	idString := result.ContainsKey("id").Value("id").String().Raw()
	id, err := strconv.ParseInt(idString, 10, 64)

	_, err = dao.GetMasterByID(id)
	if err != nil {
		t.Error("cannot get Master after signUp")
	}
	// assert(master.User != nil, "user for Master doesnot exist")
	defer dao.DeleteMaster(id)
}

func TestSignUpSalon(t *testing.T) {
	e := expect(t)

	email := getEmail()
	password := randString()

	reqParams := fmt.Sprintf("email:\"%s\", password:\"%s\", role: %d", email, password, models.SalonRole)
	respParams := "... on Salon{ id }"
	query := graphQLBody("mutation {signUp(%s){%s}}", reqParams, respParams)

	result := e.POST(graphqlURL).
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object().Value("data").
		Object().Value("signUp").Object()

	idString := result.ContainsKey("id").Value("id").String().Raw()
	id, err := strconv.ParseInt(idString, 10, 64)

	_, err = dao.GetSalonByID(id)
	if err != nil {
		t.Error("cannot get Salon after signUp")
	}
	// assert(salon.User != nil, "user for Salon doesnot exist")
	defer dao.DeleteSalon(id)
}

func TestSignInClient(t *testing.T) {
	e := expect(t)

	email := getEmail()
	password, passwordHash := getPassword()

	client, err := dao.CreateClient(email, passwordHash)
	if err != nil {
		t.Error("cannot create Client")
	}

	defer dao.DeleteClient(client.ID)

	reqParams := fmt.Sprintf("email:\"%s\", password:\"%s\"", email, password)
	respParams := "token"
	query := graphQLBody("mutation {signIn(%s){%s}}", reqParams, respParams)

	token := e.POST(graphqlURL).
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object().
		Value("data").Object().Value("signIn").Object().
		ContainsKey("token").Value("token").String().Raw()
	assert.NotEmpty(t, token, "empty token")

	query = graphQLBody("query { viewer{... on Client { id }}}")

	root := e.POST(graphqlURL).WithBytes(query).WithHeader("Authorization", "Bearer "+token).
		Expect().Status(http.StatusOK).JSON().Object().Value("data").Object().Value("viewer").Object()
	root.ContainsKey("id").Value("id").Equal(strconv.Itoa(int(client.ID)))
}

func TestSignInMaster(t *testing.T) {
	e := expect(t)

	email := getEmail()
	password, passwordHash := getPassword()
	master, err := dao.CreateMaster(email, passwordHash)
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

	query = graphQLBody("query { viewer{ ... on Master{id}}}")

	root := e.POST(graphqlURL).WithBytes(query).WithHeader("Authorization", "Bearer "+token).
		Expect().Status(http.StatusOK).JSON().Object().Value("data").Object().Value("viewer").Object()

	root.ContainsKey("id").Value("id").Equal(strconv.Itoa(int(master.ID)))
}

func TestSignInSalon(t *testing.T) {
	e := expect(t)

	email := getEmail()
	password, passwordHash := getPassword()
	salon, err := dao.CreateSalon(email, passwordHash)
	if err != nil {
		t.Error("cannot create Salon")
	}
	defer dao.DeleteSalon(salon.ID)

	reqParams := fmt.Sprintf("email:\"%s\", password:\"%s\"", email, password)
	respParams := "token"
	query := graphQLBody("mutation {signIn(%s){%s}}", reqParams, respParams)

	token := e.POST(graphqlURL).
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object().Value("data").
		Object().Value("signIn").Object().ContainsKey("token").Value("token").String().Raw()

	assert.NotEmpty(t, token, "empty token")

	query = graphQLBody("query { viewer{ ... on Salon{id}}}")

	root := e.POST(graphqlURL).WithBytes(query).WithHeader("Authorization", "Bearer "+token).
		Expect().Status(http.StatusOK).JSON().Object().Value("data").Object().Value("viewer").Object()
		root.ContainsKey("id").Value("id").Equal(strconv.Itoa(int(salon.ID)))
}