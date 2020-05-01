package main

import (
        "fmt"
        "strconv"
)

var header1 = "| 12 | 13 | 14 | 15 | 16 | 17 | BAR | 18 | 19 | 20 | 21 | 22 | 23 | OFF |\n";
var header2 = "|------- OUTER BOARD -------- |     |--------- O Home Board ------|     |\n";
var middle  = "|---------------------------- |     |-----------------------------|     |\n";
var footer1 = "|------- OUTER BOARD -------- |     |--------- P Home Board ------|     |\n";
var footer2 = "| 11 | 10 | 09 | 08 | 07 | 06 | BAR |  5 |  4 |  3 |  2 |  1 |  0 | OFF |\n";

func printBoardHalf(bg BackgammonState, upper bool) {
    for i := 0; i < 5; i++ {
        p := 0
        if upper { p += 12 }
        for k := 0; k < 14; k++ {
           fmt.Printf("| ")
           if (k < 6 || ( k > 6 && k < 13)) {
               c := bg.points[p][0] + bg.points[p][1];
               s := "O"
               if bg.points[p][1] > 0 { s = "X" }
               f := ""
               row := i + 1
               if !upper { row = 5 - i }
               if (c == 0 || c < row) {
                   f = "  "
               } else if (c >= row && ( row != 5 || c==5)) {
                   f = " " + s 
               } else if (c > 9) {
                   f = strconv.Itoa(c)   
               } else {
                   f = " " + strconv.Itoa(c)
               }
               fmt.Printf(f + " ")
               p++
           } else {
               fmt.Printf("    ")
           }
        }
        fmt.Printf("|\n")
    } 
}

func drawBgState(bg BackgammonState) {
    fmt.Printf(header1) 
    fmt.Printf(header2) 
    printBoardHalf(bg, true)
    fmt.Printf(middle)
    printBoardHalf(bg, false)
    fmt.Printf(footer1)
    fmt.Printf(footer2)
}

func main() {
    bg := BackgammonState{}
    bg.InitBeginPosition()
    drawBgState(bg)
}

