package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	mockDB       map[int]*Address
	addressJSON  string
	expectedJSON string
)

type APITestSuite struct {
	suite.Suite
}

func TestAPITestSuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}

func (suite *APITestSuite) SetupTest() {
	mockDB = map[int]*Address{
		1: &Address{1, "Justin", "Greenlee", "jgreenlee24@gmail.com", "214-422-6709"},
		2: &Address{1, "Bob", "Smith", "bobsmith@testing.com", "111-111-1111"},
		3: &Address{1, "Demo", "Address", "demo@demo.org", "222-222-2222"},
	}

	addressJSON = `{"first":"Justin","last":"Greenlee","email":"jgreenlee24@gmail.com","phone":"214-422-6709"}`
	expectedJSON = `{"id":1,"first":"Justin","last":"Greenlee","email":"jgreenlee24@gmail.com","phone":"214-422-6709"}`
}

func (suite *APITestSuite) TestAddressController_Get() {
	rec, c := createEchoContext("GET", "/", strings.NewReader(addressJSON))
	c.SetPath("/address/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	h := NewAddressController(mockDB)

	if assert.NoError(suite.T(), h.Get(c)) {
		assert.Equal(suite.T(), http.StatusOK, rec.Code)
		assert.Equal(suite.T(), expectedJSON, rec.Body.String())
	}
}

func (suite *APITestSuite) TestAddressController_Post() {
	rec, c := createEchoContext("POST", "/address", strings.NewReader(addressJSON))
	h := NewAddressController(mockDB)

	if assert.NoError(suite.T(), h.Post(c)) {
		assert.Equal(suite.T(), http.StatusCreated, rec.Code)
		assert.Equal(suite.T(), expectedJSON, rec.Body.String())
	}
}

func (suite *APITestSuite) TestAddressController_Put() {
	input := `{"first":"Test","last":"Testing","email":"test@testing.com","phone":"111-111-1111"}`
	exp := `{"id":1,"first":"Test","last":"Testing","email":"test@testing.com","phone":"111-111-1111"}`

	rec, c := createEchoContext("PUT", "/", strings.NewReader(input))
	c.SetPath("/address/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	h := NewAddressController(mockDB)

	if assert.NoError(suite.T(), h.Put(c)) {
		assert.Equal(suite.T(), http.StatusOK, rec.Code)
		assert.Equal(suite.T(), exp, rec.Body.String())
	}
}

func (suite *APITestSuite) TestAddressController_Delete() {
	rec, c := createEchoContext("DELETE", "/address", strings.NewReader(addressJSON))
	c.SetPath("/address/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	h := NewAddressController(mockDB)

	if assert.NoError(suite.T(), h.Delete(c)) {
		assert.Equal(suite.T(), http.StatusNoContent, rec.Code)
		assert.Equal(suite.T(), "", rec.Body.String())
	}
}

func (suite *APITestSuite) TestAddressController_List() {
	rec, c := createEchoContext("GET", "/address", strings.NewReader(addressJSON))
	h := NewAddressController(mockDB)

	if assert.NoError(suite.T(), h.List(c)) {
		assert.Equal(suite.T(), http.StatusOK, rec.Code)
		resp := make(map[int]*Address)

		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			suite.T().Error("Error in unmarshalling response", mockDB, rec.Body.String())

		} else {
			assert.Equal(suite.T(), mockDB, resp)
		}
	}
}

func (suite *APITestSuite) TestAddressController_ImportCSV() {
	// e := echo.New()
	// req := httptest.NewRequest(echo.POST, "/address/import", strings.NewReader(addressJSON))
	// req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	// rec := httptest.NewRecorder()
	// c := e.NewContext(req, rec)
	// h := NewAddressController(mockDB)

}

func (suite *APITestSuite) TestAddressController_ExportCSV() {
	// e := echo.New()
	// req := httptest.NewRequest(echo.GET, "/address/export", strings.NewReader(addressJSON))
	// req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	// rec := httptest.NewRecorder()
	// c := e.NewContext(req, rec)
	// h := NewAddressController(mockDB)

}

// HELPERS
func createEchoContext(method string, uri string, reader *strings.Reader) (*httptest.ResponseRecorder, echo.Context) {
	e := echo.New()
	req := httptest.NewRequest(method, uri, reader)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return rec, c
}