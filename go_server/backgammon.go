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
    allCheckersOnHomeBoard []bool
    checkersBearedOff []int
}

func (bgs *BackgammonState) FindNumBearedOffCheckerrs() {
    var numCheckersOnFieldPlayer1 int = 0
    var numCheckersOnFieldPlayer2 int = 0
    for i := 0; i<24; i++ {
        numCheckersOnFieldPlayer1 += bgs.points[i][0]
        numCheckersOnFieldPlayer2 += bgs.points[i][1]
    }
    numCheckersOnFieldPlayer1 += bgs.bar[0]
    numCheckersOnFieldPlayer2 += bgs.bar[1]
    bgs.checkersBearedOff[0] = 15 - numCheckersOnFieldPlayer1
    bgs.checkersBearedOff[1] = 15 - numCheckersOnFieldPlayer2
}

func (bgs *BackgammonState) FindIfAllCheckersOnHomeBoard() {
    var player1AllCheckersHome bool = true
    var player2AllCheckersHome bool = true
    for i := 0; i<24; i++ {
        if (i < 6 || i > 11) && bgs.points[i][1] > 0 {
            player2AllCheckersHome = false
        }
        if i < 18 && bgs.points[i][0] > 0 {
            player1AllCheckersHome = false
        }
    }
    bgs.allCheckersOnHomeBoard[0] = player1AllCheckersHome && bgs.bar[0] == 0
    bgs.allCheckersOnHomeBoard[1] = player2AllCheckersHome && bgs.bar[1] == 0
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
    bgs.FindIfAllCheckersOnHomeBoard()
}

