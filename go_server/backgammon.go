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
    lastMove1 string
    lastMove2 string
    lastMove3 string
    lastMove4 string
    playerTurn int
}

func (bgs *BackgammonState) PlayerWins(player int) bool {
    bgs.FindNumBearedOffCheckers()
    return (bgs.checkersBearedOff[player] == 15)
}

func (bgs *BackgammonState) FindNumBearedOffCheckers() {
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
    bar1s := c.FormValue("bar1")
    bar1i, _ := strconv.Atoi(bar1s)
    bar2s := c.FormValue("bar2")
    bar2i, _ := strconv.Atoi(bar2s)
    bgs.bar[0] = bar1i
    bgs.bar[1] = bar2i
    
    bgs.lastMove1 = ""
    bgs.lastMove2 = ""
    bgs.lastMove3 = ""
    bgs.lastMove4 = ""

    bgs.playerTurn, _  = strconv.Atoi(c.FormValue("playerTurn"))
    bgs.playerTurn -= 1

    bgs.FindIfAllCheckersOnHomeBoard()
}

func (bgs *BackgammonState) InitFromOtherState(other *BackgammonState) {
    bgs.Init()
    for i:=0; i<24; i++ {
        bgs.points[i][0] = other.points[i][0]
        bgs.points[i][1] = other.points[i][1]
    }
    bgs.bar[0] = other.bar[0]
    bgs.bar[1] = other.bar[1]
    bgs.allCheckersOnHomeBoard[0] = other.allCheckersOnHomeBoard[0]
    bgs.allCheckersOnHomeBoard[1] = other.allCheckersOnHomeBoard[1]

    bgs.lastMove1 = other.lastMove1
    bgs.lastMove2 = other.lastMove2
    bgs.lastMove3 = other.lastMove3
    bgs.lastMove4 = other.lastMove4
    bgs.dice1 = other.dice1
    bgs.dice1 = other.dice1
    bgs.dice1 = other.dice1
    bgs.dice2 = other.dice2
  
    bgs.playerTurn = other.playerTurn
    bgs.FindNumBearedOffCheckers()
}

func (bgs *BackgammonState) printState() {
    for i := 0; i<24; i++ {
        fmt.Printf("player 1 has %v checkers on point %v\n", bgs.points[i][0], i)
        fmt.Printf("player 2 has %v checkers on point %v\n", bgs.points[i][1], i)
    }
}

func (bgs *BackgammonState) rollDice() {
    bgs.dice1 = rand.Intn(6) + 1
    bgs.dice2 = rand.Intn(6) + 1
}

func (bgs *BackgammonState) toString() string {
    stateString := ""
    for i := 0; i<24; i++ {
        stateString += strconv.Itoa(bgs.points[i][0])
        stateString += "."
        stateString += strconv.Itoa(bgs.points[i][1])
        if i == 11 {
            stateString += "|"
        } else {
           stateString += ","
        }
    }
    stateString += strconv.Itoa(bgs.bar[0])
    stateString += "."
    stateString += strconv.Itoa(bgs.bar[1])
    stateString += ":"
    //stateString += strconv.Itoa(bgs.dice1)
    //stateString += ","
    //stateString += strconv.Itoa(bgs.dice2)
    //stateString += "lm1:" + bgs.lastMove1
    //stateString += "lm2:" + bgs.lastMove2
    return stateString
}

func (bgs *BackgammonState) findFollowUpsWithDiceRoll(playerTurn int, a int, b int) (map[string]BackgammonState, bool) {
    if bgs.PlayerWins(playerTurn) || bgs.PlayerWins((playerTurn + 1)%2) {
        followUpStates := make(map[string]BackgammonState)
        newBackGammonState := BackgammonState{}
        newBackGammonState.InitFromOtherState(bgs)
        newBackGammonState.playerTurn = ((playerTurn + 1)%2)
        followUpStates[newBackGammonState.toString()] = newBackGammonState 
        return followUpStates, true
    }
    bgs.dice1 = a
    bgs.dice2 = b
    bgs.playerTurn = playerTurn
    //fmt.Printf("----- player %v rolled %v %v\n", playerTurn, bgs.dice1, bgs.dice2)
    //doubleRoll := bgs.dice1 == bgs.dice2
    temp := bgs.findFollowUpStatesWithDiceRolled(playerTurn)
    if len(temp) == 0 {
       newBackGammonState := BackgammonState{}
       newBackGammonState.InitFromOtherState(bgs)
       temp[newBackGammonState.toString()] = newBackGammonState
    }
    return temp, false
}

