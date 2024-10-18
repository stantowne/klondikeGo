package main

import (
	"fmt"
	"sort"
	"strings"
)

func playAllMoveS(bIn board, doThisMove move, moveNum int, deckNum int) {
	if verbose > 2 {
		printBoard(bIn)
		// printMove(doThisMove)  NOT written yet
		fmt.Printf("moveNum: %v\n gameLengthLimit: %v", moveNum, gameLengthLimit)
	}

	// Make the indicated move
	var bNew board
	if moveNum == 0 {
		bNew = bIn
		if printTree == "C" {
			fmt.Printf("\n\nDeck %v\n", deckNum)
			printBoard(bIn)
			fmt.Printf("\n\n Strat #")
			for i := 1; i <= 200; i++ {
				fmt.Printf("    %3v ", i)
			}
			fmt.Printf("\n")
		}
	} else {
		bNew = moveMaker(bIn, doThisMove)
		moveNum++
		MovesTriedThisDeck++
		if printTree == "C" {
			fmt.Printf("%8v", AllMvStratNumThisDeck)
			for i := 1; i <= moveNum-1; i++ {
				fmt.Printf("        ")
			}
			fmt.Printf("%v\n", moveShortName[doThisMove.name])
		}
	}
	// Check for repetitive move

	// Check for win

	// Check for loss

	// Find Next Moves
	aMoves := detectAvailableMoves(bNew, moveNum, singleGame)
	// if more than one move is available, sort them
	if len(aMoves) > 1 { //sort them by priority if necessary
		sort.SliceStable(aMoves, func(i, j int) bool {
			return aMoves[i].priority < aMoves[j].priority
		})
	}
	//detects Loss
	if len(aMoves) == 0 { //No available moves; game lost.
		if verbose > 1 {
			fmt.Printf("****Deck %v: XXXXGame lost after %v moves\n", deckNum, moveNum)
		}
		if verbose > 2 {
			fmt.Printf("GameLost: Frequency of each moveType:\n%v\n", moveTypes)
			fmt.Printf("GameLost: aMovesNumberOf:\n%v\n", 0 /*aMovesNumberOf*/)
		}
		//lossesAtNoMoves++
		return // something
	} else {
		// Now Try all moves

		for i, move := range aMoves {
			if i != 1 {
				AllMvStratNumThisDeck++
			}
			if strings.Contains(verboseSpecial, "A") {
				fmt.Printf("\n \nDeck:%v   moveNum:%v   AllMvStratNumThisDeck:%v   PriorMove: %v\n", deckNum, moveNum, AllMvStratNumThisDeck, printMove(doThisMove, moveNum))
				printBoard(bNew)
				fmt.Printf("All Possible Next Moves: ")
				for j, _ := range aMoves {
					if j != 0 {
						fmt.Printf("                         ")
					}
					//fmt.Printf("%v", aMoves[j])
					fmt.Printf("%v", printMove(aMoves[j]), moveNum)
					if i == j {
						fmt.Printf("                <- Next Move")
					}
					fmt.Printf("\n")
				}
				fmt.Printf("\n\n****************************************\n")
			}
			playAllMoveS(bNew, move, moveNum+1, deckNum)
		}
	}

}
