package main

import (
	"net/http"
	
	"github.com/labstack/echo"
        "github.com/labstack/echo/middleware"
        "fmt"
)

type User struct {
  Name  string `json:"name" xml:"name"`
  Email string `json:"email" xml:"email"`
}

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
              u := &User{
              Name:  "Jon",
              Email: "jon@labstack.com",
              }
            fmt.Printf(c.Param("0"));
            return c.JSON(http.StatusOK, u)
	})

        e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
          Format: "method=${method}, uri=${uri}, status=${status}\n",
        }))

        e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:4200"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	e.Logger.Fatal(e.Start(":1323"))
}
