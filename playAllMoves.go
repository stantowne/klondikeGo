package main

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
)

func playAllMoves(bIn board,
	moveNum int,
	deckNum int,
	cfg *Configuration,
	vPA *variablesSpecificToPlayAll,
) (string, int) {

	/* Return Codes: NMA = Strategy Loss: No Moves Available

	   A "Strategy" is any unique set of moves on a path leading back to move 0 aka a "Branch"
	   Note: When available moves = 1 a new "Strategy"/"Branch" is NOT created
	         When available moves > 1 a new "Strategy"/"Branch" IS     created

	   RB     = Strategy Loss: Repetitive Board
	   MajSE  = Major Strategy Loss: all possible moves Exhausted
	   MinSE  = Minor Strategy Loss: all possible moves Exhausted Not Really Considered a strategy!
	   GLE    = Strategy Loss: GameLength Limit Exceeded
	   EL     = Strategy Loss: Early Loss
	   SW     = Strategy Win

	*/
	// add code for findAllWinStrats

	var aMoves []move //available Moves
	var recurReturnV1 string
	var recurReturnNum int

	if moveNum > vPA.TDotherSQL.moveNumMax {
		vPA.TDotherSQL.moveNumMax = moveNum
	}

	// Check to see if the gameLengthLimit has been exceeded.
	//If, treats this as a loss and returns with loss codes.
	if vPA.TD.mvsTried >= cfg.PlayAll.GameLengthLimit*1_000_000 {
		prntMDet(bIn,
			aMoves,
			0,
			deckNum,
			moveNum,
			"MbM_ANY",
			2,
			"GLE: Game Length of: %v exceeds limit: %v\n",
			strconv.Itoa(vPA.TD.mvsTried),
			strconv.Itoa(cfg.PlayAll.GameLengthLimit*1_000_000),
			cfg,
			vPA)
		vPA.TD.stratLossesGLE++
		prntMDetTreeReturnComment(" ==> GLE", deckNum, recurReturnNum, cfg, vPA)
		return "GLE", recurReturnNum + 1
	}

	// Find Next Moves
	aMoves = detectAvailableMoves(bIn, moveNum, cfg.General.NumberOfDecksToBePlayed == 1)

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

		/* Actually before we actually try all moves let's first: print (optionally based on cfg.PlayAll.RestrictReportingTo.pType) the incoming board
		      and check the incoming board for various end-of-strategy conditions
		   Note: This was done this way, so as to ensure that when returns backed up the moveNum, the board would reprint.
		*/

		// Possible increment to vPA.TD.stratNum
		if i != 0 {
			// Started at 0 in playAll for each deck.  Increment each time a second through nth move is made
			//     That way strategy 0 is the "All Best Moves" strategy.
			vPA.TD.stratNum++
			vPA.TD.stratTried++
		}

		// Print the incoming board EVEN IF we are returning to it to try the next available move
		//       This had to be done after possible increment to vPA.TD.stratNum so that each time a board is reprinted it shows the NEW strategy number
		//       Before when it was above the possible increment the board was printing out with the stratNum of the last failed strategy
		prntMDet(bIn, aMoves, i, deckNum, moveNum, "MbM_ANY", 1, "", "", "", cfg, vPA)

		if i == 0 {
			// Check for repetitive board
			// We only have to check if i == 0 because if i >= 0 then we are just returning to this board at which point
			//   the board will already have been checked (when i == 0 {past tense})
			bNewBcode := bIn.boardCode(deckNum) //  Just use the array itself
			// Have we seen this board before?
			if vPA.TDother.priorBoards[bNewBcode] {
				// OK we did see it before so return to try next available move (if any) in aMoves[] aka strategy
				vPA.TD.stratLossesRB++
				prntMDet(bIn, aMoves, i, deckNum, moveNum, "MbM_ANY", 2, "RB: Repetitive Board - \"Next Move\" yielded a repeat of a board.", "", "", cfg, vPA)
				prntMDetTreeReturnComment(" ==> RB", deckNum, recurReturnNum, cfg, vPA)
				return "RB", recurReturnNum + 1 // Repetitive Board
			} else {
				// Remember the board state by putting it into the map "vPA.TDother.priorBoards"
				vPA.TDother.priorBoards[bNewBcode] = true
			}
		}

		// Check if No Moves Available
		if i == 0 && aMoves[0].name == "No Moves Available" {
			vPA.TD.stratLossesNMA++
			prntMDet(bIn, aMoves, i, deckNum, moveNum, "MbM_ANY", 2, "NMA: No Moves Available: Strategy Lost %v%v\n", "", "", cfg, vPA)
			prntMDetTreeReturnComment(" ==> NMA", deckNum, recurReturnNum, cfg, vPA)
			return "NMA", recurReturnNum + 1
		}

		//Detect Win (formerly Early Win)
		if detectWinEarly(bIn) {
			vPA.TD.stratWins++
			/*  delete
			cmt := "  SW   Strategy Win: %v%v"
			if cfg.PlayAll.FindAllWinStrats {
				cmt += "  Will Continue to look for additional winning strategies for this deck"
			} else {
				cmt += "  Go to Next Deck (if any)"
			}
			prntMDet(bIn, aMoves, i, deckNum, moveNum, "DbDorMbM", 2, cmt, "", "", cfg, vPA)
			*/

			prntMDetTreeReturnComment(" ==> DECK WON", deckNum, recurReturnNum, cfg, vPA)
			vPA.TDotherSQL.moveNumAtWin = moveNum
			return "SW", recurReturnNum + 1 //  Strategy Win
		}

		// OK, done with the various end-of-strategy conditions
		// let's print out the list of available moves and make the next available move
		// The board state was already printed above
		prntMDet(bIn, aMoves, i, deckNum, moveNum, "MbM_R", 3, "", "", "", cfg, vPA) // Print Available Moves
		prntMDetTree(bIn, aMoves, i, deckNum, moveNum, cfg, vPA)

		bNew := bIn.copyBoard() // Critical Must use copyBoard!!!

		// ********** 1st of the 2 MOST IMPORTANT statements in this function:  ******************************
		bNew = moveMaker(bNew, aMoves[i])

		vPA.TD.mvsTried++

		prntMDet(bIn, aMoves, i, deckNum, moveNum, "MbM_S", 2, "\n      bIn: %v\n", "", "", cfg, vPA)
		prntMDet(bNew, aMoves, i, deckNum, moveNum, "MbM_S", 2, "     bNew: %v", "", "", cfg, vPA)

		if cfg.PlayAll.ProgressCounter > 0 && math.Mod(float64(vPA.TD.mvsTried+vPA.AD.mvsTried), float64(cfg.PlayAll.ProgressCounter)) <= 0.1 {
			//avgRepTime := time.Since(vPA.ADother.startTime) / time.Duration((vPA.TD.mvsTried+vPA.AD.mvsTried)/cfg.PlayAll.ProgressCounter)
			//estMaxRemTimeTD := time.Since(vPA.ADother.startTime) * time.Duration((cfg.PlayAll.GameLengthLimit*1_000_000 - vPA.TD.mvsTried)) / time.Duration(vPA.TD.mvsTried+vPA.AD.mvsTried)
			decksPlayed := float64(deckNum-cfg.General.FirstDeckNum) + .5
			decksCompleted := float64(vPA.ADother.decksWon + vPA.ADother.decksLost + vPA.ADother.decksLostGLE)
			estRemTimeAD := time.Duration(float64(time.Since(vPA.ADother.startTime)) * (float64(cfg.General.NumberOfDecksToBePlayed) - decksPlayed) / decksPlayed)
			// NOTE THE FOLLOWING PRINT STATEMENT NEVER GOES TO FILE - ALWAYS TO CONSOLE
			_, _ = pfmt.Printf("\rDk: %d   Mvs: %vmm   Strats: %vmm   UnqBoards: %vmm  MaxMoveNum: %v  Elapsed: %7s  estRem: %7s  W/L/GLE: %v/%v/%v  W/L/GLE %%: %3.1f/%3.1f/%3.1f\r",
				deckNum, (vPA.TD.mvsTried+vPA.AD.mvsTried)/1000000, vPA.AD.stratNum/1000000, vPA.AD.unqBoards/1000000, vPA.ADother.moveNumMax, time.Since(vPA.ADother.startTime).Round(time.Minute).String(), estRemTimeAD.Round(time.Minute), vPA.ADother.decksWon, vPA.ADother.decksLost, vPA.ADother.decksLostGLE, float64(vPA.ADother.decksWon)/decksCompleted*100.0, float64(vPA.ADother.decksLost)/decksCompleted*100.0, float64(vPA.ADother.decksLostGLE)/decksCompleted*100.0)
			//_, _ = pfmt.Printf("%v %v %v\r", vPA.ADother.decksWon, decksCompleted, float64(vPA.ADother.decksWon)/decksCompleted*100.0)
			vPA.TDother.lastPrintTime = time.Now()
		}

		// ********** 2nd of the 2 MOST IMPORTANT statements in this function:  ******************************
		recurReturnV1, recurReturnNum = playAllMoves(bNew, moveNum+1, deckNum, cfg, vPA)

		// CONSIDER DELETING prntMDet(bIn, aMoves, i, deckNum, moveNum, "DbDorMbM", 1, "  Returned: %v - %v After Call at deckNum: %v  moveNum: %v   vPA.TD.stratNum: %v   vPA.TD.mvsTried: %v   UnqBds: %v   ElTimTD: %v   ElTimADs: %v\n", recurReturnV1, recurReturnV2, cfg, vPA)

		if recurReturnV1 == "SW" {
			// save winning moves into a slice in reverse
			vPA.TDother.winningMoves = append(vPA.TDother.winningMoves, aMoves[i])
			return recurReturnV1, recurReturnNum + 1 // return up the call stack to end strategies search  if findAllWinStrats false, and we had a win
		}
		if recurReturnV1 == "GLE" {
			return recurReturnV1, recurReturnNum + 1 //
		}
	}

	if moveNum != 0 { // The initial call to playAllMoves from playALL (when MoveNum = 0) is only there to display the initial board
		if len(aMoves) == 1 { // and does not represent a move tried - Therefore returning from the initial call is not either a Maj or MinSE
			vPA.TD.stratLossesMinSE++
		} else {
			vPA.TD.stratLossesMajSE++
		}
	}
	cmt := "  MajSE   Major Strategy loss: all possible moves Exhausted%v%v"
	if recurReturnNum == 0 {
		prntMDet(bIn, aMoves, len(aMoves), deckNum, moveNum, "DbDorMbM", 2, cmt, "", "", cfg, vPA)
	}
	prntMDetTreeReturnComment("^Maj SE^", deckNum, recurReturnNum, cfg, vPA)
	return "SE", recurReturnNum + 1 //  Strategy Exhausted
}

