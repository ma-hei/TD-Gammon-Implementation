package main

import (
	"github.com/labstack/echo"
        "fmt"
        "strconv"
        "strings"
        "math/rand"
)

type BackgammonState struct {
    points [][]int
    dice1 int
    dice2 int
    bar []int
}

func (bgs *BackgammonState) InitFromString(c echo.Context) {
    bgs.Init()
    for i := 0; i<24; i++ {
        val := c.FormValue(strconv.Itoa(i))
        if val != "" {
            temp := strings.Split(val, ",")
            nCheckers, _ := strconv.Atoi(temp[0])
            player, _ := strconv.Atoi(temp[1])
            bgs.points[i][player-1] = nCheckers
        }
    }
}

func (bgs *BackgammonState) printState() {
    for i := 0; i<24; i++ {
        fmt.Printf("player 1 has %v checkers on point %v\n", bgs.points[i][0], i)
        fmt.Printf("player 2 has %v checkers on point %v\n", bgs.points[i][1], i)
    }
}

func (bgs *BackgammonState) rollDice() {
    rand.Seed(42)
    bgs.dice1 = rand.Intn(6) + 1
    bgs.dice2 = rand.Intn(6) + 1
}

func (bgs *BackgammonState) findAllPossibleFollowUpStates(diceRoll int, playerTurn int) []BackgammonState {
    possibleStates := make([]BackgammonState, 0, 100)
    otherPlayer := 2 - playerTurn
    for i, _ := range bgs.points {
        if bgs.points[i][playerTurn] > 0 {
            var targetPoint int
            if (playerTurn == 1) {
                targetPoint = i + diceRoll
            } else {
                targetPoint = i - diceRoll
            }
            targetPointOnField := targetPoint >= 0 && targetPoint < 24
            targetPointOpen := targetPointOnField && bgs.points[targetPoint][otherPlayer] < 2
            targetIsHit := targetPointOnField && bgs.points[targetPoint][otherPlayer] == 1
            if targetPointOpen {
                newBackGammonState := bgs
                newBackGammonState.points[i][playerTurn] -= 1
                newBackGammonState.points[targetPoint][playerTurn] += 1
                if targetIsHit {
                    newBackGammonState.points[targetPoint][otherPlayer] -= 1
                    newBackGammonState.bar[otherPlayer] += 1
                }
                possibleStates = append(possibleStates, *newBackGammonState)
            } 
        }
    } 
    return possibleStates 
}

func (bgs *BackgammonState) Init() {
    bgs.points = make([][]int, 24)
    for i, _ := range bgs.points {
        bgs.points[i] = make([]int, 2)
        bgs.points[i][0] = 0
        bgs.points[i][1] = 0
    }
    bgs.bar = make([]int, 2)
    bgs.bar[0] = 0
    bgs.bar[1] = 0
}

