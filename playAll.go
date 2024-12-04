package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"strconv"
	"time"
)

type Statistics struct {
	mvsTried         int
	stratNum         int // stratNum starts at 0             Both are not really needed but it was written that way and not worth changing
	stratTried       int // Strategies TRIED = stratNum + 1  Both are not really needed but it was written that way and not worth changing
	stratWins        int
	stratLosses      int
	stratLossesGLE   int // Strategy Game Length Exceeded
	stratLossesGLEAb int // Strategy Game Length Exceeded Aborted moves
	stratLossesNMA   int // Strategy No Moves Available
	stratLossesRB    int // Strategy Repetitive Board
	stratLossesMajSE int // Strategy Exhausted Major
	stratLossesMinSE int // Strategy Exhausted Minor
	stratLossesEL    int // Strategy Early Loss
	winningMovesCnt  int // = length of winning moves
	unqBoards        int
	elapsedTime      time.Duration
}

type TDotherSQL struct { // Variables in addition to TD that should be output to SQL
	moveNumMax   int
	moveNumAtWin int
}

type variablesSpecificToPlayAll struct {
	TD         Statistics
	TDotherSQL TDotherSQL
	TDother    struct { // Variables NOT common to TD and AD
		startTime     time.Time
		lastPrintTime time.Time
		treePrevMoves string // Used to retain values between calls to prntMDetTree for a single deck - Needed for when the strategy "Backs Uo"
		winningMoves  []move
		priorBoards   map[bCode]bool // NOTE: bcode is an array of 65 ints as defined in board.go
	}
	AD      Statistics
	ADother struct {
		startTime       time.Time
		moveNumMax      int
		moveNumAtWinMin int
		moveNumAtWinMax int
		decksPlayed     int
		decksWon        int
		decksLost       int
		decksLostGLE    int
	}
}

