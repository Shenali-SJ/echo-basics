package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)


func main() {
	e := echo.New()

	e.GET("/cats", GetCats)
	e.GET("/cats2/:data", GetCats2)
	e.POST("/cats", AddCat)

	e.Logger.Fatal(e.Start(":8000"))
}

// query param
// http://localhost:8000/cats?name=Miawington&type=Persian
func GetCats(c echo.Context) error {
	catName := c.QueryParam("name")
	catType := c.QueryParam("type")

	return c.String(http.StatusOK, "Cat name : " + catName + " type : " + catType)
}

// path variables and query param
// http://localhost:8000/cats2/json?name=mickyBoo&type=spinx
func GetCats2(c echo.Context) error {
	// queryparams
	catName := c.QueryParam("name")
	catType := c.QueryParam("type")
	
	//path param
	dataType := c.Param("data")

	if dataType == "string" {
		return c.String(http.StatusOK, "Cat Name : "+ catName + " type : " + catType)
	} else if dataType == "json" {
		return c.JSON(http.StatusOK, map[string]string{
			"name": catName,
			"type": catType,
		})
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Specify data type as string or JSON",
		})
	}
}

func AddCat(c echo.Context) error {
	type Cat struct {
		Name string `json"name"`
		Type string `json"type"`
	}

	cat := Cat{}
	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&cat)
	if err != nil {
		log.Fatalf("Failed to read the request body %s", err)
	}
	log.Printf("Your cat %#v", cat)
	return c.String(http.StatusOK, "We got your cat")
}