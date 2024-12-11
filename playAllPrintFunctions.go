package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

func statisticsPrint(TDorAD *Statistics, which string) {
	if which != "AD" {
		if TDorAD.elapsedTime == 0 {
			_, _ = fmt.Fprintf(oW, "\nElapsed Time: <.5ms    (Windows Minimum Resolution)")
		} else {
			_, _ = fmt.Fprintf(oW, "\nElapsed Time: %v", TDorAD.elapsedTime)
		}
	}
	_, _ = fmt.Fprintf(oW, "\n\nStrategies:")
	_, _ = pfmt.Fprintf(oW, "\n   Tried: %-15d", TDorAD.stratTried)
	_, _ = fmt.Fprintf(oW, "   Tried Detail:                   (Must sum to Strategies Tried)")
	if TDorAD.stratWins != 0 {
		_, _ = pfmt.Fprintf(oW, "\n                                     Won: %15d", TDorAD.stratWins)
	}
	if TDorAD.stratLossesNMA != 0 {
		_, _ = pfmt.Fprintf(oW, "\n                                     NMA: %15d   (No Moves Available)", TDorAD.stratLossesNMA)
	}
	if TDorAD.stratLossesRB != 0 {
		_, _ = pfmt.Fprintf(oW, "\n                                      RB: %15d   (Repetitive Board)", TDorAD.stratLossesRB)
	}
	if TDorAD.stratLossesEL != 0 {
		_, _ = pfmt.Fprintf(oW, "\n                                      EL: %15d   (Early Loss)", TDorAD.stratLossesEL)
	}
	if TDorAD.stratLossesGLE != 0 {
		_, _ = pfmt.Fprintf(oW, "\n                                     GLE: %15d   (Game Length Exceeded)", TDorAD.stratLossesGLE)
	}
	if TDorAD.stratTried != TDorAD.stratWins+TDorAD.stratLossesNMA+TDorAD.stratLossesRB+TDorAD.stratLossesEL+TDorAD.stratLossesGLE {
		_, _ = fmt.Fprintf(oW, "\n        ************* Strategies Tried != TDorAD.stratWins+TDorAD.stratLossesNMA+TDorAD.stratLossesRB+TDorAD.stratLossesEL+TDorAD.stratLossesGLE")
	}
	_, _ = fmt.Fprintf(oW, "\n\nMoves:")
	_, _ = pfmt.Fprintf(oW, "\n   Tried: %-15d", TDorAD.mvsTried)
	_, _ = fmt.Fprintf(oW, "   Tried Detail:                   (Must Sum to Moves Tried)")
	if TDorAD.stratLossesNMA != 0 {
		_, _ = pfmt.Fprintf(oW, "\n                                     NMA: %15d   (No Moves Available)", TDorAD.stratLossesNMA)
	}
	if TDorAD.stratLossesRB != 0 {
		_, _ = pfmt.Fprintf(oW, "\n                                      RB: %15d   (Repetitive Board)", TDorAD.stratLossesRB)
	}
	if TDorAD.stratLossesMajSE != 0 {
		_, _ = pfmt.Fprintf(oW, "\n                                   MajSE: %15d   (Major Exhausted)", TDorAD.stratLossesMajSE)
	}
	if TDorAD.stratLossesMinSE != 0 {
		_, _ = pfmt.Fprintf(oW, "\n                                   MinSE: %15d   (Minor Exhausted)", TDorAD.stratLossesMinSE)
	}
	if TDorAD.stratLossesEL != 0 {
		_, _ = pfmt.Fprintf(oW, "\n                                      EL: %15d   (Early Loss)", TDorAD.stratLossesEL)
	}
	if TDorAD.stratLossesGLE != 0 {
		_, _ = pfmt.Fprintf(oW, "\n                                     GLE: %15d   (Game Length Exceeded)", TDorAD.stratLossesGLE)
	}
	if TDorAD.stratLossesGLEAb != 0 {
		_, _ = pfmt.Fprintf(oW, "\n                                   GLEAb: %15d   (Game Length Exceeded Aborted moves)", TDorAD.stratLossesGLEAb)
	}
	if TDorAD.winningMovesCnt != 0 {
		_, _ = pfmt.Fprintf(oW, "\n                                   WMCnt: %15d   (Winning Moves Count)", TDorAD.winningMovesCnt)
	}
	x := TDorAD.stratLossesNMA + TDorAD.stratLossesRB + TDorAD.stratLossesMajSE + TDorAD.stratLossesMinSE + TDorAD.stratLossesEL + TDorAD.stratLossesGLE + TDorAD.stratLossesGLEAb + TDorAD.winningMovesCnt
	if TDorAD.mvsTried != x {
		_, _ = fmt.Fprintf(oW, "\n        ************* Moves Tried %v != TDorAD.stratLossesNMA+TDorAD.stratLossesRB+TDorAD.stratLossesMajSE+TDorAD.stratLossesMinSE+TDorAD.stratLossesEL+TDorAD.stratLossesGLE+TDorAD.winningMovesCnt", x)
	}
	_, _ = fmt.Fprintf(oW, "\n")
	if TDorAD.unqBoards != 0 {
		_, _ = pfmt.Fprintf(oW, "\n  UnqBds: %15d   (Unique Boards)", TDorAD.unqBoards)
	}
}