func (bgs *BackgammonState) InitFromOtherState(other *BackgammonState) {
    bgs.Init()
    for i:=0; i<24; i++ {
        bgs.points[i][0] = other.points[i][0]
        bgs.points[i][1] = other.points[i][1]
        bgs.bar[0] = other.bar[0]
        bgs.bar[1] = other.bar[1]
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

func (bgs *BackgammonState) toString() string {
    stateString := ""
    for i := 0; i<24; i++ {
        stateString += strconv.Itoa(bgs.points[i][0])
        stateString += strconv.Itoa(bgs.points[i][1])
        stateString += ","
    }
    stateString += strconv.Itoa(bgs.bar[0])
    stateString += strconv.Itoa(bgs.bar[1])
    return stateString
}

func (bgs *BackgammonState) rollDiceAndFindFollowUpStates(playerTurn int) map[string]BackgammonState {
    bgs.rollDice()
    fmt.Printf("----- rolled %v %v\n", bgs.dice1, bgs.dice2)
    //doubleRoll := bgs.dice1 == bgs.dice2
    return bgs.findFollowUpStatesWithDiceRolled(playerTurn)
}

func (bgs *BackgammonState) findFollowUpStatesWithDiceRolled(playerTurn int) map[string]BackgammonState {
    allFollowUpStates := make(map[string]BackgammonState)
    fmt.Printf("---- first dice \n")
    afterFirstRoll := bgs.findAllPossibleFollowUpStatesWithSingleDiceRoll(bgs.dice1, playerTurn)
    fmt.Printf("----- second dice \n")
    for _, s := range afterFirstRoll {
        afterSecond := s.findAllPossibleFollowUpStatesWithSingleDiceRoll(bgs.dice2, playerTurn)
        for _, s2 := range afterSecond {
            allFollowUpStates[s2.toString()] = s2
        }
    }
    return allFollowUpStates
}

func getFollowUpPoint(start int, playerTurn int, diceRoll int) int {
    var targetPoint int
    if playerTurn == 0 {
        if start < 12 {
            targetPoint = start - diceRoll
        } else {
            targetPoint = start + diceRoll
        }
        if targetPoint < 0 {
            targetPoint = 12 + (targetPoint * (-1)) - 1
        }
    } else {
        if start < 12 {
            targetPoint = start + diceRoll
        } else {
            targetPoint = start - diceRoll
        }
        if start > 11 && targetPoint <= 11 {
            targetPoint = 11 - targetPoint
        }
        if start < 12 && targetPoint >= 12 {
            targetPoint = 24 + (targetPoint - 12)
        }
    }
    
    return targetPoint
}

func getFollowUpStateFromPoint2(startPoint int, playerTurn int, diceRoll int, currentState *BackgammonState) *BackgammonState {
    otherPlayer := (playerTurn+1) % 2
    targetPoint := getFollowUpPoint(startPoint, playerTurn, diceRoll)
    targetPointOnField := targetPoint >= 0 && targetPoint < 24
    targetPointOpen := targetPointOnField && currentState.points[targetPoint][otherPlayer] < 2
    targetIsHit := targetPointOnField && currentState.points[targetPoint][otherPlayer] == 1
    if targetPointOpen {
        fmt.Printf("found possible move from %v to %v\n", startPoint, targetPoint)
        newBackGammonState := BackgammonState{}
        newBackGammonState.InitFromOtherState(currentState)
        newBackGammonState.points[startPoint][playerTurn] -= 1
        newBackGammonState.points[targetPoint][playerTurn] += 1
        if targetIsHit {
            newBackGammonState.points[targetPoint][otherPlayer] -= 1
            newBackGammonState.bar[otherPlayer] += 1
        }
        return &newBackGammonState
    }
    return nil
}

func createNewBackGammonState(startPoint int, targetPoint int, playerTurn int, targetIsHit bool, currentState *BackgammonState) *BackgammonState {
    otherPlayer := (playerTurn+1) % 2
    newBackGammonState := BackgammonState{}
    newBackGammonState.InitFromOtherState(currentState)
    newBackGammonState.points[startPoint][playerTurn] -= 1
    newBackGammonState.points[targetPoint][playerTurn] += 1
    if targetIsHit {
        newBackGammonState.points[targetPoint][otherPlayer] -= 1
        newBackGammonState.bar[otherPlayer] += 1
    }
    return &newBackGammonState
}

func getFollowUpStateFromPoint(startPoint int, playerTurn int, diceRoll int, currentState *BackgammonState) *BackgammonState {
    otherPlayer := (playerTurn+1) % 2
    targetPoint := getFollowUpPoint(startPoint, playerTurn, diceRoll)
    //playerCanBearOff := currentState.allCheckersOnHomeBoard[playerTurn]
    targetPointOnField := targetPoint >= 0 && targetPoint < 24
    targetPointOnFieldAndPointOpen := targetPointOnField && currentState.points[targetPoint][otherPlayer] < 2
    targetCanBeHit := targetPointOnFieldAndPointOpen && currentState.points[targetPoint][otherPlayer] == 1
    
    if (targetPointOnFieldAndPointOpen) {
        // create new state
        fmt.Printf("found possible move from %v to %v\n", startPoint, targetPoint)
        newBackGammonState := createNewBackGammonState(startPoint, targetPoint, playerTurn, targetCanBeHit, currentState)
        return newBackGammonState
    }
    return nil
    
    //targetIsExactBearOff := (playerTurn == 0 && targetPoint == 24) || (playerTurn == 1 && targetPoint == 12)
    //if (playerCanBearOff && targetIsExactBearOff) {
    // crete new state
    //} 
}

func (bgs *BackgammonState) findAllPossibleFollowUpStatesWithSingleDiceRoll(diceRoll int, playerTurn int) []BackgammonState {
    playerTurn--
    possibleStates := make([]BackgammonState, 0, 100)
    for i, _ := range bgs.points {
        if bgs.points[i][playerTurn] > 0 {
            newBackGammonState := getFollowUpStateFromPoint(i, playerTurn, diceRoll, bgs)
            if newBackGammonState != nil {
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
    bgs.allCheckersOnHomeBoard = make([]bool, 2)
    bgs.checkersBearedOff = make([]int, 2)
}

