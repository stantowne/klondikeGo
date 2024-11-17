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

func playAllMoveS(bIn board,
	moveNum int,
	deckNum int,
	cfg Configuration,
	varSp2PN variablesSpecificToPlayNew,
	startTimeTD time.Time) (string, string, int) {
	// if cfg.PlayNew.PrintMoveDetails {}
	/* Return Codes: SL  = Strategy Lost	 NMA = No Moves Available
	                 						 RB  = Repetitive Board
	                                         SE  = Strategy Exhausted
	                                         GLE = GameLength Limit exceeded
	                 SW  = Strategy Win      EW  = Early Win
											 SW  = Standard Win  Obsolete all wins are early
	*/
	// add code for findAllWinStrats
	// Setup pfmt to print thousands with commas
	var pfmt = message.NewPrinter(language.English)

	var recurReturnNum int

	var aMoves []move //available Moves
	var recurReturnV1, recurReturnV2 string
	if moveNum > moveNumMax {
		moveNumMax = moveNum
	}

	// Check to see if the gameLenthLimit has been exceeded.
	//If, treats this as a loss and returns with loss codes.
	if mvsTriedTD >= cfg.PlayNew.GameLengthLimit*1_000_000 {
		prntMDet(bIn,
			aMoves,
			0,
			deckNum,
			moveNum,
			"BB",
			2,
			"\n  SL-GLE: Game Length of: %v exceeds limit: %v\n",
			strconv.Itoa(mvsTriedTD),
			strconv.Itoa(cfg.PlayNew.GameLengthLimit*1_000_000),
			cfg,
			varSp2PN)
		stratLossesGLE_TD++
		prntMDetTreeReturnComment("Game Length Limit Exceeded", deckNum, recurReturnNum, cfg)
		return "SL", "GLE", 1
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

		/* Actually before we actually try all moves let's first: print (optionally based on cfg.PlayNew.RestrictReportingTo.pType) the incoming board
		      and check the incoming board for various end-of-strategy conditions
		   Note: This was done this way, so as to ensure that when returns backed up the moveNum, the board would reprint.
		*/

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
				prntMDetTreeReturnComment("RB", deckNum, recurReturnNum, cfg)
				return "SL", "RB", 1 // Repetitive Board
			} else {
				// Remember the board state by putting it into the map "varSp2PN.priorBoards"
				varSp2PN.priorBoards[bNewBcode] = true
			}
		}

		// Check if No Moves Available
		if i == 0 && aMoves[0].name == "No Moves Available" {
			stratLossesNMA_TD++
			prntMDet(bIn, aMoves, i, deckNum, moveNum, "BB", 2, "  SL-NMA: No Moves Available: Strategy Lost %v%v\n", "", "", cfg, varSp2PN)
			return "SL", "NMA", 1
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
			prntMDetTreeReturnComment("DECK WON", deckNum, recurReturnNum, cfg)
			return "SW", "EW", 1 //  Strategy Early Win
		}

		// OK, done with the various end-of-strategy conditions
		// let's print out the list of available moves and make the next available move
		// The board state was already printed above
		//if cfg.PlayNew.RestrictReporting && cfg.PlayNew.ReportingMoveByMove && cfg.PlayNew.MoveByMoveReportingOptions.Type == "regular" && prntMDetTestRange(deckNum, cfg) {   // DELETE???
		prntMDet(bIn, aMoves, i, deckNum, moveNum, "BB", 3, "", "", "", cfg, varSp2PN) // DELETE???
		//}    // DELETE???

		bNew := bIn.copyBoard() // Critical Must use copyBoard

		// ********** 1st of the 2 MOST IMPORTANT statements in this function:  ******************************
		bNew = moveMaker(bNew, aMoves[i])

		prntMDet(bIn, aMoves, i, deckNum, moveNum, "BBS", 1, "\n\nBefore Call at Deck: %v   Move: %v   Strategy #: %v  Moves Tried: %v   Unique Boards: %v   Elapsed TD: %v   Elapsed ADs: %v\n", "", "", cfg, varSp2PN)
		prntMDet(bIn, aMoves, i, deckNum, moveNum, "BBS", 2, "      bIn: %v\n", "", "", cfg, varSp2PN)
		prntMDet(bNew, aMoves, i, deckNum, moveNum, "BBS", 2, "     bNew: %v\n", "", "", cfg, varSp2PN)
		prntMDetTree(bIn, aMoves, i, deckNum, moveNum, cfg, varSp2PN)

		mvsTriedTD++

		if cfg.PlayNew.ProgressCounter > 0 && math.Mod(float64(mvsTriedTD), float64(cfg.PlayNew.ProgressCounter)) <= 0.1 {
			avgRepTime := time.Since(startTimeTD) / time.Duration(mvsTriedTD/cfg.PlayNew.ProgressCounter)
			estMaxRemTimeTD := time.Since(startTimeTD) * time.Duration((cfg.PlayNew.GameLengthLimit*1_000_000-mvsTriedTD)/mvsTriedTD)
			_, err = pfmt.Printf("\rDk: %5d   ____   MvsTried: %13v   MoveNum: %3v   Max MoveNum: %3v   StratsTried: %12v   UnqBoards: %11v  - %7s  SinceLast: %6s  Avg: %6s  estMaxRem: %7s\r", deckNum, mvsTriedTD, moveNum, moveNumMax, stratNumTD, len(varSp2PN.priorBoards), time.Since(startTimeTD).Truncate(100*time.Millisecond).String(), time.Since(startTimeTD).Truncate(100*time.Millisecond).String(), avgRepTime.Truncate(100*time.Millisecond).String(), estMaxRemTimeTD.Truncate(1*time.Minute))
			//lastPrintTime := time.Now()
		}

		// ********** 2nd of the 2 MOST IMPORTANT statements in this function:  ******************************
		recurReturnV1, recurReturnV2, recurReturnNum = playAllMoveS(bNew, moveNum+1, deckNum, cfg, varSp2PN, startTimeTD)

		// CONSIDER DELETING prntMDet(bIn, aMoves, i, deckNum, moveNum, "NOTX", 1, "  Returned: %v - %v After Call at deckNum: %v  moveNum: %v   StratNumTD: %v   MvsTriedTD: %v   UnqBds: %v   ElTimTD: %v   ElTimADs: %v\n", recurReturnV1, recurReturnV2, cfg, varSp2PN)

		if cfg.PlayNew.FindAllWinStrats != true && recurReturnV1 == "SW" {
			// save winning moves into a slice in reverse
			return recurReturnV1, recurReturnV2, recurReturnNum + 1 // return up the call stack to end strategies search  if findAllWinStrats false, and we had a win
		}
		if recurReturnV1 == "SL" && recurReturnV2 == "GLE" {
			return recurReturnV1, recurReturnV2, recurReturnNum + 1 //
		}
	}

	stratLossesSE_TD++
	prntMDetTreeReturnComment("No More Moves", deckNum, recurReturnNum, cfg)
	return "SL", "SE", recurReturnNum + 1 //  Strategy Exhausted
}

