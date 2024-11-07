package main

import (
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
)

func playAllMoveS(bIn board, moveNum int, deckNum int) (string, string) {

	/* Return Codes: SL  = Strategy Lost	 NMA = No Moves Available
	                 						 RB  = Repetitive Board
	                                         SE  = Strategy Exhausted
	                                         GML = GameLength Limit exceeded
	                 SW  = Strategy Win      EW  = Early Win
											 SW  = Standard Win  Obsolete all wins are early
	*/
	// add code for findAllWinStrats
	// Setup pfmt to print thousands with commas
	var pfmt = message.NewPrinter(language.English)
	var aMoves []move //available Moves

	if moveNum > moveNumMax {
		moveNumMax = moveNum
	}

	if mvsTriedTD >= gameLengthLimit {
		prntMDet(bIn, aMoves, 0, deckNum, moveNum, "BB", 2, "\n  SL-RB: Game Length of: %v exceeds limit: %v\n", strconv.Itoa(mvsTriedTD), strconv.Itoa(gameLengthLimit))
		stratLossesGML_TD++
		return "SL", "GML"
	}

	// Find Next Moves
	aMoves = detectAvailableMoves(bIn, moveNum, singleGame)

	if len(aMoves) == 0 {
		m := move{name: "No Moves Available"}
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

		/* Actually before we actually try all moves let's first: print (optionally based on pMD.pType) the incoming board
		      and check the incoming board for various end-of-strategy conditions
		   Note: This was done this way, so as to ensure that when returns backed up the moveNum, the board would reprint.
		*/
		// Print the incoming board if in debugging range
		prntMDet(bIn, aMoves, i, deckNum, moveNum, "BB", 1, "", "", "") //stan do this???

		if i != 0 {
			// Started at 0 in playNew for each deck.  Increment each time a second through nth move is made
			//     That way strategy 0 is the "All Best Moves" strategy.  It is also why in playNew aMStratTriedAllDecks
			//     is incremented by stratNumTD + 1 after each deck.
			stratNumTD++
		} else {
			// Check for repetitive board
			bNewBcode := bIn.boardCode(deckNum) //  consider modifying the boardCode and boardDeCode methods to produce strings
			bNewBcodeS := string(bNewBcode[:])  //  consider modifying the boardCode and boardDeCode methods to produce strings
			// Have we seen this board before?
			if _, ok := priorBoards[bNewBcodeS]; ok {
				// OK we did see it before so return to try next available move (if any) in aMoves[] aka strategy
				stratLossesRB_TD++
				prntMDet(bIn, aMoves, i, deckNum, moveNum, "BB", 2, "\n  SF-RB: Repetitive Board - \"Next Move\" yielded a repeat of the board at at MvsTriedTD: %v which was at move: %v\n", strconv.Itoa(priorBoards[bNewBcodeS].mvsTriedTD), strconv.Itoa(priorBoards[bNewBcodeS].moveNum))
				prntMDet(bIn, aMoves, i, deckNum, moveNum, "BBSS", 2, "\n  SF-RB: Repetitive Board - \"Next Move\" yielded a repeat of the board at at MvsTriedTD: %v which was at move: %v\n", strconv.Itoa(priorBoards[bNewBcodeS].mvsTriedTD), strconv.Itoa(priorBoards[bNewBcodeS].moveNum))
				return "SL", "RB" // Repetitive Move
			} else {
				// Remember the board state by putting it into the map "priorBoards"
				bInf := boardInfo{
					moveNum:    moveNum,
					mvsTriedTD: mvsTriedTD,
				}
				priorBoards[bNewBcodeS] = bInf
			}
		}

		// Check if No Moves Available
		if i == 0 && aMoves[0].name == "No Moves Available" {
			stratLossesNMA_TD++
			prntMDet(bIn, aMoves, i, deckNum, moveNum, "BB", 2, "  SL-NMA: No Moves Available: Strategy Lost %v%v\n", "", "")
			return "SL", "NMA"
		}

		//Detect Early Win
		if detectWinEarly(bIn) {
			stratWinsTD++
			cmt := "  SW-EW: Strategy Win: Early Win%v%v"
			if findAllWinStrats {
				cmt += "  Will Continue to look for additional winning strategies for this deck"
			} else {
				cmt += "  Go to Next Deck (if any)"
			}
			prntMDet(bIn, aMoves, i, deckNum, moveNum, "NOTX", 1, cmt, "", "")

			// Verbose Special "WL" Starts Here - No effect on operation
			if strings.Contains(verboseSpecial, ";WL;") { // Deck Win Loss Summary Statistics
				/*if len(deckWinLossDetail)-1 < deckNum {
					dWLDStats.winLoss = "W"
					dWLDStats.moveNumAt1stWinOrAtLoss = moveNum
					dWLDStats.moveNumMinWinIfFindAll = moveNum
					dWLDStats.moveNumMaxWinIfFindAll = moveNum
					dWLDStats.stratNumAt1stWinOrAtLoss = stratNumTD
					dWLDStats.mvsTriedAt1stWinOrAtLoss = mvsTriedTD
					deckWinLossDetail = append(deckWinLossDetail, dWLDStats)
				} else {
					if deckWinLossDetail[deckNum].moveNumMinWinIfFindAll > moveNum {
						deckWinLossDetail[deckNum].moveNumMinWinIfFindAll = moveNum
					}
					if deckWinLossDetail[deckNum].moveNumMaxWinIfFindAll < moveNum {
						deckWinLossDetail[deckNum].moveNumMaxWinIfFindAll = moveNum
					}
				}*/

			}
			// Verbose Special "WL" Ends Here - No effect on operation
			return "SW", "EW" //  Strategy Early Win
		}

		// OK, done with the various end-of-strategy conditions
		// let's print out the list of available moves and make the next available move
		// The board state was already printed above
		if pMD.pType == "BB" && prntMDetTestRange(deckNum) {
			prntMDet(bIn, aMoves, i, deckNum, moveNum, "BB", 3, "", "", "")
		}

		bNew := bIn.copyBoard() // Critical Must use copyBoard

		// ********** 1st of the 2 MOST IMPORTANT statements in this function:  ******************************
		bNew = moveMaker(bNew, aMoves[i])

		prntMDet(bIn, aMoves, i, deckNum, moveNum, "BBS", 1, "\n\nBefore Call at Deck: %v   Move: %v   Strategy #: %v  Moves Tried: %v   Unique Boards: %v   Elapsed TD: %v   Elapsed ADs: %v\n", "", "")
		prntMDet(bIn, aMoves, i, deckNum, moveNum, "BBS", 2, "      bIn: %v\n", "", "")
		prntMDet(bNew, aMoves, i, deckNum, moveNum, "BBS", 2, "     bNew: %v\n", "", "")
		mvsTriedTD++

		// verboseSpecial option PROGRESSdddd starts here      Consider Moving clause to be a clause in prntMDet
		if verboseSpecialProgressCounter > 0 && math.Mod(float64(mvsTriedTD), float64(verboseSpecialProgressCounter)) <= 0.1 {
			avgRepTime := time.Since(startTimeTD) / time.Duration(mvsTriedTD/verboseSpecialProgressCounter)
			estMaxRemTimeTD := time.Since(startTimeTD) * time.Duration((gameLengthLimit-mvsTriedTD)/mvsTriedTD)
			_, err = pfmt.Printf("\rDk: %5d   ____   MvsTried: %13v   MoveNum: %3v   Max MoveNum: %3v   StratsTried: %12v   UnqBoards: %11v   El Time: %7s   Since Last Rep: %6s   Avg Btwn Reps: %6s   estMaxRemTD: %7s\r", deckNum, mvsTriedTD, moveNum, moveNumMax, stratNumTD, len(priorBoards), time.Since(startTimeTD).Truncate(100*time.Millisecond).String(), time.Since(verboseSpecialProgressCounterLastPrintTime).Truncate(100*time.Millisecond).String(), avgRepTime.Truncate(100*time.Millisecond).String(), estMaxRemTimeTD.Truncate(1*time.Minute))
			verboseSpecialProgressCounterLastPrintTime = time.Now()
		}
		// verboseSpecial option PROGRESSdddd ends here

		// ********** 2nd of the 2 MOST IMPORTANT statements in this function:  ******************************
		recurReturnV1, recurReturnV2 := playAllMoveS(bNew, moveNum+1, deckNum)

		prntMDet(bIn, aMoves, i, deckNum, moveNum, "NOTX", 1, "  Returned: %v - %v After Call at deckNum: %v  moveNum: %v   StratNumTD: %v   MvsTriedTD: %v   UnqBds: %v   ElTimTD: %v   ElTimADs: %v\n", recurReturnV1, recurReturnV2)

		if findAllWinStrats != true && recurReturnV1 == "SW" {
			// save winning moves into a slice in reverse
			return recurReturnV1, recurReturnV2 // return up the call stack to end strategies search  if findAllWinStrats false, and we had a win
		}
		if recurReturnV1 == "SL" && recurReturnV2 == "GML" {
			return recurReturnV1, recurReturnV2 // return up the call stack to end strategies search  if findAllWinStrats false, and we had a win
		}
	}

	stratLossesSE_TD++
	return "SL", "SE" //  Strategy Exhausted
}

