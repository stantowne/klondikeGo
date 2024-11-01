package main

import (
	"fmt"
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
											 SW  = Standard Win  Obsolete all wins are aearly
	*/
	// add code for findAllSuccessfulStrategies

	if mvsTriedTD > gameLengthLimit {
		return "SL", "GML"
	}

	// Find Next Moves
	aMoves := detectAvailableMoves(bIn, moveNum, singleGame)

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

		/* Actually before we actually try all moves let's first: print (optionally based on printMoveDetail.pType) the incoming board
		      and check the incoming board for various end-of-strategy conditions
		   Note: This was done this way, so as to ensure that when returns backed up the moveNum, the board would reprint.
		*/
		// Print the incoming board if in debugging range
		pMd(bIn, deckNum, moveNum, "BB", 1, "", "", "")

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
				pMd(bIn, deckNum, moveNum, "BB", 2, "\n  SF-RB: Repetitive Board - \"Next Move\" yielded a repeat of the board at at MvsTriedTD: %v which was at move: %v\n", strconv.Itoa(priorBoards[bNewBcodeS].aMmvsTriedTD), strconv.Itoa(priorBoards[bNewBcodeS].mN))
				pMd(bIn, deckNum, moveNum, "BBSS", 2, "\n  SF-RB: Repetitive Board - \"Next Move\" yielded a repeat of the board at at MvsTriedTD: %v which was at move: %v\n", strconv.Itoa(priorBoards[bNewBcodeS].aMmvsTriedTD), strconv.Itoa(priorBoards[bNewBcodeS].mN))
				return "SL", "RB" // Repetitive Move
			} else {
				bInf := boardInfo{
					mN:           moveNum,
					aMmvsTriedTD: mvsTriedTD,
				}
				priorBoards[bNewBcodeS] = bInf
			}
		}

		// Check if No Moves Available
		if i == 0 && aMoves[0].name == "No Moves Available" {
			pMd(bIn, deckNum, moveNum, "BB", 2, "  SF-NMA: No Moves Available: Strategy Failed %v%v\n", "", "")
			return "SL", "NMA"
		}

		//Detect Early Win
		if detectWinEarly(bIn) {
			stratWinsTD++
			cmt := "  SW-EW: Strategy Win: Early Win%v%v"
			if findAllSuccessfulStrategies {
				cmt += "  Will Continue to look for additional winning strategies for this deck"
			} else {
				cmt += "  Go to Next Deck (if any)"
			}
			pMd(bIn, deckNum, moveNum, "NOTX", 1, cmt, "", "")

			// Verbose Special "WL" Starts Here - No effect on operation
			if strings.Contains(verboseSpecial, "/WL/") { // Deck Win Loss Summary Statistics
				if len(deckWinLossDetail)-1 < deckNum {
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
				}

			}
			// Verbose Special "WL" Ends Here - No effect on operation
			return "SW", "EW" //  Strategy Early Win
		}

		/*		//Detects Standard Win
				if len(bIn.piles[0])+len(bIn.piles[1])+len(bIn.piles[2])+len(bIn.piles[3]) == 52 {
					aMstandardWinCounterThisDeck++
					stratWinsTD++
					c := "  SW-SW: Strategy Win: Standard Win%v%v"
					if findAllSuccessfulStrategies {
						c = c + "  Will Continue to look for additional winning strategies for this deck"
					} else {
						c = c + "  Go to Next Deck (if any)"
					}
					pMd(bIn, deckNum, moveNum, "BB", 2, c, "", "")
					return "SW", "SW" //  Strategy Standard Win
				}
		*/

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

		/*//test
		if mvsTriedTD == 1460 || mvsTriedTD == 280 || mvsTriedTD == 1461 || mvsTriedTD == 281 {
			fmt.Printf("\n\n#########################\nBefore mvsTriedTD++ \nB4 COPY, B4 MOVEMAKER bIn BELOW mvsTriedTD == %v  move = %v", mvsTriedTD, aMoves[i])
			printBoard(bIn)
			fmt.Printf("\nB4 COPY, B4 MOVEMAKER bIn ABOVE mvsTriedTD == %v  move = %v\nAllMoves: %v\n$$$$$$$$$$$$$$$$$\n", mvsTriedTD, aMoves[i], aMoves)
		}
		//end test*/

		bNew := bIn.copyBoard() // Critical Must use copyBoard

		/*//test
		if mvsTriedTD == 1460 || mvsTriedTD == 280 || mvsTriedTD == 1461 || mvsTriedTD == 281 {
			fmt.Printf("\n\n#########################\nBefore mvsTriedTD++ \nAFTER COPY, B4 MOVEMAKER bIn BELOW mvsTriedTD == %v  move = %v", mvsTriedTD, aMoves[i])
			printBoard(bIn)
			fmt.Printf("\nAFTER COPY, B4 MOVEMAKER bIn ABOVE mvsTriedTD == %v  move = %v", mvsTriedTD, aMoves[i])
			fmt.Printf("\nAFTER COPY, B4 MOVEMAKER bNew BELOW mvsTriedTD == %v  move = %v", mvsTriedTD, aMoves[i])
			printBoard(bNew)
			fmt.Printf("\nAFTER COPY, B4 MOVEMAKER bNew ABOVE mvsTriedTD == %v  move = %v\nAllMoves: %v\n$$$$$$$$$$$$$$$$$\n", mvsTriedTD, aMoves[i], aMoves)
		}
		//end test*/

		bNew = moveMaker(bNew, aMoves[i])

		/*//test
		if mvsTriedTD == 122574 {
			fmt.Printf("\n\n#########################\nBefore mvsTriedTD++ \nAFTER COPY, AFTER MOVEMAKER bIn BELOW mvsTriedTD == %v  move = %v", mvsTriedTD, aMoves[i])
			printBoard(bIn)
			fmt.Printf("\nAFTER COPY, AFTER MOVEMAKER bIn ABOVE mvsTriedTD == %v  move = %v", mvsTriedTD, aMoves[i])
			fmt.Printf("\nAFTER COPY, AFTER MOVEMAKER bNew BELOW mvsTriedTD == %v  move = %v", mvsTriedTD, aMoves[i])
			printBoard(bNew)
			fmt.Printf("\nAFTER COPY, AFTER MOVEMAKER bNew ABOVE mvsTriedTD == %v  move = %v\nAllMoves: %v\n$$$$$$$$$$$$$$$$$\n", mvsTriedTD, aMoves[i], aMoves)
		}
		//end test*/
		pMd(bIn, deckNum, moveNum, "BBS", 1, "\n\nBefore Call at Deck: %v   Move: %v   Strategy #: %v  Moves Tried: %v   Unique Boards: %v   Elapsed TD: %v   Elapsed ADs: %v\n", "", "")
		pMd(bIn, deckNum, moveNum, "BBS", 2, "      bIn: %v\n", "", "")
		pMd(bNew, deckNum, moveNum, "BBS", 2, "     bNew: %v\n", "", "")
		mvsTriedTD++

		//	fmt.Printf("moveNum: %v   aMStratNumTD: %v   MvsTriedTD: %v   UnqBds: %v   Move: %v\n", moveNum, stratNumTD, mvsTriedTD, len(priorBoards), printMove(aMoves[i]))

		recurReturnV1, recurReturnV2 := playAllMoveS(bNew, moveNum+1, deckNum)
		pMd(bIn, deckNum, moveNum, "NOTX", 1, "  Returned: %v - %v After Call at deckNum: %v  moveNum: %v   aMStratNumTD: %v   MvsTriedTD: %v   UnqBds: %v   ElTimTD: %v   ElTimADs: %v\n", recurReturnV1, recurReturnV2)
		if findAllSuccessfulStrategies != true && recurReturnV1 == "SW" {
			return recurReturnV1, recurReturnV2 // return up the call stack to end strategies search  if findAllSuccessfulStrategies false and we had a win
		}
		if recurReturnV1 == "SF" && recurReturnV2 == "GML" {
			return recurReturnV1, recurReturnV2 // return up the call stack to end strategies search  if findAllSuccessfulStrategies false and we had a win
		}
	}

	stratLossesSE_TD++
	return "SL", "SE" //  Strategy Exhausted
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
		if mvsTriedTD >= printMoveDetail.aMvsThisDkStartVal && (printMoveDetail.aMvsThisDkContinueFor == 0 || mvsTriedTD < printMoveDetail.aMvsThisDkStartVal+printMoveDetail.aMvsThisDkContinueFor) {
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

	/*if mvsTriedTD < 300 || math.Mod(float64(mvsTriedTD), 5000) == 0 {
	 */
	if printMoveDetail.pType != "X" && pMdTestRange(dN) {
		switch {
		case pTypeIn == "BB" && printMoveDetail.pType == pTypeIn && variant == 1: // for BB
			fmt.Printf("\n****************************************\n")
			fmt.Printf("\nDeck: %v   Move: %v   Strategy #: %v  Moves Tried: %v   Unique Boards: %v   Elapsed TD: %v   Elapsed ADs: %v\n", dN, mN, stratNumTD, mvsTriedTD, len(priorBoards), time.Now().Sub(startTimeAD), time.Now().Sub(startTimeTD))
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
			fmt.Printf(comment, dN, mN, stratNumTD, mvsTriedTD, len(priorBoards), time.Now().Sub(startTimeAD), time.Now().Sub(startTimeTD))
		case pTypeIn == "BBS" && printMoveDetail.pType == pTypeIn && variant == 2: // for BBS or BBSS
			fmt.Printf(comment, b)
		case pTypeIn == "NOTX" && printMoveDetail.pType != "X" && variant == 1:
			fmt.Printf(comment, s1, s2, dN, mN, stratNumTD, mvsTriedTD, len(priorBoards), time.Now().Sub(startTimeAD), time.Now().Sub(startTimeTD))
		case (pTypeIn == "TW" || pTypeIn == "TS" || pTypeIn == "TSS") && variant == 1:
		case (pTypeIn == "TW" || pTypeIn == "TS" || pTypeIn == "TSS") && variant == 2:
		}
	}
}
