package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

func playAllMoveS(bIn board, doThisMove move, moveNum int, deckNum int) string {

	// add code for findAllSuccessfulStrategies

	if verbose > 2 {
		printBoard(bIn)
		printMove(doThisMove, moveNum)
		//fmt.Printf("moveNum: %v\n gameLengthLimit: %v NOT IMPLEMENTED YET", moveNum, gameLengthLimit)
	}

	// Make the indicated move
	var bNew board
	if moveNum == 0 {
		bNew = bIn
		if printTree == "C" {
			fmt.Printf("\n\nDeck %v\n", deckNum)
			fmt.Printf("\n\n Strat #")
			for i := 1; i <= 200; i++ {
				fmt.Printf("    %3v ", i)
			}
			fmt.Printf("\n")
		}
	} else {
		bNew = moveMaker(bIn, doThisMove)
		aMmvsTriedThisDeck++
		if printTree == "C" {
			fmt.Printf("%v  ", moveShortName[doThisMove.name])
		}
	}
	// Check for repetitive move
	bNewBcode := bNew.boardCode()      //  consider modifying the boardCode and boardDeCode methods to produce strings
	bNewBcodeS := string(bNewBcode[:]) //  consider modifying the boardCode and boardDeCode methods to produce strings
	if priorBoards[bNewBcodeS] {
		aMStratlossesAtRepMveThisDeck++
		//		if printTree == "C" {
		//			fmt.Printf("  Repetitive Move - Loop:end strategy")
		//		}
		if verbose > 1 {
			fmt.Printf("Deck %v Strategy #: %v lost after %v moves - repetitive.  Total moves tried: %v \n", deckNum, aMStratNumThisDeck, moveNum, aMmvsTriedThisDeck)
		}
		return "RM" // Repetitive Move
	} else {
		priorBoards[bNewBcodeS] = true
	}

	//Detect Early Win
	if detectWinEarly(bNew) {
		aMearlyWinCounterThisDeck++
		aMwinCounterThisDeck++

		if verbose > 0 {
			fmt.Printf("Deck %v Game won early on strategy #: %v after %v moves.  Total moves tried: %v \n", deckNum, aMStratNumThisDeck, moveNum, aMmvsTriedThisDeck)
		}
		//	if verbose > 1 {
		//		fmt.Printf("GameWon: Frequency of each moveType:\n%v\n", moveTypes)
		//	}
		return "EW" //  Early Win
	}

	//Detects Win
	if len(bNew.piles[0])+len(bNew.piles[1])+len(bNew.piles[2])+len(bNew.piles[3]) == 52 {
		aMstandardWinCounterThisDeck++
		aMwinCounterThisDeck++

		if verbose > 0 {
			fmt.Printf("Deck %v Game won on strategy #: %v after %v moves.  Total moves tried: %v \n", deckNum, aMStratNumThisDeck, moveNum, aMmvsTriedThisDeck)
		}
		//if verbose > 1 {
		//	fmt.Printf("GameWon: Frequency of each moveType:\n%v\n", moveTypes)
		//}
		return "SW" //  Standard Win
	}
	/*
		This code is now performed by Repetitive Move test above
		//Detects Loss

		if aMoves[0].name == "flipWasteToStock" {
			if moveCounter < 20 { // changed from < 20
				priorBoardNullWaste = b
			} else if reflect.DeepEqual(b, priorBoardNullWaste) {
				if verbose > 1 {
					fmt.Printf("*****Loss detected after %v moves\n", moveCounter)
				}
				regularLosses++
				continue newInitialOverrideStrategy
			} else {
				priorBoardNullWaste = b
			}
		}
	*/

	// Find Next Moves
	aMoves := detectAvailableMoves(bNew, moveNum, singleGame)

	//detects Loss
	if len(aMoves) == 0 { //No available moves; game lost.
		if verbose > 1 {
			fmt.Printf("Deck %v Strategy #: %v lost after %v moves.  Total moves tried: %v \n", deckNum, aMStratNumThisDeck, moveNum, aMmvsTriedThisDeck)
		}
		//	if verbose > 2 {
		//		fmt.Printf("Strategy Lost: Frequency of each moveType:\n%v\n", moveTypes)
		//		}
		aMStratlossesAtNoMovesThisDeck++
		return "NM" // No Moves available
	}

	// if more than one move is available, sort them
	if len(aMoves) > 1 { //sort them by priority if necessary
		sort.SliceStable(aMoves, func(i, j int) bool {
			return aMoves[i].priority < aMoves[j].priority
		})
	}

	// Now Try all moves

	for i, move := range aMoves {
		if i != 0 {
			aMStratNumThisDeck++
		}

		// Verbose Special Starts Here - No effect on operation
		if strings.Contains(verboseSpecial, "M") {
			fmt.Printf("\n \nDeck: %v   moveNum: %v   aMStratNumThisDeck: %v  aMmvsTriedThisDeck: %v   PriorMove: %v\n", deckNum, moveNum, aMStratNumThisDeck, aMmvsTriedThisDeck, printMove(doThisMove, moveNum))
			if moveNum >= 21 {
				printBoard(bNew)
			} //if test
			fmt.Printf("All Possible Moves: ")
			for j := range aMoves {
				if j != 0 {
					fmt.Printf("                         ")
				}
				//fmt.Printf("%v", aMoves[j])
				fmt.Printf("%v", printMove(aMoves[j], moveNum))
				if i == j {
					fmt.Printf("                <- Next Move")
				}
				fmt.Printf("\n")
			}
			if math.Mod(float64(aMmvsTriedThisDeck), 20) == 0 { //test
				fmt.Printf("\n") //test
			} //test
			fmt.Printf("\n\n****************************************\n")
		}
		// Verbose Special Ends Here - No effect on operation

		result := playAllMoveS(bNew, move, moveNum+1, deckNum)
		if printTree == "C" {
			if result != "Win" {
				fmt.Printf(" %v", result)
			}
			fmt.Printf("\n%8v", aMStratNumThisDeck)
			for i := 1; i <= moveNum-1; i++ {
				fmt.Printf("        ")
			}
		}
		if result == "EW" || result == "SW" || result == "Win" && findAllSuccessfulStrategies != true {
			return "Win"
		}
	}
	aMStratlossesExhaustedThisDeck++
	return "SE" //  Strategy Exhausted
}
