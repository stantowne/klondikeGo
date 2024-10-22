package main

import (
	"fmt"
	"sort"
)

func playAllMoveS(bIn board, moveNum int, deckNum int) (string, string) {

	/* Return Codes: SF  = Strategy Failed	NPM	= No Possible Moves
	                 						 RB = Repetitive Board
	                                         SE = Strategy Exhausted
	                 SW  = Strategy Win      EW = Early Win
											 SW = Standard Win
	*/

	// add code for findAllSuccessfulStrategies

	// Find Next Moves
	aMoves := detectAvailableMoves(bIn, moveNum, singleGame)

	if len(aMoves) == 0 {
		m := move{name: "No Possible Moves"}
		aMoves = append(aMoves, m)

	} else {
		// if more than one move is available, sort them
		if len(aMoves) > 1 { //sort them by priority if necessary
			sort.SliceStable(aMoves, func(i, j int) bool {
				return aMoves[i].priority < aMoves[j].priority
			})
		}
	}

	// Try all moves
	for i := range aMoves {

		/* Actually before we actually try all moves let's first: print (optionally based on printMoveDetail.pType) the incoming board
		      and check the incoming board for various end-of-strategy conditions
		   Note: This was done this way, so as to ensure that when returns backed up the moveNum, the board would reprint.
		*/

		// Print the incoming board
		if printMoveDetail.pType == "BB" && pMdTestRange(deckNum) {
			if moveNum != 0 {
				fmt.Printf("\n\n****************************************\n")
			}
			fmt.Printf("\n \nDeck: %v   moveNum: %v   aMStratNumThisDeck: %v  aMmvsTriedThisDeck: %v \n", deckNum, moveNum, aMStratNumThisDeck, aMmvsTriedThisDeck)
			printBoard(bIn)
		}

		// Check if No possible Moves
		if aMoves[0].name == "No Possible Moves" {
			if printMoveDetail.pType == "BB" && pMdTestRange(deckNum) {
				fmt.Printf("  SF-NPM: No Possible Moves: Strategy Failed\n")
			}
			return "SF", "NPM"
		}
		// Check for repetitive move
		bNewBcode := bIn.boardCode() //  consider modifying the boardCode and boardDeCode methods to produce strings

		bNewBcodeS := string(bNewBcode[:]) //  consider modifying the boardCode and boardDeCode methods to produce strings
		if priorBoards[bNewBcodeS] != 0 {
			aMStratlossesAtRepMveThisDeck++
			if printMoveDetail.pType == "BB" && pMdTestRange(deckNum) {
				fmt.Printf("  SF-RM: Repetitive Board - Loop:end strategy - see board at aMmvsTriedThisDeck: %v", priorBoards[bNewBcodeS])
			}
			return "SF", "RM" // Repetitive Move
		} else {
			priorBoards[bNewBcodeS] = aMmvsTriedThisDeck
		}

		//Detect Early Win
		if detectWinEarly(bIn) {
			aMearlyWinCounterThisDeck++
			aMwinCounterThisDeck++
			if printMoveDetail.pType == "BB" && pMdTestRange(deckNum) {
				fmt.Printf("  SW-EW: Strategy Win: Early Win")
				if findAllSuccessfulStrategies {
					fmt.Printf("  Will Continue to look for additional winning strategies for this deck")
				} else {
					fmt.Printf("  Go to Next Deck (if any)")
				}
			}
			return "SW", "EW" //  Early Win
		}

		//Detects Standard Win
		if len(bIn.piles[0])+len(bIn.piles[1])+len(bIn.piles[2])+len(bIn.piles[3]) == 52 {
			aMstandardWinCounterThisDeck++
			aMwinCounterThisDeck++
			if printMoveDetail.pType == "BB" && pMdTestRange(deckNum) {
				fmt.Printf("  SW-SW: Strategy Win: Standard Win")
				if findAllSuccessfulStrategies {
					fmt.Printf("  Will Continue to look for additional winning strategies for this deck")
				} else {
					fmt.Printf("  Go to Next Deck (if any)")
				}
			}
			return "SW", "SW" //  Standard Win
		}

		// OK, done with the various end-of-strategy conditions
		// let's print out the list of available moves and make the next available move
		if printMoveDetail.pType == "BB" && pMdTestRange(deckNum) {
			fmt.Printf("     All Possible Moves: ")
			for j := range aMoves {
				if j != 0 {
					fmt.Printf("                         ")
				}
				fmt.Printf("%v", printMove(aMoves[j]))
				if i == j {
					fmt.Printf("                <- Next Move")
				}
				fmt.Printf("\n")
			}
		}
		bNew := moveMaker(bIn, aMoves[i])
		aMmvsTriedThisDeck++
		playAllMoveS(bNew, moveNum+1, deckNum)
	}

	/*
				if printTree == "C" {
				if result != "Win" {
					fmt.Printf(" %v", result)
				}
				fmt.Printf("\n%8v", aMStratNumThisDeck)
				for i := 1; i <= moveNum-1; i++ {
					fmt.Printf("        ")
				}

			if result == "EW" || result == "SW" || result == "Win" && findAllSuccessfulStrategies != true {
				return "Win"
			}
		}
	*/
	aMStratlossesExhaustedThisDeck++
	return "SF", "SE" //  Strategy Exhausted
}

func pMdTestRange(deckNum int) bool {
	if printMoveDetail.startVal == 0 && printMoveDetail.continueFor == 0 {
		return true
	}
	if printMoveDetail.startType == "DECK" && deckNum >= printMoveDetail.startVal {
		if printMoveDetail.continueFor == 0 || deckNum < printMoveDetail.startVal+printMoveDetail.continueFor {
			return true
		}
	}
	if printMoveDetail.startType == "MvsT" && aMmvsTriedThisDeck >= printMoveDetail.startVal {
		if printMoveDetail.continueFor == 0 || aMmvsTriedThisDeck < printMoveDetail.startVal+printMoveDetail.continueFor {
			return true
		}
	}
	return false
}
