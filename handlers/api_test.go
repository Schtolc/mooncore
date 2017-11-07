package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Schtolc/mooncore/dependencies"
	"github.com/gavv/httpexpect"

	"github.com/Schtolc/mooncore/models"
	"math/rand"
	"net/http"
	"net/url"
	"testing"
)

var (
	conf      = dependencies.ConfigInstance()
	localhost = url.URL{Scheme: "http", Host: conf.Server.Hostbase.Host + ":" + conf.Server.Hostbase.Port}
	testUser  = &models.User{
		Email:    "newEmail",
		Password: "passPass",
		Role:     1,
	}
	testUserDetails = &models.UserDetails{
		UserID:    1,
		Name:      "user_name",
		AddressID: 1,
		PhotoID:   1,
	}
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

// [[ TEST CREATE USER ]]

func TestCreateUser(t *testing.T) {
	e := expect(t)

	reqParams := fmt.Sprintf("email:\"%s\", password: \"%s\", role: %d", testUser.Email, testUser.Password, testUser.Role)
	respParams := "id, email, role"
	query := graphQLBody("mutation{createUser(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")

	body := resp.Value("body").Object().Value("createUser").Object()
	body.Value("email").Equal(testUser.Email)
	body.Value("role").Equal(testUser.Role)
}

func TestCreateSecondUser(t *testing.T) {
	e := expect(t)

	reqParams := fmt.Sprintf("email:\"%s\", password: \"%s\", role: %d", testUser.Email, testUser.Password, testUser.Role)
	respParams := "id, email, role"
	query := graphQLBody("mutation{createUser(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")

	body := resp.Value("body").Object().Value("createUser").Object()
	body.Value("email").Equal(testUser.Email)
	body.Value("role").Equal(testUser.Role)
}

func TestCreateUserBadParams(t *testing.T) {
	e := expect(t)

	reqParams := fmt.Sprintf("email:\"%s\", password: %s, role: %d", testUser.Email, testUser.Password, testUser.Role)
	respParams := "id, email, role"
	query := graphQLBody("mutation{createUser(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")

	errorMessage := "Argument \"password\" has invalid value passPass.\nExpected type \"String\", found passPass."
	resp.Value("code").Number().Equal(http.StatusNotFound)
	resp.Value("body").Array().First().Object().Value("message").Equal(errorMessage)
}

func TestCreateUserWithoutParams(t *testing.T) {
	e := expect(t)

	reqParams := fmt.Sprintf("email:\"%s\", password: \"%s\"", testUser.Email, testUser.Password)
	respParams := "id, email, role"
	query := graphQLBody("mutation{createUser(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")

	errorMessage := "Field \"createUser\" argument \"role\" of type \"Int!\" is required but not provided."
	resp.Value("code").Number().Equal(http.StatusNotFound)
	resp.Value("body").Array().First().Object().Value("message").Equal(errorMessage)
}

// [[ QUERY USER ]]

func TestGetUser(t *testing.T) {
	e := expect(t)

	reqParams := fmt.Sprintf("email:\"%s\"", testUser.Email)
	respParams := "name, address{lat, lon}, avatar{path}"
	query := graphQLBody("{getUser(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")

	body := resp.Value("body").Object().Value("getUser").Object()

	body.Value("address").Object().Value("lat").Equal(0)
	body.Value("address").Object().Value("lon").Equal(0)

	body.Value("avatar").Object().Value("path").Equal("default")
}

func TestListUsers(t *testing.T) {
	e := expect(t)

	respParams := "id, name, address{lat, lon}, avatar{path}"
	query := graphQLBody("{listUsers{%s}}", respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")
	body := resp.Value("body").Object().Value("listUsers").Array().Element(1).Object()
	body.Value("name").Equal("")
	body.Value("address").Object().Value("lat").Equal(0)
	body.Value("address").Object().Value("lon").Equal(0)
	body.Value("avatar").Object().Value("path").Equal("default")
}

// [[ TEST USER DETAILS EDIT]]

func TestEditUserProfile(t *testing.T) {
	e := expect(t)

	reqParams := fmt.Sprintf("email:\"%s\", name:\"%s\", address_id: %d, avatar_id: %d",
		testUser.Email, testUserDetails.Name, testUserDetails.AddressID, testUserDetails.PhotoID)
	respParams := fmt.Sprintf("id, name, address{lat,lon,id}, avatar{id, path}, photos{id, path}")

	query := graphQLBody("mutation{editUserProfile(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")

	resp.Value("code").Number().Equal(http.StatusOK)
	body := resp.Value("body").Object().Value("editUserProfile").Object()

	body.Value("name").Equal(testUserDetails.Name)
	body.Value("address").Object().Value("id").Equal(testUserDetails.AddressID)
	body.Value("address").Object().Value("lat").NotNull()
	body.Value("address").Object().Value("lon").NotNull()

	body.Value("avatar").Object().Value("id").Equal(testUserDetails.PhotoID)
	body.Value("avatar").Object().Value("path").NotNull()
}

func TestEditUserProfileNotAllParams(t *testing.T) {
	e := expect(t)

	reqParams := fmt.Sprintf("email:\"%s\", avatar_id: %d",
		testUser.Email, testUserDetails.PhotoID)
	respParams := fmt.Sprintf("id, name, address{lat,lon,id}, avatar{id, path}, photos{id, path}")

	query := graphQLBody("mutation{editUserProfile(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")

	resp.Value("code").Number().Equal(http.StatusOK)
	body := resp.Value("body").Object().Value("editUserProfile").Object()

	body.Value("name").Equal(testUserDetails.Name)
	body.Value("address").Object().Value("id").Equal(testUserDetails.AddressID)
	body.Value("address").Object().Value("lat").NotNull()
	body.Value("address").Object().Value("lon").NotNull()

	body.Value("avatar").Object().Value("id").Equal(testUserDetails.PhotoID)
	body.Value("avatar").Object().Value("path").NotNull()
}
func TestCreateAddress(t *testing.T) {
	e := expect(t)

	address := &models.Address{
		Lat:         rand.Float64(),
		Lon:         rand.Float64(),
		Description: "description",
	}

	reqParams := fmt.Sprintf("lat:\"%.20v\", lon:\"%.20v\", description:\"%s\"", address.Lat, address.Lon, address.Description)
	respParams := "id, lat, lon, description"
	query := graphQLBody("mutation{createAddress(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")

	resp.Value("code").Number().Equal(http.StatusOK)

	body := resp.Value("body").Object().Value("createAddress")
	body.Object().Value("lat").Equal(address.Lat)
	body.Object().Value("lon").Equal(address.Lon)
	body.Object().Value("description").Equal(address.Description)
}

func TestCreateAddressBadParamLat(t *testing.T) {
	e := expect(t)

	address := &models.Address{
		Lat:         rand.Float64(),
		Lon:         rand.Float64(),
		Description: "description",
	}

	reqParams := fmt.Sprintf("lat:\"string\", lon:\"string\", description:\"%s\"", address.Description)
	respParams := "id, lat, lon, description"
	query := graphQLBody("mutation{createAddress(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Value("code").Number().Equal(http.StatusNotFound)
	errorMessage := "strconv.ParseFloat: parsing \"string\": invalid syntax"
	resp.Value("body").Array().First().Object().Value("message").Equal(errorMessage)
}

func TestCreateAddressBadParamDescription(t *testing.T) {
	e := expect(t)

	address := &models.Address{
		Lat:         rand.Float64(),
		Lon:         rand.Float64(),
		Description: "description",
	}

	reqParams := fmt.Sprintf("lat:\"%.16v\", lon:\"%.16v\", description:%.16v", address.Lat, address.Lon, address.Lon)
	respParams := "id, lat, lon, description"
	query := graphQLBody("mutation{createAddress(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	errorMessage := fmt.Sprintf("Argument \"description\" has invalid value %.16v.\nExpected type \"String\", found %.16v.", address.Lon, address.Lon)

	resp.Value("code").Number().Equal(http.StatusNotFound)
	resp.Value("body").Array().First().Object().Value("message").Equal(errorMessage)
}

func TestCreateAddressWithoutParams(t *testing.T) {
	e := expect(t)

	reqParams := fmt.Sprintf("lat:\"string\", lon:\"string\"")
	respParams := "id, lat, lon, description"
	query := graphQLBody("mutation{createAddress(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	errorMessage := "Field \"createAddress\" argument \"description\" of type \"String!\" is required but not provided."
	resp.Value("code").Number().Equal(http.StatusNotFound)
	resp.Value("body").Array().First().Object().Value("message").Equal(errorMessage)

}

func TestCreateAvatar(t *testing.T) {
	e := expect(t)

	path := "random_path"
	tags := []int{1, 2}

	reqParams := fmt.Sprintf("path:\"%s\", tags:[%d,%d]", path, tags[0], tags[1])
	respParams := "path, tags{id, name}"
	query := graphQLBody("mutation{createAvatar(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")
	resp.Value("code").Number().Equal(http.StatusOK)

	body := resp.Value("body").Object().Value("createAvatar").Object()
	body.Value("path").Equal(path)
	body.Value("tags").Array().First().Object().Value("id").Equal(tags[0])
	body.Value("tags").Array().Last().Object().Value("id").Equal(tags[1])
}

func TestCreateAvatarWithNotExistingTags(t *testing.T) {
	e := expect(t)

	path := "random_path"
	tags := []int{1000, 2000}
	reqParams := fmt.Sprintf("path:\"%s\", tags:[%d,%d]", path, tags[0], tags[1])
	respParams := "path, tags{id, name}"
	query := graphQLBody("mutation{createAvatar(%s){%s}}", reqParams, respParams)
	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")
	resp.Value("code").Number().Equal(http.StatusNotFound)
	resp.Value("body").Array().First().Object().Value("message").Equal("No such tags")
}

func TestCreateAvatarWithBadParamsTags(t *testing.T) {
	e := expect(t)

	path := "random_path"
	tags := []float64{
		rand.Float64(), rand.Float64(),
	}

	reqParams := fmt.Sprintf("path:\"%s\", tags:[%0.16v,%0.16v]", path, tags[0], tags[1])
	respParams := "path, tags{id, name}"
	query := graphQLBody("mutation{createAvatar(%s){%s}}", reqParams, respParams)
	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")

	errorMessage := fmt.Sprintf("Argument \"tags\" has invalid value [%0.16v, %0.16v].\nIn element #1: Expected type \"Int\", found %0.16v.\nIn element #1: Expected type \"Int\", found %0.16v.", tags[0], tags[1], tags[0], tags[1])

	resp.Value("code").Number().Equal(http.StatusNotFound)
	resp.Value("body").Array().First().Object().Value("message").Equal(errorMessage)
}

func TestCreateAvatarWithoutTags(t *testing.T) {
	e := expect(t)

	path := "random_path"

	reqParams := fmt.Sprintf("path:\"%s\"", path)
	respParams := "path"
	query := graphQLBody("mutation{createAvatar(%s){%s}}", reqParams, respParams)
	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")

	resp.Value("code").Number().Equal(http.StatusOK)
	body := resp.Value("body").Object().Value("createAvatar").Object()
	body.Value("path").Equal(path)
}

// [[ TEST PHOTOS CREATE]]

func TestAddPhoto(t *testing.T) {
	e := expect(t)

	path := "random_path123"
	tags := []int{1, 2}

	reqParams := fmt.Sprintf("email:\"%s\", path:\"%s\", tags:[%d,%d]", testUser.Email, path, tags[0], tags[1])
	respParams := "path, tags{id, name}"
	query := graphQLBody("mutation{addPhoto(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")
	resp.Value("code").Number().Equal(http.StatusOK)

	body := resp.Value("body").Object().Value("addPhoto").Object()
	body.Value("path").Equal(path)
	body.Value("tags").Array().First().Object().Value("id").Equal(tags[0])
	body.Value("tags").Array().Last().Object().Value("id").Equal(tags[1])
}

// [[ TEST PHOTOS GET]]

func TestGetPhoto(t *testing.T) {
	e := expect(t)

	reqParams := fmt.Sprintf("id:2")
	respParams := "path, tags{id, name}"
	query := graphQLBody("{getPhoto(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")
	resp.Value("code").Number().Equal(http.StatusOK)

	body := resp.Value("body").Object().Value("getPhoto").Object()

	body.Value("path").NotNull()
	body.Value("tags").Array().First().Object().Value("id").NotNull()
	body.Value("tags").Array().Last().Object().Value("id").NotNull()
}

func TestGetUserPhotos(t *testing.T) {
	e := expect(t)

	reqParams := fmt.Sprintf("email:\"%s\"", testUser.Email)
	respParams := "path, tags{id, name}"
	query := graphQLBody("{getUserPhotos(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")
	resp.Value("code").Number().Equal(http.StatusOK)

	body := resp.Value("body").Object().Value("getUserPhotos").Array()

	body.First().Object().Value("path").NotNull()
	body.First().Object().Value("tags").NotNull()
}

// [[ TEST ADD SIGN ]]

func TestCreateSign(t *testing.T) {
	e := expect(t)

	signs := []int{1, 2}
	reqParams := fmt.Sprintf("email:\"%s\", signs:[%d, %d]", testUser.Email, signs[0], signs[1])
	respParams := "id, signs{id, name, photo{path}, description}"
	query := graphQLBody("mutation{addSigns(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")

	resp.Value("code").Number().Equal(http.StatusOK)

	body := resp.Value("body").Object().Value("addSigns").Object()
	firstElem := body.Value("signs").Array().First().Object()
	firstElem.Value("description").NotNull()
	firstElem.Value("id").NotNull()
	firstElem.Value("name").NotNull()
	firstElem.Value("photo").Object().Value("path").NotNull()

	lastElem := body.Value("signs").Array().Last().Object()
	lastElem.Value("description").NotNull()
	lastElem.Value("id").NotNull()
	lastElem.Value("name").NotNull()
	lastElem.Value("photo").Object().Value("path").NotNull()

}

func TestCreateSignWithBadParamsSigns(t *testing.T) {
	e := expect(t)

	reqParams := fmt.Sprintf("email:\"%s\"", testUser.Email)
	respParams := "id, signs{id, name}"
	query := graphQLBody("mutation{addSigns(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")

	errorMessage := "Field \"addSigns\" argument \"signs\" of type \"[Int]!\" is required but not provided."
	resp.Value("code").Number().Equal(http.StatusNotFound)
	resp.Value("body").Array().First().Object().Value("message").Equal(errorMessage)

}

func TestCreateSignWithoutSigns(t *testing.T) {
	e := expect(t)

	reqParams := fmt.Sprintf("email:\"%s\"", testUser.Email)
	respParams := "id, signs{id, name}"
	query := graphQLBody("mutation{addSigns(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")

	errorMessage := "Field \"addSigns\" argument \"signs\" of type \"[Int]!\" is required but not provided."
	resp.Value("code").Number().Equal(http.StatusNotFound)
	resp.Value("body").Array().First().Object().Value("message").Equal(errorMessage)
}

// [[ TEST GET SIGN ]]

func TestGetSign(t *testing.T) {
	e := expect(t)

	reqParams := fmt.Sprintf("email:\"%s\"", testUser.Email)
	respParams := "id, name, photo{path}, description"
	query := graphQLBody("{getSigns(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")

	resp.Value("code").Number().Equal(http.StatusOK)

	body := resp.Value("body").Object().Value("getSigns").Array()
	body.First().Object().Value("name").Equal("accuracy")
	body.First().Object().Value("description").Equal("means accuracy")

	body.First().Object().Value("photo").Object().Value("path").Equal("default")
}

//// [[ GET FEED ]]

func TestGetFeed(t *testing.T) {
	e := expect(t)
	query := graphQLBody("{getFeed(limit:3){id, name, avatar{id, path, tags{id, name}}, photos{id, path, tags{name}}, signs{id, name, description, photo{id, path, tags{id, name}}}}}")
	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()
	resp.Keys().ContainsOnly("code", "body")

	body := resp.Value("body").Object().Value("getFeed")

	photo := body.Array().First().Object().Value("photos").Array().First()
	photo.Object().Value("id").NotNull()
	photo.Object().Value("path").NotNull()
	photo.Object().Value("tags").Array()

	sign := body.Array().First().Object().Value("signs").Array().First()
	sign.Object().Value("description").NotNull()
	sign.Object().Value("name").NotNull()
	sign.Object().Value("id").NotNull()
	sign.Object().Value("photo").Object().Value("path").NotNull()
	sign.Object().Value("photo").Object().Value("id").NotNull()
	//sign.Object().Value("photo").Object().Value("tags").Array()

	pho := body.Array().First().Object().Value("avatar")
	pho.Object().Value("id").NotNull()
	pho.Object().Value("path").NotNull()
	//pho.Object().Value("tags").Array().First()
}

//// [[ GET ADDRESSES IN GIVEN AREA ]]

func TestGetAddressInGivenArea(t *testing.T) {
	e := expect(t)

	db := dependencies.DBInstance()

	boundary1 := models.Address{
		Lat: rand.Float64() + 10, // to avoid ranges close to zero
		Lon: rand.Float64() + 10,
	}
	boundary2 := models.Address{
		Lat: rand.Float64() + 10,
		Lon: rand.Float64() + 10,
	}

	if boundary1.Lat > boundary2.Lat { // To reach correct order
		boundary1.Lat, boundary2.Lat = boundary2.Lat, boundary1.Lat
	}
	if boundary1.Lon < boundary2.Lon { // To reach correct order
		boundary1.Lon, boundary2.Lon = boundary2.Lon, boundary1.Lon
	}

	objectInArea := models.Address{
		Lat: (boundary1.Lat + boundary2.Lat) / 2,
		Lon: (boundary1.Lon + boundary2.Lon) / 2,
	}

	objectOutOfArea1 := models.Address{
		Lat: boundary1.Lat - rand.Float64(),
		Lon: (boundary1.Lon + boundary2.Lon) / 2,
	}
	objectOutOfArea2 := models.Address{
		Lat: boundary2.Lat + rand.Float64(),
		Lon: (boundary1.Lon + boundary2.Lon) / 2,
	}
	objectOutOfArea3 := models.Address{
		Lat: (boundary1.Lat + boundary2.Lat) / 2,
		Lon: boundary1.Lon + rand.Float64(),
	}
	objectOutOfArea4 := models.Address{
		Lat: (boundary1.Lat + boundary2.Lat) / 2,
		Lon: boundary2.Lon - rand.Float64(),
	}

	db.Create(&objectInArea)
	db.Create(&objectOutOfArea1)
	db.Create(&objectOutOfArea2)
	db.Create(&objectOutOfArea3)
	db.Create(&objectOutOfArea4)

	query := graphQLBody(
		"{addressListInArea(lat1:\"%f\", lon1:\"%f\", lat2:\"%f\", lon2:\"%f\"){lat, lon, id}}",
		boundary1.Lat,
		boundary1.Lon,
		boundary2.Lat,
		boundary2.Lon,
	)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	db.Delete(&objectInArea)
	db.Delete(&objectOutOfArea1)
	db.Delete(&objectOutOfArea2)
	db.Delete(&objectOutOfArea3)
	db.Delete(&objectOutOfArea4)

	resp.Value("code").Number().Equal(http.StatusOK)

	obj := resp.Value("body").
		Object().Value("addressListInArea").Array()
	obj.Length().Equal(1)
	obj.Element(0).Object().Value("id").Equal(objectInArea.ID)
}