func prntMDetTestRange(deckNum int, cfg *Configuration, vPA *variablesSpecificToPlayAll) bool {
	deckRangeOK := false
	if cfg.PlayAll.RestrictReportingTo.DeckStartVal == 0 && cfg.PlayAll.RestrictReportingTo.DeckContinueFor == 0 {
		deckRangeOK = true
	} else {
		if deckNum >= cfg.PlayAll.RestrictReportingTo.DeckStartVal && (cfg.PlayAll.RestrictReportingTo.DeckContinueFor == 0 || deckNum < cfg.PlayAll.RestrictReportingTo.DeckStartVal+cfg.PlayAll.RestrictReportingTo.DeckContinueFor) {
			deckRangeOK = true
		}
	}
	aMvsThisDkRangeOK := false
	if cfg.PlayAll.RestrictReportingTo.MovesTriedStartVal == 0 && cfg.PlayAll.RestrictReportingTo.MovesTriedContinueFor == 0 {
		aMvsThisDkRangeOK = true
	} else {
		if vPA.TD.mvsTried >= cfg.PlayAll.RestrictReportingTo.MovesTriedStartVal && (cfg.PlayAll.RestrictReportingTo.MovesTriedContinueFor == 0 || vPA.TD.mvsTried < cfg.PlayAll.RestrictReportingTo.MovesTriedStartVal+cfg.PlayAll.RestrictReportingTo.MovesTriedContinueFor) {
			aMvsThisDkRangeOK = true
		}
	}
	if deckRangeOK && aMvsThisDkRangeOK {
		return true
	} else {
		return false
	}
}

