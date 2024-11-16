package main

import (
	"encoding/csv"
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"io"
	"log"
	"math"
	"strconv"
	"time"
)

// type cumulativeByDeckVariables map[string] int v,
type variablesSpecificToPlayNew struct {
	priorBoards                                map[bCode]bool // NOTE: bcode is an array of 65 ints as defined in board.go
	verboseSpecialProgressCounterLastPrintTime time.Time
	treePrevMovesTD                            string
}

//var treePrevMovesTD string

var stratWinsTD = 0
var stratLossesTD = 0

var stratLossesGLE_TD = 0
var stratLossesNMA_TD = 0
var stratLossesRB_TD = 0
var stratLossesSE_TD = 0
var stratNumTD = 0
var mvsTriedTD = 0

//var mvsTriedAD = 0

var startTimeTD time.Time
var startTimeAD time.Time

type deckWinLossDetailStats struct {
	wonLost                                        string //
	moveNum                                        int    // Do Not Report if wonTotal > 1  OR  if wonLost == "Lost"
	mvsTried                                       int    // mvsTriedTD
	stratNum                                       int    // stratNumTD
	unqBoards                                      int
	elapsedTime                                    time.Duration // time.Since(startTimeTD)
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
}

var dWLDStats deckWinLossDetailStats
var deckWinLossDetail []deckWinLossDetailStats

