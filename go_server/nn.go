package main

import (
    "gonum.org/v1/gonum/mat"
    "math"
    "math/rand"
    "fmt"
    "os"
    "time"
    "encoding/binary"
) 

type NeuralNetwork struct {
    nNodesInLayer []int
    nLayers int
    weights []*mat.Dense
    biases []*mat.Dense
    activations []*mat.Dense
    z []*mat.Dense
    errors []*mat.Dense
    etracesWeights []*mat.Dense
    etracesBiases []*mat.Dense
    derivativeWeights []*mat.Dense
    derivativeBiases []*mat.Dense
    
    etracesWeightsSplit [][]*mat.Dense
    etracesBiasesSplit [][]*mat.Dense
    derivativeWeightsSplit [][]*mat.Dense
    derivativeBiasesSplit [][]*mat.Dense

}

func (nn *NeuralNetwork) InitBiases() {
    nn.biases = make([]*mat.Dense, nn.nLayers)
    for i := 1; i < nn.nLayers; i++ {
        nNodes := nn.nNodesInLayer[i]
        data := GenerateRandomValues(nNodes)
        nn.biases[i-1] = mat.NewDense(nNodes, 1, data)
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
    nn.errors = make([]*mat.Dense, nn.nLayers)
    for i := 0; i < nn.nLayers; i++ {
        nNodesInLayer := nn.nNodesInLayer[i]
        nn.activations[i] = mat.NewDense(nNodesInLayer, 1, nil)
        nn.z[i] = mat.NewDense(nNodesInLayer, 1, nil)
        nn.errors[i] = mat.NewDense(nNodesInLayer, 1, nil)
    }
}

func (nn *NeuralNetwork) Init(nNodesInLayer []int) {
    nn.nNodesInLayer = nNodesInLayer
    nn.nLayers = len(nNodesInLayer)
    nn.InitWeights()
    nn.InitBiases()
    nn.InitActivations()
    nn.InitETracesAndDerivatives()
}

func (nn *NeuralNetwork) InitETracesAndDerivatives() {
    nn.etracesWeightsSplit = make([][]*mat.Dense, 2)
    nn.etracesWeightsSplit[0] = make([]*mat.Dense, nn.nLayers)
    nn.etracesWeightsSplit[1] = make([]*mat.Dense, nn.nLayers)
    nn.etracesWeights = make([]*mat.Dense, nn.nLayers)
    nn.derivativeWeightsSplit = make([][]*mat.Dense, 2)
    nn.derivativeWeightsSplit[0] = make([]*mat.Dense, nn.nLayers)
    nn.derivativeWeightsSplit[1] = make([]*mat.Dense, nn.nLayers)
    nn.derivativeWeights = make([]*mat.Dense, nn.nLayers)
    for i := 1; i < nn.nLayers; i++ {
        nNodesPrevious := nn.nNodesInLayer[i-1]
        nNodesCurrent := nn.nNodesInLayer[i]
        nn.etracesWeights[i-1] = mat.NewDense(nNodesCurrent, nNodesPrevious, nil)
        nn.derivativeWeights[i-1] = mat.NewDense(nNodesCurrent, nNodesPrevious, nil)
        nn.derivativeWeights[i-1].Zero()
        nn.etracesWeights[i-1].Zero()
        nn.etracesWeightsSplit[0][i-1] = mat.NewDense(nNodesCurrent, nNodesPrevious, nil)
        nn.etracesWeightsSplit[0][i-1].Zero()
        nn.etracesWeightsSplit[1][i-1] = mat.NewDense(nNodesCurrent, nNodesPrevious, nil)
        nn.etracesWeightsSplit[1][i-1].Zero()
        nn.derivativeWeightsSplit[0][i-1] =  mat.NewDense(nNodesCurrent, nNodesPrevious, nil)
        nn.derivativeWeightsSplit[0][i-1].Zero()
        nn.derivativeWeightsSplit[1][i-1] = mat.NewDense(nNodesCurrent, nNodesPrevious, nil)
        nn.derivativeWeightsSplit[1][i-1].Zero()
    }
    nn.etracesBiasesSplit = make([][]*mat.Dense, 2)
    nn.etracesBiasesSplit[0] = make([]*mat.Dense, nn.nLayers)
    nn.etracesBiasesSplit[1] = make([]*mat.Dense, nn.nLayers)
    nn.etracesBiases = make([]*mat.Dense, nn.nLayers)
    nn.derivativeBiasesSplit = make([][]*mat.Dense, 2)
    nn.derivativeBiasesSplit[0] = make([]*mat.Dense, nn.nLayers) 
    nn.derivativeBiasesSplit[1] = make([]*mat.Dense, nn.nLayers)
    nn.derivativeBiases = make([]*mat.Dense, nn.nLayers)
    for i := 1; i < nn.nLayers; i++ {
        nNodes := nn.nNodesInLayer[i]
        nn.derivativeBiases[i-1] = mat.NewDense(nNodes, 1, nil)
        nn.etracesBiases[i-1] = mat.NewDense(nNodes, 1, nil)
        nn.derivativeBiases[i-1].Zero()
        nn.etracesBiases[i-1].Zero()
        nn.etracesBiasesSplit[0][i-1] = mat.NewDense(nNodes, 1, nil)
        nn.etracesBiasesSplit[0][i-1].Zero()
        nn.etracesBiasesSplit[1][i-1] = mat.NewDense(nNodes, 1, nil)
        nn.etracesBiasesSplit[1][i-1].Zero()
        nn.derivativeBiasesSplit[0][i-1] = mat.NewDense(nNodes, 1, nil)
        nn.derivativeBiasesSplit[0][i-1].Zero()
        nn.derivativeBiasesSplit[1][i-1] = mat.NewDense(nNodes, 1, nil)
        nn.derivativeBiasesSplit[1][i-1].Zero()
    }
}

func GenerateRandomValues(length int) []float64 {
    data := make([]float64, length)
    for k := range data {
        rand.Seed(time.Now().UnixNano())
        temp := rand.NormFloat64()
        temp = temp/10000.0
        data[k] = temp
    }
    return data
}

func (nn *NeuralNetwork) FeedForward(input []float64) ([]*mat.Dense, []*mat.Dense) {
    nn.activations[0].SetCol(0, input)
    for i := 1; i < nn.nLayers; i++ {
       nn.z[i].Mul(nn.weights[i-1], nn.activations[i-1])
       nn.z[i].Add(nn.z[i], nn.biases[i-1])
       nn.activations[i].Apply(SigmoidWrapper, nn.z[i])
    }
    return nn.z, nn.activations
}

//func (nn *NeuralNetwork) BackpropagateLastActivation2() []*mat.Dense {
//    f1 := mat.NewDense(nn.nNodesInLayer[
//}

// computes derivative of activation in final layer wrt. weights
func (nn *NeuralNetwork) BackpropagateLastActivation() []*mat.Dense{
     f1 := mat.NewDense(nn.nNodesInLayer[nn.nLayers - 1], 1, nil)
     for i := 0; i < nn.nNodesInLayer[nn.nLayers - 1]; i++ {
         f1.Set(i,0,1)
     }
     f2 := mat.DenseCopyOf(nn.z[nn.nLayers - 1])
     f2.Apply(SigmoidDiffWrapper, f2) 
     nn.errors[nn.nLayers - 1] = HadamardProduct(f1, f2)
     // compute derivative of last layer weights
     nn.derivativeWeights[nn.nLayers - 2].Mul(nn.errors[nn.nLayers - 1], nn.activations[nn.nLayers-2].T())
     nn.derivativeBiases[nn.nLayers - 2].Copy(nn.errors[nn.nLayers - 1])
     for i := nn.nLayers - 2; i > 0; i-- {
         temp := nn.weights[i].T()
         Wtranspose := mat.DenseCopyOf(temp)
         nn.errors[i].Mul(Wtranspose, nn.errors[i+1])
         f2 := mat.DenseCopyOf(nn.z[i])
         f2.Apply(SigmoidDiffWrapper, f2) 
         nn.errors[i] = HadamardProduct(nn.errors[i], f2)
         nn.derivativeWeights[i - 1].Mul(nn.errors[i], nn.activations[i-1].T())
         nn.derivativeBiases[i - 1].Copy(nn.errors[i])
     }
     return nn.errors
}

func (nn *NeuralNetwork) UpdateEligibilityTraceWithLastDerivative() {
    learnRate := 0.2
    for i := 0; i < nn.nLayers-1; i++ {
        nn.etracesWeights[i].Scale(learnRate, nn.etracesWeights[i])
        nn.etracesWeights[i].Add(nn.etracesWeights[i], nn.derivativeWeights[i])
        nn.etracesBiases[i].Scale(learnRate, nn.etracesBiases[i])
        nn.etracesBiases[i].Add(nn.etracesBiases[i], nn.derivativeBiases[i])
    }
}

func (nn *NeuralNetwork) TdUpdate(reward float64, newStateValue float64, oldStateValue float64) {
    learnRate := 0.05
    for i := 0; i < nn.nLayers-1; i++ {
        temp := learnRate * (reward + newStateValue - oldStateValue)
        temp2 := mat.DenseCopyOf(nn.etracesWeights[i])
        temp2.Scale(temp, temp2)
        nn.weights[i].Add(nn.weights[i], temp2)
        temp3 := mat.DenseCopyOf(nn.etracesBiases[i])
        temp3.Scale(temp, temp3)
        nn.biases[i].Add(nn.biases[i], temp3)
    }
}
// computes derivative of error function 1/2(y-y')^2 wrt.weights
func (nn *NeuralNetwork) Backpropagate(input []float64) []*mat.Dense{
     f1 := mat.NewDense(len(input), 1, nil)
     for i := 0; i < len(input); i++ {
         f1.Set(i,0,input[i])
     }
     f1.Sub(nn.activations[nn.nLayers - 1], f1)
     f2 := mat.DenseCopyOf(nn.z[nn.nLayers - 1])
     f2.Apply(SigmoidDiffWrapper, f2) 
     nn.errors[nn.nLayers - 1] = HadamardProduct(f1, f2)
     for i := nn.nLayers - 2; i > 0; i-- {
         temp := nn.weights[i].T()
         Wtranspose := mat.DenseCopyOf(temp)
         nn.errors[i].Mul(Wtranspose, nn.errors[i+1])
         f2 := mat.DenseCopyOf(nn.z[i])
         f2.Apply(SigmoidDiffWrapper, f2) 
         nn.errors[i] = HadamardProduct(nn.errors[i], f2)
     }
     return nn.errors
}

func PrintMatrix(m *mat.Dense) {
    a := mat.Formatted(m, mat.Prefix(""), mat.Squeeze())
    fmt.Printf("\n%v\n", a) 
}

func SigmoidWrapper(i, j int, v float64) float64{
    return Sigmoid(v)
}

func Sigmoid(val float64) float64{
    return 1/(1+(math.Pow(math.E, -1 * val)))
}

func SigmoidDiffWrapper(i, j int, v float64) float64{
    return SigmoidDiff(v)
}

func SigmoidDiff(val float64) float64 {
    return Sigmoid(val) * (1 - Sigmoid(val))
}

func HadamardProduct(a, b *mat.Dense) *mat.Dense {
    nRow, _ := a.Dims()
    result := mat.NewDense(nRow, 1, nil)
    for i := 0; i < nRow; i++ {
        result.Set(i, 0, a.At(i, 0) * b.At(i, 0))
    }
    return result
}

func ReadImagesFromFile(fileName string) [][]float64 {
    f, _ := os.Open(fileName)
    temp := make([]byte, 4)
    f.Read(temp)
    f.Read(temp)
    nTrainImages := int(binary.BigEndian.Uint32(temp))
    trainingImages := make([][]float64, nTrainImages)
    f.Read(temp)
    nCols := int(binary.BigEndian.Uint32(temp))
    f.Read(temp)
    nRows := int(binary.BigEndian.Uint32(temp))
    for i := 0; i < nTrainImages; i++ {
        trainingImages[i] = make([]float64, nRows * nCols)
        image := make([]byte, nRows * nCols)
        f.Read(image)
        for k := 0; k < nRows * nCols; k++ {
            trainingImages[i][k] = float64(image[k])/255.0
        }
    }
    return trainingImages
}

func ReadLabelsFromFile(fileName string) []int {
    f, _ := os.Open(fileName)
    temp := make([]byte, 4)
    f.Read(temp)
    f.Read(temp)
    nLabels := int(binary.BigEndian.Uint32(temp))
    labels := make([]int, nLabels)
    label := make([]byte, 1)
    for i := 0; i < nLabels; i++ {
       f.Read(label) 
       labels[i] = int(label[0])
    }
    return labels
}  

// Takes a slice of labels [1,3,9,5,9,...] and
// creates a slice of slices where each slice
// represents a single label:
// [
// [0,1,0,0,0,0,0,0,0,0],
// [0,0,0,1,0,0,0,0,0,0],
// [0,0,0,0,0,0,0,0,0,1],
// [0,0,0,0,0,1,0,0,0,0],
// ...]
func MakeTrainingLabels(labels []int) [][]float64 {
    nLabels := len(labels)
    result := make([][]float64, nLabels)
    for i := 0; i < nLabels; i++ {
        result[i] = make([]float64, 10)
        result[i][labels[i]] = 1
    }
    return result
}

func PrintImage(i int, t[][]float64) {
    for l := 0; l<28; l++ {
        for k := 0; k<28; k++ {
            temp := 1
            if t[i][l*28+k] == 0 {
                temp = 0
            }
            fmt.Printf("%v ", temp)
        }
        fmt.Printf("\n")
    }
}

func (nn *NeuralNetwork) UpdateWeights(learnRate float64) {
    for i := nn.nLayers-1; i > 0; i-- {
        r,c := nn.weights[i-1].Dims()
        temp := mat.NewDense(r, c, nil)    
        temp.Mul(nn.errors[i], nn.activations[i-1].T())
        temp.Scale(learnRate, temp)
        nn.weights[i-1].Sub(nn.weights[i-1], temp)
        temp2 := mat.NewDense(nn.nNodesInLayer[i], 1, nil)
        temp2.Scale(learnRate, nn.errors[i])
        nn.biases[i-1].Sub(nn.biases[i-1], temp2)
    }
}

func (nn *NeuralNetwork) Classify(in []float64) int {
    nn.FeedForward(in)
    max := 0.0
    idx := -1
    for i := 0; i < 10; i++ {
        if (nn.activations[2].At(i, 0) > max) {
            max = nn.activations[2].At(i, 0)
            idx = i
        }
    }
    return idx
}

func (nn *NeuralNetwork) Evaluate(test [][]float64, labels []int) float64{
     correct := 0
     nTest := len(test)
     for i := 0; i < nTest; i++ {
        pred := nn.Classify(test[i])
        trueLabel := labels[i]
        if pred == trueLabel {
            correct++
        }
    }
    return float64(correct)/float64(nTest)
}

// n = 10 -> [0,3,5,9,1,2,4,8,6,7]
func GetShuffledIndexes(n int) []int {
    temp := make([]int, n)
    for i := 0; i < n; i++ {
        temp[i] = i
    }
    rand.Shuffle(len(temp), func(i, j int) { temp[i], temp[j] = temp[j], temp[i] })
    return temp
}

//func main() {
//    nn := NeuralNetwork{} 
//    nn.Init([]int{28*28, 70, 10})
//    trainImages := ReadImagesFromFile("train-images-idx3-ubyte")
//    trainLabels := MakeTrainingLabels(ReadLabelsFromFile("train-labels-idx1-ubyte"))
//    testImages := ReadImagesFromFile("t10k-images-idx3-ubyte")
//    testLabels := ReadLabelsFromFile("t10k-labels-idx1-ubyte")
//    nTrain := len(trainImages)
//
//   nEpochs := 100
//    learnRate := 1.2
//    for k := 0; k < nEpochs; k++ {
//        fmt.Printf("Epoch %v..", k)
//        idx := GetShuffledIndexes(nTrain)
//        for i := 0; i < nTrain; i++ {
//            nn.FeedForward(trainImages[idx[i]])
//            nn.Backpropagate(trainLabels[idx[i]])
//            nn.UpdateWeights(learnRate)
//        }
//        score := nn.Evaluate(testImages, testLabels)
//        fmt.Printf(" Score: %v\n", score)
//    }
//}


