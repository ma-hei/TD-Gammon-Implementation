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
        //fmt.Printf("%v", len(b.points))
	e := echo.New()
	e.GET("/GetBoardUpdate", func(c echo.Context) error {
              u := &User{
              Name:  "Jon",
              Email: "jon@labstack.com",
              }
            b := BackgammonState{}
            b.InitFromString(c)
            b.printState()
            followUpStates := b.rollDiceAndFindFollowUpStates(2)
            for _, v := range followUpStates {
                fmt.Printf("--------------------\n")
                v.rollDice()
                fmt.Printf("found another state\n")
            }
            //followUpStates := b.findAllPossibleFollowUpStates(4, 2)
            //fmt.Printf("n follow up states %v\n", len(followUpStates))
            //for _, s := range followUpStates {
            //    fmt.Printf("------------------\n")
            //    s.printState()
            //}
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