func prntMDetTestDeckRange(deckNum int, cfg *Configuration) bool {
	deckRangeOK := false
	if cfg.PlayAll.RestrictReportingTo.DeckStartVal == 0 && cfg.PlayAll.RestrictReportingTo.DeckContinueFor == 0 {
		deckRangeOK = true
	} else {
		if deckNum >= cfg.PlayAll.RestrictReportingTo.DeckStartVal && (cfg.PlayAll.RestrictReportingTo.DeckContinueFor == 0 || deckNum < cfg.PlayAll.RestrictReportingTo.DeckStartVal+cfg.PlayAll.RestrictReportingTo.DeckContinueFor) {
			deckRangeOK = true
		}
	}
	return deckRangeOK
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
	cfg *Configuration,
	vPA *variablesSpecificToPlayAll) {
	// Done here just to clean up mainline logic of playAllMoves
	// Do some repetitive printing to track progress
	// This function will use the struct pMD
	//      variant will be used for different outputs under the same pType

	if !cfg.PlayAll.RestrictReporting || prntMDetTestRange(dN, cfg, vPA) {
		switch {
		case pTypeIn == "MbM_ANY" && cfg.PlayAll.ReportingType.MoveByMove && variant == 1: // for "MbM_R", "MbM_S", "MbM_VS"
			if cfg.PlayAll.MoveByMoveReportingOptions.Type == "regular" && !(mN == 0 && nextMove == 0) {
				_, _ = fmt.Fprintf(oW, "\n****************************************\n")
			}
			_, _ = pfmt.Fprintf(oW, "\nDeck: %v   Move: %3v   Strategy #: %v  Moves Tried: %v   Unique Boards: %v   Elapsed TD: %v",
				dN,
				mN,
				vPA.TD.stratNum,
				vPA.TD.mvsTried,
				len(vPA.TDother.priorBoards),
				time.Since(vPA.TDother.startTime),
			)
			if cfg.PlayAll.MoveByMoveReportingOptions.Type != "regular" {
				pM1, pM2 := printMove(aMoves[nextMove], true)
				_, _ = fmt.Fprintf(oW, "    Next Move: %s   %s", pM1, strings.TrimSpace(pM2))
			}
			if cfg.PlayAll.MoveByMoveReportingOptions.Type == "regular" || (mN == 0 && nextMove == 0) {
				_, _ = fmt.Fprintf(oW, "\n")
				printBoard(b)
			}
		case pTypeIn == "MbM_R" && cfg.PlayAll.ReportingType.MoveByMove && variant == 2 && cfg.PlayAll.MoveByMoveReportingOptions.Type == "regular": // for "MbM_R"
			_, _ = pfmt.Fprintf(oW, "\n   "+comment, s1, s2)
		case pTypeIn == "MbM_SorMBM_VS" && cfg.PlayAll.ReportingType.MoveByMove && variant == 2: // for "MbM_R", "MbM_S", "MbM_VS"
			// comment must have 2 %v in it
			if cfg.PlayAll.MoveByMoveReportingOptions.Type == "regular" {
				_, _ = fmt.Fprintf(oW, "\n")
			}
			_, _ = pfmt.Fprintf(oW, "   "+comment, s1, s2)

		case pTypeIn == "MbM_SorMBM_VS" && cfg.PlayAll.ReportingType.MoveByMove && cfg.PlayAll.MoveByMoveReportingOptions.Type != "regular" && variant == 2: // for "MbM_R"
			// comment must have 2 %v in it
			_, _ = pfmt.Fprintf(oW, comment, s1, s2)
		case pTypeIn == "MbM_R" && cfg.PlayAll.ReportingType.MoveByMove && cfg.PlayAll.MoveByMoveReportingOptions.Type == "regular" && variant == 3:
			_, _ = fmt.Fprintf(oW, "\n All Possible Moves: \n")
			for j := range aMoves {
				pM1, pM2 := printMove(aMoves[j], false)
				if nextMove == j {
					_, _ = fmt.Fprintf(oW, "     Next Move ->  ")
				} else {
					_, _ = fmt.Fprintf(oW, "                   ")
				}
				_, _ = fmt.Fprintf(oW, "%s", pM1)
				if nextMove == j {
					pad := strings.Repeat(" ", 110-printableLength(pM1))
					_, _ = fmt.Fprintf(oW, pad+"<- Next Move")
				}
				if len(pM2) > 0 {
					_, _ = fmt.Fprintf(oW, "\n                           %s", pM2)
				}

				_, _ = fmt.Fprintf(oW, "\n")
			}
		/*case cfg.PlayAll.ReportingType.MoveByMove && pTypeIn == "MbM_SorMBM_VS" && (cfg.PlayAll.MoveByMoveReportingOptions.Type == "short" || cfg.PlayAll.MoveByMoveReportingOptions.Type == "very short") && variant == 1: // for "MbM_S" or "MbM_VS"
		_, _ = pfmt.Fprintf( oW,comment, dN, mN, vPA.TD.stratNum, vPA.TD.mvsTried, len(vPA.TDother.priorBoards), time.Since(vPA.ADother.startTime), time.Since(vPA.TDother.startTime))
		_, _ = pfmt.Fprintf( oW,comment, dN, mN, vPA.TD.stratNum, vPA.TD.mvsTried, len(vPA.TDother.priorBoards), time.Since(vPA.ADother.startTime), time.Since(vPA.TDother.startTime))
		*/
		case cfg.PlayAll.ReportingType.MoveByMove && pTypeIn == "MbM_S" && cfg.PlayAll.MoveByMoveReportingOptions.Type == "short" && variant == 2:
			_, _ = pfmt.Fprintf(oW, comment, b)
		//case pTypeIn == "DbDorMbM" && (cfg.PlayAll.ReportingType.MoveByMove || cfg.PlayAll.ReportingType.DeckByDeck) && variant == 1: //   formerly "NOTX"
		//	_, _ = pfmt.Fprintf( oW,comment, s1, s2, dN, mN, vPA.TD.stratNum, vPA.TD.mvsTried, len(vPA.TDother.priorBoards), time.Since(vPA.ADother.startTime), time.Since(vPA.TDother.startTime))
		case pTypeIn == "DbDorMbM" && (cfg.PlayAll.ReportingType.MoveByMove || cfg.PlayAll.ReportingType.DeckByDeck) && variant == 2: //   formerly "NOTX"
			_, _ = pfmt.Fprintf(oW, comment, s1, s2)
		}
	}
	// End of Deck Reporting Restrictions only on Deck Number
	if pTypeIn == "ANY" && prntMDetTestDeckRange(dN, cfg) && !cfg.PlayAll.ReportingType.NoReporting {
		if cfg.PlayAll.ReportingType.DeckByDeck && cfg.PlayAll.DeckByDeckReportingOptions.Type == "regular" ||
			cfg.PlayAll.ReportingType.MoveByMove && cfg.PlayAll.MoveByMoveReportingOptions.Type == "regular" ||
			cfg.PlayAll.ReportingType.Tree && cfg.PlayAll.TreeReportingOptions.Type == "regular" {
			_, _ = pfmt.Fprintf(oW, " \n\nDECK: %v%s    Result Codes: %v", dN, s1, s2) // !!! Initial Blank Required in case ProgressCounter last printed
			statisticsPrint(&vPA.TD, "")
			TDotherSQLPrint(&vPA.TDotherSQL)
		} else {
			statisticsPrintOneLine(vPA, dN, s1, cfg.General.FirstDeckNum, cfg.General.NumberOfDecksToBePlayed)
		}

	}

}

