package main

import (
	"net/http"
        "math/rand"
	"github.com/labstack/echo"
        "github.com/labstack/echo/middleware"
        "fmt"

        "strconv"
)

type User struct {
  Name  string `json:"name" xml:"name"`
  Email string `json:"email" xml:"email"`
}

type Response struct {
  ReturnState  string `json:"returnState" xml:"returnState"`
  //  newState string `json:"state" xml:"state"`
}

func main() {
        //fmt.Printf("%v", len(b.points))
	e := echo.New()
	e.GET("/GetBoardUpdate", func(c echo.Context) error {
            b := BackgammonState{}
            b.InitFromString(c)
            b.printState()
            followUpStates := b.rollDiceAndFindFollowUpStates(b.playerTurn)
            fmt.Printf("found %v states\n", len(followUpStates))
            asList := make([]string, 0, len(followUpStates))
            if len(followUpStates) == 0 {
               temp := b.toString()
               temp += strconv.Itoa(b.dice1)
               temp += ","
               temp += strconv.Itoa(b.dice2)
               temp += "lm1:" + b.lastMove1
               temp += "lm2:" + b.lastMove2
               temp += "lm3:" + b.lastMove3
               temp += "lm4:" + b.lastMove4
               asList = append(asList,temp)
            }
            for k, v := range followUpStates {
                fmt.Printf("--------------------\n")
                fmt.Printf("%v\n", k)
                temp := k
                temp += strconv.Itoa(v.dice1)
                temp += ","
                temp += strconv.Itoa(v.dice2)
                temp += "lm1:" + v.lastMove1
                temp += "lm2:" + v.lastMove2
                temp += "lm3:" + v.lastMove3
                temp += "lm4:" + v.lastMove4
                asList = append(asList, temp)
            }
            stateToReturn := rand.Intn(len(asList))
            fmt.Printf("-----> returning %v\n", asList[stateToReturn])
            u := &Response{
              ReturnState:  asList[stateToReturn],
            }
            return c.JSON(http.StatusOK, u)
            //followUpStates := b.findAllPossibleFollowUpStates(4, 2)
            //fmt.Printf("n follow up states %v\n", len(followUpStates))
            //for _, s := range followUpStates {
            //    fmt.Printf("------------------\n")
            //    s.printState()
            //}

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