func prntMDetTestRange(deckNum int) bool {
	deckRangeOK := false
	if pMD.deckStartVal == 0 && pMD.deckContinueFor == 0 {
		deckRangeOK = true
	} else {
		if deckNum >= pMD.deckStartVal && (pMD.deckContinueFor == 0 || deckNum < pMD.deckStartVal+pMD.deckContinueFor) {
			deckRangeOK = true
		}
	}
	aMvsThisDkRangeOK := false
	if pMD.movesTriedTDStartVal == 0 && pMD.movesTriedTDContinueFor == 0 {
		aMvsThisDkRangeOK = true
	} else {
		if mvsTriedTD >= pMD.movesTriedTDStartVal && (pMD.movesTriedTDContinueFor == 0 || mvsTriedTD < pMD.movesTriedTDStartVal+pMD.movesTriedTDContinueFor) {
			aMvsThisDkRangeOK = true
		}
	}
	if deckRangeOK && aMvsThisDkRangeOK {
		return true
	} else {
		return false
	}
}

func prntMDet(b board, aMoves []move, nextMove int, dN int, mN int, pTypeIn string, variant int, comment string, s1 string, s2 string) {
	// Done here just to clean up mainline logic of playAllMoves
	// Do some repetitive printing to track progress
	// This function will use the struct pMD
	//      variant will be used for different outputs under the same pType

	// Setup pfmt to print thousands with commas
	var pfmt = message.NewPrinter(language.English)

	if pMD.pType != "X" && prntMDetTestRange(dN) {
		switch {
		case pTypeIn == "BB" && pMD.pType == pTypeIn && variant == 1: // for BB
			fmt.Printf("\n****************************************\n")
			_, err = pfmt.Printf("\nDeck: %v   Move: %v   Strategy #: %v  Moves Tried: %v   Unique Boards: %v   Elapsed TD: %v   Elapsed ADs: %v\n", dN, mN, stratNumTD, mvsTriedTD, len(priorBoards), time.Now().Sub(startTimeAD), time.Now().Sub(startTimeTD))
			printBoard(b)
		case pTypeIn == "BB" && pMD.pType == pTypeIn && variant == 2: // for BB
			// comment must have 2 %v in it
			_, err = pfmt.Printf(comment, s1, s2)
		case pMD.pType == "BB" && variant == 3:
			fmt.Printf("\n     All Possible Moves: ")
			for j := range aMoves {
				if j != 0 {
					fmt.Printf("                         ")
				}
				fmt.Printf("%v", printMove(aMoves[j]))
				if nextMove == j {
					fmt.Printf("                <- Next Move")
				}
				fmt.Printf("\n")
			}
		case strings.HasPrefix(pTypeIn, "BBS") && strings.HasPrefix(pMD.pType, "BBS") && variant == 1: // for BBS or BBSS
			_, err = pfmt.Printf(comment, dN, mN, stratNumTD, mvsTriedTD, len(priorBoards), time.Now().Sub(startTimeAD), time.Now().Sub(startTimeTD))
		case pTypeIn == "BBS" && pMD.pType == pTypeIn && variant == 2:
			_, err = pfmt.Printf(comment, b)
		case pTypeIn == "NOTX" && pMD.pType != "X" && variant == 1:
			_, err = pfmt.Printf(comment, s1, s2, dN, mN, stratNumTD, mvsTriedTD, len(priorBoards), time.Now().Sub(startTimeAD), time.Now().Sub(startTimeTD))
		case pTypeIn == "NOTX" && pMD.pType != "X" && variant == 2:
			_, err = pfmt.Printf(comment, s1, s2)
		case (pTypeIn == "TW" || pTypeIn == "TS" || pTypeIn == "TSS") && variant == 1:
		case (pTypeIn == "TW" || pTypeIn == "TS" || pTypeIn == "TSS") && variant == 2:
		}
	}
}