func prntMDetTree(b board, aMoves []move, nextMove int, dN int, mN int, cfg *Configuration, vPA *variablesSpecificToPlayAll) {

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

	if cfg.PlayAll.ReportingType.Tree && (!cfg.PlayAll.RestrictReporting || prntMDetTestRange(dN, cfg, vPA)) {
		if mN == 0 && nextMove == 0 {
			_, _ = pfmt.Fprintf(oW, "\n\n Deck: %v\n", dN)
			printBoard(b)
			if cfg.PlayAll.TreeReportingOptions.Type == "regular" {
				_, _ = fmt.Fprintf(oW, "\n       Move # ==>")
				for i := 1; i <= 150; i++ {
					if i == 1 {
						_, _ = fmt.Fprintf(oW, "%4s", strconv.Itoa(i)+"  ")
					} else {
						_, _ = fmt.Fprintf(oW, "%8s", strconv.Itoa(i)+"  ")
					}
				}
			}
			_, _ = fmt.Fprintf(oW, "\n   Strategy # ")
			_, _ = fmt.Fprintf(oW, "\n            0  ")
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

		switch cfg.PlayAll.TreeReportingOptions.Type {
		case "very narrow":
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
			if cfg.General.OutputTo == "console" { // No need to sleep if not printing to console where
				time.Sleep(cfg.PlayAll.TreeReportingOptions.TreeSleepBetwnMovesDur) // someone would be watching in real time
			}
			_, _ = fmt.Fprintf(oW, "%s", treeThisMove)
		} else {
			if cfg.General.OutputTo == "console" { // No need to sleep if not printing to console where
				time.Sleep(cfg.PlayAll.TreeReportingOptions.TreeSleepBetwnStrategiesDur) // someone would be watching in real time
			}
			x := []rune(vPA.TDother.treePrevMoves)
			x = x[0 : mN*treeMoveWidth]
			vPA.TDother.treePrevMoves = string(x)
			_, _ = pfmt.Fprintf(oW, "\n%13s  %s%s", strconv.Itoa(vPA.TD.stratNum), vPA.TDother.treePrevMoves, treeThisMove)
		}
		vPA.TDother.treePrevMoves += treeAddToPrev
	}
}

func prntMDetTreeReturnComment(c string, dN int, recurReturnNum int, cfg *Configuration, vPA *variablesSpecificToPlayAll) {
	if cfg.PlayAll.ReportingType.Tree && recurReturnNum == 0 && prntMDetTestRange(dN, cfg, vPA) {
		_, _ = fmt.Fprintf(oW, c)
	}
}
func printableLength(s string) int {
	withinEscapeSequence := false
	runeCount := 0
	for _, r := range s {
		if withinEscapeSequence {
			if r == 'm' {
				withinEscapeSequence = false
			}
		} else {
			if r == '\x1b' {
				withinEscapeSequence = true
			} else {
				runeCount++
			}
		}
	}
	return runeCount
}
