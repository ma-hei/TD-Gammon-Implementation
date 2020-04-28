package main

import (
    "fmt"
    "math"
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

func findBestFollowUpState(selection map[string]BackgammonState, nn NeuralNetwork, playerTurn int) (float64, []float64, bool, *BackgammonState) {
    max := math.Inf(-1)
    bestState := make([]float64, 198)
    isWin := false
    var BgState *BackgammonState
    for _, v := range selection { 
        nnInput := translateToNNInput(v)
        nn.FeedForward(nnInput)
        goodNess := nn.activations[2].At(0,0)
        if goodNess > max {
            max = goodNess
            copy(nnInput, bestState)
            if v.PlayerWins(playerTurn) {
                isWin = true
            }
            BgState = &v
        } 
    }
    return max, bestState, isWin, BgState
}

func main() {
    nn := NeuralNetwork{} 
    nn.Init([]int{192 + 6, 40, 1})
    for i := 0; i< 1000; i++ {
        b := BackgammonState{}
        b.InitBeginPosition()        
        statePrev := translateToNNInput(b) 
        nn.FeedForward(statePrev)
        scorePrev := nn.activations[2].At(0,0)
        playerTurn := 0 
        gameOver := false    
        turn := 0
        for !gameOver {
            followUpStates := b.rollDiceAndFindFollowUpStates(playerTurn)
            if len(followUpStates) == 0 {
                followUpStates[b.toString()] = b
            }
            scoreNext, bestNextState, isWin, newBgState := findBestFollowUpState(followUpStates, nn, playerTurn)
            nn.FeedForward(statePrev)
            nn.BackpropagateLastActivation()
            nn.UpdateEligibilityTraceWithLastDerivative()
            reward := 0
            if isWin {
                reward = 1
                gameOver = true
                fmt.Printf("win: %v\n", playerTurn)
            }
            nn.TdUpdate(reward, scoreNext, scorePrev)
            copy(statePrev, bestNextState)
            turn++
            playerTurn = (playerTurn + 1) % 2
            b = *newBgState 
            fmt.Printf("done %v: %v\n", i, turn)
        }
    }
}

