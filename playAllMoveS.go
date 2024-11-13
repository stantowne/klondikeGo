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

func playAllMoveS(bIn board, moveNum int, deckNum int, cfg Configuration, varSp2PN variablesSpecificToPlayNew, startTimeTD time.Time) (string, string) {
	if cfg.PlayNew.PrintMoveDetails {

	}
	/* Return Codes: SL  = Strategy Lost	 NMA = No Moves Available
	                 						 RB  = Repetitive Board
	                                         SE  = Strategy Exhausted
	                                         GML = GameLength Limit exceeded
	                 SW  = Strategy Win      EW  = Early Win
											 SW  = Standard Win  Obsolete all wins are early
	*/

	// Setup pfmt to print thousands with commas
	var pfmt = message.NewPrinter(language.English)
	var aMoves []move //available Moves
	var recurReturnV1, recurReturnV2 string
	if moveNum > moveNumMax {
		moveNumMax = moveNum
	}

	if mvsTriedTD >= gameLengthLimit {
		prntMDet(bIn, aMoves, 0, deckNum, moveNum, "BB", 2, "\n  SL-RB: Game Length of: %v exceeds limit: %v\n", strconv.Itoa(mvsTriedTD), strconv.Itoa(gameLengthLimit), cfg, varSp2PN)
		stratLossesGML_TD++
		return "SL", "GML"
	}

	// Find Next Moves
	aMoves = detectAvailableMoves(bIn, moveNum, singleGame)

	if len(aMoves) == 0 {
		m := move{name: "No Moves Available"} // This is a pseudo move not created by detectAvailable Moves it exists to remember
		aMoves = append(aMoves, m)            //      this state and for the various printing routines that come below.  If a return
		//      was made from this point no history of it would exist.  See code at "// Check if No Moves Available"
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

		// Print the incoming board EVEN IF we are returning to it to try the next available move
		prntMDet(bIn, aMoves, i, deckNum, moveNum, "BB", 1, "", "", "", cfg, varSp2PN)

		// Possible increment to stratNumTD
		if i != 0 {
			// Started at 0 in playNew for each deck.  Increment each time a second through nth move is made
			//     That way strategy 0 is the "All Best Moves" strategy.  It is also why in playNew aMStratTriedAllDecks
			//     is incremented by stratNumTD + 1 after each deck.
			stratNumTD++
		}

		// Print the incoming board EVEN IF we are returning to it to try the next available move
		//       This had to be done after possible increment to stratNumTD so that each time a board is reprinted it shows the NEW strategy number
		//       Before when it was above the possible increment the board was printing out with the stratNum of the last failed strategy
		prntMDet(bIn, aMoves, i, deckNum, moveNum, "BB", 1, "", "", "", cfg, varSp2PN)

		if i == 0 {
			// Check for repetitive board
			// We only have to check if i == 0 because if i >= 0 then we are just returning to this board at which point
			//   the board will already have been checked (when i == 0 {past tense})
			bNewBcode := bIn.boardCode(deckNum) //  Just use the array itself
			// Have we seen this board before?
			if varSp2PN.priorBoards[bNewBcode] {
				// OK we did see it before so return to try next available move (if any) in aMoves[] aka strategy
				stratLossesRB_TD++
				//  	prntMDet(b board, aMoves []move, nextMove int, dN int, mN int, pTypeIn string, variant int, comment string, s1 string, s2 string) {
				prntMDet(bIn, aMoves, i, deckNum, moveNum, "BB", 2, "\n  SF-RB: Repetitive Board - \"Next Move\" yielded a repeat of a board.\n", "", "", cfg, varSp2PN)
				prntMDet(bIn, aMoves, i, deckNum, moveNum, "BBSS", 2, "\n  SF-RB: Repetitive Board - \"Next Move\" yielded a repeat of a board.\n", "", "", cfg, varSp2PN)
				return "SL", "RB" // Repetitive Move
			} else {
				// Remember the board state by putting it into the map "varSp2PN.priorBoards"
				varSp2PN.priorBoards[bNewBcode] = true
			}
		}

		// Check if No Moves Available
		if i == 0 && aMoves[0].name == "No Moves Available" {
			stratLossesNMA_TD++
			prntMDet(bIn, aMoves, i, deckNum, moveNum, "BB", 2, "  SL-NMA: No Moves Available: Strategy Lost %v%v\n", "", "", cfg, varSp2PN)
			return "SL", "NMA"
		}

		//Detect Early Win
		if detectWinEarly(bIn) {
			stratWinsTD++
			cmt := "  SW-EW: Strategy Win: Early Win%v%v"
			if cfg.PlayNew.FindAllWinStrats {
				cmt += "  Will Continue to look for additional winning strategies for this deck"
			} else {
				cmt += "  Go to Next Deck (if any)"
			}
			prntMDet(bIn, aMoves, i, deckNum, moveNum, "NOTX", 1, cmt, "", "", cfg, varSp2PN)

			// Verbose Special "WL" Starts Here - No effect on operation
			if cfg.PlayNew.WinLossReport { // Deck Win Loss Summary Statistics   MOVE THIS!!!!!!!!!
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
		if cfg.PlayNew.PrintMoveDetails && cfg.PlayNew.PrintMoveDetailOptions.Type == "BB" && prntMDetTestRange(deckNum, cfg) {
			prntMDet(bIn, aMoves, i, deckNum, moveNum, "BB", 3, "", "", "", cfg, varSp2PN)
		}

		bNew := bIn.copyBoard() // Critical Must use copyBoard

		// ********** 1st of the 2 MOST IMPORTANT statements in this function:  ******************************
		bNew = moveMaker(bNew, aMoves[i])

		prntMDet(bIn, aMoves, i, deckNum, moveNum, "BBS", 1, "\n\nBefore Call at Deck: %v   Move: %v   Strategy #: %v  Moves Tried: %v   Unique Boards: %v   Elapsed TD: %v   Elapsed ADs: %v\n", "", "", cfg, varSp2PN)
		prntMDet(bIn, aMoves, i, deckNum, moveNum, "BBS", 2, "      bIn: %v\n", "", "", cfg, varSp2PN)
		prntMDet(bNew, aMoves, i, deckNum, moveNum, "BBS", 2, "     bNew: %v\n", "", "", cfg, varSp2PN)
		mvsTriedTD++

		// verboseSpecial option PROGRESSdddd starts here      Consider Moving clause to be a clause in prntMDet
		if cfg.PlayNew.ProgressCounter > 0 && math.Mod(float64(mvsTriedTD), float64(cfg.PlayNew.ProgressCounter)) <= 0.1 {
			avgRepTime := time.Since(startTimeTD) / time.Duration(mvsTriedTD/cfg.PlayNew.ProgressCounter)
			estMaxRemTimeTD := time.Since(startTimeTD) * time.Duration((gameLengthLimit-mvsTriedTD)/mvsTriedTD)
			_, err = pfmt.Printf("\rDk: %5d   ____   MvsTried: %13v   MoveNum: %3v   Max MoveNum: %3v   StratsTried: %12v   UnqBoards: %11v  - %7s  SinceLast: %6s  Avg: %6s  estMaxRem: %7s\r", deckNum, mvsTriedTD, moveNum, moveNumMax, stratNumTD, len(varSp2PN.priorBoards), time.Since(startTimeTD).Truncate(100*time.Millisecond).String(), time.Since(startTimeTD).Truncate(100*time.Millisecond).String(), avgRepTime.Truncate(100*time.Millisecond).String(), estMaxRemTimeTD.Truncate(1*time.Minute))
			//lastPrintTime := time.Now()
		}
		// verboseSpecial option PROGRESSdddd ends here

		// ********** 2nd of the 2 MOST IMPORTANT statements in this function:  ******************************
		recurReturnV1, recurReturnV2 = playAllMoveS(bNew, moveNum+1, deckNum, cfg, varSp2PN, startTimeTD)

		prntMDet(bIn, aMoves, i, deckNum, moveNum, "NOTX", 1, "  Returned: %v - %v After Call at deckNum: %v  moveNum: %v   StratNumTD: %v   MvsTriedTD: %v   UnqBds: %v   ElTimTD: %v   ElTimADs: %v\n", recurReturnV1, recurReturnV2, cfg, varSp2PN)

		if cfg.PlayNew.FindAllWinStrats != true && recurReturnV1 == "SW" {
			// save winning moves into a slice in reverse
			return recurReturnV1, recurReturnV2 // return up the call stack to end strategies search  if findAllWinStrats false, and we had a win
		}
		if recurReturnV1 == "SL" && recurReturnV2 == "GML" {
			return recurReturnV1, recurReturnV2 //
		}
	}

	stratLossesSE_TD++
	return "SL", "SE" //  Strategy Exhausted
}

func prntMDetTestRange(deckNum int, cfg Configuration) bool {
	deckRangeOK := false
	pMD := cfg.PlayNew.PrintMoveDetailOptions
	if pMD.DeckStartVal == 0 && pMD.DeckContinueFor == 0 {
		deckRangeOK = true
	} else {
		if deckNum >= pMD.DeckStartVal && (pMD.DeckContinueFor == 0 || deckNum < pMD.DeckStartVal+pMD.DeckContinueFor) {
			deckRangeOK = true
		}
	}
	aMvsThisDkRangeOK := false
	if pMD.MovesTriedStartVal == 0 && pMD.MovesTriedContinueFor == 0 {
		aMvsThisDkRangeOK = true
	} else {
		if mvsTriedTD >= pMD.MovesTriedStartVal && (pMD.MovesTriedContinueFor == 0 || mvsTriedTD < pMD.MovesTriedStartVal+pMD.MovesTriedContinueFor) {
			aMvsThisDkRangeOK = true
		}
	}
	if deckRangeOK && aMvsThisDkRangeOK {
		return true
	} else {
		return false
	}
}

func prntMDet(b board, aMoves []move, nextMove int, dN int, mN int, pTypeIn string, variant int, comment string, s1 string, s2 string, cfg Configuration, varSp2PN variablesSpecificToPlayNew) {
	pMD := cfg.PlayNew.PrintMoveDetailOptions
	// Done here just to clean up mainline logic of playAllMoves
	// Do some repetitive printing to track progress
	// This function will use the struct pMD
	//      variant will be used for different outputs under the same pType

	// Setup pfmt to print thousands with commas
	var pfmt = message.NewPrinter(language.English)

	if pMD.Type != "none" && prntMDetTestRange(dN, cfg) {
		switch {
		case pTypeIn == "BB" && pMD.Type == pTypeIn && variant == 1: // for BB
			fmt.Printf("\n****************************************\n")
			_, err = pfmt.Printf("\nDeck: %v   Move: %v   Strategy #: %v  Moves Tried: %v   Unique Boards: %v   Elapsed TD: %v   Elapsed ADs: %v\n", dN, mN, stratNumTD, mvsTriedTD, len(varSp2PN.priorBoards), time.Since(startTimeAD), time.Since(startTimeTD))
			printBoard(b)
		case pTypeIn == "BB" && pMD.Type == pTypeIn && variant == 2: // for BB
			// comment must have 2 %v in it
			_, err = pfmt.Printf(comment, s1, s2)
		case pMD.Type == "BB" && variant == 3:
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
		case strings.HasPrefix(pTypeIn, "BBS") && strings.HasPrefix(pMD.Type, "BBS") && variant == 1: // for BBS or BBSS
			_, err = pfmt.Printf(comment, dN, mN, stratNumTD, mvsTriedTD, len(varSp2PN.priorBoards), time.Since(startTimeAD), time.Since(startTimeTD))
		case pTypeIn == "BBS" && pMD.Type == pTypeIn && variant == 2:
			_, err = pfmt.Printf(comment, b)
		case pTypeIn == "NOTX" && pMD.Type != "X" && variant == 1:
			_, err = pfmt.Printf(comment, s1, s2, dN, mN, stratNumTD, mvsTriedTD, len(varSp2PN.priorBoards), time.Since(startTimeAD), time.Since(startTimeTD))
		case pTypeIn == "NOTX" && pMD.Type != "X" && variant == 2:
			_, err = pfmt.Printf(comment, s1, s2)
		case (pTypeIn == "TW" || pTypeIn == "TS" || pTypeIn == "TSS") && variant == 1:
		case (pTypeIn == "TW" || pTypeIn == "TS" || pTypeIn == "TSS") && variant == 2:
		}
	}
}
