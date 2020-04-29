package main

import (
    "fmt"
    "math"
    "math/rand"
)

func fillInNNInputForPointAndPlayer(point int, player int, b BackgammonState, nnInput []float64) {
    offset := point * 8 + player * 4
    for i := 0; i < 4; i++ {
        nnInput[offset + i] = 0   
    }
    if b.points[point][player] > 0 {
        nnInput[offset + 0] = 1
    }
    if b.points[point][player] > 1 {
    }
    for i := 0; i < 3; i++ {
        if b.points[point][player] > i {
            nnInput[offset + i] = 1
        }
    }
    if b.points[point][player] > 3 {
        nnInput[offset + 3] = float64(b.points[point][player] - 3) / 2.0
    }
}

func translateToNNInput(b BackgammonState) []float64 {
    nnInput := make([]float64, 198)
    for i := 0; i < 24; i++ {
        for k := 0; k < 2; k++ {
            fillInNNInputForPointAndPlayer(i, k, b, nnInput)
        }
    }
    nnInput [192] = float64(b.bar[0]) / 2.0
    nnInput [193] = float64(b.bar[1]) / 2.0
    nnInput [194] = float64(b.checkersBearedOff[0]) / 15.0
    nnInput [195] = float64(b.checkersBearedOff[1]) / 15.0
    if b.playerTurn == 0 {
        nnInput[196] = 1
        nnInput[197] = 0
    } else {
        nnInput[196] = 0
        nnInput[197] = 1
    }
    return nnInput
}

func findBestFollowUpState(selection map[string]BackgammonState, nn NeuralNetwork, playerTurn int) (float64, []float64, BackgammonState) {
    bestState := make([]float64, 198)
    if playerTurn == 0 {
        randomState := selectRandomFollowUp(selection)
        nnInput := translateToNNInput(randomState)
        goodness := nn.activations[2].At(0,0)
        copy(bestState, nnInput)
        return goodness, nnInput, randomState
    }
    //fmt.Printf("playerTurn %v\n", playerTurn)
    max := 0.0
    //if playerTurn == 1 {
    max = math.Inf(-1)
    //} else {
    //   max = math.Inf(1)
    //}
    var BgState BackgammonState
    for _, v := range selection { 
        //v.playerTurn = (v.playerTurn + 1) % 2
        nnInput := translateToNNInput(v)
        nn.FeedForward(nnInput)
        //goodness := nn.activations[2].At(0,0) - nn.activations[2].At(1,0)
        goodness := nn.activations[2].At(0,0)
        //fmt.Printf("%v\n", k)
        //fmt.Printf("goodness: %v\n", goodness)
        chose := false
        chose = goodness > max
        if chose {
            max = goodness
            copy(bestState, nnInput)
            BgState = v
        } 
    }
    //fmt.Printf("chsoing: %v\n", BgState.toString())
    return max, bestState, BgState
}

func getControlState() BackgammonState {
    c := BackgammonState{}
    c.InitBeginPosition()
    c.points[0][0] = 0
    c.points[11][0] = 0
    c.points[16][0] = 0
    c.points[18][0] = 0
    c.points[23][1] = 0
    c.points[0][1] = 0
    c.points[23][0] = 8
    c.points[22][0] = 3
    c.FindNumBearedOffCheckers() 
    c.playerTurn = 0
    return c
}

func controlNN(nn NeuralNetwork, c_ []float64) {
    c_[196] = 1
    c_[197] = 0
    nn.FeedForward(c_)
    temp := nn.activations[2].At(0,0)
    fmt.Printf("control score: %v\n", temp)
}

func selectRandomFollowUp(selection map[string]BackgammonState) BackgammonState {
   asList := make([]BackgammonState, 0, len(selection)) 
   for _, v := range selection {
       asList = append(asList,v)
   }
   stateToReturn := rand.Intn(len(asList))
   return asList[stateToReturn]
}

func compareStates(a []float64, b []float64) {
    for i:=1;i<198;i++ {
        if a[i] != b[i] {
           fmt.Printf("warning %v\n", i)
        }
    }
}