func (bgs *BackgammonState) putRandomOnBar(playerTurn int) bool {
    for i := 0; i<24; i++ {
        if bgs.points[i][playerTurn] > 0 {
            bgs.points[i][playerTurn]--;
            bgs.bar[playerTurn]++ 
            return true
        }
    }
    return false
}

func (bgs *BackgammonState) rollDiceAndFindFollowUpStates(playerTurn int) (map[string]BackgammonState, bool) {
    if (bgs.PlayerWins(playerTurn) || bgs.PlayerWins((playerTurn + 1)%2)) {
        //fmt.Printf("Current State is a win\n")
        followUpStates := make(map[string]BackgammonState)
        newBackGammonState := BackgammonState{}
        newBackGammonState.InitFromOtherState(bgs)
        newBackGammonState.playerTurn = ((playerTurn + 1) % 2)
        followUpStates[newBackGammonState.toString()] = newBackGammonState 
        return followUpStates, true
    }
    bgs.rollDice()
    bgs.lastMove1=""
    bgs.lastMove2=""
    bgs.lastMove3=""
    bgs.lastMove4=""
    //bgs.dice1 = 3
    //bgs.dice2 = 3
    bgs.playerTurn = playerTurn
    //fmt.Printf("----- player %v rolled %v %v\n", playerTurn, bgs.dice1, bgs.dice2)
    //doubleRoll := bgs.dice1 == bgs.dice2
    temp := bgs.findFollowUpStatesWithDiceRolled(playerTurn)
    //fmt.Printf("--------------\n") 
    fmt.Printf("found %v\n", len(temp))
    if len(temp) == 0 {
       fmt.Printf("no state\n")
       newBackGammonState := BackgammonState{}
       newBackGammonState.InitFromOtherState(bgs)
       newBackGammonState.playerTurn = ((playerTurn + 1) % 2)
       temp[newBackGammonState.toString()] = newBackGammonState
    }
    return temp, false
}

func (bgs *BackgammonState) findFollowUpStatesWithTwoSameDiceRolled(playerTurn int) map[string]BackgammonState {
    allFollowUpStates := make(map[string]BackgammonState)
    diceRolled := bgs.dice1
    // bug remove *bgs from start here
    start := []BackgammonState{*bgs}

    for i := 0; i<4; i++ {
        //fmt.Printf("---- using dice %v for %v time\n", diceRolled, (i+1))
        followUpStatesAfterDiceUsed := findAllPossibleFollowUpStatesFromStatesWithSingleDiceRoll(diceRolled, playerTurn, start)
        if len(followUpStatesAfterDiceUsed) == 0 {
            //fmt.Printf("---- can't use dice %v times\n", (i+1))
            break;
        } 
        followUpStatesAfterDiceUsed = removeDuplicateStates(followUpStatesAfterDiceUsed)
        start = followUpStatesAfterDiceUsed
        //fmt.Printf("---- find %v states\n", len(start))
    }
    
    addFollowUpStatesToAllStates(start, allFollowUpStates)
    return allFollowUpStates
}

