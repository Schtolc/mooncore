package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Schtolc/mooncore/dependencies"
	"github.com/Schtolc/mooncore/models"
	"github.com/gavv/httpexpect"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net/http"
	"net/url"
	"testing"
)

var (
	conf      = dependencies.ConfigInstance()
	localhost = url.URL{Scheme: "http", Host: conf.Server.Hostbase.Host + ":" + conf.Server.Hostbase.Port}
	db        = dependencies.DBInstance()
)

type graphQLQuery struct {
	Query string `json:"query"`
}

func graphQLBody(query string, a ...interface{}) []byte {
	body, _ := json.Marshal(graphQLQuery{
		fmt.Sprintf(query, a...),
	})
	return body
}

func expect(t *testing.T) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  localhost.String(),
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: nil,
	})
}

func randString() string {
	var result string
	l := rand.Int() % 50
	for i := 0; i < l; i++ {
		result += string(rand.Int()%('z'-'a') + 'a')
	}
	return result
}

func createAddress(t *testing.T) *models.Address {
	address := &models.Address{
		Lat: rand.Float64(),
		Lon: rand.Float64(),
	}
	assert.Nil(t, db.Create(&address).Error, "address was not created")
	return address
}

func createPhoto(t *testing.T) *models.Photo {
	photo := &models.Photo{
		Path: randString(),
	}
	assert.Nil(t, db.Create(&photo).Error, "photo was not created")
	return photo
}

func TestCreateAddress(t *testing.T) {
	e := expect(t)

	lat := rand.Float64()
	lon := rand.Float64()

	query := graphQLBody("mutation{createAddress(lat:%f,lon:%f){id}}", lat, lon)

	root := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON()

	root.Object().Value("code").Number().Equal(http.StatusOK)

	id := root.Object().Value("body").
		Object().Value("createAddress").
		Object().Value("id").Number().Raw()

	address := models.Address{}

	assert.Nil(t, db.First(&address, int(id)).Error, "address was not created")
	assert.Equal(t, int(address.Lat), int(lat), "created lat doesn't equal to returned")
	assert.Equal(t, int(address.Lon), int(lon), "created lon doesn't equal to returned")

	db.Delete(&address)
}

func TestGetAddress(t *testing.T) {
	e := expect(t)

	address := createAddress(t)

	query := graphQLBody("{address(id:%d){lat, lon}}", address.ID)

	root := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON()

	root.Object().Value("code").Number().Equal(http.StatusOK)

	obj := root.Object().Value("body").
		Object().Value("address")

	lat := obj.Object().Value("lat").Number().Raw()
	lon := obj.Object().Value("lon").Number().Raw()

	assert.Equal(t, int(address.Lat), int(lat), "created lat doesn't equal to returned")
	assert.Equal(t, int(address.Lon), int(lon), "created lon doesn't equal to returned")

	db.Delete(&address)

	root = e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON()

	root.Object().Value("code").Number().Equal(http.StatusNotFound)
	root.Object().Value("body").NotNull()
}

func TestCreatePhoto(t *testing.T) {
	e := expect(t)

	path := randString()

	query := graphQLBody("mutation{createPhoto(path:\"%s\"){id}}", path)

	root := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON()

	root.Object().Value("code").Number().Equal(http.StatusOK)

	id := root.Object().Value("body").
		Object().Value("createPhoto").
		Object().Value("id").Number().Raw()

	photo := models.Photo{}

	assert.Nil(t, db.First(&photo, int(id)).Error, "address was not created")
	assert.Equal(t, photo.Path, path, "created path doesn't equal to returned")

	db.Delete(&photo)
}

func TestGetPhoto(t *testing.T) {
	e := expect(t)

	photo := createPhoto(t)

	query := graphQLBody("{photo(id:%d){path}}", photo.ID)

	root := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON()

	root.Object().Value("code").Number().Equal(http.StatusOK)

	root = root.Object().Value("body").
		Object().Value("photo")

	path := root.Object().Value("path").Raw()

	assert.Equal(t, photo.Path, path, "created path doesn't equal to returned")

	db.Delete(&photo)

	root = e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON()

	root.Object().Value("code").Number().Equal(http.StatusNotFound)
	root.Object().Value("body").NotNull()
}

func TestCreateUser(t *testing.T) {
	e := expect(t)

	address := createAddress(t)
	photo := createPhoto(t)

	name := randString()
	password := randString()
	email := randString()
	addressID := address.ID
	photoID := photo.ID

	query := graphQLBody("mutation{createUser(name:\"%s\", "+
		"password:\"%s\", "+
		"email:\"%s\", "+
		"address_id:%d, "+
		"photo_id:%d){id}}", name, password, email, addressID, photoID)

	root := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON()

	root.Object().Value("code").Number().Equal(http.StatusOK)

	id := root.Object().Value("body").
		Object().Value("createUser").
		Object().Value("id").Number().Raw()

	user := models.User{}

	assert.Nil(t, db.First(&user, int(id)).Error, "user was not created")

	assert.Equal(t, user.Name, name, "created name doesn't equal to returned")
	assert.Equal(t, user.Password, password, "created password doesn't equal to returned")
	assert.Equal(t, user.Email, email, "created email doesn't equal to returned")
	assert.Equal(t, user.AddressID, addressID, "created addressID doesn't equal to returned")
	assert.Equal(t, user.PhotoID, photoID, "created photoID doesn't equal to returned")

	db.Delete(&user)
	db.Delete(&photo)
	db.Delete(&address)
}

func TestGetUser(t *testing.T) {
	e := expect(t)

	address := createAddress(t)
	photo := createPhoto(t)

	user := models.User{
		Name:      randString(),
		Password:  randString(),
		Email:     randString(),
		AddressID: address.ID,
		PhotoID:   photo.ID,
	}

	assert.Nil(t, db.Create(&user).Error, "user was not created")

	query := graphQLBody("{user(id:%d){name, email, address{id}, photo{id}}}", user.ID)

	root := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON()

	root.Object().Value("code").Number().Equal(http.StatusOK)

	root = root.Object().Value("body").
		Object().Value("user")

	root.Object().Value("name").Equal(user.Name)
	root.Object().Value("email").Equal(user.Email)
	root.Object().Value("address").Object().Value("id").Equal(user.AddressID)
	root.Object().Value("photo").Object().Value("id").Equal(user.PhotoID)

	db.Delete(&user)
	db.Delete(&photo)
	db.Delete(&address)

	root = e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON()

	root.Object().Value("code").Number().Equal(http.StatusNotFound)
	root.Object().Value("body").NotNull()
}
