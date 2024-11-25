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

type ThisDeck struct {
	boardCodeOfDeck  [65]byte // The boardCode of the deck read
	mvsTried         int
	stratNum         int // stratNum starts at 0
	stratTried       int // This Deck Strategies TRIED = stratNum + 1
	stratWins        int
	stratLosses      int
	stratLossesGLE   int // Strategy Game Length Exceeded
	stratLossesGLEAb int // Strategy Game Length Exceeded Aborted moves
	stratLossesNMA   int // Strategy No Moves Available
	stratLossesRB    int // Strategy Repetitive Board
	stratLossesMajSE int // Strategy Exhausted Major
	stratLossesMinSE int // Strategy Exhausted Minor
	stratLossesEL    int // Strategy Early Loss
	unqBoards        int
	moveNumAtWin     int
	moveNumMax       int
	elapsedTime      time.Duration
	winningMoves     []move
	winningMovesCnt  int // = length of winning moves
}

type AllDecks struct {
	decksPlayed int
	mvsTried    int // NOTE: Removed various min and max versions of variables here in AD struct as they would only pertain to the specific decks included in this run
	//                      if info across any set of decks is wanted it should be derived from the saved deck history to be found in SQL
	stratTried       int
	stratWins        int
	stratLosses      int
	stratLossesGLE   int // Strategy Game Length Exceeded
	stratLossesGLEAb int // Strategy Game Length Exceeded Aborted moves
	stratLossesNMA   int // Strategy No Moves Available
	stratLossesRB    int // Strategy Repetitive Board
	stratLossesMajSE int // Strategy Exhausted
	stratLossesMinSE int // Strategy Exhausted Minor
	stratLossesEL    int // Strategy Early Loss
	unqBoards        int
	startTime        time.Time
	deckWins         int
	deckLosses       int
	winningMovesCnt  int // = length of winning moves
}

type variablesSpecificToPlayAll struct {
	priorBoards map[bCode]bool // NOTE: bcode is an array of 65 ints as defined in board.go
	TD          ThisDeck
	AD          AllDecks
	/*TD          struct {       // TD = This Deck
		boardCodeOfDeck  [65]byte // The boardCode of the deck read
		mvsTried         int
		stratNum         int // stratNum starts at 0
		stratTried       int // This Deck Strategies TRIED = stratNum + 1
		stratWins        int
		stratLosses      int
		stratLossesGLE   int // Strategy Game Length Exceeded
		stratLossesGLEAb int // Strategy Game Length Exceeded Aborted moves
		stratLossesNMA   int // Strategy No Moves Available
		stratLossesRB    int // Strategy Repetitive Board
		stratLossesMajSE int // Strategy Exhausted Major
		stratLossesMinSE int // Strategy Exhausted Minor
		stratLossesEL    int // Strategy Early Loss
		unqBoards        int
		moveNumAtWin     int
		moveNumMax       int
		elapsedTime      time.Duration
		winningMoves     []move
		winningMovesCnt  int // = length of winning moves
	}*/
	TDother struct { // Variables NOT needed in SQL output
		startTime     time.Time
		treePrevMoves string // Used to retain values between calls to prntMDetTree for a single deck - Needed for when the strategy "Backs Uo"
	}
	/*AD struct { // AD = All Decks This run
		decksPlayed int
		mvsTried    int // NOTE: Removed various min and max versions of variables here in AD struct as they would only pertain to the specific decks included in this run
		//                      if info across any set of decks is wanted it should be derived from the saved deck history to be found in SQL
		stratTried       int
		stratWins        int
		stratLosses      int
		stratLossesGLE   int // Strategy Game Length Exceeded
		stratLossesGLEAb int // Strategy Game Length Exceeded Aborted moves
		stratLossesNMA   int // Strategy No Moves Available
		stratLossesRB    int // Strategy Repetitive Board
		stratLossesMajSE int // Strategy Exhausted
		stratLossesMinSE int // Strategy Exhausted Minor
		stratLossesEL    int // Strategy Early Loss
		unqBoards        int
		startTime        time.Time
		deckWins         int
		deckLosses       int
		winningMovesCnt  int // = length of winning moves
	}*/
}