func playAll(reader csv.Reader, cfg *Configuration) {
	firstDeckNum := cfg.General.FirstDeckNum                       // Shorthand name but really is a copy - OK since never changed (but would Pointer or address be better?)
	numberOfDecksToBePlayed := cfg.General.NumberOfDecksToBePlayed // Shorthand name but really is a copy - OK since never changed (but would Pointer or address be better?)
	verbose := cfg.General.Verbose                                 // Shorthand name but really is a copy - OK since never changed (but would Pointer or address be better?)
	var vPA variablesSpecificToPlayAll
	vPA.TDother.priorBoards = map[bCode]bool{}
	vPA.TDother.treePrevMoves = ""
	vPA.TD.stratTried = 1
	vPA.TDother.startTime = time.Now()
	vPA.TDother.lastPrintTime = time.Now()
	vPA.ADother.startTime = time.Now()

	for deckNum := firstDeckNum; deckNum < (firstDeckNum + numberOfDecksToBePlayed); deckNum++ {
		if !cfg.PlayAll.ReportingType.NoReporting && !(cfg.PlayAll.ReportingType.DeckByDeck && cfg.PlayAll.DeckByDeckReportingOptions.Type != "regular") {
			fmt.Printf("\n\n******************************************************************************************************\n")
		}
		protoDeck, err := reader.Read() // protoDeck is a slice of strings: rank, suit, rank, suit, etc.
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Cannot read from inputFileName:", err)
		}

		if verbose > 1 {
			_, _ = pfmt.Printf("\nDeck #%d:\n", deckNum)
		}
		var d Deck

		for i := 0; i < 52; i++ {
			rank, _ := strconv.Atoi(protoDeck[i*2])
			suit, _ := strconv.Atoi(protoDeck[i*2+1])
			c := Card{
				Rank:   rank,
				Suit:   suit,
				FaceUp: false,
			}
			d = append(d, c)

		}
		//deal Deck onto board
		//temp		AllMvStratNum := 0
		var b = dealDeck(d)

		// This statement is executed once per deck and transfers program execution.
		// When this statement returns the deck has been played.
		result1, result2 := playAllMoves(b, 0, deckNum, cfg, &vPA)

		if vPA.TD.stratLossesGLE > 0 {
			vPA.TD.stratLossesGLEAb = result2 - 2
		}
		vPA.TD.elapsedTime = time.Since(vPA.TDother.startTime)
		vPA.TD.winningMovesCnt = len(vPA.TDother.winningMoves)
		var dummy []move

		vPA.TD.unqBoards = len(vPA.TDother.priorBoards)
		var s1 string
		s2 := "   DECK"
		s1 = strconv.Itoa(deckNum)
		if result1 == "SW" {
			vPA.ADother.decksWon += 1
			s2 += " WON "
		} else {
			if vPA.TD.stratLossesGLE > 0 {
				s2 += " GLE "
				vPA.ADother.decksLostGLE += 1
			} else {
				vPA.TD.stratLosses++
				vPA.ADother.decksLost += 1
				s2 += " LOST"
			}
		}
		prntMDet(b, dummy, 1, deckNum, 1, "MbM_ANY", 2, "\n   "+s2+"\n", "", "", cfg, &vPA)
		prntMDetTreeReturnComment("\n   "+s2+"\n", deckNum, 0, cfg, &vPA)

		// End of Deck Statistics Reporting
		prntMDet(b, dummy, 0, deckNum, 1, "ANY", 0, s1, s2, result1, cfg, &vPA)
		/*if !cfg.PlayAll.ReportingType.NoReporting && cfg.PlayAll.DeckByDeckReportingOptions.Type == "regular" { // End of deck Statistics
			_, _ = pfmt.Printf("\n\nDECK: %s%s    Result Codes: %v", s1, s2, result1)

			statisticsPrint(&vPA.TD, "")
			TDotherSQLPrint(&vPA.TDotherSQL)
		}*/
		// if cfg.PlayAll.SaveResultsToSQL == true OR cfg.PlayAll.PrintWinningMoves == true
		if (cfg.PlayAll.SaveResultsToSQL || cfg.PlayAll.PrintWinningMoves) && vPA.TD.winningMovesCnt != 0 {
			// First Reverse the slice (which was collected in reverse as we backed up the call chain)
			for i := 0; i < len(vPA.TDother.winningMoves)/2; i++ {
				vPA.TDother.winningMoves[i], vPA.TDother.winningMoves[len(vPA.TDother.winningMoves)-i-1] = vPA.TDother.winningMoves[len(vPA.TDother.winningMoves)-i-1], vPA.TDother.winningMoves[i]
			}
			if cfg.PlayAll.PrintWinningMoves {
				fmt.Printf("\n\n     Winning Moves:\n")
				for mN := range vPA.TDother.winningMoves {
					m1, m2 := printMove(vPA.TDother.winningMoves[mN], true)
					fmt.Printf("        %3v.  %s\n", mN+1, m1)
					if len(m2) != 0 {
						fmt.Printf("          %s\n", m2)
					}
				}
				//fmt.Printf("\n")
			}
			if cfg.PlayAll.SaveResultsToSQL {
				// write ConfigurationSubsetOnlyForSQLWriting and vPA.TD out to sql/csv here
			}
		}

		/*	// This If Block is Print Only for DbD_S or DbD_VS
			if cfg.PlayAll.ReportingType.DeckByDeck && cfg.PlayAll.DeckByDeckReportingOptions.Type != "regular" {
				var est time.Duration
				//                      nanosecondsTD   / Decks Played So Far         * remaining decks [remaining decks = numbertobeplayed - decksplayed so far
				est = time.Duration(float64(time.Since(vPA.ADother.startTime))/float64(deckNum+1-firstDeckNum)*float64(numberOfDecksToBePlayed-(deckNum+1-firstDeckNum))) * time.Nanosecond
				wL := ""
				if vPA.TD.stratWins > 0 {
					wL = "WON " // Note additional space -- for alignment
				} else {
					wL = "LOST"
				}
				elTimeSinceStartTimeADFormatted := time.Since(vPA.ADother.startTime).Truncate(100 * time.Millisecond).String()
				if time.Since(vPA.ADother.startTime) > time.Duration(5*time.Minute) {
					elTimeSinceStartTimeADFormatted = time.Since(vPA.ADother.startTime).Truncate(time.Second).String()
				}
				_, _ = pfmt.Printf("Dk: %5d   "+wL+"   MvsTried: %13v   MoveNum: %3v   Max MoveNum: %3v   StratsTried: %12v   UnqBoards: %11v   Won: %5v   Lost: %5v   GLE: %5v   Won: %5.1f%%   Lost: %5.1f%%   GLE: %5.1f%%   ElTime TD: %9s   ElTime ADs: %9s  Rem Time: %11s   ResCodes: %2s %3s   Time Now: %8s\n", deckNum, vPA.TD.mvsTried, vPA.TDotherSQL.moveNumAtWin, vPA.TDotherSQL.moveNumMax, vPA.TD.stratNum, vPA.TD.unqBoards, vPA.ADother.decksWon, vPA.ADother.decksLost, vPA.AD.stratLossesGLE, roundFloatIntDiv(vPA.ADother.decksWon*100, deckNum+1-firstDeckNum, 1), roundFloatIntDiv(vPA.ADother.decksLost*100, deckNum+1-firstDeckNum, 1), roundFloatIntDiv(vPA.AD.stratLossesGLE*100, deckNum+1-firstDeckNum, 1), vPA.TD.elapsedTime.Truncate(100*time.Millisecond).String(), elTimeSinceStartTimeADFormatted, est.Truncate(time.Second).String(), result1, "", time.Now().Format(" 3:04 pm"))
			}*/
		vPA.AD.mvsTried += vPA.TD.mvsTried
		vPA.TD.mvsTried = 0
		vPA.AD.stratNum += vPA.TD.stratNum
		vPA.TD.stratNum = 0
		vPA.AD.stratTried += vPA.TD.stratTried
		vPA.TD.stratTried = 1 // NOTE: Starts at 1 not 0
		vPA.AD.stratWins += vPA.TD.stratWins
		vPA.TD.stratWins = 0
		vPA.AD.stratLosses += vPA.TD.stratLosses
		vPA.TD.stratLosses = 0
		vPA.AD.stratLossesGLE += vPA.TD.stratLossesGLE
		vPA.TD.stratLossesGLE = 0
		vPA.AD.stratLossesGLEAb += vPA.TD.stratLossesGLEAb
		vPA.TD.stratLossesGLEAb = 0
		vPA.AD.stratLossesNMA += vPA.TD.stratLossesNMA
		vPA.TD.stratLossesNMA = 0
		vPA.AD.stratLossesRB += vPA.TD.stratLossesRB
		vPA.TD.stratLossesRB = 0
		vPA.AD.stratLossesMajSE += vPA.TD.stratLossesMajSE
		vPA.TD.stratLossesMajSE = 0
		vPA.AD.stratLossesMinSE += vPA.TD.stratLossesMinSE
		vPA.TD.stratLossesMinSE = 0
		vPA.AD.stratLossesEL += vPA.TD.stratLossesEL
		vPA.TD.stratLossesEL = 0
		vPA.AD.unqBoards += vPA.TD.unqBoards
		vPA.TD.unqBoards = 0
		vPA.AD.winningMovesCnt += vPA.TD.winningMovesCnt
		vPA.TD.winningMovesCnt = 0
		vPA.AD.elapsedTime += vPA.TD.elapsedTime

		if vPA.ADother.moveNumMax == 0 || vPA.ADother.moveNumMax < vPA.TDotherSQL.moveNumMax {
			vPA.ADother.moveNumMax = vPA.TDotherSQL.moveNumMax
		}
		if vPA.ADother.moveNumAtWinMin == 0 || vPA.ADother.moveNumAtWinMin > vPA.TDotherSQL.moveNumAtWin {
			vPA.ADother.moveNumAtWinMin = vPA.TDotherSQL.moveNumAtWin
		}
		if vPA.ADother.moveNumAtWinMax == 0 || vPA.ADother.moveNumAtWinMax < vPA.TDotherSQL.moveNumAtWin {
			vPA.ADother.moveNumAtWinMax = vPA.TDotherSQL.moveNumAtWin
		}

		vPA.TDotherSQL.moveNumAtWin = 0
		vPA.TDotherSQL.moveNumMax = 0 //to keep track of length of the longest strategy so far

		vPA.ADother.decksPlayed++
		vPA.TDother.treePrevMoves = ""
		vPA.TDother.winningMoves = nil
		vPA.TDother.startTime = time.Now()
		vPA.TDother.lastPrintTime = time.Now()
		clear(vPA.TDother.priorBoards)
		vPA.TD.stratNum = 0
		vPA.TD.stratTried = 1 // NOTE: Starts at 1 not 0

	}

	// At this point, all decks to be played have been played.  Time to report aggregate won loss.
	// From this point on, the program only prints.

	fmt.Printf("\n\n******************   Summary Statistics   ******************\n")
	averageElapsedTimePerDeck := time.Duration(float64(time.Since(vPA.ADother.startTime)) / float64(numberOfDecksToBePlayed))
	fmt.Printf("\n     Elapsed Time: %v", time.Since(vPA.ADother.startTime).Truncate(100*time.Millisecond).String())
	fmt.Printf("\nAvg Time per Deck: %v\n", averageElapsedTimePerDeck.Truncate(100*time.Millisecond).String())
	_, _ = pfmt.Printf("\n          Decks Played: %6d", vPA.ADother.decksPlayed)
	_, _ = pfmt.Printf("\n             Decks Won: %6d   %4v%%", vPA.ADother.decksWon, roundFloatIntDiv(vPA.ADother.decksWon*100, vPA.ADother.decksPlayed, 1))
	_, _ = pfmt.Printf("\n            Decks Lost: %6d   %4v%%", vPA.ADother.decksLost, roundFloatIntDiv(vPA.ADother.decksLost*100, vPA.ADother.decksPlayed, 1))
	_, _ = pfmt.Printf("\n         Decks LostGLE: %6d   %4v%%", vPA.ADother.decksLostGLE, roundFloatIntDiv(vPA.ADother.decksLostGLE*100, vPA.ADother.decksPlayed, 1))
	_, _ = pfmt.Printf("\n  Decks Lost + LostGLE: %6d   %4v%%", vPA.ADother.decksLostGLE, roundFloatIntDiv((vPA.ADother.decksLost+vPA.ADother.decksLostGLE)*100, vPA.ADother.decksPlayed, 1))

	statisticsPrint(&vPA.AD, "AD")
	ADotherSQLPrint(&vPA)

	if cfg.PlayAll.WinLossReport { // Deck Win Loss Summary Statistics
		// Close sql/csv file for writing and open it for reading and report it here
	}
}

