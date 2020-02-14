package main

import (
    "gonum.org/v1/gonum/mat"
    "math/rand"
    "fmt"
) 

type NeuralNetwork struct {
    nNodesInLayer []int
    nLayers int
    weights []*mat.Dense
    biases []*mat.Dense
    activations []*mat.Dense
    z []*mat.Dense
}

func (nn *NeuralNetwork) InitBiases() {
    nn.biases = make([]*mat.Dense, nn.nLayers)
    for i := 0; i < nn.nLayers; i++ {
        nNodes := nn.nNodesInLayer[i]
        data := GenerateRandomValues(nNodes)
        nn.biases[i] = mat.NewDense(nNodes, 1, data)
    }
}

func (nn *NeuralNetwork) InitWeights() {
    nn.weights = make([]*mat.Dense, nn.nLayers)
    for i := 1; i < nn.nLayers; i++ {
        nNodesPrevious := nn.nNodesInLayer[i-1]
        nNodesCurrent := nn.nNodesInLayer[i]
        data := GenerateRandomValues(nNodesPrevious * nNodesCurrent)
        nn.weights[i-1] = mat.NewDense(nNodesCurrent, nNodesPrevious, data)
    }
}

func (nn *NeuralNetwork) InitActivations() {
    nn.activations = make([]*mat.Dense, nn.nLayers)
    nn.z = make([]*mat.Dense, nn.nLayers)
    for i := 0; i < nn.nLayers; i++ {
        nNodesInLayer := nn.nNodesInLayer[i]
        nn.activations[i] = mat.NewDense(nNodesInLayer, 1, nil)
        nn.z[i] = mat.NewDense(nNodesInLayer, 1, nil)
    }
}

func (nn *NeuralNetwork) Init(nNodesInLayer []int) {
    nn.nNodesInLayer = nNodesInLayer
    nn.nLayers = len(nNodesInLayer)
    nn.InitWeights()
    nn.InitBiases()
    nn.InitActivations()
}

func GenerateRandomValues(length int) []float64 {
    data := make([]float64, length)
    for k := range data {
        data[k] = rand.NormFloat64()
    }
    return data
}

func (nn *NeuralNetwork) FeedForward(input []float64) {
    nn.activations[0].SetCol(0, input)
    for i := 1; i < nn.nLayers; i++ {
       PrintMatrix(nn.weights[i-1])
       PrintMatrix(nn.activations[i-1])
       nn.z[i].Mul(nn.weights[i-1], nn.activations[i-1])
       nn.z[i].Add(nn.z[i], nn.biases[i])
    }
}

func PrintMatrix(m *mat.Dense) {
    a := mat.Formatted(m, mat.Prefix(""), mat.Squeeze())
    fmt.Printf("\n%v\n", a) 
}

func main() {
    nn := NeuralNetwork{} 
    nn.Init([]int{3,4,1})
    nn.FeedForward([]float64{1,2,3})
}


