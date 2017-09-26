package handlers

import (
	"fmt"
	"github.com/Schtolc/mooncore/config"
	"github.com/Schtolc/mooncore/database"
	"github.com/Schtolc/mooncore/models"
	"github.com/gavv/httpexpect"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"testing"
)

var (
	conf      = config.Instance()
	localhost = url.URL{Scheme: "http", Host: conf.Server.Hostbase.Host + ":" + conf.Server.Hostbase.Port}
	db        = database.Instance()
)

func randString() string {
	var result string
	l := rand.Int() % 50
	for i := 0; i < l; i++ {
		result += string(rand.Int()%('z'-'a') + 'a')
	}
	return result
}

func TestCreateAddress(t *testing.T) {
	e := httpexpect.New(t, localhost.String())

	lat := rand.Float64()
	lon := rand.Float64()

	query := fmt.Sprintf("mutation{createAddress(lat:%f,lon:%f){id}}", lat, lon)

	root := e.GET("/graphql").
		WithQuery("query", query).Expect().
		Status(http.StatusOK).JSON()

	root.Object().Value("code").Number().Equal(OK)

	id := root.Object().Value("body").
		Object().Value("createAddress").
		Object().Value("id").Number().Raw()

	address := models.Address{}
	db.First(&address, int(id))

	if math.Abs(address.Lat-lat) > 1 || math.Abs(address.Lon-lon) > 1 {
		t.Fail()
	}

	db.Delete(&address)
}

func TestGetAddress(t *testing.T) {
	e := httpexpect.New(t, localhost.String())

	address := &models.Address{
		Lat: rand.Float64(),
		Lon: rand.Float64(),
	}
	db.Create(&address)

	query := fmt.Sprintf("{address(id:%d){lat, lon}}", address.ID)

	root := e.GET("/graphql").
		WithQuery("query", query).Expect().
		Status(http.StatusOK).JSON()

	root.Object().Value("code").Number().Equal(OK)

	obj := root.Object().Value("body").
		Object().Value("address")

	lat := obj.Object().Value("lat").Number().Raw()
	lon := obj.Object().Value("lon").Number().Raw()

	if math.Abs(address.Lat-lat) > 1 || math.Abs(address.Lon-lon) > 1 {
		t.Fail()
	}

	db.Delete(&address)

	root = e.GET("/graphql").
		WithQuery("query", query).Expect().
		Status(http.StatusOK).JSON()

	root.Object().Value("code").Number().Equal(ERROR)
	root.Object().Value("body").NotNull()
}

func TestCreatePhoto(t *testing.T) {
	e := httpexpect.New(t, localhost.String())

	path := randString()

	query := fmt.Sprintf("mutation{createPhoto(path:\"%s\"){id}}", path)

	root := e.GET("/graphql").
		WithQuery("query", query).Expect().
		Status(http.StatusOK).JSON()

	root.Object().Value("code").Number().Equal(OK)

	id := root.Object().Value("body").
		Object().Value("createPhoto").
		Object().Value("id").Number().Raw()

	photo := models.Photo{}
	db.First(&photo, int(id))

	if photo.Path != path {
		t.Fail()
	}

	db.Delete(&photo)
}

func TestGetPhoto(t *testing.T) {
	e := httpexpect.New(t, localhost.String())

	photo := &models.Photo{
		Path: randString(),
	}
	db.Create(&photo)

	query := fmt.Sprintf("{photo(id:%d){path}}", photo.ID)

	root := e.GET("/graphql").
		WithQuery("query", query).Expect().
		Status(http.StatusOK).JSON()

	root.Object().Value("code").Number().Equal(OK)

	root = root.Object().Value("body").
		Object().Value("photo")

	path := root.Object().Value("path").Raw()

	if path != photo.Path {
		t.Fail()
	}

	db.Delete(&photo)

	root = e.GET("/graphql").
		WithQuery("query", query).Expect().
		Status(http.StatusOK).JSON()

	root.Object().Value("code").Number().Equal(ERROR)
	root.Object().Value("body").NotNull()
}

func TestCreateUser(t *testing.T) {
	e := httpexpect.New(t, localhost.String())

	address := &models.Address{
		Lat: rand.Float64(),
		Lon: rand.Float64(),
	}
	db.Create(&address)

	path := randString()
	photo := &models.Photo{
		Path: path,
	}
	db.Create(&photo)

	name := randString()
	password := randString()
	email := randString()
	addressID := address.ID
	photoID := photo.ID

	query := fmt.Sprintf("mutation{createUser(name:\"%s\", "+
		"password:\"%s\", "+
		"email:\"%s\", "+
		"address_id:%d, "+
		"photo_id:%d){id}}", name, password, email, addressID, photoID)

	root := e.GET("/graphql").
		WithQuery("query", query).Expect().
		Status(http.StatusOK).JSON()

	root.Object().Value("code").Number().Equal(OK)

	id := root.Object().Value("body").
		Object().Value("createUser").
		Object().Value("id").Number().Raw()

	user := models.User{}
	db.First(&user, int(id))

	if name != user.Name || password != user.Password || email != user.Email || addressID != user.AddressID || photoID != user.PhotoID {
		t.Fail()
	}

	db.Delete(&user)
	db.Delete(&photo)
	db.Delete(&address)
}

func TestGetUser(t *testing.T) {
	e := httpexpect.New(t, localhost.String())

	address := &models.Address{
		Lat: rand.Float64(),
		Lon: rand.Float64(),
	}
	db.Create(&address)

	path := randString()
	photo := &models.Photo{
		Path: path,
	}
	db.Create(&photo)

	user := models.User{
		Name:      randString(),
		Password:  randString(),
		Email:     randString(),
		AddressID: address.ID,
		PhotoID:   photo.ID,
	}

	db.Create(&user)

	query := fmt.Sprintf("{user(id:%d){name, email, address{id}, photo{id}}}", user.ID)

	root := e.GET("/graphql").
		WithQuery("query", query).Expect().
		Status(http.StatusOK).JSON()

	root.Object().Value("code").Number().Equal(OK)

	root = root.Object().Value("body").
		Object().Value("user")

	root.Object().Value("name").Equal(user.Name)
	root.Object().Value("email").Equal(user.Email)
	root.Object().Value("address").Object().Value("id").Equal(user.AddressID)
	root.Object().Value("photo").Object().Value("id").Equal(user.PhotoID)

	db.Delete(&user)
	db.Delete(&photo)
	db.Delete(&address)

	root = e.GET("/graphql").
		WithQuery("query", query).Expect().
		Status(http.StatusOK).JSON()

	root.Object().Value("code").Number().Equal(ERROR)
	root.Object().Value("body").NotNull()
}
