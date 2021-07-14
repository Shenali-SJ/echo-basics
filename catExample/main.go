package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// -------------------------------custom context--------------------------
type CustomContext struct {
	echo.Context
}

func (c *CustomContext) Foo() {
	println("foo")
}

func (c *CustomContext) Bar() {
	println("bar")
}

// ---------------------------------cookie-------------------------------
func writeCookie(c echo.Context) error {
	//creating a cookie
	cookie := new(http.Cookie)
	cookie.Name = "username"
	cookie.Value = "jon"
	cookie.Expires = time.Now().Add(24 * time.Hour)

	//adds a Set-Cookie header in HTTP response
	c.SetCookie(cookie)

	return c.String(http.StatusOK, "Write a cookie")
}

func readCookie(c echo.Context) error {
	//cookie is read by name
	cookie, err := c.Cookie("username")
	if err != nil {
		return err
	}

	fmt.Println(cookie.Name)
	fmt.Println(cookie.Value)

	return c.String(http.StatusOK, "Read a cookie")
}

func readAllCookies(c echo.Context) error {
	for _, cookie := range c.Cookies() {
		fmt.Println(cookie.Name)
		fmt.Println(cookie.Value)
	}
	return c.String(http.StatusOK, "Read all the cookies")
}

// ---------------------------custome error handler----------------------
func customErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	errorPage := fmt.Sprintf("%d.html", code)
	if err := c.File(errorPage); err != nil {
		c.Logger().Error(err)
	}
	c.Logger().Error(err)
}

func main() {
	e := echo.New()

	e.GET("/cats", GetCats)
	e.GET("/cats2/:data", GetCats2)
	e.POST("/cats", AddCat)

	//custom context
	//create a middleware to extend default context
	//this middleware should be registered before any other middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &CustomContext{c}
			return next(cc)
		}
	})

e.GET("/", func(c echo.Context) error {
	cc := c.(*CustomContext)
	cc.Foo()
	cc.Bar()

	//cookies
	writeCookie(cc)
	readCookie(cc)
	readAllCookies(cc)

	return cc.String(200, "OK")
})

// for error handling
// e.GET("/cats", func (c echo.Context) error {
// 	type Cat struct {
// 		Name string `json"name"`
// 		Type string `json"type"`
// 	}

// 	cat := Cat{}
// 	defer c.Request().Body.Close()

// 	err := json.NewDecoder(c.Request().Body).Decode(&cat)
// 	if err != nil {
// 		e.HTTPErrorHandler = customErrorHandler
// 		customErrorHandler(err, c)
// 		log.Fatalf("Failed to read the request body %s", err)
// 	}
// 	log.Printf("Your cat %#v", cat)
// 	return c.String(http.StatusOK, "We got your cat")
// })

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
	writeCookie(c)  //cookie

	// queryparams
	catName := c.QueryParam("name")
	catType := c.QueryParam("type")
	
	//path param
	dataType := c.Param("data")

	readAllCookies(c)

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