func main() {
    nn := NeuralNetwork{} 
    nn.Init([]int{192 + 6, 80, 1})
    wins0 := 0
    wins1 := 0
    //for i := 0; i < 10000; i++ {
    //    b := BackgammonState {}
    //    b.InitBeginPosition2()
    //    playerTurn := i%2
    //    gameOver := false    
    //    nTurns := 0
    //    for !gameOver {
    //        followUpStates, won := b.rollDiceAndFindFollowUpStates(playerTurn)
    //        nextState := selectRandomFollowUp(followUpStates)
    //        if (won) {
    //            gameOver = true
    //            if nextState.playerTurn == 1 {
    //                wins1++
    //            } else {
    //               wins0++
    //            }
                //fmt.Printf("win %v\n", nextState.playerTurn) 
                //fmt.Printf("%v\n", nextState.toString())
    //        }
    //        b = nextState
    //        playerTurn = (playerTurn + 1)%2
    //        nTurns++
    //    }
    //    fmt.Printf("%v %v %v %v\n", i, wins0, wins1, nTurns)
    //}
    //c.InitBeginPosition3()  
    //followUpStates, _ := c.rollDiceAndFindFollowUpStates(0)
    //fmt.Printf("n follow ups %v\n", len(followUpStates))
    //for k, _ := range followUpStates {
    //   fmt.Printf("%v\n", k)
    //}
    //for i :=0 ; i < 1--; i++ {
    //    b := BackgammonState{}
    //    b.InitBeginPosition3()        
    //    playerTurn := i % 2
    //    b.playerTurn = playerTurn
    //    nn.InitETracesAndDerivatives()
        
    //}

    //b := BackgammonState{}
    //b.InitBeginPosition3()        
    //followUpStates, _ := b.rollDiceAndFindFollowUpStates(1)
    //for k, _ := range followUpStates {
    //    fmt.Printf("%v\n", k)
    //}
    //asList := make([]BackgammonState, 0, len(followUpStates)) 
    //for _, v := range followUpStates {
    //    asList = append(asList,v)
    //}
    //fmt.Printf("%v\n", asList[2].toString())
    //c1 := translateToNNInput(asList[1])
    //c2 := translateToNNInput(asList[0])
    //s := translateToNNInput(asList[2])
    //for i := 0; i < 0; i++ {
    //    nn.FeedForward(c1)
    //    scorec1 := nn.activations[2].At(0,0)
    //    nn.FeedForward(c2)
    //    scorec2 := nn.activations[2].At(0,0)
    //    nn.FeedForward(s)
    //    scores := nn.activations[2].At(0,0)
    //    if (i%200) == 0 {
    //        fmt.Printf("%v %v %v\n", scorec1, scorec2, scores)
    //    }
        //nn.InitETracesAndDerivatives() 
    //    nn.FeedForward(s)
    //    nn.BackpropagateLastActivation()
    //    nn.UpdateEligibilityTraceWithLastDerivative()
    //    nn.TdUpdate(1.0, 0.0, 0.0)
    //}
    for i := 0; i < 100000; i++ {
    //    fmt.Printf("------------\n")
        //controlNN(nn, c_)
        b := BackgammonState{}
        b.InitBeginPosition8()        
        playerTurn := i%2
        b.playerTurn = playerTurn
        gameOver := false    
        turn := 1
        nn.InitETracesAndDerivatives()
        //for k:=0;k<20000;k++{
        for !gameOver {
            followUpStates, won := b.rollDiceAndFindFollowUpStates(playerTurn)
            currentState := translateToNNInput(b)
            nn.FeedForward(currentState)
            //scoreCurrent := nn.activations[2].At(0,0) - nn.activations[2].At(1,0)
            scoreCurrent := nn.activations[2].At(0,0)
            nn.BackpropagateLastActivation()
            nn.UpdateEligibilityTraceWithLastDerivative()
            //fmt.Printf("player turn: %v\n", playerTurn)
            //fmt.Printf("n follow ups %v\n", len(followUpStates))
            //for k, _ := range followUpStates {
            //    fmt.Printf("%v\n", k)
            //}
            isWin := won
            scoreNext, bestNextState, newBgState := findBestFollowUpState(followUpStates, nn, playerTurn)
            reward := 0.0
            if isWin {
                playerWon := (playerTurn + 1) % 2
                if playerWon == 1 {
                    reward = 1.0
                    fmt.Printf("player 1 won\n")
                    //reward = 1.0
                    wins1++
                } else {
                    fmt.Printf("player 0 won\n")
                    //fmt.Printf("%v\n", newBgState.toString())
                    wins0++
                }
                compareStates(currentState, bestNextState)
                gameOver = true
    //            fmt.Printf("win: %v\n", playerWon)
    //            fmt.Printf("last state: %v \n", b.toString())
            }
            nn.TdUpdate(reward, scoreNext, scoreCurrent)
            if isWin {
                fmt.Printf("the new score                               %v\n", scoreNext)
                fmt.Printf("what I previously thought about this state: %v\n", scoreCurrent)
                nn.FeedForward(currentState)
                scoreCurrent := nn.activations[2].At(0,0)
                fmt.Printf("what I think now                          : %v\n", scoreCurrent)
                //nn.FeedForward(c1)
                //scorec1 := nn.activations[2].At(0,0)
                //nn.FeedForward(c2)
                //scorec2 := nn.activations[2].At(0,0)
                //nn.FeedForward(s)
                //scores := nn.activations[2].At(0,0)
                //fmt.Printf("%v %v %v\n", scorec1, scorec2, scores)
            }
    //        if (isWin) {
    //            fmt.Printf("next minus previous                  : %v\n", (scoreNext - scorePrev))
    //            nn.FeedForward(statePrev)
    //            temp := nn.activations[2].At(0,0) - nn.activations[2].At(1,0)
    //            fmt.Printf("difference of score before after     : %v\n", (temp - scorePrev))
    //        }
            copy(currentState, bestNextState)
            playerTurn = (playerTurn + 1) % 2
            b = newBgState 
            turn++
            //gameOver = true
        }
        fmt.Printf("done %v: %v, %v, %v\n", i, wins0, wins1, turn)
    }
    //b := BackgammonState{}
    //b.InitBeginPosition3()        
    //currentState := translateToNNInput(b)
    //nn.FeedForward(currentState)
    //fmt.Printf("--> %v %v\n", nn.activations[2].At(0,0), nn.activations[2].At(1,0))
    //nn.BackpropagateLastActivation()
    //nn.UpdateEligibilityTraceWithLastDerivative()
    //nn.TdUpdate(1.0, 0.0, 0.0)
    //nn.FeedForward(currentState)
    //fmt.Printf("--> %v %v\n", nn.activations[2].At(0,0), nn.activations[2].At(1,0))
}

