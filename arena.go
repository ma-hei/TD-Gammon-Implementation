package main

import (
    "fmt"
    "strings"
    "strconv"
) 

func replacePointNumberInMove(n int) string {
    if (n <= 11) {
        return strconv.Itoa(11 - n)
    }
    return strconv.Itoa(n)
}

func isPoint(s string) bool {
    return s != "off" && s != "b"
}

func replacePointNumbersInMove(move string) string {
    if move == "" {
        return move
    }
    temp := strings.Split(move, ".")
    a := ""
    if isPoint(temp[0]) {
        isHit := strings.Contains(temp[0], "h")
        asInt := 0
        if isHit {
            asInt, _ = strconv.Atoi(temp[0][0:(len(temp[0])-1)])
        } else {
            asInt, _ = strconv.Atoi(temp[0])
        }
        a = replacePointNumberInMove(asInt)
        if isHit {
            a = a + "h"
        }
    } else {
        a = temp[0]
    }

    b := ""
    if isPoint(temp[1]) {
        isHit := strings.Contains(temp[1], "h")
        asInt := 0
        if isHit {
            asInt, _ = strconv.Atoi(temp[1][0:(len(temp[1])-1)])
        } else {
            asInt, _ = strconv.Atoi(temp[1])
        }
        b = replacePointNumberInMove(asInt)
        if isHit {
            b = b + "h"
        }
    } else {
        b = temp[1]
    }
    return a + "." + b
}

func appendMove(move string, fullString string) string {
    if (move != "") {
        return fullString + " " + move
    }
    return fullString
}

func getStateFromUser(possibleStates map[string]BackgammonState) BackgammonState {
    count := 1
    selection := make(map[int]BackgammonState)
    for _, v := range possibleStates { 
        //fmt.Printf("%v\n", v.toString())
        move1 := replacePointNumbersInMove(v.lastMove1)
        move2 := replacePointNumbersInMove(v.lastMove2)
        move3 := replacePointNumbersInMove(v.lastMove3)
        move4 := replacePointNumbersInMove(v.lastMove4)
        fullString := ""
        fullString = appendMove(move1, fullString) 
        fullString = appendMove(move2, fullString) 
        fullString = appendMove(move3, fullString) 
        fullString = appendMove(move4, fullString) 
        selection[count] = v
        fill := 30 - len(fullString)
        for k :=0; k< fill; k++ {
            fullString += " "
        }
        fmt.Printf("%v: %v \t %v\n", count, fullString, v.toString())
        count++
    }
    var i int
    fmt.Scanf("%d", &i)
    return selection[i] 
}

func printLastMove(v BackgammonState) {
     move1 := replacePointNumbersInMove(v.lastMove1)
     move2 := replacePointNumbersInMove(v.lastMove2)
     move3 := replacePointNumbersInMove(v.lastMove3)
     move4 := replacePointNumbersInMove(v.lastMove4)
     fullString := ""
     fullString = appendMove(move1, fullString) 
     fullString = appendMove(move2, fullString) 
     fullString = appendMove(move3, fullString) 
     fullString = appendMove(move4, fullString) 
     fmt.Printf("nn move: %v\n", fullString)
}

func competeAgainst(nn NeuralNetwork) {
    b := BackgammonState{}
    b.InitBeginPosition()        
    gameOver := false    
    turn := 1
    for !gameOver {
        drawBgState(b) 
        followUpStates, won := b.rollDiceAndFindFollowUpStates(b.playerTurn)
        gameOver = won
        var newBgState BackgammonState
        if (b.playerTurn == 0) {
            fmt.Printf("dice roll: %v %v\n", b.dice1, b.dice2)
            _, _, newBgState = findBestFollowUpState(followUpStates, nn, b.playerTurn, turn, 0)
            printLastMove(newBgState)
        } else {
            fmt.Printf("dice roll: %v %v\n", b.dice1, b.dice2)
            newBgState = getStateFromUser(followUpStates)
        }
        b = newBgState
    }
}

func letCompete(nn1 NeuralNetwork, nn2 NeuralNetwork, ngames int) {
    win0 := 0
    win1 := 0
    for i := 0; i < ngames; i++ {
        b := BackgammonState{}
        b.InitBeginPosition()        
        b.playerTurn = i%2
        gameOver := false    
        turn := 1
        for !gameOver {
            followUpStates, won := b.rollDiceAndFindFollowUpStates(b.playerTurn)
            if won {
                b.playerTurn = ((b.playerTurn + 1) %2)
            }
            var newBgState BackgammonState
            if (b.playerTurn == 0) {
                _, _, newBgState = findBestFollowUpState(followUpStates, nn1, b.playerTurn, turn, i)
            } else {
                _, _, newBgState = findBestFollowUpState(followUpStates, nn2, b.playerTurn, turn, i)
            }
            if won {
                if b.playerTurn == 0 {
                    win0++
                } else {
                    win1++
                }
                gameOver = true
            }
            b = newBgState
        }
        fmt.Printf("%v %v\n", win0, win1)
    } 
}

func testingStuff() {
    b := BackgammonState{}
    b.InitBeginPosition()        
    b.playerTurn = 1
    
    a, _ := b.rollDiceAndFindFollowUpStates(b.playerTurn)
    getStateFromUser(a)
}