// Use at best both dice. If not use only a single one. If you can't use both dice, but each single one use the larger
func (bgs *BackgammonState) findFollowUpStatesWithTwoDifferentDiceRolled(playerTurn int) map[string]BackgammonState {
    allFollowUpStates := make(map[string]BackgammonState)
    dice1 := bgs.dice1
    dice2 := bgs.dice2
    start := []BackgammonState{*bgs}
    //fmt.Printf("---- from current state, using first dice %v\n", dice1)
    usingOnlyFirstDice := findAllPossibleFollowUpStatesFromStatesWithSingleDiceRoll(dice1, playerTurn, start)
    //fmt.Printf("---- found  %v states when using only first dice\n", len(usingOnlyFirstDice))
    
    //fmt.Printf("---- after first dice used, using second dice %v\n", dice2)    
    usingFirstDiceThenSecondDice := findAllPossibleFollowUpStatesFromStatesWithSingleDiceRoll(dice2, playerTurn, usingOnlyFirstDice)
    //fmt.Printf("---- found  %v states when using first then second dice\n", len(usingFirstDiceThenSecondDice))

    //fmt.Printf("---- from current state, using second dice %v\n", dice2)
    usingOnlySecondDice := findAllPossibleFollowUpStatesFromStatesWithSingleDiceRoll(dice2, playerTurn, start)
    //fmt.Printf("---- found  %v states when using only second dice\n", len(usingOnlySecondDice))

    //fmt.Printf("---- after second dice used, using first dice %v\n", dice1)    
    usingSecondDiceThenFirstDice := findAllPossibleFollowUpStatesFromStatesWithSingleDiceRoll(dice1, playerTurn, usingOnlySecondDice)
    //fmt.Printf("---- found  %v states when using second dice then first dice\n", len(usingSecondDiceThenFirstDice))

    // If it was possible to use both dice, then return that as follow up states
    if (len(usingFirstDiceThenSecondDice) > 0 || len(usingSecondDiceThenFirstDice) > 0) {
        addFollowUpStatesToAllStates(usingFirstDiceThenSecondDice, allFollowUpStates)
        addFollowUpStatesToAllStates(usingSecondDiceThenFirstDice, allFollowUpStates)
        return allFollowUpStates
    }

    // It was not possible to use both dice. If both single dice could be used, use the larger one (that's the rules)
    if len(usingOnlyFirstDice) > 0 && len(usingOnlySecondDice) > 0 {
        if dice1 > dice2 {
            addFollowUpStatesToAllStates(usingOnlyFirstDice, allFollowUpStates)
            return allFollowUpStates
        }
        addFollowUpStatesToAllStates(usingOnlySecondDice, allFollowUpStates)
        return allFollowUpStates
    }
    // Only one of the single dice could be used
    addFollowUpStatesToAllStates(usingOnlyFirstDice, allFollowUpStates)
    addFollowUpStatesToAllStates(usingOnlySecondDice, allFollowUpStates)
    return allFollowUpStates
}

func (bgs *BackgammonState) findFollowUpStatesWithDiceRolled(playerTurn int) map[string]BackgammonState {
    if bgs.dice1 == bgs.dice2 {
        return bgs.findFollowUpStatesWithTwoSameDiceRolled(playerTurn)        
    }
    return bgs.findFollowUpStatesWithTwoDifferentDiceRolled(playerTurn)
}

func addFollowUpStatesToAllStates(followUpStates []BackgammonState, allStates map[string]BackgammonState) {
    for _, s := range followUpStates {
        //fmt.Printf("%v\n", s.toString())
        //fmt.Printf("%v\n", s.lastMove1)
        //fmt.Printf("%v\n", s.lastMove2)
        //fmt.Printf("%v\n", s.lastMove3)
        //fmt.Printf("%v\n", s.lastMove4)
        allStates[s.toString()] = s
    }
}

func removeDuplicateStates(states []BackgammonState) []BackgammonState {
    statesMap := make(map[string]BackgammonState)
    deduplicatedStates := make([]BackgammonState, 0, 100)
    for _, s := range states {
        _, alreadyPresent := statesMap[s.toString()]
        if !alreadyPresent {
            statesMap[s.toString()] = s
            deduplicatedStates = append(deduplicatedStates, s)
        }
    } 
    return deduplicatedStates
}