func statisticsPrintOneLine(vPA *variablesSpecificToPlayAll, dN int, s1 string, firstDeckNum int, numberOfDecksToBePlayed int) {
	var est time.Duration
	//                      nanosecondsTD   / Decks Played So Far         * remaining decks [remaining decks = numbertobeplayed - decksplayed so far
	est = time.Duration(float64(time.Since(vPA.ADother.startTime))/float64(dN+1-firstDeckNum)*float64(numberOfDecksToBePlayed-(dN+1-firstDeckNum))) * time.Nanosecond
	elTimeSinceStartTimeADFormatted := time.Since(vPA.ADother.startTime).Round(100 * time.Millisecond).String()
	if time.Since(vPA.ADother.startTime) > time.Duration(5*time.Minute) {
		elTimeSinceStartTimeADFormatted = time.Since(vPA.ADother.startTime).Round(time.Second).String()
	}
	_, _ = pfmt.Fprintf(oW, "Deck: %7v%9s   Strategy #: %11v   Moves Tried: %11v   Unique Boards: %10v   Elapsed TD: %10v"+
		"   stratTried: %11v   stratWins: %1v   stratLossesNMA: %9v   stratLossesRB: %10v   stratLossesEL: %1v   stratLossesGLE: %5v"+
		"   stratLossesMajSE: %11v   stratLossesMinSE: %11v   stratLossesEL: %7v   stratLossesGLEAb: %3v"+
		"   winningMovesCnt: %3v   moveNumMax: %3v   moveNumAtWin:%3v   Elapsed AD: %10s   Est Rem: %10s   Now: %8s\n",
		dN,
		s1,
		vPA.TD.stratNum,
		vPA.TD.mvsTried,
		vPA.TD.unqBoards,
		vPA.TD.elapsedTime.Round(10*time.Millisecond).String(),
		vPA.TD.stratTried, vPA.TD.stratWins, vPA.TD.stratLossesNMA, vPA.TD.stratLossesRB, vPA.TD.stratLossesEL, vPA.TD.stratLossesGLE,
		vPA.TD.stratLossesMajSE, vPA.TD.stratLossesMinSE, vPA.TD.stratLossesEL, vPA.TD.stratLossesGLEAb,
		vPA.TD.winningMovesCnt, vPA.TDotherSQL.moveNumMax, vPA.TDotherSQL.moveNumAtWin,
		elTimeSinceStartTimeADFormatted, est.Round(100*time.Millisecond).String(), time.Now().Format(" 3:04 pm"))
}

func TDotherSQLPrint(x *TDotherSQL) {
	if x.moveNumMax != 0 {
		_, _ = pfmt.Fprintf(oW, "\n   MNMax: %15d   (Move Number Max)", x.moveNumMax)
	}
	if x.moveNumAtWin != 0 {
		_, _ = pfmt.Fprintf(oW, "\n   MNWin: %15d   (Move Number at Win)", x.moveNumAtWin)
	}
}

