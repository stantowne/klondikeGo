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

// type cumulativeByDeckVariables map[string] int v, ???????????????

type variablesSpecificToPlayAll struct {
	priorBoards map[bCode]bool // NOTE: bcode is an array of 65 ints as defined in board.go
	TD          struct {
		mvsTried       int
		stratWins      int
		stratLosses    int
		stratLossesGLE int // Strategy Game Length Exceeded
		stratLossesNMA int // Strategy No Moves Available
		stratLossesRB  int // Strategy Repetitive Board
		stratLossesSE  int // Strategy Exhausted
		stratLossesEL  int // Strategy Early Loss
		stratNum       int
		startTime      time.Time
		treePrevMoves  string
	}
	AD struct {
		mvsTried       int
		stratWins      int
		stratLosses    int
		stratLossesGLE int // Strategy Game Length Exceeded
		stratLossesNMA int // Strategy No Moves Available
		stratLossesRB  int // Strategy Repetitive Board
		stratLossesSE  int // Strategy Exhausted
		stratLossesEL  int // Strategy Early Loss
		stratNum       int
		startTime      time.Time
		deckWins       int
		deckLosses     int
	}
}

// var stratWinsTD = 0
// var stratLossesTD = 0
// var stratLossesGLE_TD = 0
// var stratLossesNMA_TD = 0
// var stratLossesRB_TD = 0
// var stratLossesSE_TD = 0
//var stratNumTD = 0

//var mvsTriedTD = 0
//var mvsTriedAD = 0

/*type deckWinLossDetailStats struct {
	wonLost                                        string //
	moveNum                                        int    // Do Not Report if wonTotal > 1  OR  if wonLost == "Lost"
	stratNum                                       int    // vPA.TD.stratNum
	unqBoards                                      int
	elapsedTime                                    time.Duration // time.Since(vPA.TD.startTime)
	ifWonAndFindAllSuccessfulStrategiesAndGLE_flag bool          // true to indicate that Elapsed time This Deck may not be relevant
	resultCode1                                    string        // result1
	resultCode2                                    string        // result2
	moveNumMin                                     int           //
	moveNumMax                                     int           //
	wonTotal                                       int
	elapsedTimeAt1stWin                            time.Duration
	moveNumAt1stWin                                int //
	mvsTriedAt1stWin                               int //
	stratNumAt1stWin                               int //
	moveNumMinAt1stWin                             int //
	moveNumMaxAt1stWin                             int //
}*/

//var dWLDStats deckWinLossDetailStats
//var deckWinLossDetail []deckWinLossDetailStats

