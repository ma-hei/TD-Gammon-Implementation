package main

import (
    "fmt"
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

func main() {
    b := BackgammonState{}
    b.InitBeginPosition()        
    nn := NeuralNetwork{} 
    nn.Init([]int{192 + 6, 40, 1})
    temp := translateToNNInput(b) 
    for i :=0 ; i < 198; i++ {
        fmt.Printf("%v: %v\n", (i + 1), temp[i])
    }
}