// Divide 2 integers and round to precision digits
func roundFloatIntDiv(numer int, denom int, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(float64(numer)/float64(denom)*ratio) / ratio
}

func statisticsPrint(TDorAD *Statistics, which string) {
	if which != "AD" {
		if TDorAD.elapsedTime == 0 {
			fmt.Printf("\nElapsed Time: <.5ms    (Windows Minimum Resolution)")
		} else {
			fmt.Printf("\nElapsed Time: %v", TDorAD.elapsedTime)
		}
	}
	fmt.Printf("\n\nStrategies:")
	_, _ = pfmt.Printf("\n   Tried: %13d", TDorAD.stratTried)
	fmt.Printf("   Tried Detail:                 (Must sum to Strategies Tried)")
	if TDorAD.stratWins != 0 {
		_, _ = pfmt.Printf("\n                                   Won: %13d", TDorAD.stratWins)
	}
	if TDorAD.stratLossesNMA != 0 {
		_, _ = pfmt.Printf("\n                                   NMA: %13d   (No Moves Available)", TDorAD.stratLossesNMA)
	}
	if TDorAD.stratLossesRB != 0 {
		_, _ = pfmt.Printf("\n                                    RB: %13d   (Repetitive Board)", TDorAD.stratLossesRB)
	}
	if TDorAD.stratLossesEL != 0 {
		_, _ = pfmt.Printf("\n                                   EL: %13d   (Early Loss)", TDorAD.stratLossesEL)
	}
	if TDorAD.stratLossesGLE != 0 {
		_, _ = pfmt.Printf("\n                                   GLE: %13d   (Game Length Exceeded)", TDorAD.stratLossesGLE)
	}
	if TDorAD.stratTried != TDorAD.stratWins+TDorAD.stratLossesNMA+TDorAD.stratLossesRB+TDorAD.stratLossesEL+TDorAD.stratLossesGLE {
		fmt.Printf("\n        ************* Strategies Tried != TDorAD.stratWins+TDorAD.stratLossesNMA+TDorAD.stratLossesRB+TDorAD.stratLossesEL+TDorAD.stratLossesGLE")
	}
	fmt.Printf("\n\nMoves:")
	_, _ = pfmt.Printf("\n   Tried: %13d", TDorAD.mvsTried)
	fmt.Printf("   Tried Detail:                 (Must Sum to Moves Tried)")
	if TDorAD.stratLossesNMA != 0 {
		_, _ = pfmt.Printf("\n                                   NMA: %13d   (No Moves Available)", TDorAD.stratLossesNMA)
	}
	if TDorAD.stratLossesRB != 0 {
		_, _ = pfmt.Printf("\n                                    RB: %13d   (Repetitive Board)", TDorAD.stratLossesRB)
	}
	if TDorAD.stratLossesMajSE != 0 {
		_, _ = pfmt.Printf("\n                                 MajSE: %13d   (Major Exhausted)", TDorAD.stratLossesMajSE)
	}
	if TDorAD.stratLossesMinSE != 0 {
		_, _ = pfmt.Printf("\n                                 MinSE: %13d   (Minor Exhausted)", TDorAD.stratLossesMinSE)
	}
	if TDorAD.stratLossesEL != 0 {
		_, _ = pfmt.Printf("\n                                    EL: %13d   (Early Loss)", TDorAD.stratLossesEL)
	}
	if TDorAD.stratLossesGLE != 0 {
		_, _ = pfmt.Printf("\n                                   GLE: %13d   (Game Length Exceeded)", TDorAD.stratLossesGLE)
	}
	if TDorAD.stratLossesGLEAb != 0 {
		_, _ = pfmt.Printf("\n                                 GLEAb: %13d   (Game Length Exceeded Aborted moves)", TDorAD.stratLossesGLEAb)
	}
	if TDorAD.winningMovesCnt != 0 {
		_, _ = pfmt.Printf("\n                                 WMCnt: %13d   (Winning Moves Count)", TDorAD.winningMovesCnt)
	}
	x := TDorAD.stratLossesNMA + TDorAD.stratLossesRB + TDorAD.stratLossesMajSE + TDorAD.stratLossesMinSE + TDorAD.stratLossesEL + TDorAD.stratLossesGLE + TDorAD.stratLossesGLEAb + TDorAD.winningMovesCnt
	if TDorAD.mvsTried != x {
		fmt.Printf("\n        ************* Moves Tried %v != TDorAD.stratLossesNMA+TDorAD.stratLossesRB+TDorAD.stratLossesMajSE+TDorAD.stratLossesMinSE+TDorAD.stratLossesEL+TDorAD.stratLossesGLE+TDorAD.winningMovesCnt", x)
	}
	fmt.Printf("\n")
	if TDorAD.unqBoards != 0 {
		_, _ = pfmt.Printf("\n  UnqBds: %13d   (Unique Boards)", TDorAD.unqBoards)
	}
}