func prntMDetTestRange(deckNum int, cfg Configuration) bool {
	deckRangeOK := false
	if cfg.PlayNew.RestrictReportingTo.DeckStartVal == 0 && cfg.PlayNew.RestrictReportingTo.DeckContinueFor == 0 {
		deckRangeOK = true
	} else {
		if deckNum >= cfg.PlayNew.RestrictReportingTo.DeckStartVal && (cfg.PlayNew.RestrictReportingTo.DeckContinueFor == 0 || deckNum < cfg.PlayNew.RestrictReportingTo.DeckStartVal+cfg.PlayNew.RestrictReportingTo.DeckContinueFor) {
			deckRangeOK = true
		}
	}
	aMvsThisDkRangeOK := false
	if cfg.PlayNew.RestrictReportingTo.MovesTriedStartVal == 0 && cfg.PlayNew.RestrictReportingTo.MovesTriedContinueFor == 0 {
		aMvsThisDkRangeOK = true
	} else {
		if mvsTriedTD >= cfg.PlayNew.RestrictReportingTo.MovesTriedStartVal && (cfg.PlayNew.RestrictReportingTo.MovesTriedContinueFor == 0 || mvsTriedTD < cfg.PlayNew.RestrictReportingTo.MovesTriedStartVal+cfg.PlayNew.RestrictReportingTo.MovesTriedContinueFor) {
			aMvsThisDkRangeOK = true
		}
	}
	if deckRangeOK && aMvsThisDkRangeOK {
		return true
	} else {
		return false
	}
}