func playNew(reader csv.Reader, cfg Configuration) {
	firstDeckNum := cfg.General.FirstDeckNum
	numberOfDecksToBePlayed := cfg.General.NumberOfDecksToBePlayed
	verbose := cfg.General.Verbose
	var varSp2PN variablesSpecificToPlayNew
	varSp2PN.priorBoards = map[bCode]bool{}
	varSp2PN.treePrevMovesTD = ""
	//treePrevMovesTD = ""

	varSp2PN.verboseSpecialProgressCounterLastPrintTime = time.Now()
	var deckWinsAD = 0
	var deckLossesAD = 0
	var stratWinsAD = 0
	var stratLossesAD = 0
	var stratLossesGLE_AD = 0
	var stratLossesNMA_AD = 0
	var stratLossesRB_AD = 0
	var stratLossesSE_AD = 0
	var stratNumAD = 0
	var mvsTriedAD = 0
	startTimeAD = time.Now()

	/*	var treeMoveLen int       //commented out to elim warning
		var treeVert string
		var treeHoriz string
		switch pMD.pType {
		case "TW":
			treeMoveLen = 8
			treeVert = vert8   // Looks Like: ->    ┃   <-
			treeHoriz = horiz8 // Looks Like: ->━━━━━━━━<-
		case "TS":
			treeMoveLen = 5
			treeVert = vert5   // Looks Like: ->  ┃  <-
			treeHoriz = horiz5 // Looks Like: ->━━━━━<-
		case "TSS":
			treeMoveLen = 3
			treeVert = vert3   // Looks Like: -> ┃ <-
			treeHoriz = horiz3 // Looks Like: ->━━━<-
			// Following used only for "TSS" so no generic equivalent variable is needed
			 const horiz3NewFirstStrat    // Looks Like: -> ┳━<-
			 const horiz3NewLastStrat     // Looks Like: -> ┗━<-
			 const horiz3NewMidStrat      // Looks Like: ->━┳━<-
		}
	*/

	// Setup pfmt to print thousands with commas
	var pfmt = message.NewPrinter(language.English)

	for deckNum := firstDeckNum; deckNum < (firstDeckNum + numberOfDecksToBePlayed); deckNum++ {

		startTimeTD = time.Now()
		moveNumMax = 0                  //to keep track of length of the longest strategy so far
		protoDeck, err := reader.Read() // protoDeck is a slice of strings: rank, suit, rank, suit, etc.
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Cannot read from inputFileName:", err)
		}

		if verbose > 1 {
			pfmt.Printf("\nDeck #%d:\n", deckNum)
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
		result1, result2, _ := playAllMoveS(b, 0, deckNum, cfg, varSp2PN, startTimeTD)

		if cfg.PlayNew.ReportingMoveByMove && cfg.PlayNew.ReportingType.Tree && cfg.PlayNew.TreeReportingOptions.Type != "narrow" {
			_, err = pfmt.Printf("\n\nDeck %v\n", deckNum)
			fmt.Printf("\n\n Strat #")
			if cfg.PlayNew.TreeReportingOptions.Type == "wide" {
				for i := 1; i <= 200; i++ {
					fmt.Printf("    %3v ", i)
				}
				fmt.Printf("\n")
			}
		}

		var dummy []move
		if stratWinsTD > 0 {
			deckWinsAD += 1
			prntMDet(b, dummy, 1, deckNum, 1, "NOTX", 2, "\n   DECK WON\n", "", "", cfg, varSp2PN)
			prntMDetTreeReturnComment("\n               DECK WON\n", deckNum, 0, cfg)
		} else {
			prntMDet(b, dummy, 1, deckNum, 1, "NOTX", 2, "\n   DECK LOST\n", "", "", cfg, varSp2PN)
			prntMDetTreeReturnComment("\n               DECK LOST\n", deckNum, 0, cfg)
			deckLossesAD += 1
		}

		/*
			if cfg.PlayNew.WinLossReport { // Deck Win Loss Summary Statistics
				if stratWinsTD == 0 {
					dWLDStats.winLoss = "L"
					dWLDStats.moveNumAt1stWinOrAtLoss = 0
					dWLDStats.moveNumMinWinIfFindAll = 0
					dWLDStats.moveNumMaxWinIfFindAll = 0
					dWLDStats.stratNumAt1stWinOrAtLoss = stratWinsTD
					dWLDStats.mvsTriedAt1stWinOrAtLoss = mvsTriedTD
					// Add maxUnique Boards
					//add max movenum
					// add ElapsedTimeAt1stWin
					dWLDStats.ElapsedTimeAt1stWinOrAtLoss = time.Now().Sub(startTimeTD)
				} else {
					dWLDStats.winLoss = "W"
					dWLDStats.ElapsedTimeAt1stWinOrAtLoss = time.Now().Sub(startTimeTD)
				}
				deckWinLossDetail = append(deckWinLossDetail, dWLDStats)

			}*/

		// This If Block is Print Only
		if cfg.PlayNew.ReportingType.DeckByDeck && cfg.PlayNew.DeckByDeckReportingOptions.Type == "regular" { // Deck-by-deck Statistics
			if stratWinsTD > 0 {
				fmt.Printf("\n\n*************************\n\nDeck: %d  WON    Result Codes: %v %v", deckNum, result1, result2)
			} else {
				fmt.Printf("\n\n*************************\n\nDeck: %d  LOST   Result Codes: %v %v", deckNum, result1, result2)
			}
			fmt.Printf("\nElapsed Time is %v.", time.Since(startTimeTD))
			fmt.Printf("\n\nStrategies:")
			_, err = pfmt.Printf("\n   Tried: %d", stratNumTD)
			_, err = pfmt.Printf("\n     Won: %d", stratWinsTD)
			_, err = pfmt.Printf("\n    Lost: %d", stratLossesTD)
			fmt.Printf("\n\nStrategies Lost Detail:")
			_, err = pfmt.Printf("\n   NMA: %d   (No Moves Available)", stratLossesNMA_TD)
			_, err = pfmt.Printf("\n    RB: %d   (Repetitive Board)", stratLossesRB_TD)
			_, err = pfmt.Printf("\n    SE: %d   (Strategy Exhausted)", stratLossesSE_TD)
			_, err = pfmt.Printf("\n   GML: %d   (Game Length Limit)", stratLossesGLE_TD)
			if stratLossesNMA_TD+stratLossesRB_TD+stratLossesSE_TD+stratLossesGLE_TD != stratLossesTD {
				fmt.Printf("\n     *********** Total Strategy Losses != Sum of strategy detail")
			}
			if stratLossesTD+stratWinsTD != stratNumTD {
				fmt.Printf("\n     *********** Strategies Tried != Strategies Lost + Strategies Won")
			}
			if cfg.PlayNew.FindAllWinStrats {
				fmt.Printf("\n\n Multiple Successful Strategies were found in some wining decks.")
				_, err = pfmt.Printf("   Total winning strategies found: %d\n", stratWinsTD)
			}
		}

		// This If Block is Print Only
		if cfg.PlayNew.ReportingType.DeckByDeck && cfg.PlayNew.DeckByDeckReportingOptions.Type != "regular" {
			var est time.Duration
			//                      nanosecondsTD   / Decks Played So Far         * remaining decks [remaining decks = numbertobeplayed - decksplayed so far
			est = time.Duration(float64(time.Since(startTimeAD))/float64(deckNum+1-firstDeckNum)*float64(numberOfDecksToBePlayed-(deckNum+1-firstDeckNum))) * time.Nanosecond
			/*			est = time.Duration(float64(elapsedTimeAD)           / float64(deckNum-firstDeckNum+1)*float64(numberOfDecksToBePlayed-(deckNum-firstDeckNum+1))) * time.Nanosecond
						fmt.Printf("       (deckNum+1-firstDeckNum)                                                           = %v\n", (deckNum + 1 - firstDeckNum))
						fmt.Printf("float64(deckNum+1-firstDeckNum)                                                           = %v\n", float64(deckNum+1-firstDeckNum))
						fmt.Printf("                                                                (deckNum+1-firstDeckNum)  = %v\n", (deckNum + 1 - firstDeckNum))
						fmt.Printf("                                        numberOfDecksToBePlayed                           = %v\n", numberOfDecksToBePlayed)
						fmt.Printf("                                       (numberOfDecksToBePlayed-(deckNum+1-firstDeckNum)) = %v\n", (numberOfDecksToBePlayed - (deckNum + 1 - firstDeckNum)))
						fmt.Printf("                                float64(numberOfDecksToBePlayed-(deckNum+1-firstDeckNum)) = %v\n", float64(numberOfDecksToBePlayed-(deckNum+1-firstDeckNum)))
						fmt.Printf("float64(deckNum+1-firstDeckNum)*float64(numberOfDecksToBePlayed-(deckNum+1-firstDeckNum)) = %v\n", float64(deckNum+1-firstDeckNum)*float64(numberOfDecksToBePlayed-(deckNum+1-firstDeckNum)))
						fmt.Printf("float64(deckNum+1-firstDeckNum)*float64(numberOfDecksToBePlayed-(deckNum+1-firstDeckNum)) = %v\n", float64(deckNum+1-firstDeckNum)*float64(numberOfDecksToBePlayed-(deckNum+1-firstDeckNum)))
						fmt.Printf("time.Since(startTimeAD)                                                                                                            = %v\n", time.Since(startTimeAD))
						fmt.Printf("time.Since(startTimeAD)                                                                                                            = %v\n", time.Since(startTimeAD))
						fmt.Printf("                          time.Duration(float64(deckNum+1-firstDeckNum)*float64(numberOfDecksToBePlayed-(deckNum+1-firstDeckNum))) = %v\n", time.Duration(float64(deckNum+1-firstDeckNum)*float64(numberOfDecksToBePlayed-(deckNum+1-firstDeckNum))))
						fmt.Printf("time.Since(startTimeAD) / time.Duration(float64(deckNum+1-firstDeckNum)*float64(numberOfDecksToBePlayed-(deckNum+1-firstDeckNum))) = %v\n", time.Since(startTimeAD)/time.Duration(float64(deckNum+1-firstDeckNum)*float64(numberOfDecksToBePlayed-(deckNum+1-firstDeckNum))))
			*/
			wL := ""
			if stratWinsTD > 0 {
				wL = "WON " // Note additional space -- for alignment
			} else {
				wL = "LOST"
			}
			elTimeSinceStartTimeADFormatted := time.Since(startTimeAD).Truncate(100 * time.Millisecond).String()
			if time.Since(startTimeAD) > time.Duration(5*time.Minute) {
				elTimeSinceStartTimeADFormatted = time.Since(startTimeAD).Truncate(time.Second).String()
			}
			pfmt.Printf("Dk: %5d   "+wL+"   MvsTried: %13v   MoveNum: xxx   Max MoveNum: xxx   StratsTried: %12v   UnqBoards: %11v   Won: %5v   Lost: %5v   GLE: %5v   Won: %5.1f%%   Lost: %5.1f%%   GLE: %5.1f%%   ElTime TD: %9s   ElTime ADs: %9s  Rem Time: %11s   ResCodes: %2s %3s   Time Now: %8s\n", deckNum, mvsTriedTD /*moveNum, maxMoveNum, */, stratNumTD, len(varSp2PN.priorBoards), deckWinsAD, deckLossesAD, stratLossesGLE_AD, roundFloatIntDiv(deckWinsAD*100, deckNum+1-firstDeckNum, 1), roundFloatIntDiv(deckLossesAD*100, deckNum+1-firstDeckNum, 1), roundFloatIntDiv(stratLossesGLE_AD*100, deckNum+1-firstDeckNum, 1), time.Since(startTimeTD).Truncate(100*time.Millisecond).String(), elTimeSinceStartTimeADFormatted, est.Truncate(time.Second).String(), result1, result2, time.Now().Format(" 3:04 pm"))
		}

		// Verbose Special "BELL" Starts Here - No effect on operation
		/*if strings.Contains(verboseSpecial, ";BELL;") && time.Since(startTimeTD) > 1*time.Second { // changed 5*time.Minute to 1*time.Second for test change it back
			fmt.Printf("\a") // Ring Bell
		}*/
		// Verbose Special "BELL" Ends Here - No effect on operation

		stratWinsAD += stratWinsTD
		stratWinsTD = 0
		stratLossesAD += stratLossesTD
		stratLossesTD = 0
		stratLossesGLE_AD += stratLossesGLE_TD
		stratLossesGLE_TD = 0
		stratLossesNMA_AD += stratLossesNMA_TD
		stratLossesNMA_TD = 0
		stratLossesRB_AD += stratLossesRB_TD
		stratLossesRB_TD = 0
		stratLossesSE_AD += stratLossesSE_TD
		stratLossesSE_TD = 0
		stratNumAD += stratNumTD + 1 // Because we start at strategy 0 which is all best moves
		stratNumTD = 0
		mvsTriedAD += mvsTriedTD + 1
		mvsTriedTD = 0
		varSp2PN.treePrevMovesTD = ""
		//treePrevMovesTD = ""
		clear(varSp2PN.priorBoards)
	}

	// At this point, all decks to be played have been played.  Time to report aggregate won loss.
	// From this point on, the program only prints.

	fmt.Printf("\n******************\n\n" + "Decks:")
	_, err = pfmt.Printf("\n   Played: %d", numberOfDecksToBePlayed)
	_, err = pfmt.Printf("\n      Won: %d", deckWinsAD)
	_, err = pfmt.Printf("\n     Lost: %d", deckLossesAD)
	averageElapsedTimePerDeck := time.Duration(float64(time.Since(startTimeAD)) / float64(numberOfDecksToBePlayed))
	fmt.Printf("\nElapsed Time is %v.", time.Since(startTimeAD))
	fmt.Printf("\nAverage Elapsed Time per Deck is %s", averageElapsedTimePerDeck.Truncate(100*time.Millisecond).String())
	fmt.Printf("\n\nStrategies:")
	_, err = pfmt.Printf("\n   Tried: %d", stratNumAD)
	_, err = pfmt.Printf("\n     Won: %d", stratLossesAD)
	_, err = pfmt.Printf("\n    Lost: %d", stratWinsAD)
	fmt.Printf("\n\nStrategies Lost Detail:")
	_, err = pfmt.Printf("\n   NMA: %d   (No Moves Available)", stratLossesNMA_AD)
	_, err = pfmt.Printf("\n    RB: %d   (Repetitive Board)", stratLossesRB_AD)
	_, err = pfmt.Printf("\n    SE: %d   (Strategy Exhausted)", stratLossesSE_AD)
	_, err = pfmt.Printf("\nStrategy Losses at Game Length Limit is: %d", stratLossesGLE_AD)
	if stratLossesNMA_AD+stratLossesRB_AD+stratLossesSE_AD+stratLossesGLE_AD != stratLossesAD {
		fmt.Printf("\n     *********** Total Strategy Losses != Sum of strategy detail")
	}
	if stratLossesAD+stratWinsAD != stratNumAD {
		fmt.Printf("\n     *********** Strategies Tried != Strategies Lost + Strategies Won")
	}
	if cfg.PlayNew.FindAllWinStrats {
		fmt.Printf("\n\n Multiple Successful Strategies were found in some winng decks.")
		_, err = pfmt.Printf("   Decks Won: %d\n", deckWinsAD)
		_, err = pfmt.Printf("   Total winning strategies found: %d\n", stratWinsAD)
		_, err = pfmt.Printf("   Average winning strategies found: %d\n", stratWinsAD/deckWinsAD)
	}

	if cfg.PlayNew.WinLossReport { // Deck Win Loss Summary Statistics
		fmt.Printf("\n\n\n Deck-by Deck Win/Loss Detail   (Copy to Excel to get headings to line up with the columns)")
		fmt.Printf("\n\n Deck\tW/L\tMoveNum 1ST-Win\tStratNum At 1st-Win Or At-Loss\tMvsTried At 1st-Win Or At-Loss\tMoveNum Min-Win If-Find-All\tMoveNum Max-Win If-Find-All\tElapsed Time At 1st-Win Or At Loss\n")
		/*for dN, detail := range deckWinLossDetail {
			_, err = pfmt.Printf("\n  %5v\t  %v\t%4v\t%8v\t%8v\t%4v\t%4v", dN, detail.winLoss, detail.moveNumAt1stWinOrAtLoss, detail.stratNumAt1stWinOrAtLoss, detail.mvsTriedAt1stWinOrAtLoss, detail.moveNumMinWinIfFindAll, detail.moveNumMaxWinIfFindAll, detail.ElapsedTimeAt1stWinOrAtLoss)
		}*/
	}

}

// Divide 2 integers and round to precision digits
func roundFloatIntDiv(numer int, denom int, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(float64(numer)/float64(denom)*ratio) / ratio
}
