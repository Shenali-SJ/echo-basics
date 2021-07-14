package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)


var (
	mockDB = map[string]*User{
		"shenali@gmail.com": &User{"Shenali", "shenali@gmail.com"},
	}
	userJSON = `{"name":"Shenali","email":"shenali@gmail.com"}`
)

//test failed
func TestCreateUser(t *testing.T) {
	e := echo.New()

	// returns a new incoming server Request, suitable for passing to an http.Handler for testing.
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// returns an initialized ResponseRecorder.
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	h := &handler{mockDB}

	//assertions
	if assert.NoError(t, h.createUser(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, userJSON, rec.Body.String())
	}

}