func playAll(reader csv.Reader, cfg *Configuration) {
	firstDeckNum := cfg.General.FirstDeckNum
	numberOfDecksToBePlayed := cfg.General.NumberOfDecksToBePlayed
	verbose := cfg.General.Verbose
	var vPA variablesSpecificToPlayAll
	vPA.priorBoards = map[bCode]bool{}
	vPA.TD.treePrevMoves = ""

	//var deckWinsAD = 0
	//var deckLossesAD = 0
	//var stratWinsAD = 0
	//var stratLossesAD = 0
	//var stratLossesGLE_AD = 0
	//var stratLossesNMA_AD = 0
	//var stratLossesRB_AD = 0
	//var stratLossesSE_AD = 0
	//var stratNumAD = 0
	vPA.TD.startTime = time.Now()
	vPA.AD.startTime = time.Now()

	for deckNum := firstDeckNum; deckNum < (firstDeckNum + numberOfDecksToBePlayed); deckNum++ {

		vPA.TD.startTime = time.Now()
		moveNumMax = 0                   //to keep track of length of the longest strategy so far
		protoDeck, err5 := reader.Read() // protoDeck is a slice of strings: rank, suit, rank, suit, etc.  // err5 used to avoid shadowing err
		if err5 == io.EOF {
			break
		}
		if err5 != nil {
			log.Println("Cannot read from inputFileName:", err5)
		}

		if verbose > 1 {
			_, err = pfmt.Printf("\nDeck #%d:\n", deckNum)
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
		result1, result2, _ := playAllMoves(b, 0, deckNum, cfg, &vPA)

		var dummy []move
		var s string
		if result2 == "EW" {
			vPA.AD.deckWins += 1
			s = "DECK WON"
		} else {
			vPA.AD.deckLosses += 1
			s = "DECK LOST"
		}
		prntMDet(b, dummy, 1, deckNum, 1, "DbDorMbM", 2, "\n   "+s+"\n", "", "", cfg, &vPA) // "DbDorMbM" was formerly "NOTX"
		prntMDetTreeReturnComment("\n   "+s+"\n", deckNum, 0, cfg, &vPA)

		/*
			if cfg.PlayAll.WinLossReport { // Deck Win Loss Summary Statistics
				if vPA.TD.stratWins == 0 {
					dWLDStats.winLoss = "L"
					dWLDStats.moveNumAt1stWinOrAtLoss = 0
					dWLDStats.moveNumMinWinIfFindAll = 0
					dWLDStats.moveNumMaxWinIfFindAll = 0
					dWLDStats.stratNumAt1stWinOrAtLoss = vPA.TD.stratWins
					dWLDStats.mvsTriedAt1stWinOrAtLoss = vPA.TD.mvsTried
					// Add maxUnique Boards
					//add max movenum
					// add ElapsedTimeAt1stWin
					dWLDStats.ElapsedTimeAt1stWinOrAtLoss = time.Now().Sub(vPA.TD.startTime)
				} else {
					dWLDStats.winLoss = "W"
					dWLDStats.ElapsedTimeAt1stWinOrAtLoss = time.Now().Sub(vPA.TD.startTime)
				}
				deckWinLossDetail = append(deckWinLossDetail, dWLDStats)

			}*/

		// This If Block is Print Only
		if cfg.PlayAll.ReportingType.DeckByDeck && cfg.PlayAll.DeckByDeckReportingOptions.Type == "regular" { // Deck-by-deck Statistics
			if vPA.TD.stratWins > 0 {
				fmt.Printf("\n\n*************************\n\nDeck: %d  WON    Result Codes: %v %v", deckNum, result1, result2)
			} else {
				fmt.Printf("\n\n*************************\n\nDeck: %d  LOST   Result Codes: %v %v", deckNum, result1, result2)
			}
			fmt.Printf("\nElapsed Time is %v.", time.Since(vPA.TD.startTime))
			fmt.Printf("\n\nStrategies:")
			_, err = pfmt.Printf("\n Tried: %13d", vPA.TD.stratNum)
			_, err = pfmt.Printf("\n   Won: %13d", vPA.TD.stratWins)
			_, err = pfmt.Printf("\n  Lost: %13d", vPA.TD.stratLosses)
			fmt.Printf("\n\nStrategies Lost Detail:")
			_, err = pfmt.Printf("\n   NMA: %13d   (No Moves Available)", vPA.TD.stratLossesNMA)
			_, err = pfmt.Printf("\n    RB: %13d   (Repetitive Board)", vPA.TD.stratLossesRB)
			_, err = pfmt.Printf("\n    SE: %13d   (Strategy Exhausted)", vPA.TD.stratLossesSE)
			_, err = pfmt.Printf("\n   GLE: %13d   (Game Length Exceeded)", vPA.TD.stratLossesGLE)
			if vPA.TD.stratLossesNMA+vPA.TD.stratLossesRB+vPA.TD.stratLossesSE+vPA.TD.stratLossesGLE != vPA.TD.stratLosses {
				fmt.Printf("\n        ************* Total Strategy Losses != Sum of strategy detail")
			}
			if vPA.TD.stratLosses+vPA.TD.stratWins != vPA.TD.stratNum {
				fmt.Printf("\n        ************* Strategies Tried != Strategies Lost + Strategies Won")
			}
			if cfg.PlayAll.FindAllWinStrats && vPA.TD.stratWins != 0 {
				fmt.Printf("\n\n Multiple Successful Strategies were found in this deck.")
				_, err = pfmt.Printf("   Total winning strategies found: %d\n", vPA.TD.stratWins)
			}
		}

		// This If Block is Print Only   ??????????????  what was this for ????????       PROGRESS?????
		if cfg.PlayAll.ReportingType.DeckByDeck && cfg.PlayAll.DeckByDeckReportingOptions.Type != "regular" {
			var est time.Duration
			//                      nanosecondsTD   / Decks Played So Far         * remaining decks [remaining decks = numbertobeplayed - decksplayed so far
			est = time.Duration(float64(time.Since(vPA.AD.startTime))/float64(deckNum+1-firstDeckNum)*float64(numberOfDecksToBePlayed-(deckNum+1-firstDeckNum))) * time.Nanosecond
			/*			est = time.Duration(float64(elapsedTimeAD)           / float64(deckNum-firstDeckNum+1)*float64(numberOfDecksToBePlayed-(deckNum-firstDeckNum+1))) * time.Nanosecond
						fmt.Printf("       (deckNum+1-firstDeckNum)                                                           = %v\n", (deckNum + 1 - firstDeckNum))
						fmt.Printf("float64(deckNum+1-firstDeckNum)                                                           = %v\n", float64(deckNum+1-firstDeckNum))
						fmt.Printf("                                                                (deckNum+1-firstDeckNum)  = %v\n", (deckNum + 1 - firstDeckNum))
						fmt.Printf("                                        numberOfDecksToBePlayed                           = %v\n", numberOfDecksToBePlayed)
						fmt.Printf("                                       (numberOfDecksToBePlayed-(deckNum+1-firstDeckNum)) = %v\n", (numberOfDecksToBePlayed - (deckNum + 1 - firstDeckNum)))
						fmt.Printf("                                float64(numberOfDecksToBePlayed-(deckNum+1-firstDeckNum)) = %v\n", float64(numberOfDecksToBePlayed-(deckNum+1-firstDeckNum)))
						fmt.Printf("float64(deckNum+1-firstDeckNum)*float64(numberOfDecksToBePlayed-(deckNum+1-firstDeckNum)) = %v\n", float64(deckNum+1-firstDeckNum)*float64(numberOfDecksToBePlayed-(deckNum+1-firstDeckNum)))
						fmt.Printf("float64(deckNum+1-firstDeckNum)*float64(numberOfDecksToBePlayed-(deckNum+1-firstDeckNum)) = %v\n", float64(deckNum+1-firstDeckNum)*float64(numberOfDecksToBePlayed-(deckNum+1-firstDeckNum)))
						fmt.Printf("time.Since(vPA.AD.startTime)                                                                                                            = %v\n", time.Since(vPA.AD.startTime))
						fmt.Printf("time.Since(vPA.AD.startTime)                                                                                                            = %v\n", time.Since(vPA.AD.startTime))
						fmt.Printf("                          time.Duration(float64(deckNum+1-firstDeckNum)*float64(numberOfDecksToBePlayed-(deckNum+1-firstDeckNum))) = %v\n", time.Duration(float64(deckNum+1-firstDeckNum)*float64(numberOfDecksToBePlayed-(deckNum+1-firstDeckNum))))
						fmt.Printf("time.Since(vPA.AD.startTime) / time.Duration(float64(deckNum+1-firstDeckNum)*float64(numberOfDecksToBePlayed-(deckNum+1-firstDeckNum))) = %v\n", time.Since(vPA.AD.startTime)/time.Duration(float64(deckNum+1-firstDeckNum)*float64(numberOfDecksToBePlayed-(deckNum+1-firstDeckNum))))
			*/
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
			_, err = pfmt.Printf("Dk: %5d   "+wL+"   MvsTried: %13v   MoveNum: xxx   Max MoveNum: xxx   StratsTried: %12v   UnqBoards: %11v   Won: %5v   Lost: %5v   GLE: %5v   Won: %5.1f%%   Lost: %5.1f%%   GLE: %5.1f%%   ElTime TD: %9s   ElTime ADs: %9s  Rem Time: %11s   ResCodes: %2s %3s   Time Now: %8s\n", deckNum, vPA.TD.mvsTried /*moveNum, maxMoveNum, */, vPA.TD.stratNum, len(vPA.priorBoards), vPA.AD.deckWins, vPA.AD.deckLosses, vPA.AD.stratLossesGLE, roundFloatIntDiv(vPA.AD.deckWins*100, deckNum+1-firstDeckNum, 1), roundFloatIntDiv(vPA.AD.deckLosses*100, deckNum+1-firstDeckNum, 1), roundFloatIntDiv(vPA.AD.stratLossesGLE*100, deckNum+1-firstDeckNum, 1), time.Since(vPA.TD.startTime).Truncate(100*time.Millisecond).String(), elTimeSinceStartTimeADFormatted, est.Truncate(time.Second).String(), result1, result2, time.Now().Format(" 3:04 pm"))
		}

		// Verbose Special "BELL" Starts Here - No effect on operation
		/*if strings.Contains(verboseSpecial, ";BELL;") && time.Since(vPA.TD.startTime) > 1*time.Second { // changed 5*time.Minute to 1*time.Second for test change it back
			fmt.Printf("\a") // Ring Bell
		}*/
		// Verbose Special "BELL" Ends Here - No effect on operation

		vPA.AD.stratWins += vPA.TD.stratWins
		vPA.TD.stratWins = 0
		vPA.AD.stratLosses += vPA.TD.stratLosses
		vPA.TD.stratLosses = 0
		vPA.AD.stratLossesGLE += vPA.TD.stratLossesGLE
		vPA.TD.stratLossesGLE = 0
		vPA.AD.stratLossesNMA += vPA.TD.stratLossesNMA
		vPA.TD.stratLossesNMA = 0
		vPA.AD.stratLossesRB += vPA.TD.stratLossesRB
		vPA.TD.stratLossesRB = 0
		vPA.AD.stratLossesSE += vPA.TD.stratLossesSE
		vPA.TD.stratLossesSE = 0
		vPA.AD.stratNum += vPA.TD.stratNum + 1 // Because we start at strategy 0 which is all best moves
		vPA.TD.stratNum = 0
		vPA.AD.mvsTried += vPA.TD.mvsTried + 1
		vPA.TD.mvsTried = 0
		vPA.TD.treePrevMoves = ""
		clear(vPA.priorBoards)
	}

	// At this point, all decks to be played have been played.  Time to report aggregate won loss.
	// From this point on, the program only prints.

	fmt.Printf("\n\n******************   Summary Statistics   ******************\n" + "Decks:")
	_, err = pfmt.Printf("\n Played: %6d", numberOfDecksToBePlayed)
	_, err = pfmt.Printf("\n    Won: %6d", vPA.AD.deckWins)
	_, err = pfmt.Printf("\n   Lost: %6d", vPA.AD.deckLosses)
	averageElapsedTimePerDeck := time.Duration(float64(time.Since(vPA.AD.startTime)) / float64(numberOfDecksToBePlayed))
	fmt.Printf("\nElapsed Time is %v.", time.Since(vPA.AD.startTime))
	fmt.Printf("\nAverage Elapsed Time per Deck is %s", averageElapsedTimePerDeck.Truncate(100*time.Millisecond).String())
	fmt.Printf("\n\nStrategies:")
	_, err = pfmt.Printf("\n  Tried: %13d", vPA.AD.stratNum)
	_, err = pfmt.Printf("\n    Won: %13d", vPA.AD.stratLosses)
	_, err = pfmt.Printf("\n   Lost: %13d", vPA.AD.stratWins)
	fmt.Printf("\n\nStrategies Lost Detail:")
	_, err = pfmt.Printf("\n    NMA: %13d   (No Moves Available)", vPA.AD.stratLossesNMA)
	_, err = pfmt.Printf("\n     RB: %13d   (Repetitive Board)", vPA.AD.stratLossesRB)
	_, err = pfmt.Printf("\n     SE: %13d   (Strategy Exhausted)", vPA.AD.stratLossesSE)
	_, err = pfmt.Printf("\n    GLE: %13d   (Game Length Exceeded)", vPA.AD.stratLossesGLE)
	if vPA.AD.stratLossesNMA+vPA.AD.stratLossesRB+vPA.AD.stratLossesSE+vPA.AD.stratLossesGLE != vPA.AD.stratLosses {
		fmt.Printf("\n        ************* Total Strategy Losses != Sum of strategy detail")
	}
	if vPA.AD.stratLosses+vPA.AD.stratWins != vPA.AD.stratNum {
		fmt.Printf("\n        ************* Strategies Tried != Strategies Lost + Strategies Won")
	}
	if cfg.PlayAll.FindAllWinStrats {
		fmt.Printf("\n\n Multiple Successful Strategies were found in some winng decks.")
		_, err = pfmt.Printf("   Decks Won: %d\n", vPA.AD.deckWins)
		_, err = pfmt.Printf("   Total winning strategies found: %d\n", vPA.AD.stratWins)
		_, err = pfmt.Printf("   Average winning strategies found: %d\n", vPA.AD.stratWins/vPA.AD.deckWins)
	}

	if cfg.PlayAll.WinLossReport { // Deck Win Loss Summary Statistics
		fmt.Printf("\n\n\n Deck-by Deck Win/Loss Detail   (Copy to Excel to get headings to line up with the columns)")
		fmt.Printf("\n\n Deck\tW/L\tMoveNum 1ST-Win\tStratNum At 1st-Win Or At-Loss\tMvsTried At 1st-Win Or At-Loss\tMoveNum Min-Win If-Find-All\tMoveNum Max-Win If-Find-All\tElapsed Time At 1st-Win Or At Loss\n")
		/*for dN, detail := range deckWinLossDetail {
			_, err = pfmt.Printf("\n  %5v\t  %v\t%4v\t%8v\t%8v\t%4v\t%4v", dN, detail.winLoss, detail.moveNumAt1stWinOrAtLoss, detail.stratNumAt1stWinOrAtLoss, detail.mvsTriedAt1stWinOrAtLoss, detail.moveNumMinWinIfFindAll, detail.moveNumMaxWinIfFindAll, detail.ElapsedTimeAt1stWinOrAtLoss)
		}*/
	}
	fmt.Printf("\n\n\n")
}

// Divide 2 integers and round to precision digits
func roundFloatIntDiv(numer int, denom int, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(float64(numer)/float64(denom)*ratio) / ratio
}