func prntMDet(b board,
	aMoves []move,
	nextMove int,
	dN int, //deck number
	mN int, //move number
	pTypeIn string,
	variant int,
	comment string,
	s1 string,
	s2 string,
	cfg Configuration,
	varSp2PN variablesSpecificToPlayNew) {
	// Done here just to clean up mainline logic of playAllMoves
	// Do some repetitive printing to track progress
	// This function will use the struct pMD
	//      variant will be used for different outputs under the same pType

	// Setup pfmt to print thousands with commas
	var pfmt = message.NewPrinter(language.English)

	if cfg.PlayNew.RestrictReporting && prntMDetTestRange(dN, cfg) {
		switch {
		case pTypeIn == "BB" && cfg.PlayNew.ReportingMoveByMove && cfg.PlayNew.MoveByMoveReportingOptions.Type == "regular" && variant == 1: // for BB
			fmt.Printf("\n****************************************\n")
			_, err = pfmt.Printf("\nDeck: %v   Move: %v   Strategy #: %v  Moves Tried: %v   Unique Boards: %v   Elapsed TD: %v   Elapsed ADs: %v\n",
				dN,
				mN,
				stratNumTD,
				mvsTriedTD,
				len(varSp2PN.priorBoards),
				time.Since(startTimeAD),
				time.Since(startTimeTD))
			printBoard(b)
		case pTypeIn == "BB" && cfg.PlayNew.ReportingMoveByMove && cfg.PlayNew.MoveByMoveReportingOptions.Type == "regular" && variant == 2: // for BB
			// comment must have 2 %v in it
			_, err = pfmt.Printf(comment, s1, s2)
		case pTypeIn == "BB" && cfg.PlayNew.ReportingMoveByMove && cfg.PlayNew.MoveByMoveReportingOptions.Type == "regular" && variant == 3:
			fmt.Printf("\n     All Possible Moves: ")
			for j := range aMoves {
				pM1, pM2 := printMove(aMoves[j])
				if nextMove == j {
					fmt.Printf("   Next Move ->    ")
				} else {
					fmt.Printf("                   ")
				}
				fmt.Printf("%s", pM1)
				if nextMove == j {
					pad := strings.Repeat(" ", 110-len([]rune(pM1)))
					fmt.Printf(pad + "<- Next Move")
				}
				if len(pM2) > 0 {
					fmt.Printf("%s", pM2)
				}

				fmt.Printf("\n")
			}
		case strings.HasPrefix(pTypeIn, "BBS") && cfg.PlayNew.ReportingMoveByMove && (cfg.PlayNew.MoveByMoveReportingOptions.Type == "short" || cfg.PlayNew.MoveByMoveReportingOptions.Type == "very short") && variant == 1: // for BBS or BBSS
			_, err = pfmt.Printf(comment, dN, mN, stratNumTD, mvsTriedTD, len(varSp2PN.priorBoards), time.Since(startTimeAD), time.Since(startTimeTD))
		case pTypeIn == "BBS" && cfg.PlayNew.ReportingMoveByMove && cfg.PlayNew.MoveByMoveReportingOptions.Type == "short" && variant == 2:
			_, err = pfmt.Printf(comment, b)
		case pTypeIn == "NOTX" && cfg.PlayNew.ReportingMoveByMove && variant == 1:
			_, err = pfmt.Printf(comment, s1, s2, dN, mN, stratNumTD, mvsTriedTD, len(varSp2PN.priorBoards), time.Since(startTimeAD), time.Since(startTimeTD))
		case pTypeIn == "NOTX" && cfg.PlayNew.ReportingMoveByMove && variant == 2:
			_, err = pfmt.Printf(comment, s1, s2)
		case (pTypeIn == "TW" || pTypeIn == "TS" || pTypeIn == "TSS") && variant == 1:
		case (pTypeIn == "TW" || pTypeIn == "TS" || pTypeIn == "TSS") && variant == 2:
		}
	}
}

