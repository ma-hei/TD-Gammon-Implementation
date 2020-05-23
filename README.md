# Backgammon

The goal of this project was to implement TD gammon from scratch (Gerald Tesauro, 1995). 

The layout of the backgammon board on the command line was copied from https://github.com/dellalibera/gym-backgammon.

## How to play

Install go version 1.13 and checkout this repo. Then run

<code>go run nn.go arena.go backgammon.go readingWriting.go trainTdGammon.go boardDrawingCommandLine.go</code>

The CPU begins. After the CPU did their move, you are presented with the current board state, your dice roll and your options:
<code>
nn move:  11.16 11.16 16.21 16.21                                                                                                                                                                              
| 12 | 13 | 14 | 15 | 16 | 17 | BAR | 18 | 19 | 20 | 21 | 22 | 23 | OFF |                                                                                                                                      
|------- OUTER BOARD -------- |     |--------- O Home Board ------|     |
|  X |    |    |    |  O |    |   0 |  O |    |    |  O |    |  X |   0 |
|  X |    |    |    |  O |    |     |  O |    |    |  O |    |  X |     |
|  X |    |    |    |  O |    |     |  O |    |    |    |    |    |     |
|  X |    |    |    |    |    |     |  O |    |    |    |    |    |     |
|  X |    |    |    |    |    |     |  O |    |    |    |    |    |     |
|---------------------------- |     |-----------------------------|     |
|    |    |    |    |    |    |     |  X |    |    |    |    |    |     |
|    |    |    |    |    |    |     |  X |    |    |    |    |    |     |
|  O |    |    |    |  X |    |     |  X |    |    |    |    |    |     |
|  O |    |    |    |  X |    |     |  X |    |    |    |    |  O |     |
|  O |    |    |    |  X |    |   0 |  X |    |    |    |    |  O |   0 |
|------- OUTER BOARD -------- |     |--------- P Home Board ------|     |
| 11 | 10 | 09 | 08 | 07 | 06 | BAR |  5 |  4 |  3 |  2 |  1 |  0 | OFF |
found 33
dice roll: 2 2
1:  7.5 5.3 12.10 10.8                   3.0,0.0,0.0,0.1,0.2,0.0,0.5,0.0,0.1,0.0,0.0,2.0|0.4,0.0,0.0,0.0,3.0,0.0,5.0,0.0,0.0,2.0,0.0,0.2,0.0:
2:  7.5 5.3 12.10 12.10                  3.0,0.2,0.0,0.0,0.2,0.0,0.5,0.0,0.1,0.0,0.0,2.0|0.3,0.0,0.0,0.0,3.0,0.0,5.0,0.0,0.0,2.0,0.0,0.2,0.0:
3:  5.3 5.3 3.1 12.10                    3.0,0.1,0.0,0.0,0.3,0.0,0.3,0.0,0.1,0.0,0.1,2.0|0.4,0.0,0.0,0.0,3.0,0.0,5.0,0.0,0.0,2.0,0.0,0.2,0.0:
4:  5.3 3.1 12.10 12.10                  3.0,0.2,0.0,0.0,0.3,0.0,0.4,0.0,0.0,0.0,0.1,2.0|0.3,0.0,0.0,0.0,3.0,0.0,5.0,0.0,0.0,2.0,0.0,0.2,0.0:
5:  5.3 12.10 10.8 12.10                 3.0,0.1,0.0,0.1,0.3,0.0,0.4,0.0,0.1,0.0,0.0,2.0|0.3,0.0,0.0,0.0,3.0,0.0,5.0,0.0,0.0,2.0,0.0,0.2,0.0:
6:  12.10 10.8 8.6 6.4                   3.0,0.0,0.0,0.0,0.3,0.0,0.5,0.1,0.0,0.0,0.0,2.0|0.4,0.0,0.0,0.0,3.0,0.0,5.0,0.0,0.0,2.0,0.0,0.2,0.0:
</code>

You can now enter 1, 2, 3,... . For example, if you enter 2, you will move one checker from point 7 to point 5, one checker from point 5 to point 3, one checker from point 12 to 10 and another checker from point 12 to 10. Note that each possible "post move state" is represented by exactly one move-option. So if two different moves lead to the same "post move state", then only one of them is presented. (Sorry for bad UX).

## Some Details

This project was an exercise to get a deeper understanding of reinforcement learning. This exercise was concluded after observing that the model learned to play backgammon quite well. By "quite well", I mean that the model beat me in ~38/40 test games. No further parameter tuning was attempted after this. For example, the model was not trained with an extra reward for a "gammon" as suggested in the paper.

Golang was used in this project just because I wanted to try out the language. It is probably not the ideal language to build a neural network. The code hasn't been cleaned up yet.
