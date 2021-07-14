package main

//Remove-item alias:curl

import (
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

type User struct {
	Name string `json:"name" xml:"name" form:"name" query:"name"`
	Email string `json:"email" xml:"email" form:"email" query:"email"`
}

func main() {
	//create a new echo instance
	e := echo.New()

	e.GET("/", getHomePage)
	e.GET("/users/:id", getUser)
	e.GET("/show", show)
	e.POST("/save", save)
	e.POST("/users", getUsers)

	//start the echo server
	e.Logger.Fatal(e.Start(":1323"))

}

func getHomePage(c echo.Context) error {
	// String - sends a string response with a status code
	return c.String(http.StatusOK, "Hello, Echo World!")
}

// path parameters
// http://localhost/users/Shena
func getUser(c echo.Context) error {
	// Param - returns path parameter by name
	id := c.Param("id")
	return c.String(http.StatusOK, id)
}

// query parameters
// http://localhost/show?team=x-men&member=wolverine
func show(c echo.Context) error {
	// QueryParam -  returns the query parameter for the given name
	team := c.QueryParam("team")
	member := c.QueryParam("member")
	return c.String(http.StatusOK, "team:" + team + ", member:" + member)
}

// Content-type = application/x-www-form-urlencoded -->
// func save(c echo.Context) error {
// 	name := c.FormValue("name")
// 	email := c.FormValue("email")
// 	return c.String(http.StatusOK, "name:" + name + ", email:" + email)
// }

// Content-type = multipart/form-data
func save(c echo.Context) error {
	//get name and avatar as form values
	// FormValue - returns the form field value for the given name
	name := c.FormValue("name")
	avatar, err := c.FormFile("avatar")

	if err != nil {
		return err
	}

	//source
	src, err := avatar.Open()
	if err != nil {
		return err
	}

	//destination
	dst, err := os.Create(avatar.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	//HTML -  sends a http response wwith a status code
	return c.HTML(http.StatusOK, "<b>Thank you! " + name + "</b>")
}

func getUsers(c echo.Context) error {
	u := new(User)

	// Bind - Binds req body to provided type
	// default binder binds based on the content-type header
	if err := c.Bind(u); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, u)
	// return c.XML(http.StatusCreated, u)
}