func statisticsPrintOneLine(vPA *variablesSpecificToPlayAll, dN int, s1 string, firstDeckNum int, numberOfDecksToBePlayed int) {
	var est time.Duration
	//                      nanosecondsTD   / Decks Played So Far         * remaining decks [remaining decks = numbertobeplayed - decksplayed so far
	est = time.Duration(float64(time.Since(vPA.ADother.startTime))/float64(dN+1-firstDeckNum)*float64(numberOfDecksToBePlayed-(dN+1-firstDeckNum))) * time.Nanosecond
	elTimeSinceStartTimeADFormatted := time.Since(vPA.ADother.startTime).Truncate(100 * time.Millisecond).String()
	if time.Since(vPA.ADother.startTime) > time.Duration(5*time.Minute) {
		elTimeSinceStartTimeADFormatted = time.Since(vPA.ADother.startTime).Truncate(time.Second).String()
	}
	_, _ = pfmt.Printf("Deck: %7v%9s   Strategy #: %11v  Moves Tried: %11v   Unique Boards: %10v   Elapsed TD: %10v"+
		"   stratTried: %11v   stratWins: %1v    stratLossesNMA: %9v    stratLossesRB: %10v    stratLossesEL: %1v    stratLossesGLE: %5v"+
		"   stratLossesMajSE: %11v    stratLossesMinSE: %11v    stratLossesEL: %7v    stratLossesGLEAb: %3v"+
		"   winningMovesCnt: %3v    moveNumMax: %3v    moveNumAtWin:%3v   Elapsed AD: %10s   Est Rem: %10s   Now: %8s\n",
		dN,
		s1,
		vPA.TD.stratNum,
		vPA.TD.mvsTried,
		vPA.TD.unqBoards,
		vPA.TD.elapsedTime.Truncate(10*time.Millisecond).String(),
		vPA.TD.stratTried, vPA.TD.stratWins, vPA.TD.stratLossesNMA, vPA.TD.stratLossesRB, vPA.TD.stratLossesEL, vPA.TD.stratLossesGLE,
		vPA.TD.stratLossesMajSE, vPA.TD.stratLossesMinSE, vPA.TD.stratLossesEL, vPA.TD.stratLossesGLEAb,
		vPA.TD.winningMovesCnt, vPA.TDotherSQL.moveNumMax, vPA.TDotherSQL.moveNumAtWin,
		elTimeSinceStartTimeADFormatted, est.Truncate(100*time.Millisecond).String(), time.Now().Format(" 3:04 pm"))
}

func TDotherSQLPrint(x *TDotherSQL) {
	if x.moveNumMax != 0 {
		_, _ = pfmt.Printf("\n   MNMax: %13d   (Move Number Max)", x.moveNumMax)
	}
	if x.moveNumAtWin != 0 {
		_, _ = pfmt.Printf("\n   MNWin: %13d   (Move Number at Win)", x.moveNumAtWin)
	}
}

func ADotherSQLPrint(vPA *variablesSpecificToPlayAll) {
	if vPA.ADother.moveNumMax != 0 {
		_, _ = pfmt.Printf("\n   MNMax: %13d   (Move Number Max)", vPA.ADother.moveNumMax)
	}
	if vPA.ADother.moveNumAtWinMin != 0 {
		_, _ = pfmt.Printf("\nMinMNWin: %13d   (Minimum Move Number at Win)", vPA.ADother.moveNumAtWinMin)
	}
	if vPA.ADother.moveNumAtWinMax != 0 {
		_, _ = pfmt.Printf("\nMaxMNWin: %13d   (Maximum Move Number at Win)\n\n\n", vPA.ADother.moveNumAtWinMax)
	}
	fmt.Printf("\n\n\n")
}
