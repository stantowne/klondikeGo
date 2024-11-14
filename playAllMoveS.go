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

func playAllMoveS(bIn board, moveNum int, deckNum int, cLArgs commandLineArgs, varSp2PN variablesSpecificToPlayNew) (string, string, int) {

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

	if mvsTriedTD >= gameLengthLimit {
		prntMDet(bIn, aMoves, 0, deckNum, moveNum, "BB", 2, "\n  SL-GLE: Game Length of: %v exceeds limit: %v\n", strconv.Itoa(mvsTriedTD), strconv.Itoa(gameLengthLimit), cLArgs, varSp2PN)
		stratLossesGLE_TD++
		prntMDetTreeReturnComment("Game Length Limit Exceeded", deckNum, recurReturnNum, cLArgs)
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

		/* Actually before we actually try all moves let's first: print (optionally based on pMD.pType) the incoming board
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
		prntMDet(bIn, aMoves, i, deckNum, moveNum, "BB", 1, "", "", "", cLArgs, varSp2PN)

		if i == 0 {
			// Check for repetitive board
			// We only have to check if i == 0 because if i >= 0 then we are just returning to this board at which point
			//   the board will already have been checked (when i == 0 {past tense})
			bNewBcode := bIn.boardCode(deckNum) //  Just use the array itself
			// Have we seen this board before?
			if varSp2PN.priorBoards[bNewBcode] {
				// OK we did see it before so return to try next available move (if any) in aMoves[] aka strategy
				stratLossesRB_TD++
				prntMDet(bIn, aMoves, i, deckNum, moveNum, "BB", 2, "\n  SF-RB: Repetitive Board - \"Next Move\" yielded a repeat of a board.\n", "", "", cLArgs, varSp2PN)
				prntMDet(bIn, aMoves, i, deckNum, moveNum, "BBSS", 2, "\n  SF-RB: Repetitive Board - \"Next Move\" yielded a repeat of a board.\n", "", "", cLArgs, varSp2PN)
				prntMDetTreeReturnComment("Repetitive Board", deckNum, recurReturnNum, cLArgs)
				return "SL", "RB", 1 // Repetitive Board
			} else {
				// Remember the board state by putting it into the map "varSp2PN.priorBoards"
				varSp2PN.priorBoards[bNewBcode] = true
			}
		}

		// Check if No Moves Available
		if i == 0 && aMoves[0].name == "No Moves Available" {
			stratLossesNMA_TD++
			prntMDet(bIn, aMoves, i, deckNum, moveNum, "BB", 2, "  SL-NMA: No Moves Available: Strategy Lost %v%v\n", "", "", cLArgs, varSp2PN)
			prntMDetTreeReturnComment("No Moves Available", deckNum, recurReturnNum, cLArgs)
			return "SL", "NMA", 1
		}

		//Detect Early Win
		if detectWinEarly(bIn) {
			stratWinsTD++
			cmt := "  SW-EW: Strategy Win: Early Win%v%v"
			if cLArgs.findAllWinStrats {
				cmt += "  Will Continue to look for additional winning strategies for this deck"
			} else {
				cmt += "  Go to Next Deck (if any)"
			}
			prntMDet(bIn, aMoves, i, deckNum, moveNum, "NOTX", 1, cmt, "", "", cLArgs, varSp2PN)

			// Verbose Special "WL" Starts Here - No effect on operation
			if strings.Contains(cLArgs.verboseSpecial, ";WL;") { // Deck Win Loss Summary Statistics   MOVE THIS!!!!!!!!!
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
			prntMDetTreeReturnComment("DECK WON", deckNum, recurReturnNum, cLArgs)
			return "SW", "EW", 1 //  Strategy Early Win
		}

		// OK, done with the various end-of-strategy conditions
		// let's print out the list of available moves and make the next available move
		// The board state was already printed above
		//if pMD.pType == "BB" && prntMDetTestRange(deckNum, cLArgs) {  // DELETE???
		prntMDet(bIn, aMoves, i, deckNum, moveNum, "BB", 3, "", "", "", cLArgs, varSp2PN)
		//}                                                             // DELETE???

		bNew := bIn.copyBoard() // Critical Must use copyBoard

		// ********** 1st of the 2 MOST IMPORTANT statements in this function:  ******************************
		bNew = moveMaker(bNew, aMoves[i])

		prntMDet(bIn, aMoves, i, deckNum, moveNum, "BBS", 1, "\n\nBefore Call at Deck: %v   Move: %v   Strategy #: %v  Moves Tried: %v   Unique Boards: %v   Elapsed TD: %v   Elapsed ADs: %v\n", "", "", cLArgs, varSp2PN)
		prntMDet(bIn, aMoves, i, deckNum, moveNum, "BBS", 2, "      bIn: %v\n", "", "", cLArgs, varSp2PN)
		prntMDetTree(bIn, aMoves, i, deckNum, moveNum, cLArgs, varSp2PN)

		mvsTriedTD++ // Two above printed before bumping mVStried because they print bIN  The next prnt prints bNew

		prntMDet(bNew, aMoves, i, deckNum, moveNum, "BBS", 2, "     bNew: %v\n", "", "", cLArgs, varSp2PN)

		// verboseSpecial option PROGRESSdddd starts here      Consider Moving clause to be a clause in prntMDet
		if cLArgs.verboseSpecialProgressCounter > 0 && math.Mod(float64(mvsTriedTD), float64(cLArgs.verboseSpecialProgressCounter)) <= 0.1 {
			avgRepTime := time.Since(startTimeTD) / time.Duration(mvsTriedTD/cLArgs.verboseSpecialProgressCounter)
			estMaxRemTimeTD := time.Since(startTimeTD) * time.Duration((gameLengthLimit-mvsTriedTD)/mvsTriedTD)
			pfmt.Printf("\rDk: %5d   ____   MvsTried: %13v   MoveNum: %3v   Max MoveNum: %3v   StratsTried: %12v   UnqBoards: %11v  - %7s  SinceLast: %6s  Avg: %6s  estMaxRem: %7s\r", deckNum, mvsTriedTD, moveNum, moveNumMax, stratNumTD, len(varSp2PN.priorBoards), time.Since(startTimeTD).Truncate(100*time.Millisecond).String(), time.Since(cLArgs.verboseSpecialProgressCounterLastPrintTime).Truncate(100*time.Millisecond).String(), avgRepTime.Truncate(100*time.Millisecond).String(), estMaxRemTimeTD.Truncate(1*time.Minute))
			cLArgs.verboseSpecialProgressCounterLastPrintTime = time.Now()
		}
		// verboseSpecial option PROGRESSdddd ends here

		// ********** 2nd of the 2 MOST IMPORTANT statements in this function:  ******************************
		recurReturnV1, recurReturnV2, recurReturnNum = playAllMoveS(bNew, moveNum+1, deckNum, cLArgs, varSp2PN)

		// CONSIDER DELETING prntMDet(bIn, aMoves, i, deckNum, moveNum, "NOTX", 1, "  Returned: %v - %v After Call at deckNum: %v  moveNum: %v   StratNumTD: %v   MvsTriedTD: %v   UnqBds: %v   ElTimTD: %v   ElTimADs: %v\n", recurReturnV1, recurReturnV2, cLArgs, varSp2PN)

		if cLArgs.findAllWinStrats != true && recurReturnV1 == "SW" {
			// save winning moves into a slice in reverse
			return recurReturnV1, recurReturnV2, recurReturnNum + 1 // return up the call stack to end strategies search  if findAllWinStrats false, and we had a win
		}
		if recurReturnV1 == "SL" && recurReturnV2 == "GLE" {
			return recurReturnV1, recurReturnV2, recurReturnNum + 1 //
		}
	}

	stratLossesSE_TD++
	prntMDetTreeReturnComment("No More Moves", deckNum, recurReturnNum, cLArgs)
	return "SL", "SE", recurReturnNum + 1 //  Strategy Exhausted
}

func prntMDetTestRange(deckNum int, cLArgs commandLineArgs) bool {
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

func prntMDet(b board, aMoves []move, nextMove int, dN int, mN int, pTypeIn string, variant int, comment string, s1 string, s2 string, cLArgs commandLineArgs, varSp2PN variablesSpecificToPlayNew) {
	// Done here just to clean up mainline logic of playAllMoves
	// Do some repetitive printing to track progress
	// This function will use the struct pMD
	//      variant will be used for different outputs under the same pType

	// Setup pfmt to print thousands with commas
	var pfmt = message.NewPrinter(language.English)

	if pMD.pType != "X" && prntMDetTestRange(dN, cLArgs) {
		switch {
		case pTypeIn == "BB" && pMD.pType == pTypeIn && variant == 1: // for BB
			fmt.Printf("\n****************************************\n")
			pfmt.Printf("\nDeck: %v   Move: %v   Strategy #: %v  Moves Tried: %v   Unique Boards: %v   Elapsed TD: %v   Elapsed ADs: %v\n", dN, mN, stratNumTD, mvsTriedTD, len(varSp2PN.priorBoards), time.Since(startTimeAD), time.Since(startTimeTD))
			printBoard(b)
		case pTypeIn == "BB" && pMD.pType == pTypeIn && variant == 2: // for BB
			// comment must have 2 %v in it
			pfmt.Printf(comment, s1, s2)
		case pMD.pType == "BB" && variant == 3:
			fmt.Printf("\n                   ALL POSSIBLE MOVES:\n")
			if mN == 55 {
				time.Sleep(1 * time.Second)
			}
			for j := range aMoves {
				if nextMove == j {
					fmt.Printf("   Next Move ->    ")
				} else {
					fmt.Printf("                   ")
				}
				fmt.Printf("%s", printMove(aMoves[j]))
				if nextMove == j {
					pad := strings.Repeat(" ", 110-len([]rune(printMove(aMoves[j]))))
					fmt.Printf(pad + "<- Next Move")
				}
				fmt.Printf("\n")
			}
		case strings.HasPrefix(pTypeIn, "BBS") && strings.HasPrefix(pMD.pType, "BBS") && variant == 1: // for BBS or BBSS
			pfmt.Printf(comment, dN, mN, stratNumTD, mvsTriedTD, len(varSp2PN.priorBoards), time.Since(startTimeAD), time.Since(startTimeTD))
		case pTypeIn == "BBS" && pMD.pType == pTypeIn && variant == 2:
			pfmt.Printf(comment, b)
		case pTypeIn == "NOTX" && pMD.pType != "X" && variant == 1 && !strings.Contains(";TW;TS;TSS;", ";"+pMD.pType+";"):
			pfmt.Printf(comment, s1, s2, dN, mN, stratNumTD, mvsTriedTD, len(varSp2PN.priorBoards), time.Since(startTimeAD), time.Since(startTimeTD))
		case pTypeIn == "NOTX" && pMD.pType != "X" && variant == 2 && !strings.Contains(";TW;TS;TSS;", ";"+pMD.pType+";"):
			pfmt.Printf(comment, s1, s2)
		}
	}
}

func prntMDetTree(b board, aMoves []move, nextMove int, dN int, mN int, cLArgs commandLineArgs, varSp2PN variablesSpecificToPlayNew) {
	//

	// Setup pfmt to print thousands with commas
	var pfmt = message.NewPrinter(language.English)

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

	if prntMDetTestRange(dN, cLArgs) && strings.Contains(";TW;TS;TSS;", ";"+pMD.pType+";") {
		if mN == 0 && nextMove == 0 {
			pfmt.Printf("\n\n Deck: %v\n\n", dN)
			printBoard(b)
			fmt.Printf("\n\n Strategy #   ")
			if pMD.pType == "TW" {
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
		switch pMD.pType {
		case "TSS":
			//varSp2PN.treePrevMovesTD += treeAddToPrev
			treeMoveWidth = 1
		case "TS":
			treeThisMove += strings.Repeat(horiz1, 2)
			treeMoveWidth = 3
			treeAddToPrev += strings.Repeat(treeRepeatChar, 2)
		case "TW":
			treeThisMove += moveShortName[aMoves[nextMove].name] + horiz1
			treeMoveWidth = 8
			treeAddToPrev += strings.Repeat(treeRepeatChar, 7)
		}
		if nextMove == 0 {
			time.Sleep(treeSleepBetwnMoves)
			fmt.Printf("%s", treeThisMove)
		} else {
			time.Sleep(treeSleepBetwnStrategies)
			//x := []rune(varSp2PN.treePrevMovesTD)
			x := []rune(treePrevMovesTD)
			x = x[0 : mN*treeMoveWidth]
			//varSp2PN.treePrevMovesTD = string(x)
			treePrevMovesTD = string(x)
			//pfmt.Printf("\n%13s  %s%s", strconv.Itoa(stratNumTD), varSp2PN.treePrevMovesTD, treeThisMove)
			pfmt.Printf("\n%13s  %s%s", strconv.Itoa(stratNumTD), treePrevMovesTD, treeThisMove)
		}
		//varSp2PN.treePrevMovesTD += treeAddToPrev
		treePrevMovesTD += treeAddToPrev
	}
}

func prntMDetTreeReturnComment(c string, dN int, recurReturnNum int, cLArgs commandLineArgs) {
	if prntMDetTestRange(dN, cLArgs) && strings.Contains(";TW;TS;TSS;", ";"+pMD.pType+";") && recurReturnNum == 0 {
		fmt.Printf(" ==> " + c)
	}
}