func ADotherSQLPrint(vPA *variablesSpecificToPlayAll) {
	if vPA.ADother.moveNumMax != 0 {
		_, _ = pfmt.Fprintf(oW, "\n   MNMax: %15d   (Move Number Max)", vPA.ADother.moveNumMax)
	}
	if vPA.ADother.moveNumAtWinMin != 0 {
		_, _ = pfmt.Fprintf(oW, "\nMinMNWin: %15d   (Minimum Move Number at Win)", vPA.ADother.moveNumAtWinMin)
	}
	if vPA.ADother.moveNumAtWinMax != 0 {
		_, _ = pfmt.Fprintf(oW, "\nMaxMNWin: %15d   (Maximum Move Number at Win)\n\n\n", vPA.ADother.moveNumAtWinMax)
	}
	_, _ = fmt.Fprintf(oW, "\n\n\n")
}

func printSummaryStats(cfg *Configuration, vPA *variablesSpecificToPlayAll) {
	_, _ = fmt.Fprintf(oW, "\n\n******************   Summary Statistics   ******************\n")
	averageElapsedTimePerDeck := time.Duration(float64(time.Since(vPA.ADother.startTime)) / float64(cfg.General.NumberOfDecksToBePlayed))
	_, _ = fmt.Fprintf(oW, "\n     Elapsed Time: %5v", time.Since(vPA.ADother.startTime).Round(6*time.Second).String())
	_, _ = fmt.Fprintf(oW, "\nAvg Time per Deck: %5v\n", averageElapsedTimePerDeck.Round(100*time.Millisecond).String())
	_, _ = pfmt.Fprintf(oW, "\n          Decks Played: %-7d", vPA.ADother.decksPlayed)
	_, _ = pfmt.Fprintf(oW, "\n             Decks Won: %-7d   %4v%%     Ignoring GLE: %4v%%", vPA.ADother.decksWon, roundFloatIntDiv(vPA.ADother.decksWon*100, vPA.ADother.decksPlayed, 1), roundFloatIntDiv(vPA.ADother.decksWon*100, vPA.ADother.decksPlayed-vPA.ADother.decksLostGLE, 1))
	_, _ = pfmt.Fprintf(oW, "\n            Decks Lost: %-7d   %4v%%     Ignoring GLE: %4v%%", vPA.ADother.decksLost, roundFloatIntDiv(vPA.ADother.decksLost*100, vPA.ADother.decksPlayed, 1), roundFloatIntDiv(vPA.ADother.decksLost*100, vPA.ADother.decksPlayed-vPA.ADother.decksLostGLE, 1))
	_, _ = pfmt.Fprintf(oW, "\n         Decks LostGLE: %-7d   %4v%%", vPA.ADother.decksLostGLE, roundFloatIntDiv(vPA.ADother.decksLostGLE*100, vPA.ADother.decksPlayed, 1))
	_, _ = pfmt.Fprintf(oW, "\n  Decks Lost + LostGLE: %-7d   %4v%%", vPA.ADother.decksLostGLE, roundFloatIntDiv((vPA.ADother.decksLost+vPA.ADother.decksLostGLE)*100, vPA.ADother.decksPlayed, 1))

	statisticsPrint(&vPA.AD, "AD")
	ADotherSQLPrint(vPA)
}

func PrintWinningMoves(cfg *Configuration, vPA *variablesSpecificToPlayAll) {
	// First Reverse the slice (which were collected in reverse as we backed up the call chain)
	for i := 0; i < len(vPA.TDother.winningMoves)/2; i++ {
		vPA.TDother.winningMoves[i], vPA.TDother.winningMoves[len(vPA.TDother.winningMoves)-i-1] = vPA.TDother.winningMoves[len(vPA.TDother.winningMoves)-i-1], vPA.TDother.winningMoves[i]
	}
	// Now print them
	if cfg.PlayAll.PrintWinningMoves {
		_, _ = fmt.Fprintf(oW, "\n\n     Winning Moves:\n")
		for mN := range vPA.TDother.winningMoves {
			m1, m2 := printMove(vPA.TDother.winningMoves[mN], true)
			_, _ = fmt.Fprintf(oW, "        %3v.  %s\n", mN+1, m1)
			if len(m2) != 0 {
				_, _ = fmt.Fprintf(oW, "          %s\n", m2)
			}
		}
	}
}

// Divide 2 integers and round to precision digits
func roundFloatIntDiv(numer int, denom int, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(float64(numer)/float64(denom)*ratio) / ratio
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
