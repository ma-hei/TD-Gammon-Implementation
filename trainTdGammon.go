package main

import (
    "fmt"
    "math"
    "math/rand"
    "strconv"
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

func getScoreOfActivationForPlayer(nn NeuralNetwork, playerTurn int) float64 {
    a := playerTurn
    b := (playerTurn+1) % 2
    return (nn.activations[2].At(a,0) - nn.activations[2].At(b,0))
}

func findBestFollowUpState(selection map[string]BackgammonState, nn NeuralNetwork, playerTurn int, turn int, episode int) (float64, []float64, BackgammonState) {
    bestState := make([]float64, 198)
    if playerTurn == 3 {
        randomState := selectRandomFollowUp(selection)
        nnInput := translateToNNInput(randomState)
        nn.FeedForward(nnInput)
        goodness := getScoreOfActivationForPlayer(nn, playerTurn)
        copy(bestState, nnInput)
        return goodness, nnInput, randomState
    }
    max := 0.0
    max = math.Inf(-1)
    var BgState BackgammonState
    for _, v := range selection { 
        nnInput := translateToNNInput(v)
        nn.FeedForward(nnInput)
        goodness := getScoreOfActivationForPlayer(nn, playerTurn)
        chose := false
        chose = goodness > max
        if chose {
            max = goodness
            copy(bestState, nnInput)
            BgState = v
        } 
    }
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

func trainNN() {
    nn := NeuralNetwork{} 
    nn.Init([]int{192 + 6, 40, 2})
    nn = readNNFromFile("nn_train_plus_199000_episodes.bin")
    wins0 := 0
    wins1 := 0
       for i:=0; i<200000; i++ {
        b := BackgammonState{}
        b.InitBeginPosition()
        b.playerTurn = i%2
        gameOver := false    
        turn := 1
        nn.reInitEtracesAndDerivatives()
        for !gameOver {
            followUpStates, won := b.rollDiceAndFindFollowUpStates(b.playerTurn)
            if won {
                b.playerTurn = ((b.playerTurn + 1) %2)
            }
            currentState := translateToNNInput(b)
            nn.FeedForward(currentState)
            scoreCurrent := getScoreOfActivationForPlayer(nn, b.playerTurn)
            nn.BackpropagateLastActivationPerOutputUnit()
            nn.UpdateEligibilityTraceWithLastDerivativePerOutputUnit()
            scoreNext, bestNextState, newBgState := findBestFollowUpState(followUpStates, nn, b.playerTurn, turn, i)
            if i>200 && (i%10) == 0 {
                fmt.Printf("player turn %v\n", b.playerTurn)
                fmt.Printf("dice1 %v dice2 %v\n", b.dice1, b.dice2)
                fmt.Printf("lm1: %v lm2: %v lm3: %v lm4: %v\n", newBgState.lastMove1, newBgState.lastMove2, newBgState.lastMove3, newBgState.lastMove4)
                drawBgState(newBgState) 
            }
            reward := 0.0
            if won {
                reward = 1.0
                if b.playerTurn == 1 {
                    wins1++
                } else {
                    wins0++
                }
                compareStates(currentState, bestNextState)
                gameOver = true
            }
            nn.TdUpdatePerOutputUnit(reward, scoreNext, scoreCurrent, b.playerTurn, ((b.playerTurn+1)%2))
            copy(currentState, bestNextState)
            b = newBgState 
            if i > 0 && (i%1000 == 0) {
                fileName := "nn_train_plus_plus_" + strconv.Itoa(i) + "_episodes.bin"
                writeNNToFile(fileName, nn)
            }
            turn++
        }
        fmt.Printf("done game %v: %v %v %v %v\n", i, wins0, wins1, float64(wins1)/float64(wins0), turn)
    }
}

func main() {
    //trainNN()
    //testingStuff()
    nn1 := readNNFromFile("nn_train_plus_199000_episodes.bin")
    //nn2 := readNNFromFile("nn_train_plus_plus_150000_episodes.bin")
    //letCompete(nn1, nn2, 50000)
    competeAgainst(nn1) 
}
