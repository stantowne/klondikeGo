package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
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
		if i != 0 {
			// Started at 0 in playNew for each deck.  Increment each time a second through nth move is made
			//     That way strategy 0 is the "All Best Moves" strategy.  It is also why in playNew aMStratTriedAllDecks
			//     is incremented by aMStratNumThisDeck + 1 after each deck.
			aMStratNumThisDeck++
		}

		/* Actually before we actually try all moves let's first: print (optionally based on printMoveDetail.pType) the incoming board
		      and check the incoming board for various end-of-strategy conditions
		   Note: This was done this way, so as to ensure that when returns backed up the moveNum, the board would reprint.
		*/
		// Print the incoming board if in debugging range
		pMd(bIn, deckNum, moveNum, "BB", 1, "", "", "")

		// Check if No possible Moves
		if i == 0 && aMoves[0].name == "No Possible Moves" {
			pMd(bIn, deckNum, moveNum, "BB", 2, "  SF-NPM: No Possible Moves: Strategy Failed %v%v\n", "", "")
			return "SF", "NPM"
		}

		// Check for repetitive board
		bNewBcode := bIn.boardCode()       //  consider modifying the boardCode and boardDeCode methods to produce strings
		bNewBcodeS := string(bNewBcode[:]) //  consider modifying the boardCode and boardDeCode methods to produce strings
		// Have we seen this board before?
		if _, ok := priorBoards[bNewBcodeS]; ok {
			// OK we did see it before but lets check if we have just returned to try the next available move (if any)
			if priorBoards[bNewBcodeS].mN == moveNum {
				// Do Nothing We are just back up at this board checking for the next available move
			} else {
				aMStratlossesAtRepMveThisDeck++
				pMd(bIn, deckNum, moveNum, "BB", 2, "  SF-RB: Repetitive Board - Loop:end strategy - see board at aMmvsTriedThisDeck: %v%v\n", strconv.Itoa(priorBoards[bNewBcodeS].aMmvsTriedTD), "")
				return "SF", "RB" // Repetitive Move
			}
		} else {
			bI := boardInfo{
				mN:           moveNum,
				aMmvsTriedTD: aMmvsTriedThisDeck,
			}
			priorBoards[bNewBcodeS] = bI
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
			pMd(bIn, deckNum, moveNum, "BB", 2, c, "", "")
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
			pMd(bIn, deckNum, moveNum, "BB", 2, c, "", "")
			return "SW", "SW" //  Standard Win
		}

		// OK, done with the various end-of-strategy conditions
		// let's print out the list of available moves and make the next available move
		if printMoveDetail.pType == "BB" && pMdTestRange(deckNum) {
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

		bNew := bIn.copyBoard() // Critical Must use copyBoard
		bNew = moveMaker(bNew, aMoves[i])
		pMd(bIn, deckNum, moveNum, "BBS", 1, "\n\nBefore Call at deckNum: %v  moveNum: %v   aMStratNumThisDeck: %v   aMmvsTriedThisDeck: %v\n", "", "")
		pMd(bIn, deckNum, moveNum, "BBS", 2, "      bIn: %v\n", "", "")
		pMd(bNew, deckNum, moveNum, "BBS", 2, "     bNew: %v\n", "", "")
		aMmvsTriedThisDeck++

		r1, r2 := playAllMoveS(bNew, moveNum+1, deckNum)

		//pMd(bIn, deckNum, moveNum, "BBS", 1, " After Call at deckNum: %v  moveNum: %v   aMStratNumThisDeck: %v   aMmvsTriedThisDeck: %v\n", "", "")
		//pMd(bIn, deckNum, moveNum, "BBS", 2, "      bIn: %v\n", "", "")
		pMd(bIn, deckNum, moveNum, "NOTX", 1, "  Returned r1: %v, r2: %v After Call at deckNum: %v  moveNum: %v   aMStratNumThisDeck: %v   aMmvsTriedThisDeck: %v\n", r1, r2)
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
	}*/

	aMStratlossesExhaustedThisDeck++
	return "SF", "SE" //  Strategy Exhausted
}

func pMdTestRange(deckNum int) bool {
	deckRangeOK := false
	if printMoveDetail.deckStartVal == 0 && printMoveDetail.deckContinueFor == 0 {
		deckRangeOK = true
	} else {
		if deckNum >= printMoveDetail.deckStartVal && (printMoveDetail.deckContinueFor == 0 || deckNum < printMoveDetail.deckStartVal+printMoveDetail.deckContinueFor) {
			deckRangeOK = true
		}
	}
	aMvsThisDkRangeOK := false
	if printMoveDetail.aMvsThisDkStartVal == 0 && printMoveDetail.aMvsThisDkContinueFor == 0 {
		aMvsThisDkRangeOK = true
	} else {
		if aMmvsTriedThisDeck >= printMoveDetail.aMvsThisDkStartVal && (printMoveDetail.aMvsThisDkContinueFor == 0 || aMmvsTriedThisDeck < printMoveDetail.aMvsThisDkStartVal+printMoveDetail.aMvsThisDkContinueFor) {
			aMvsThisDkRangeOK = true
		}
	}
	if deckRangeOK && aMvsThisDkRangeOK {
		return true
	} else {
		return false
	}
}

func pMd(b board, dN int, mN int, pTypeIn string, variant int, comment string, s1 string, s2 string) {
	// Done here just to clean up mainline logic of playAllMoves
	// Do some repetitive printing to track progress
	// This function will use the struct printMoveDetail
	//      variant will be used for different outputs under the same pType
	if printMoveDetail.pType != "X" && pMdTestRange(dN) {
		switch {
		case pTypeIn == "BB" && printMoveDetail.pType == pTypeIn && variant == 1: // for BB
			if mN != 0 {
				fmt.Printf("\n\n****************************************\n")
			}
			fmt.Printf("\n \nDeck: %v   mN: %v   aMStratNumThisDeck: %v  aMmvsTriedThisDeck: %v \n", dN, mN, aMStratNumThisDeck, aMmvsTriedThisDeck)
			printBoard(b)
		case pTypeIn == "BB" && printMoveDetail.pType == pTypeIn && variant == 2: // for BB
			// comment must have 2 %v in it
			fmt.Printf(comment, s1, s2)
			/*case printMoveDetail.pType == "BB" && variant == 3:     // for BB FIX THIS AT SOME POINT
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
		case strings.HasPrefix(pTypeIn, "BBS") && strings.HasPrefix(printMoveDetail.pType, "BBS") && variant == 1: // for BBS or BBSS
			fmt.Printf(comment, dN, mN, aMStratNumThisDeck, aMmvsTriedThisDeck)
		case pTypeIn == "BBS" && printMoveDetail.pType == pTypeIn && variant == 2: // for BBS or BBSS
			fmt.Printf(comment, b)
		case pTypeIn == "NOTX" && printMoveDetail.pType != "X" && variant == 1:
			fmt.Printf(comment, s1, s2, dN, mN, aMStratNumThisDeck, aMmvsTriedThisDeck)
		}
	}
}
