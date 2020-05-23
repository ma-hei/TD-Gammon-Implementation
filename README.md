# Backgammon

The goal of this project was to implement TD gammon from scratch (Gerald Tesauro, 1995). 

## How to play

Install go version 1.13 and checkout this repo. Then run

<code>go run nn.go arena.go backgammon.go readingWriting.go trainTdGammon.go boardDrawingCommandLine.go</code>

## Some Details

This project was an exercise to get a deeper understanding of reinforcement learning. This exercise was concluded after observing that the model learned to play backgammon quite well. By "quite well", I mean that the model beat me in ~38/40 test games. No further parameter tuning was attempted after this. For example, the model was not trained with an extra reward for a "gammon" as suggested in the paper.

Golang was used in this project just because I wanted to try out the language. It is probably not the ideal language to build a neural network. The code hasn't been cleaned up yet.