func findAllPossibleFollowUpStatesFromStatesWithSingleDiceRoll(diceRoll int, playerTurn int, startStates []BackgammonState) []BackgammonState{
    possibleFollowUpStates := make([]BackgammonState, 0, 100)
    for _, s := range startStates {
        // If a checker is on the bar move that first
        if s.bar[playerTurn] > 0 {
            newBackGammonState := getFollowUpStateFromPoint(-1, playerTurn, diceRoll, &s)
            if newBackGammonState != nil {
                possibleFollowUpStates = append(possibleFollowUpStates, *newBackGammonState)
            }
        } else {
            for i, _ := range s.points {
                if s.points[i][playerTurn] > 0 {
                    newBackGammonState := getFollowUpStateFromPoint(i, playerTurn, diceRoll, &s)
                    if newBackGammonState != nil {
                        possibleFollowUpStates = append(possibleFollowUpStates, *newBackGammonState)
                    }
                }
            } 
        }
    }
    return possibleFollowUpStates 
}

func getFollowUpPointFromBar(playerTurn int, diceRoll int) int {
    if playerTurn == 0 {
        return 12 - diceRoll
    } else {
        return 24 - diceRoll
    }
}

// gets the point to which we would put a checker for the player playerTurn
// and dice roll diceRoll. If the point is >=24, this means the checker would
// be beared off. 24 means its exactly beared off. 25 means its beared off +1, etc.
func getFollowUpPoint(start int, playerTurn int, diceRoll int) int {
    if start == -1 {
        return getFollowUpPointFromBar(playerTurn, diceRoll)
    }
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

func (bgs *BackgammonState) AddMoveAsString(move string) {
    if bgs.lastMove1 == "" {
        bgs.lastMove1 = move
    } else if bgs.lastMove2 == "" {
        bgs.lastMove2 = move
    } else if bgs.lastMove3 == "" {
        bgs.lastMove3 = move
    } else if bgs.lastMove4 == "" {
        bgs.lastMove4 = move
    }
}

func createNewBackGammonState(startPoint int, targetPoint int, playerTurn int, targetIsHit bool, currentState *BackgammonState) *BackgammonState {
    newBackGammonState := BackgammonState{}
    newBackGammonState.InitFromOtherState(currentState)
    newBackGammonState.playerTurn = ((playerTurn + 1) % 2)
    moveAsString := ""
    if startPoint != -1 {
        moveAsString = strconv.Itoa(startPoint) + "."
        newBackGammonState.points[startPoint][playerTurn] -= 1
    } else {
        moveAsString = "b."
        newBackGammonState.bar[playerTurn] -= 1
    }
    if targetPoint < 24 {
        moveAsString += strconv.Itoa(targetPoint)
        newBackGammonState.points[targetPoint][playerTurn] += 1
        otherPlayer := (playerTurn+1) % 2
        if targetIsHit {
            moveAsString += "h"
            newBackGammonState.points[targetPoint][otherPlayer] -= 1
            newBackGammonState.bar[otherPlayer] += 1
        }
    } else {
        moveAsString += "off"
    }
    //fmt.Printf("---> move as string: %v\n", moveAsString)
    newBackGammonState.AddMoveAsString(moveAsString)
    newBackGammonState.FindIfAllCheckersOnHomeBoard()
    newBackGammonState.FindNumBearedOffCheckers()
    return &newBackGammonState
}

func (bgs *BackgammonState) getHighestPointWithChecker(playerTurn int) int {
    lowest := 24;
    for i := 0; i<24; i++ {
        if bgs.points[i][playerTurn] > 0 && lowest == 24 {
            lowest = i
        }
    }
    if playerTurn == 0 {
        return 24 - lowest    
    } else {
        return 12 - lowest
    }
}

func getFollowUpStateFromPoint(startPoint int, playerTurn int, diceRoll int, currentState *BackgammonState) *BackgammonState {
    //fmt.Printf("dice roll %v\n", diceRoll)
    otherPlayer := (playerTurn+1) % 2
    //fmt.Printf("start point %v\n", startPoint)
    targetPoint := getFollowUpPoint(startPoint, playerTurn, diceRoll)
    //fmt.Printf("target point %v\n", targetPoint)
    targetPointOnField := targetPoint >= 0 && targetPoint < 24
    targetPointOnFieldAndPointOpen := targetPointOnField && currentState.points[targetPoint][otherPlayer] < 2
    targetCanBeHit := targetPointOnFieldAndPointOpen && currentState.points[targetPoint][otherPlayer] == 1
    
    if (targetPointOnFieldAndPointOpen) {
        // create new state
        //fmt.Printf("found possible move from %v to %v\n", startPoint, targetPoint)
        newBackGammonState := createNewBackGammonState(startPoint, targetPoint, playerTurn, targetCanBeHit, currentState)
        return newBackGammonState
    }

    //fmt.Printf("%v\n",currentState.allCheckersOnHomeBoard[playerTurn])
    if !targetPointOnField && currentState.allCheckersOnHomeBoard[playerTurn] == true {
        if targetPoint == 24 {
            newBackGammonState := createNewBackGammonState(startPoint, targetPoint, playerTurn, targetCanBeHit, currentState)
            return newBackGammonState
        }
        highestPointWithChecker := currentState.getHighestPointWithChecker(playerTurn)
        if (highestPointWithChecker >= 1 && highestPointWithChecker <= 6 && diceRoll > highestPointWithChecker) {
            newBackGammonState := createNewBackGammonState(startPoint, targetPoint, playerTurn, targetCanBeHit, currentState)
            return newBackGammonState
        }
    }  
    return nil
    
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
    bgs.allCheckersOnHomeBoard[0] = false;
    bgs.allCheckersOnHomeBoard[1] = false;
    bgs.checkersBearedOff = make([]int, 2)
    bgs.checkersBearedOff[0] = 0;
    bgs.checkersBearedOff[1] = 0;
}

func (bgs *BackgammonState) InitBeginPosition() {
   bgs.Init()
   //this.putNewCheckersOnPoint(5, 0, 1);
   bgs.points[0][0] =  5;
   //this.putNewCheckersOnPoint(3, 4, 2);
   bgs.points[4][1] = 3;
   //this.putNewCheckersOnPoint(5, 6, 2);
   bgs.points[6][1] = 5;
   //this.putNewCheckersOnPoint(5, 12, 2);
   bgs.points[12][1] = 5;
   //this.putNewCheckersOnPoint(2, 11, 1);
   bgs.points[11][0] = 2;
   //this.putNewCheckersOnPoint(3, 16, 1);
   bgs.points[16][0] = 3;
   //this.putNewCheckersOnPoint(5, 18, 1);
   bgs.points[18][0] = 5;
   //this.putNewCheckersOnPoint(2, 23, 2);
   bgs.points[23][1] = 2;
   bgs.playerTurn = 0
}

func (bgs *BackgammonState) InitBeginPosition2() {
   bgs.Init()
   //this.putNewCheckersOnPoint(5, 0, 1);
   bgs.points[23][0] = 0;
   //this.putNewCheckersOnPoint(3, 4, 2);
   bgs.points[22][0] = 0;
   //this.putNewCheckersOnPoint(5, 6, 2);
   bgs.points[21][0] = 0;
   //this.putNewCheckersOnPoint(5, 12, 2);
   bgs.points[20][0] = 0;
   //this.putNewCheckersOnPoint(2, 11, 1);
   bgs.points[19][0] = 0;
   //this.putNewCheckersOnPoint(3, 16, 1);
   bgs.points[18][0] = 15;
   //this.putNewCheckersOnPoint(5, 18, 1);
   //this.putNewCheckersOnPoint(2, 23, 2);
   bgs.points[6][1] = 15;
   bgs.points[7][1] = 0;
   bgs.points[8][1] = 0;
   bgs.points[9][1] = 0;
   bgs.points[10][1] = 0;
   bgs.points[11][1] = 0;
   bgs.playerTurn = 0
   bgs.FindIfAllCheckersOnHomeBoard() 
}

func (bgs *BackgammonState) InitBeginPosition3() {
   bgs.Init()
   //this.putNewCheckersOnPoint(5, 0, 1);
   bgs.points[23][0] = 0;
   //this.putNewCheckersOnPoint(3, 4, 2);
   bgs.points[22][0] = 0;
   //this.putNewCheckersOnPoint(5, 6, 2);
   bgs.points[21][0] = 1;
   //this.putNewCheckersOnPoint(5, 12, 2);
   bgs.points[20][0] = 0;
   //this.putNewCheckersOnPoint(2, 11, 1);
   bgs.points[19][0] = 1;
   //this.putNewCheckersOnPoint(3, 16, 1);
   bgs.points[18][0] = 0;
   bgs.points[12][0] = 0;
   //this.putNewCheckersOnPoint(5, 18, 1);
   //this.putNewCheckersOnPoint(2, 23, 2);
   bgs.points[0][1] = 0;
   bgs.points[6][1] = 0;
   bgs.points[7][1] = 1;
   bgs.points[8][1] = 0;
   bgs.points[9][1] = 1;
   bgs.points[10][1] = 0;
   bgs.points[11][1] = 0;
   bgs.playerTurn = 0
   bgs.FindIfAllCheckersOnHomeBoard() 
}

func (bgs *BackgammonState) InitBeginPosition4() {
   bgs.Init()
   //this.putNewCheckersOnPoint(5, 0, 1);
   bgs.points[23][0] = 0;
   //this.putNewCheckersOnPoint(3, 4, 2);
   bgs.points[22][0] = 0;
   //this.putNewCheckersOnPoint(5, 6, 2);
   bgs.points[21][0] = 3;
   //this.putNewCheckersOnPoint(5, 12, 2);
   bgs.points[20][0] = 4;
   //this.putNewCheckersOnPoint(2, 11, 1);
   bgs.points[19][0] = 4;
   //this.putNewCheckersOnPoint(3, 16, 1);
   bgs.points[18][0] = 4;
   bgs.points[12][0] = 0;
   //this.putNewCheckersOnPoint(5, 18, 1);
   //this.putNewCheckersOnPoint(2, 23, 2);
   bgs.points[0][1] = 0;
   bgs.points[6][1] = 4;
   bgs.points[7][1] = 4;
   bgs.points[8][1] = 4;
   bgs.points[9][1] = 3;
   bgs.points[10][1] = 0;
   bgs.points[11][1] = 0;
   bgs.playerTurn = 0
   bgs.FindIfAllCheckersOnHomeBoard() 
}

func (bgs *BackgammonState) InitBeginPosition5() {
   bgs.Init()
   //this.putNewCheckersOnPoint(5, 0, 1);
   bgs.points[23][0] = 0;
   //this.putNewCheckersOnPoint(3, 4, 2);
   bgs.points[22][0] = 0;
   //this.putNewCheckersOnPoint(5, 6, 2);
   bgs.points[21][0] = 0;
   //this.putNewCheckersOnPoint(5, 12, 2);
   bgs.points[20][0] = 0;
   //this.putNewCheckersOnPoint(2, 11, 1);
   bgs.points[19][0] = 0;
   //this.putNewCheckersOnPoint(3, 16, 1);
   bgs.points[18][0] = 0;
   bgs.points[12][0] = 15;
   //this.putNewCheckersOnPoint(5, 18, 1);
   //this.putNewCheckersOnPoint(2, 23, 2);
   bgs.points[0][1] = 15;
   bgs.points[6][1] = 0;
   bgs.points[7][1] = 0;
   bgs.points[8][1] = 0;
   bgs.points[9][1] = 0;
   bgs.points[10][1] = 0;
   bgs.points[11][1] = 0;
   bgs.playerTurn = 0
   bgs.FindIfAllCheckersOnHomeBoard() 
}

func (bgs *BackgammonState) InitBeginPosition6() {
   bgs.Init()
   //this.putNewCheckersOnPoint(5, 0, 1);
   bgs.points[23][0] = 0;
   //this.putNewCheckersOnPoint(3, 4, 2);
   bgs.points[22][0] = 0;
   //this.putNewCheckersOnPoint(5, 6, 2);
   bgs.points[21][0] = 0;
   //this.putNewCheckersOnPoint(5, 12, 2);
   bgs.points[20][0] = 0;
   //this.putNewCheckersOnPoint(2, 11, 1);
   bgs.points[19][0] = 0;
   //this.putNewCheckersOnPoint(3, 16, 1);
   bgs.points[18][0] = 0;
   bgs.points[17][0] = 3;
   bgs.points[16][0] = 3;
   bgs.points[15][0] = 3;
   bgs.points[14][0] = 2;
   bgs.points[13][0] = 2;
   bgs.points[12][0] = 2;
   //this.putNewCheckersOnPoint(5, 18, 1);
   //this.putNewCheckersOnPoint(2, 23, 2);
   bgs.points[0][1] = 2;
   bgs.points[1][1] = 2;
   bgs.points[2][1] = 2;
   bgs.points[3][1] = 3;
   bgs.points[4][1] = 3;
   bgs.points[5][1] = 3;
   bgs.points[6][1] = 0;
   bgs.points[7][1] = 0;
   bgs.points[8][1] = 0;
   bgs.points[9][1] = 0;
   bgs.points[10][1] = 0;
   bgs.points[11][1] = 0;
   bgs.playerTurn = 0
   bgs.FindIfAllCheckersOnHomeBoard() 
}

func (bgs *BackgammonState) InitBeginPosition7() {
   bgs.Init()
   //this.putNewCheckersOnPoint(5, 0, 1);
   bgs.points[23][0] = 0;
   //this.putNewCheckersOnPoint(3, 4, 2);
   bgs.points[22][0] = 0;
   //this.putNewCheckersOnPoint(5, 6, 2);
   bgs.points[21][0] = 0;
   //this.putNewCheckersOnPoint(5, 12, 2);
   bgs.points[20][0] = 0;
   //this.putNewCheckersOnPoint(2, 11, 1);
   bgs.points[19][0] = 0;
   //this.putNewCheckersOnPoint(3, 16, 1);
   bgs.points[18][0] = 0;
   bgs.points[17][0] = 0;
   bgs.points[16][0] = 0;
   bgs.points[15][0] = 0;
   bgs.points[14][0] = 0;
   bgs.points[13][0] = 0;
   bgs.points[12][1] = 2;
   //this.putNewCheckersOnPoint(5, 18, 1);
   //this.putNewCheckersOnPoint(2, 23, 2);
   bgs.points[0][1] = 0;
   bgs.points[1][1] = 0;
   bgs.points[2][1] = 0;
   bgs.points[3][1] = 0;
   bgs.points[4][1] = 0;
   bgs.points[5][0] = 2;
   bgs.points[6][0] = 3;
   bgs.points[7][1] = 5;
   bgs.points[8][0] = 5;
   bgs.points[9][1] = 5;
   bgs.points[10][0] = 5;
   bgs.points[11][1] = 3;
   bgs.playerTurn = 0
   bgs.FindIfAllCheckersOnHomeBoard() 
}

func (bgs *BackgammonState) InitBeginPosition8() {
   bgs.Init()
   //this.putNewCheckersOnPoint(5, 0, 1);
   bgs.points[23][0] = 0;
   //this.putNewCheckersOnPoint(3, 4, 2);
   bgs.points[22][0] = 2;
   //this.putNewCheckersOnPoint(5, 6, 2);
   bgs.points[21][0] = 2;
   //this.putNewCheckersOnPoint(5, 12, 2);
   bgs.points[20][0] = 3;
   //this.putNewCheckersOnPoint(2, 11, 1);
   bgs.points[19][0] = 3;
   //this.putNewCheckersOnPoint(3, 16, 1);
   bgs.points[18][0] = 5;
   //this.putNewCheckersOnPoint(5, 18, 1);
   //this.putNewCheckersOnPoint(2, 23, 2);
   bgs.points[6][1] = 5;
   bgs.points[7][1] = 3;
   bgs.points[8][1] = 3;
   bgs.points[9][1] = 2;
   bgs.points[10][1] = 2;
   bgs.points[11][1] = 0;
   bgs.playerTurn = 0
   bgs.FindIfAllCheckersOnHomeBoard() 
}

func (bgs *BackgammonState) InitBeginPosition9() {
   bgs.Init()
   //this.putNewCheckersOnPoint(5, 0, 1);
   bgs.points[12][1] = 0;
   //this.putNewCheckersOnPoint(3, 4, 2);
   bgs.points[13][1] = 0;
   //this.putNewCheckersOnPoint(5, 6, 2);
   bgs.points[14][1] = 3;
   //this.putNewCheckersOnPoint(5, 12, 2);
   bgs.points[15][1] = 3;
   //this.putNewCheckersOnPoint(2, 11, 1);
   bgs.points[16][1] = 3;
   //this.putNewCheckersOnPoint(3, 16, 1);
   bgs.points[17][1] = 3;
   bgs.points[18][1] = 3;
   bgs.points[19][0] = 0;

   //this.putNewCheckersOnPoint(5, 18, 1);
   //this.putNewCheckersOnPoint(2, 23, 2);
   bgs.points[0][0] = 0;
   bgs.points[1][0] = 0;
   bgs.points[2][0] = 3;
   bgs.points[3][0] = 3;
   bgs.points[4][0] = 3;
   bgs.points[5][0] = 3;
   bgs.points[6][0] = 3;
   bgs.points[7][0] = 0;
   bgs.playerTurn = 0
   bgs.FindIfAllCheckersOnHomeBoard() 
}

func (bgs *BackgammonState) InitBeginPosition10() {
   bgs.Init()
   //this.putNewCheckersOnPoint(5, 0, 1);
   bgs.points[12][1] = 0;
   //this.putNewCheckersOnPoint(3, 4, 2);
   bgs.points[13][1] = 0;
   //this.putNewCheckersOnPoint(5, 6, 2);
   bgs.points[14][1] = 3;
   //this.putNewCheckersOnPoint(5, 12, 2);
   bgs.points[15][1] = 3;
   //this.putNewCheckersOnPoint(2, 11, 1);
   bgs.points[16][1] = 3;
   //this.putNewCheckersOnPoint(3, 16, 1);
   bgs.points[17][1] = 3;
   bgs.points[18][1] = 3;
   bgs.points[19][0] = 0;

   //this.putNewCheckersOnPoint(5, 18, 1);
   //this.putNewCheckersOnPoint(2, 23, 2);
   bgs.points[0][0] = 0;
   bgs.points[1][0] = 0;
   bgs.points[2][0] = 3;
   bgs.points[3][0] = 3;
   bgs.points[4][0] = 3;
   bgs.points[5][0] = 3;
   bgs.points[6][0] = 3;
   bgs.points[7][0] = 0;
   bgs.playerTurn = 0
   bgs.FindIfAllCheckersOnHomeBoard() 
}

func (bgs *BackgammonState) InitBeginPosition11() {
   bgs.Init()
   //this.putNewCheckersOnPoint(5, 0, 1);
   bgs.points[12][1] = 3;
   //this.putNewCheckersOnPoint(3, 4, 2);
   bgs.points[13][1] = 0;
   //this.putNewCheckersOnPoint(5, 6, 2);
   bgs.points[14][1] = 0;
   //this.putNewCheckerOnPoint(5, 12, 2);
   bgs.points[15][1] = 0;
   //this.putNewCheckersOnPoint(2, 11, 1);
   bgs.points[16][0] = 3;
   //this.putNewCheckersOnPoint(3, 16, 1);
   //bgs.points[17][1] = 3;
   bgs.points[18][0] = 5;
   bgs.points[19][1] = 2;

   //this.putNewCheckersOnPoint(5, 18, 1);
   //this.putNewCheckersOnPoint(2, 23, 2);
   bgs.points[0][0] = 5;
   bgs.points[1][0] = 0;
   bgs.points[2][0] = 3;
   bgs.points[3][1] = 2;
   bgs.points[4][1] = 1;
   bgs.points[5][0] = 3;
   bgs.points[6][1] = 3;
   bgs.points[7][1] = 2;
   bgs.points[8][0] = 1;
   bgs.points[9][1] = 2;
   bgs.bar[0]=1;
   bgs.playerTurn = 0
   bgs.FindIfAllCheckersOnHomeBoard() 
}

