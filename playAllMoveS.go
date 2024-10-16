package main

import (
	"fmt"
	"sort"
)

func playAllMoveS(bIn board, doThisMove move, moveNum int, deckNum int, verbose int, findAllSuccessfulStrategies bool, printTree string, AllMvStratNum int) {
	if verbose > 2 {
		printBoard(bIn)
		// printMove(doThisMove)  NOT written yet
		fmt.Printf("moveNum: %v, deckNum: %v, verbose: %v, findAllSuccessfulStrategies: %v, printTree: %v, AllMvStratNum: %v\n", moveNum, deckNum, verbose, findAllSuccessfulStrategies, printTree, AllMvStratNum)
		fmt.Printf("gameLengthLimitNew: %v", gameLengthLimitNew)
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
				fmt.Printf("    %3i ", i)
			}
			fmt.Printf("\n")
		} else {
			bNew = moveMaker(bIn, doThisMove)
			moveNum++
			AllMvStratNum++
			fmt.Printf("%8i", AllMvStratNum)
			for i := 1; i <= moveNum-1; i++ {
				fmt.Printf("        ")
			}
			fmt.Printf(" %v", "x" /*moveShortName[doThisMove.name]*/)
		}
	}

	// Check for repetitive move

	// Check for win

	// Check for loss

	// Find Next Moves

	// Now Try all moves
	aMoves := detectAvailableMoves(bNew, moveNum, false /*singleGame*/)
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
		/*	for i, move := range aMoves {
				playAllMoveS(bNew, aMoves[i], moveNum+1, deckNum, verbose, findAllSuccessfulStrategies, printTree, AllMvStratNum)
			}
		*/
	}

}
