package main

import (
	"log"
	"os"
	"encoding/binary"
	"bytes"
	//"math/rand"
	//"time"
)

type intPayload struct {
    n int
}

func readNNFromFile(fileName string) NeuralNetwork {
    file, err := os.Open(fileName)
    defer file.Close()
    if err != nil {
        log.Fatal(err)
    }
    var nLayers int32
    data := readNextBytes(file, 4)
    buffer := bytes.NewBuffer(data)
    err = binary.Read(buffer, binary.BigEndian, &nLayers)
    nodesPerLayer := make([]int , nLayers)
    for l := 0; l < int(nLayers); l++ {
        var nNodesInLayer int32
        data := readNextBytes(file, 4)
        buffer := bytes.NewBuffer(data)
        err = binary.Read(buffer, binary.BigEndian, &nNodesInLayer)
        nodesPerLayer[l] = int(nNodesInLayer)
    }
    nn := NeuralNetwork{}
    nn.Init(nodesPerLayer)
    for l := 0; l < int(nLayers)-1; l++ {
        rows := nodesPerLayer[l+1]
        cols := nodesPerLayer[l]
        for r := 0; r < rows; r ++ {
            for c := 0; c < cols; c++ {
                var weight float64
                data := readNextBytes(file, 8)
                buffer := bytes.NewBuffer(data)
                err = binary.Read(buffer, binary.BigEndian, &weight)
                nn.weights[l].Set(r, c, weight)
            }
        }
        for c := 0; c < nodesPerLayer[l+1]; c++ {
            var bias float64
            data := readNextBytes(file, 8)
            buffer := bytes.NewBuffer(data)
            err = binary.Read(buffer, binary.BigEndian, &bias)
            nn.biases[l].Set(c, 0, bias)
        }
    }
    return nn
}

func writeNNToFile(fileName string, nn NeuralNetwork) {
    file, err := os.Create(fileName)
    defer file.Close()
    if err != nil {
        log.Fatal(err)
    }
    var buf bytes.Buffer
    layers := int32(nn.nLayers)
    binary.Write(&buf, binary.BigEndian, layers)
    writeNextBytes(file, buf.Bytes())
    for l := 0; l < nn.nLayers; l++ {
        nNodesLayer := int32(nn.nNodesInLayer[l])
        var buf bytes.Buffer
        binary.Write(&buf, binary.BigEndian, nNodesLayer)
        writeNextBytes(file, buf.Bytes())
    }
    for l := 0; l < nn.nLayers-1; l++ {
        cols := nn.nNodesInLayer[l]
        rows := nn.nNodesInLayer[l+1]
        for r := 0; r < rows; r++ {
            for c := 0; c < cols; c++ {
                weight := nn.weights[l].At(r,c)
                var buf bytes.Buffer
                binary.Write(&buf, binary.BigEndian, weight)
                writeNextBytes(file, buf.Bytes())
            }
        }
        nNodes := nn.nNodesInLayer[l+1]
        for c := 0; c < nNodes; c++ {
            bias := nn.biases[l].At(c,0)
            var buf bytes.Buffer
            binary.Write(&buf, binary.BigEndian, bias)
            writeNextBytes(file, buf.Bytes())
        }
    }
}

//this type represnts a record with three fields
type payload struct {
	Two float64
}

//func main() {
//    b := BackgammonState{}
//    b.InitBeginPosition()
//    currentState := translateToNNInput(b)
//    nn := NeuralNetwork{} 
//    nn.Init([]int{192 + 6, 40, 2}) 
//    nn.FeedForward(currentState)
//    fmt.Printf("%v\n",nn.activations[2].At(0,0))
//    writeNNToFile("bgState.bin", nn)
//    nn2 := readNNFromFile("bgState.bin")
//    nn2.FeedForward(currentState)
//    fmt.Printf("%v\n",nn2.activations[2].At(0,0))
    //writeFile()
    //readFile()
    //buf := new(bytes.Buffer)
    //var n int32 = 10
    //binary.Write(buf, binary.LittleEndian, n)
    //fmt.Printf("%v",len(buf.Bytes()))
//}

func writeNextBytes(file *os.File, bytes []byte) {

	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}

}

func readNextBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}