func prntMDetTree(b board, aMoves []move, nextMove int, dN int, mN int, cfg Configuration, varSp2PN variablesSpecificToPlayNew) {
	//

	// Setup pfmt to print thousands with commas
	//var pfmt = message.NewPrinter(language.English)

	const (
		vert1      = string('\u2503') // Looks Like: ->┃<-
		horiz1     = string('\u2501') // Looks Like: ->━<-
		firstStrat = string('\u2533') // Looks Like: ->┳<-
		lastStrat  = string('\u2517') // Looks Like: ->┗<-
		midStrat   = string('\u2523') // Looks Like: ->┣<-
	)
	var treeThisMove string
	var treeAddToPrev string
	var treeMoveWidth int
	var treeRepeatChar string

	if prntMDetTestRange(dN, cfg) && cfg.PlayNew.ReportingType.Tree {
		if mN == 0 && nextMove == 0 {
			_, err = pfmt.Printf("\n\n Deck: %v\n\n", dN)
			printBoard(b)
			fmt.Printf("\n\n Strategy #   ")
			if cfg.PlayNew.TreeReportingOptions.Type == "regular" {
				fmt.Printf("\n             ")
				for i := 1; i <= 150; i++ {
					fmt.Printf("%8s", strconv.Itoa(i)+"  ")
				}
			}
			fmt.Printf("\n\n            0  ")
		}
		switch {
		case len(aMoves) == 1:
			treeThisMove = horiz1
			treeAddToPrev = " "
			treeRepeatChar = " "
		case nextMove == 0:
			treeThisMove = firstStrat
			treeAddToPrev = vert1
			treeRepeatChar = " "
		case nextMove == len(aMoves)-1:
			treeThisMove = lastStrat
			treeAddToPrev = " "
			treeRepeatChar = " "
		default:
			treeThisMove = midStrat
			treeAddToPrev = vert1
			treeRepeatChar = " "
		}
		switch cfg.PlayNew.TreeReportingOptions.Type {
		case "very narrow":
			//varSp2PN.treePrevMovesTD += treeAddToPrev
			treeMoveWidth = 1
		case "narrow":
			treeThisMove += strings.Repeat(horiz1, 2)
			treeMoveWidth = 3
			treeAddToPrev += strings.Repeat(treeRepeatChar, 2)
		case "regular":
			treeThisMove += moveShortName[aMoves[nextMove].name] + horiz1
			treeMoveWidth = 8
			treeAddToPrev += strings.Repeat(treeRepeatChar, 7)
		}
		if nextMove == 0 {
			time.Sleep(cfg.PlayNew.TreeReportingOptions.TreeSleepBetwnMoves)
			fmt.Printf("%s", treeThisMove)
		} else {
			time.Sleep(cfg.PlayNew.TreeReportingOptions.TreeSleepBetwnStrategies)
			//x := []rune(varSp2PN.treePrevMovesTD)
			x := []rune(varSp2PN.treePrevMovesTD)
			x = x[0 : mN*treeMoveWidth]
			//varSp2PN.treePrevMovesTD = string(x)
			varSp2PN.treePrevMovesTD = string(x)
			//pfmt.Printf("\n%13s  %s%s", strconv.Itoa(stratNumTD), varSp2PN.treePrevMovesTD, treeThisMove)
			_, err = pfmt.Printf("\n%13s  %s%s", strconv.Itoa(stratNumTD), varSp2PN.treePrevMovesTD, treeThisMove)
		}
		//varSp2PN.treePrevMovesTD += treeAddToPrev
		varSp2PN.treePrevMovesTD += treeAddToPrev
	}
}

func prntMDetTreeReturnComment(c string, dN int, recurReturnNum int, cfg Configuration) {
	if prntMDetTestRange(dN, cfg) && cfg.PlayNew.ReportingType.Tree && recurReturnNum == 0 {
		fmt.Printf(" ==> " + c)
	}
}
