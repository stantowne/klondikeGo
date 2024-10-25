package main

import (
	"fmt"
	"sort"
	"strconv"
)

func playAllMoveS(bIn board, moveNum int, deckNum int) (string, string) {

	/* Return Codes: SF  = Strategy Failed	NPM	= No Possible Moves
	                 						 RB = Repetitive Board
	                                         SE = Strategy Exhausted
	                 SW  = Strategy Win      EW = Early Win
											 SW = Standard Win
	*/

	// add code for findAllSuccessfulStrategies

	// Closures START:

	pMdTestRange := func( /*deckNum int*/ ) bool {
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

	pMd := func(pTypeIn string, variant int, comment string, s1 string, s2 string) {
		// Do some repetitive printing to track
		// This function will use the struct printMoveDetail
		//      variant will be used for different outputs under the same pType
		if printMoveDetail.pType != "X" && pMdTestRange() {
			switch {
			case printMoveDetail.pType == "BB" && variant == 1:
				if moveNum != 0 {
					fmt.Printf("\n\n****************************************\n")
				}
				fmt.Printf("\n \nDeck: %v   moveNum: %v   aMStratNumThisDeck: %v  aMmvsTriedThisDeck: %v \n", deckNum, moveNum, aMStratNumThisDeck, aMmvsTriedThisDeck)
				printBoard(bIn)
			case printMoveDetail.pType == "BB" && variant == 2:
				// comment must have 2 %v in it
				fmt.Printf(comment, s1, s2)
				/*case printMoveDetail.pType == "BB" && variant == 3:
					fmt.Printf("\n     All Possible Moves: ")
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
				}*/
			}

		}
	}
	// Closures END:

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
		pMd("BB", 1, "", "", "")

		// We only have to do these checks if it is the first time we have been at this board
		// We know it is not the first time at this board if we are executing the move aMoves[i] where i > 0
		if i == 0 {
			// Check if No possible Moves
			if aMoves[0].name == "No Possible Moves" {
				pMd("BB", 2, "  SF-NPM: No Possible Moves: Strategy Failed%v%v\n", "", "")
				return "SF", "NPM"
			}
			// Check for repetitive move
			bNewBcode := bIn.boardCode() //  consider modifying the boardCode and boardDeCode methods to produce strings

			bNewBcodeS := string(bNewBcode[:]) //  consider modifying the boardCode and boardDeCode methods to produce strings
			if priorBoards[bNewBcodeS] != 0 {
				aMStratlossesAtRepMveThisDeck++
				pMd("BB", 2, "  SF-RM: Repetitive Board - Loop:end strategy - see board at aMmvsTriedThisDeck: %v%v", strconv.Itoa(priorBoards[bNewBcodeS]), "")
				return "SF", "RM" // Repetitive Move
			} else {
				priorBoards[bNewBcodeS] = aMmvsTriedThisDeck
			}

			//Detect Early Win
			if detectWinEarly(bIn) {
				aMearlyWinCounterThisDeck++
				aMwinCounterThisDeck++
				c := "  SW-EW: Strategy Win: Early Win%v%v"
				if findAllSuccessfulStrategies {
					c += "  Will Continue to look for additional winning strategies for this deck"
				} else {
					c += "  Go to Next Deck (if any)"
				}
				pMd("BB", 2, c, "", "")
				return "SW", "EW" //  Early Win
			}

			//Detects Standard Win
			if len(bIn.piles[0])+len(bIn.piles[1])+len(bIn.piles[2])+len(bIn.piles[3]) == 52 {
				aMstandardWinCounterThisDeck++
				aMwinCounterThisDeck++
				c := "  SW-SW: Strategy Win: Standard Win%v%v"
				if findAllSuccessfulStrategies {
					c = c + "  Will Continue to look for additional winning strategies for this deck"
				} else {
					c = c + "  Go to Next Deck (if any)"
				}
				pMd("BB", 2, c, "", "")
				return "SW", "SW" //  Standard Win
			}
		}
		// OK, done with the various end-of-strategy conditions
		// let's print out the list of available moves and make the next available move
		if printMoveDetail.pType == "BB" && pMdTestRange() {
			fmt.Printf("\n     All Possible Moves: ")
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
		//pMd("BB",3,"","","")
		bNew := bIn.copyBoard()
		bNew = moveMaker(bNew, aMoves[i])
		aMmvsTriedThisDeck++
		/*if pMdTestRange() {
			fmt.Printf("\n\nBefore Call at deckNum: %v  moveNum: %v   aMmvsTriedThisDeck: %v\n      bIn: %v\n  bInOrig: %v\n     bNew: %v\n", deckNum, moveNum, aMmvsTriedThisDeck, bIn, bInOrig, bNew)
		}*/
		r1, r2 := playAllMoveS(bNew, moveNum+1, deckNum)

		//fmt.Printf("\n@@@@@@@@@@@After Call at deckNum: %v  moveNum: %v   aMmvsTriedThisDeck: %v\n      r1: %v\n  r2: %v\n ", deckNum, moveNum, aMmvsTriedThisDeck, r1, r2)
		if deckNum == 0 {
			x := 2
			x += 3
		}
		if pMdTestRange() {
			fmt.Printf("\nAfter Call at deckNum: %v  moveNum: %v   aMmvsTriedThisDeck: %v\n      bIn: %v\n    bNew: %v\n", deckNum, moveNum, aMmvsTriedThisDeck, bIn, bNew)
		}
		x := 1
		x++
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
