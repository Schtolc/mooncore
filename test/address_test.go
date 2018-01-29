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

func TestCreateAddress(t *testing.T) {
	e := expect(t)

	lat := rand.Float64()
	lon := rand.Float64()
	desc := "description"

	reqParams := fmt.Sprintf("lat:\"%.20v\", lon:\"%.20v\", description:\"%s\"", lat, lon, desc)
	respParams := "id, lat, lon, description"
	query := graphQLBody("mutation{createAddress(%s){%s}}", reqParams, respParams)

	resp := e.POST(graphqlUrl).
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object().Value("data").Object().Value("createAddress").Object()

	resp.Value("lat").Equal(lat)
	resp.Value("lon").Equal(lon)
	resp.Value("description").Equal(desc)

	id, err := strconv.ParseInt(resp.Value("id").String().Raw(), 10, 64)

	if err != nil {
		t.Error("Cannot parse id from response")
	}

	dbAddress, err := dao.GetAddressById(id)
	if err != nil {
		t.Error("address hasn't been created")
	}

	assert.Equal(t, lat, dbAddress.Lat, "lat differs")
	assert.Equal(t, lon, dbAddress.Lon, "lon differs")
	assert.Equal(t, desc, dbAddress.Description, "description differs")

	if err := dao.DeleteAddress(dbAddress.ID); err != nil {
		t.Error("address cannot be deleted")
	}
}

func TestGetAddress(t *testing.T) {
	e := expect(t)

	lat := rand.Float64()
	lon := rand.Float64()
	desc := "description"

	address, err := dao.CreateAddress(lat, lon, desc)

	if err != nil {
		t.Error("address cannot be created")
	}

	reqParams := fmt.Sprintf("id:\"%v\"", address.ID)
	respParams := "id, lat, lon, description"
	query := graphQLBody("{address(%s){%s}}", reqParams, respParams)

	resp := e.POST(graphqlUrl).
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object().Value("data").Object().Value("address").Object()
	resp.Value("lat").Equal(lat)
	resp.Value("lon").Equal(lon)
	resp.Value("description").Equal(desc)

	if err := dao.DeleteAddress(address.ID); err != nil {
		t.Error("address cannot be deleted")
	}
}
