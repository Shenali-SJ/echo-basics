package main

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
	e := echo.New()

	e.GET("/", getHomePage)
	e.GET("/users/:id", getUser)
	e.GET("/show", show)
	e.POST("/save", save)
	e.GET("/users", getUsers)

	e.Logger.Fatal(e.Start(":1323"))

}

func getHomePage(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, Echo World!")
}

func getUser(c echo.Context) error {
	id := c.Param("id")
	return c.String(http.StatusOK, id)
}

func show(c echo.Context) error {
	team := c.QueryParam("team")
	member := c.QueryParam("member")
	return c.String(http.StatusOK, "team:" + team + ", member:" + member)
}

// func save(c echo.Context) error {
// 	name := c.FormValue("name")
// 	email := c.FormValue("email")
// 	return c.String(http.StatusOK, "name:" + name + ", email:" + email)
// }

func save(c echo.Context) error {
	//get name and avatar as form values
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

	return c.HTML(http.StatusOK, "<b>Thank you! " + name + "</b>")
}

func getUsers(c echo.Context) error {
	u := new(User)

	if err := c.Bind(u); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, u)
	// return c.XML(http.StatusCreated, u)
}