func playAll(reader csv.Reader, cfg *Configuration) {
	firstDeckNum := cfg.General.FirstDeckNum                       // Shorthand name but really is a copy - OK since never changed (but would Pointer or address be better?)
	numberOfDecksToBePlayed := cfg.General.NumberOfDecksToBePlayed // Shorthand name but really is a copy - OK since never changed (but would Pointer or address be better?)
	verbose := cfg.General.Verbose                                 // Shorthand name but really is a copy - OK since never changed (but would Pointer or address be better?)
	var vPA variablesSpecificToPlayAll
	vPA.priorBoards = map[bCode]bool{}
	vPA.TDother.treePrevMoves = ""
	vPA.TD.stratTried = 1
	vPA.TDother.startTime = time.Now()
	vPA.AD.startTime = time.Now()

	for deckNum := firstDeckNum; deckNum < (firstDeckNum + numberOfDecksToBePlayed); deckNum++ {

		vPA.TD.moveNumMax = 0           //to keep track of length of the longest strategy so far
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
		vPA.TD.boardCodeOfDeck = b.boardCode(deckNum)
		// This statement is executed once per deck and transfers program execution.
		// When this statement returns the deck has been played.
		result1, result2 := playAllMoves(b, 0, deckNum, cfg, &vPA)

		vPA.TD.elapsedTime = time.Since(vPA.TDother.startTime)
		vPA.TD.winningMovesCnt = len(vPA.TD.winningMoves)
		var dummy []move
		var s string
		if result1 == "SW" {
			vPA.AD.deckWins += 1
			s = "DECK WON"
		} else {
			vPA.TD.stratLosses++
			vPA.AD.deckLosses += 1
			s = "DECK LOST"
		}
		prntMDet(b, dummy, 1, deckNum, 1, "MbM_R", 2, "\n   "+s+"\n", "", "", cfg, &vPA)
		prntMDetTreeReturnComment("\n   "+s+"\n", deckNum, 0, cfg, &vPA)

		// This If Block is Print Only for DbD_R  ????
		if cfg.PlayAll.ReportingType.DeckByDeck && cfg.PlayAll.DeckByDeckReportingOptions.Type == "regular" { // Deck-by-deck Statistics
			if vPA.TD.stratWins > 0 {
				_, _ = pfmt.Printf(" \n*************************\n\nDeck: %d  WON    Result Codes: %v", deckNum, result1)
			} else {
				_, _ = pfmt.Printf(" \n*************************\n\nDeck: %d  LOST   Result Codes: %v", deckNum, result1)
			}
			if vPA.TD.elapsedTime == 0 {
				fmt.Printf("\nElapsed Time: <.5ms    (Windows Minimum Resolution)")
			} else {
				fmt.Printf("\nElapsed Time: %v", vPA.TD.elapsedTime)
			}
			fmt.Printf("\n\nStrategies:")
			_, _ = pfmt.Printf("\n Tried: %13d", vPA.TD.stratTried)
			fmt.Printf("\n\n Tried Detail:   (Must sum to Strategies Tried)")
			if vPA.TD.stratWins != 0 {
				_, _ = pfmt.Printf("\n     Won: %13d", vPA.TD.stratWins)
			}
			if vPA.TD.stratLossesNMA != 0 {
				_, _ = pfmt.Printf("\n     NMA: %13d   (No Moves Available)", vPA.TD.stratLossesRB)
			}
			if vPA.TD.stratLossesRB != 0 {
				_, _ = pfmt.Printf("\n      RB: %13d   (Repetitive Board)", vPA.TD.stratLossesRB)
			}
			if vPA.TD.stratLossesEL != 0 {
				_, _ = pfmt.Printf("\n      EL: %13d   (Early Loss)", vPA.TD.stratLossesEL)
			}
			if vPA.TD.stratLossesGLE != 0 {
				_, _ = pfmt.Printf("\n     GLE: %13d   (Game Length Exceeded)", vPA.TD.stratLossesGLE)
			}
			if vPA.TD.stratTried != vPA.TD.stratWins+vPA.TD.stratLossesNMA+vPA.TD.stratLossesRB+vPA.TD.stratLossesEL+vPA.TD.stratLossesGLE {
				fmt.Printf("\n        ************* Strategies Tried != vPA.TD.stratWins+vPA.TD.stratLossesNMA+vPA.TD.stratLossesRB+vPA.TD.stratLossesEL+vPA.TD.stratLossesGLE")
			}
			fmt.Printf("\n\nMoves:")
			_, _ = pfmt.Printf("\n Tried: %13d", vPA.TD.mvsTried)
			fmt.Printf("\n Tried Detail:   (Must Sum to Moves Tried)")
			if vPA.TD.stratLossesNMA != 0 {
				_, _ = pfmt.Printf("\n     NMA: %13d   (No Moves Available)", vPA.TD.stratLossesNMA)
			}
			if vPA.TD.stratLossesRB != 0 {
				_, _ = pfmt.Printf("\n      RB: %13d   (Repetitive Board)", vPA.TD.stratLossesRB)
			}
			if vPA.TD.stratLossesMajSE != 0 {
				_, _ = pfmt.Printf("\n   MajSE: %13d   (Major Exhausted)", vPA.TD.stratLossesMajSE)
			}
			if vPA.TD.stratLossesMinSE != 0 {
				_, _ = pfmt.Printf("\n   MinSE: %13d   (Minor Exhausted)", vPA.TD.stratLossesMinSE)
			}
			if vPA.TD.stratLossesEL != 0 {
				_, _ = pfmt.Printf("\n      EL: %13d   (Early Loss)", vPA.TD.stratLossesEL)
			}
			if vPA.TD.stratLossesGLE != 0 {
				vPA.TD.stratLossesGLEAb = result2 - 2
				_, _ = pfmt.Printf("\n     GLE: %13d   (Game Length Exceeded)", vPA.TD.stratLossesGLE)
				_, _ = pfmt.Printf("\n   GLEAb: %13d   (Game Length Exceeded Aborted moves)", vPA.TD.stratLossesGLEAb)
			}
			if vPA.TD.winningMovesCnt != 0 {
				_, _ = pfmt.Printf("\n   WMCnt: %13d   (Winning Moves Count)", vPA.TD.winningMovesCnt)
			}
			if vPA.TD.mvsTried != vPA.TD.stratLossesNMA+vPA.TD.stratLossesRB+vPA.TD.stratLossesMajSE+vPA.TD.stratLossesMinSE+vPA.TD.stratLossesEL+vPA.TD.stratLossesGLE+vPA.TD.stratLossesGLEAb+vPA.TD.winningMovesCnt {
				fmt.Printf("\n        ************* Moves Tried != vPA.TD.mvsTried != vPA.TD.stratLossesNMA+vPA.TD.stratLossesRB+vPA.TD.stratLossesMajSE+vPA.TD.stratLossesMinSE+vPA.TD.stratLossesEL+vPA.TD.stratLossesGLE+vPA.TD.winningMovesCnt")
			}
			fmt.Printf("\n")
		}

		// This If Block is Print Only for DbD_S or DbD_VS
		if cfg.PlayAll.ReportingType.DeckByDeck && cfg.PlayAll.DeckByDeckReportingOptions.Type != "regular" {
			var est time.Duration
			//                      nanosecondsTD   / Decks Played So Far         * remaining decks [remaining decks = numbertobeplayed - decksplayed so far
			est = time.Duration(float64(time.Since(vPA.AD.startTime))/float64(deckNum+1-firstDeckNum)*float64(numberOfDecksToBePlayed-(deckNum+1-firstDeckNum))) * time.Nanosecond
			wL := ""
			if vPA.TD.stratWins > 0 {
				wL = "WON " // Note additional space -- for alignment
			} else {
				wL = "LOST"
			}
			elTimeSinceStartTimeADFormatted := time.Since(vPA.AD.startTime).Truncate(100 * time.Millisecond).String()
			if time.Since(vPA.AD.startTime) > time.Duration(5*time.Minute) {
				elTimeSinceStartTimeADFormatted = time.Since(vPA.AD.startTime).Truncate(time.Second).String()
			}
			_, _ = pfmt.Printf("Dk: %5d   "+wL+"   MvsTried: %13v   MoveNum: %3v   Max MoveNum: %3v   StratsTried: %12v   UnqBoards: %11v   Won: %5v   Lost: %5v   GLE: %5v   Won: %5.1f%%   Lost: %5.1f%%   GLE: %5.1f%%   ElTime TD: %9s   ElTime ADs: %9s  Rem Time: %11s   ResCodes: %2s %3s   Time Now: %8s\n", deckNum, vPA.TD.mvsTried, vPA.TD.moveNumAtWin, vPA.TD.moveNumMax, vPA.TD.stratNum, len(vPA.priorBoards), vPA.AD.deckWins, vPA.AD.deckLosses, vPA.AD.stratLossesGLE, roundFloatIntDiv(vPA.AD.deckWins*100, deckNum+1-firstDeckNum, 1), roundFloatIntDiv(vPA.AD.deckLosses*100, deckNum+1-firstDeckNum, 1), roundFloatIntDiv(vPA.AD.stratLossesGLE*100, deckNum+1-firstDeckNum, 1), vPA.TD.elapsedTime.Truncate(100*time.Millisecond).String(), elTimeSinceStartTimeADFormatted, est.Truncate(time.Second).String(), result1, "", time.Now().Format(" 3:04 pm"))
		}
		// if Winning moves to be printed or output of Deck is to be reported to SQL
		if cfg.PlayAll.SaveResultsToSQL || cfg.PlayAll.PrintWinningMoves {
			// First Reverse the slice (which was collected in reverse as we backed up the call chain)
			for i := 0; i < len(vPA.TD.winningMoves)/2; i++ {
				vPA.TD.winningMoves[i], vPA.TD.winningMoves[len(vPA.TD.winningMoves)-i-1] = vPA.TD.winningMoves[len(vPA.TD.winningMoves)-i-1], vPA.TD.winningMoves[i]
			}
			if cfg.PlayAll.PrintWinningMoves {
				fmt.Printf("\n     Winning Moves:\n")
				for mN := range vPA.TD.winningMoves {
					m1, m2 := printMove(vPA.TD.winningMoves[mN], true)
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
		vPA.AD.stratTried += vPA.TD.stratTried
		vPA.TD.stratTried = 1
		vPA.AD.mvsTried += vPA.TD.mvsTried + 1
		vPA.AD.decksPlayed++
		vPA.TD.mvsTried = 0
		vPA.TDother.treePrevMoves = ""
		vPA.TD.moveNumMax = 0
		vPA.TD.moveNumAtWin = 0
		vPA.TD.winningMoves = nil
		vPA.TD.winningMovesCnt = 0
		clear(vPA.priorBoards)
		vPA.TD.stratNum = 0
		vPA.TDother.startTime = time.Now()
	}

	// At this point, all decks to be played have been played.  Time to report aggregate won loss.
	// From this point on, the program only prints.

	fmt.Printf("\n\n******************   Summary Statistics   ******************\n" + "Decks:")
	_, _ = pfmt.Printf("\nDecks Played: %6d", vPA.AD.decksPlayed)
	_, _ = pfmt.Printf("\n   Decks Won: %6d", vPA.AD.deckWins)
	_, _ = pfmt.Printf("\n  Decks Lost: %6d", vPA.AD.deckLosses)
	_, _ = pfmt.Printf("\n\n Moves Tried: %6d", vPA.AD.mvsTried)
	averageElapsedTimePerDeck := time.Duration(float64(time.Since(vPA.AD.startTime)) / float64(numberOfDecksToBePlayed))
	fmt.Printf("\nElapsed Time: %v", time.Since(vPA.AD.startTime))
	fmt.Printf("\nAverage Elapsed Time per Deck is %v", averageElapsedTimePerDeck)
	fmt.Printf("\n\nStrategies:")
	_, _ = pfmt.Printf("\n Tried: %13d", vPA.AD.stratTried)
	_, _ = pfmt.Printf("\n   Won: %13d", vPA.AD.stratWins)
	_, _ = pfmt.Printf("\n  Lost: %13d", vPA.AD.stratLosses)
	fmt.Printf("\n\nStrategies Lost Detail:")
	_, _ = pfmt.Printf("\n   NMA: %13d   (No Moves Available)", vPA.AD.stratLossesNMA)
	_, _ = pfmt.Printf("\n    RB: %13d   (Repetitive Board)", vPA.AD.stratLossesRB)
	_, _ = pfmt.Printf("\n MajSE: %13d   (Major Exhausted)", vPA.AD.stratLossesMajSE)
	_, _ = pfmt.Printf("\n MinSE: %13d   (Minor Exhausted)", vPA.AD.stratLossesMinSE)
	_, _ = pfmt.Printf("\n    EL: %13d   (Early Loss)", vPA.AD.stratLossesEL)
	_, _ = pfmt.Printf("\n   GLE: %13d   (Game Length Exceeded)", vPA.AD.stratLossesGLE)
	if vPA.AD.stratLossesNMA+vPA.AD.stratLossesRB+vPA.AD.stratLossesMajSE+vPA.AD.stratLossesEL+vPA.AD.stratLossesGLE != vPA.AD.stratLosses {
		fmt.Printf("\n        ************* Total Strategy Losses != Sum of strategy detail NOT including stratLossesMinSE")
	}
	if vPA.AD.stratLosses+vPA.AD.stratWins != vPA.AD.decksPlayed {
		fmt.Printf("\n        ************* Strategies Tried != Strategies Lost + Strategies Won")
	}
	if vPA.AD.stratLosses+vPA.AD.stratWins != vPA.AD.stratTried {
		fmt.Printf("\n        ************* Strategies Tried != Strategies Lost + Strategies Won")
	}
	if vPA.AD.mvsTried != vPA.AD.stratLosses+vPA.AD.stratLossesMinSE+vPA.AD.stratWins {
		fmt.Printf("\n\nMoves Tried != strat Losses + strat LossesMinSE + strat Wins")
	}

	if cfg.PlayAll.WinLossReport { // Deck Win Loss Summary Statistics
		// Close sql/csv file for writing and open it for reading and report it here
	}
}

// Divide 2 integers and round to precision digits
func roundFloatIntDiv(numer int, denom int, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(float64(numer)/float64(denom)*ratio) / ratio
}
