package test

import (
	"fmt"
	"github.com/Schtolc/mooncore/dao"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net/http"
	"strconv"
	"testing"
)

func TestCreateMaster(t *testing.T) {
	e := expect(t)

	email := randString()
	password := randString()
	name := randString()

	lat := rand.Float64() + 55
	lon := rand.Float64() + 37
	description := randString()

	reqParams := fmt.Sprintf("email:\"%s\", password:\"%s\", name:\"%s\", lat:\"%f\", lon:\"%f\", description:\"%s\"", email, password, name, lat, lon, description)
	respParams := "id, user {id, username, email, role}, name, address {id, lat, lon, description}, avatar {id, path, tags { id, name } }, photos {id, path, tags { id, name } }, stars, signs {id, name, description, icon}, services {id, name, description, price }"

	query := graphQLBody("mutation{createMaster(%s){%s}}", reqParams, respParams)

	root := e.POST(graphqlURL).
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object().Value("data").Object().Value("createMaster").Object()

	id, err := strconv.ParseInt(root.ContainsKey("id").Value("id").String().Raw(), 10, 64)

	if err != nil {
		t.Error("cannot parse id")
	}

	root.ContainsKey("name").Value("name").String().Equal(name)

	user := root.ContainsKey("user").Value("user").Object()
	user.ContainsKey("email").Value("email").String().Equal(email)

	responseAddress := root.ContainsKey("address").Value("address").Object()

	responseAddress.ContainsKey("lat").Value("lat").Number().EqualDelta(lat, 0.0001)
	responseAddress.ContainsKey("lon").Value("lon").Number().EqualDelta(lon, 0.0001)
	responseAddress.ContainsKey("description").Value("description").Equal(description)

	master, err := dao.GetMasterByID(id)
	if err != nil {
		t.Error("cannot load master from database")
	}

	assert.Equal(t, master.Name, name, "name in database differs")

	dao.DeleteMaster(id)
}

func TestGetMaster(t *testing.T) {
	e := expect(t)

	address, err := dao.CreateAddress(rand.Float64(), rand.Float64())

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

	reqParams := fmt.Sprintf("id:\"%d\"", master.ID)
	respParams := "id, user {id, username, email, role}, name, address {id, lat, lon, description}, avatar {id, path, tags { id, name } }, photos {id, path, tags { id, name } }, stars, signs {id, name, description, icon}, services {id, name, description, price }"

	query := graphQLBody("query {master(%s){%s}}", reqParams, respParams)

	root := e.POST(graphqlURL).
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object().Value("data").Object().Value("master").Object()

	root.ContainsKey("name").Value("name").String().Equal(name)

	user := root.ContainsKey("user").Value("user").Object()
	user.ContainsKey("id").Value("id").String().Equal(strconv.FormatInt(master.UserID, 10))
	user.ContainsKey("email").Value("email").String().Equal(email)

	responseAddress := root.ContainsKey("address").Value("address").Object()
	addressID, err := strconv.ParseInt(responseAddress.ContainsKey("id").Value("id").String().Raw(), 10, 64)
	if err != nil {
		t.Error("cannot parse address.id")
	}

	assert.Equal(t, address.ID, addressID, "address in response id differs")
	responseAddress.ContainsKey("lat").Value("lat").Equal(address.Lat)
	responseAddress.ContainsKey("lon").Value("lon").Equal(address.Lon)
	responseAddress.ContainsKey("description").Value("description").Equal(address.Description)
}

func TestFeed(t *testing.T) {
	address, err := dao.CreateAddress(rand.Float64(), rand.Float64())

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

	e := expect(t)

	realCount, err := dao.MasterCount()
	if err != nil {
		t.Error("cannot count master")
	}

	reqParams := fmt.Sprintf("limit:%d, offset:%d", realCount, 0)
	respParams := "id, user {id, username, email, role}, name, address {id, lat, lon, description}, avatar {id, path, tags { id, name } }, photos {id, path, tags { id, name } }, stars, signs {id, name, description, icon}, services {id, name, description, price }"

	query := graphQLBody("query {feed(%s){%s}}", reqParams, respParams)

	root := e.POST(graphqlURL).
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object().Value("data").Object().Value("feed").Array()

	root.Length().Ge(1)

	first := root.First().Object()
	first.ContainsKey("id").
		ContainsKey("user").
		ContainsKey("name").
		ContainsKey("address").
		ContainsKey("avatar").
		ContainsKey("photos").
		ContainsKey("stars").
		ContainsKey("signs").
		ContainsKey("services")
}
