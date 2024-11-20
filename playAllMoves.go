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
) (string, string, int) {
	// if cfg.PlayAll.PrintMoveDetails {}
	/* Return Codes: SL  = Strategy Lost	 NMA = No Moves Available
	                 						 RB  = Repetitive Board
	                                         SE  = Strategy Exhausted
	                                         GLE = GameLength Limit exceeded
	                 SW  = Strategy Win      EW  = Early Win
											 SW  = Standard Win  Obsolete all wins are early
	*/
	// add code for findAllWinStrats

	var recurReturnNum int

	var aMoves []move //available Moves
	var recurReturnV1, recurReturnV2 string
	if moveNum > moveNumMax {
		moveNumMax = moveNum
	}

	// Check to see if the gameLenthLimit has been exceeded.
	//If, treats this as a loss and returns with loss codes.
	if vPA.TD.mvsTried >= cfg.PlayAll.GameLengthLimit*1_000_000 {
		prntMDet(bIn,
			aMoves,
			0,
			deckNum,
			moveNum,
			"MbM_R",
			2,
			"\n  SL-GLE: Game Length of: %v exceeds limit: %v\n",
			strconv.Itoa(vPA.TD.mvsTried),
			strconv.Itoa(cfg.PlayAll.GameLengthLimit*1_000_000),
			cfg,
			vPA)
		vPA.TD.stratLossesGLE++
		prntMDetTreeReturnComment(" ==> GLE", deckNum, recurReturnNum, cfg, vPA)
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

		/* Actually before we actually try all moves let's first: print (optionally based on cfg.PlayAll.RestrictReportingTo.pType) the incoming board
		      and check the incoming board for various end-of-strategy conditions
		   Note: This was done this way, so as to ensure that when returns backed up the moveNum, the board would reprint.
		*/

		// Possible increment to vPA.TD.stratNum
		if i != 0 {
			// Started at 0 in playNew for each deck.  Increment each time a second through nth move is made
			//     That way strategy 0 is the "All Best Moves" strategy.  It is also why in playNew aMStratTriedAllDecks
			//     is incremented by vPA.TD.stratNum + 1 after each deck.
			vPA.TD.stratNum++
		}

		// Print the incoming board EVEN IF we are returning to it to try the next available move
		//       This had to be done after possible increment to vPA.TD.stratNum so that each time a board is reprinted it shows the NEW strategy number
		//       Before when it was above the possible increment the board was printing out with the stratNum of the last failed strategy
		prntMDet(bIn, aMoves, i, deckNum, moveNum, "MbM_R", 1, "", "", "", cfg, vPA)

		if i == 0 {
			// Check for repetitive board
			// We only have to check if i == 0 because if i >= 0 then we are just returning to this board at which point
			//   the board will already have been checked (when i == 0 {past tense})
			bNewBcode := bIn.boardCode(deckNum) //  Just use the array itself
			// Have we seen this board before?
			if vPA.priorBoards[bNewBcode] {
				// OK we did see it before so return to try next available move (if any) in aMoves[] aka strategy
				vPA.TD.stratLossesRB++
				//  	prntMDet(b board, aMoves []move, nextMove int, dN int, mN int, pTypeIn string, variant int, comment string, s1 string, s2 string) {
				prntMDet(bIn, aMoves, i, deckNum, moveNum, "MbM_R", 2, "\n  SF-RB: Repetitive Board - \"Next Move\" yielded a repeat of a board.\n", "", "", cfg, vPA)
				prntMDet(bIn, aMoves, i, deckNum, moveNum, "MbM_VS", 2, "\n  SF-RB: Repetitive Board - \"Next Move\" yielded a repeat of a board.\n", "", "", cfg, vPA)
				prntMDetTreeReturnComment(" ==> RB", deckNum, recurReturnNum, cfg, vPA)
				return "SL", "RB", 1 // Repetitive Board
			} else {
				// Remember the board state by putting it into the map "vPA.priorBoards"
				vPA.priorBoards[bNewBcode] = true
			}
		}

		// Check if No Moves Available
		if i == 0 && aMoves[0].name == "No Moves Available" {
			vPA.TD.stratLossesNMA++
			prntMDet(bIn, aMoves, i, deckNum, moveNum, "MbM_R", 2, "  SL-NMA: No Moves Available: Strategy Lost %v%v\n", "", "", cfg, vPA)
			prntMDetTreeReturnComment(" ==> NMA", deckNum, recurReturnNum, cfg, vPA)
			return "SL", "NMA", 1
		}

		//Detect Early Win
		if detectWinEarly(bIn) {
			vPA.TD.stratWins++
			cmt := "  SW-EW: Strategy Win: Early Win%v%v"
			if cfg.PlayAll.FindAllWinStrats {
				cmt += "  Will Continue to look for additional winning strategies for this deck"
			} else {
				cmt += "  Go to Next Deck (if any)"
			}
			prntMDet(bIn, aMoves, i, deckNum, moveNum, "DbDorMbM", 2, cmt, "", "", cfg, vPA)

			// Verbose Special "WL" Starts Here - No effect on operation
			if cfg.PlayAll.WinLossReport { // Deck Win Loss Summary Statistics   MOVE THIS!!!!!!!!!
				/*if len(deckWinLossDetail)-1 < deckNum {
					dWLDStats.winLoss = "W"
					dWLDStats.moveNumAt1stWinOrAtLoss = moveNum
					dWLDStats.moveNumMinWinIfFindAll = moveNum
					dWLDStats.moveNumMaxWinIfFindAll = moveNum
					dWLDStats.stratNumAt1stWinOrAtLoss = vPA.TD.stratNum
					dWLDStats.mvsTriedAt1stWinOrAtLoss = vPA.TD.mvsTried
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
			prntMDetTreeReturnComment(" ==> DECK WON", deckNum, recurReturnNum, cfg, vPA)
			return "SW", "EW", 1 //  Strategy Early Win
		}

		// OK, done with the various end-of-strategy conditions
		// let's print out the list of available moves and make the next available move
		// The board state was already printed above
		//if cfg.PlayAll.RestrictReporting && cfg.PlayAll.ReportingMoveByMove && cfg.PlayAll.MoveByMoveReportingOptions.Type == "regular" && prntMDetTestRange(deckNum, &cfg) {   // DELETE???
		prntMDet(bIn, aMoves, i, deckNum, moveNum, "MbM_R", 3, "", "", "", cfg, vPA) // DELETE???
		//}    // DELETE???

		bNew := bIn.copyBoard() // Critical Must use copyBoard

		// ********** 1st of the 2 MOST IMPORTANT statements in this function:  ******************************
		bNew = moveMaker(bNew, aMoves[i])

		prntMDet(bIn, aMoves, i, deckNum, moveNum, "MbM_SorMBM_VS", 1, "\nBefore Call at Deck: %v   Move: %v   Strategy #: %v  Moves Tried: %v   Unique Boards: %v   Elapsed TD: %v   Elapsed ADs: %v", "", "", cfg, vPA)
		prntMDet(bIn, aMoves, i, deckNum, moveNum, "MbM_S", 2, "\n      bIn: %v\n", "", "", cfg, vPA)
		prntMDet(bNew, aMoves, i, deckNum, moveNum, "MbM_S", 2, "     bNew: %v\n", "", "", cfg, vPA)
		prntMDetTree(bIn, aMoves, i, deckNum, moveNum, cfg, vPA)

		vPA.TD.mvsTried++

		if cfg.General.ProgressCounter > 0 && math.Mod(float64(vPA.TD.mvsTried+vPA.AD.mvsTried), float64(cfg.General.ProgressCounter)) <= 0.1 {
			avgRepTime := time.Since(vPA.TD.startTime) / time.Duration(vPA.TD.mvsTried+vPA.AD.mvsTried/cfg.General.ProgressCounter)
			estMaxRemTimeTD := time.Since(vPA.TD.startTime) * time.Duration((cfg.PlayAll.GameLengthLimit*1_000_000-vPA.TD.mvsTried+vPA.AD.mvsTried)/vPA.TD.mvsTried+vPA.AD.mvsTried)
			_, err = pfmt.Printf("\rDk: %5d   ____   MvsTried: %13v   MoveNum: %3v   Max MoveNum: %3v   StratsTried: %12v   UnqBoards: %11v  - %7s  SinceLast: %6s  Avg: %6s  estMaxRem: %7s\r", deckNum, vPA.TD.mvsTried+vPA.AD.mvsTried, moveNum, moveNumMax, vPA.TD.stratNum, len(vPA.priorBoards), time.Since(vPA.TD.startTime).Truncate(100*time.Millisecond).String(), time.Since(vPA.TD.startTime).Truncate(100*time.Millisecond).String(), avgRepTime.Truncate(100*time.Millisecond).String(), estMaxRemTimeTD.Truncate(1*time.Minute))
			//lastPrintTime := time.Now()
		}

		// ********** 2nd of the 2 MOST IMPORTANT statements in this function:  ******************************
		recurReturnV1, recurReturnV2, recurReturnNum = playAllMoves(bNew, moveNum+1, deckNum, cfg, vPA)

		// CONSIDER DELETING prntMDet(bIn, aMoves, i, deckNum, moveNum, "DbDorMbM", 1, "  Returned: %v - %v After Call at deckNum: %v  moveNum: %v   vPA.TD.stratNum: %v   vPA.TD.mvsTried: %v   UnqBds: %v   ElTimTD: %v   ElTimADs: %v\n", recurReturnV1, recurReturnV2, cfg, vPA)

		if cfg.PlayAll.FindAllWinStrats != true && recurReturnV1 == "SW" {
			// save winning moves into a slice in reverse
			return recurReturnV1, recurReturnV2, recurReturnNum + 1 // return up the call stack to end strategies search  if findAllWinStrats false, and we had a win
		}
		if recurReturnV1 == "SL" && recurReturnV2 == "GLE" {
			return recurReturnV1, recurReturnV2, recurReturnNum + 1 //
		}
	}

	vPA.TD.stratLossesSE++
	prntMDetTreeReturnComment(" ==> SE", deckNum, recurReturnNum, cfg, vPA)
	return "SL", "SE", recurReturnNum + 1 //  Strategy Exhausted
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
		case pTypeIn == "MbM_R" && cfg.PlayAll.ReportingType.MoveByMove && cfg.PlayAll.MoveByMoveReportingOptions.Type == "regular" && variant == 1: // for "MbM_R"
			fmt.Printf("\n****************************************\n")
			_, err = pfmt.Printf("\nDeck: %v   Move: %v   Strategy #: %v  Moves Tried: %v   Unique Boards: %v   Elapsed TD: %v   Elapsed ADs: %v\n",
				dN,
				mN,
				vPA.TD.stratNum,
				vPA.TD.mvsTried,
				len(vPA.priorBoards),
				time.Since(vPA.AD.startTime),
				time.Since(vPA.TD.startTime))
			printBoard(b)
		case pTypeIn == "MbM_R" && cfg.PlayAll.ReportingType.MoveByMove && cfg.PlayAll.MoveByMoveReportingOptions.Type == "regular" && variant == 2: // for "MbM_R"
			// comment must have 2 %v in it
			_, err = pfmt.Printf(comment, s1, s2)
		case pTypeIn == "MbM_R" && cfg.PlayAll.ReportingType.MoveByMove && cfg.PlayAll.MoveByMoveReportingOptions.Type == "regular" && variant == 3:
			fmt.Printf("\n All Possible Moves: \n")
			for j := range aMoves {
				pM1, pM2 := printMove(aMoves[j])
				if nextMove == j {
					fmt.Printf("     Next Move ->  ")
				} else {
					fmt.Printf("                   ")
				}
				fmt.Printf("%s", pM1)
				if nextMove == j {
					pad := strings.Repeat(" ", 110-printableLength(pM1))
					fmt.Printf(pad + "<- Next Move")
				}
				if len(pM2) > 0 {
					fmt.Printf("%s", pM2)
				}

				fmt.Printf("\n")
			}
		case cfg.PlayAll.ReportingType.MoveByMove && pTypeIn == "MbM_SorMBM_VS" && (cfg.PlayAll.MoveByMoveReportingOptions.Type == "short" || cfg.PlayAll.MoveByMoveReportingOptions.Type == "very short") && variant == 1: // for "MbM_S" or "MbM_VS"
			_, err = pfmt.Printf(comment, dN, mN, vPA.TD.stratNum, vPA.TD.mvsTried, len(vPA.priorBoards), time.Since(vPA.AD.startTime), time.Since(vPA.TD.startTime))
		case cfg.PlayAll.ReportingType.MoveByMove && pTypeIn == "MbM_S" && cfg.PlayAll.MoveByMoveReportingOptions.Type == "short" && variant == 2:
			_, err = pfmt.Printf(comment, b)
			//	case pTypeIn == "DbDorMbM" && cfg.PlayAll.ReportingMoveByMove && variant == 1://   formerly "NOTX"
		case pTypeIn == "DbDorMbM" && (cfg.PlayAll.ReportingType.MoveByMove || cfg.PlayAll.ReportingType.DeckByDeck) && variant == 1: //   formerly "NOTX"
			_, err = pfmt.Printf(comment, s1, s2, dN, mN, vPA.TD.stratNum, vPA.TD.mvsTried, len(vPA.priorBoards), time.Since(vPA.AD.startTime), time.Since(vPA.TD.startTime))
			//	case pTypeIn == "DbDorMbM" && cfg.PlayAll.ReportingMoveByMove && variant == 2://   formerly "NOTX"
		case pTypeIn == "DbDorMbM" && (cfg.PlayAll.ReportingType.MoveByMove || cfg.PlayAll.ReportingType.DeckByDeck) && variant == 2: //   formerly "NOTX"
			_, err = pfmt.Printf(comment, s1, s2)
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
			_, err = pfmt.Printf("\n\n Deck: %v\n\n", dN)
			printBoard(b)
			if cfg.PlayAll.TreeReportingOptions.Type == "regular" {
				fmt.Printf("\n       Move # ==>")
				for i := 1; i <= 150; i++ {
					if i == 1 {
						fmt.Printf("%4s", strconv.Itoa(i)+"  ")
					} else {
						fmt.Printf("%8s", strconv.Itoa(i)+"  ")
					}
				}
			}
			fmt.Printf("\n   Strategy # ")
			fmt.Printf("\n            0  ")
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
			time.Sleep(cfg.PlayAll.TreeReportingOptions.TreeSleepBetwnMovesDur)
			fmt.Printf("%s", treeThisMove)
		} else {
			time.Sleep(cfg.PlayAll.TreeReportingOptions.TreeSleepBetwnStrategiesDur)
			x := []rune(vPA.TD.treePrevMoves)
			x = x[0 : mN*treeMoveWidth]
			vPA.TD.treePrevMoves = string(x)
			_, err = pfmt.Printf("\n%13s  %s%s", strconv.Itoa(vPA.TD.stratNum), vPA.TD.treePrevMoves, treeThisMove)
		}
		vPA.TD.treePrevMoves += treeAddToPrev
	}
}

func prntMDetTreeReturnComment(c string, dN int, recurReturnNum int, cfg *Configuration, vPA *variablesSpecificToPlayAll) {
	if cfg.PlayAll.ReportingType.Tree && recurReturnNum == 0 && prntMDetTestRange(dN, cfg, vPA) {
		fmt.Printf(c)
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
