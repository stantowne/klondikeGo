package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
)

//var priorPause = 0 // Remove Pauser

/*
// These are used in the Tree printing subroutine pmd(......) in playAllMoves    //commented out to elim warning
const vert1 = string('\u2503')          // Looks Like: ->┃<-
const horiz1 = string('\u2501')         // Looks Like: ->━<-
const vert5 = "  " + vert1 + "  "       // Looks Like: ->  ┃  <-
const vert8 = "  " + vert5 + " "        // Looks Like: ->    ┃   <-
const horiz5 = horiz3 + horiz1 + horiz1 // Looks Like: ->━━━━━<-
const horiz8 = horiz3 + horiz5          // Looks Like: ->━━━━━━━━<-

const vert3 = " " + vert1 + " "              // Looks Like: -> ┃ <-
const horiz3 = horiz1 + horiz1 + horiz1      // Looks Like: ->━━━<-
const horiz1NewFirstStrat = string('\u2533') // Looks Like: ->┳<-
const horiz3NewFirstStrat = horiz1 + horiz1NewFirstStrat + horiz1 // Looks Like: ->━┳━<-
const horiz1NewStratLastStrat = string('\u2517') // Looks Like: ->┗<-
const horiz3NewLastStrat = " " + horiz1NewStratLastStrat + horiz1 // Looks Like: -> ┗━<-
const horiz1NewMidStrat = string('\u2523') // Looks Like: ->┣<-
const horiz3NewMidStrat = horiz1 + horiz1NewMidStrat + horiz1     // Looks Like: -> ┣━<-
*/
var stratWinsTD = 0
var stratLossesTD = 0

var stratLossesGML_TD = 0
var stratLossesNMA_TD = 0
var stratLossesRB_TD = 0
var stratLossesSE_TD = 0
var stratNumTD = 0
var mvsTriedTD = 0
var startTimeAD = time.Now()
var startTimeTD = time.Now()

type boardInfo struct {
	mN           int
	aMmvsTriedTD int
	// May add linked list and stats later
}

type deckWinLossDetailStats struct {
	winLoss                     string
	moveNumAt1stWinOrAtLoss     int
	moveNumMinWinIfFindAll      int
	moveNumMaxWinIfFindAll      int
	stratNumAt1stWinOrAtLoss    int
	mvsTriedAt1stWinOrAtLoss    int
	ElapsedTimeAt1stWinOrAtLoss time.Duration
}

var dWLDStats deckWinLossDetailStats
var deckWinLossDetail []deckWinLossDetailStats

var priorBoards = make(map[string]boardInfo)

func playNew(reader csv.Reader) {

	var deckWinsAD = 0
	var deckLossesAD = 0
	var stratWinsAD = 0
	var stratLossesAD = 0
	var stratLossesGML_AD = 0
	var stratLossesNMA_AD = 0
	var stratLossesRB_AD = 0
	var stratLossesSE_AD = 0
	var stratNumAD int = 0
	var mvsTriedAD int = 0

	/*	var treeMoveLen int       //commented out to elim warning
		var treeVert string
		var treeHoriz string
		switch printMoveDetail.pType {
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

	for deckNum := firstDeckNum; deckNum < (firstDeckNum + numberOfDecksToBePlayed); deckNum++ {

		startTimeTD = time.Now()

		protoDeck, err := reader.Read() // protoDeck is a slice of strings: rank, suit, rank, suit, etc.
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Cannot read from inputFileName:", err)
		}

		if verbose > 1 {
			fmt.Printf("\nDeck #%d:\n", deckNum)
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

		if printMoveDetail.pType == "TW" || printMoveDetail.pType == "TS" {
			fmt.Printf("\n\nDeck %v\n", deckNum)
			fmt.Printf("\n\n Strat #")
			if printMoveDetail.pType == "TW" {
				for i := 1; i <= 200; i++ {
					fmt.Printf("    %3v ", i)
				}
				fmt.Printf("\n")
			}
		}

		result1, result2 := playAllMoveS(b, 0, deckNum)

		// Verbose Special "WL" Starts Here - No effect on operation
		if strings.Contains(verboseSpecial, "/WL/") { // Deck Win Loss Summary Statistics
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

		}
		// Verbose Special "WL" Ends Here - No effect on operation

		//Verbose Special "DBD" Starts Here - No effect on operation
		if strings.Contains(verboseSpecial, "/DBD/") { // Deck-by-deck Statistics
			if stratWinsTD > 0 {
				fmt.Printf("\n\n*************************\n\nDeck: %d  WON    Result Codes: %v %v", deckNum, result1, result2)
			} else {
				fmt.Printf("\n\n*************************\n\nDeck: %d  LOST   Result Codes: %v %v", deckNum, result1, result2)
			}
			endTime := time.Now()
			elapsedTime := endTime.Sub(startTimeTD)
			fmt.Printf("\nElapsed Time is %v.", elapsedTime)
			fmt.Printf("\n\nStrategies:")
			fmt.Printf("\n   Tried: %d", stratNumTD)
			fmt.Printf("\n     Won: %d", stratLossesTD)
			fmt.Printf("\n    Lost: %d", stratWinsTD)
			fmt.Printf("\n\nStrategies Lost Detail:")
			fmt.Printf("\n   NMA: %d   (No Moves Available)", stratLossesNMA_TD)
			fmt.Printf("\n    RB: %d   (Repetitive Board)", stratLossesRB_TD)
			fmt.Printf("\n    SE: %d   (Strategy Exhausted)", stratLossesSE_TD)
			fmt.Printf("Strategy Losses at Game Length Limit is: %d\n", stratLossesGML_TD)
			if stratLossesNMA_TD+stratLossesRB_TD+stratLossesSE_TD+stratLossesGML_TD != stratLossesTD {
				fmt.Printf("\n     *********** Total Strategy Losses != Sum of strategy detail")
			}
			if stratLossesTD+stratWinsTD != stratNumTD {
				fmt.Printf("\n     *********** Strategies Tried != Strategies Lost + Strategies Won")
			}
			if findAllSuccessfulStrategies {
				fmt.Printf("\n\n Multiple Successful Startegies were found in some wining decks.")
				fmt.Printf("   Total winning strategies found: %d\n", stratWinsTD)
			}
		}
		// Verbose Special "DBD" Ends Here - No effect on operation

		//Verbose Special "DBDS" Starts Here - No effect on operation
		if strings.Contains(verboseSpecial, "/DBDS/") { // Deck-by-deck SHORT Statistics
			var est time.Duration
			est = time.Duration(int64(time.Now().Sub(startTimeAD))/(int64(deckNum)+int64(1))*(int64(printMoveDetail.deckContinueFor)-(int64(deckNum)+int64(1)))) * time.Nanosecond

			if stratWinsTD > 0 {
				fmt.Printf("\nDeck: %5d  WON    Result Codes: %v %v   Moves Tried: %9v   Elapsed Time This Deck: %v   Elapsed Time All Decks: %v  Est. Remaining Time: %v", deckNum, result1, result2, mvsTriedTD, time.Now().Sub(startTimeTD), time.Now().Sub(startTimeAD), est)
			} else {
				fmt.Printf("\nDeck: %5d  LOST   Result Codes: %v %v   Moves Tried: %9v   Elapsed Time This Deck: %v   Elapsed Time All Decks: %v  Est. Remaining Time: %v", deckNum, result1, result2, mvsTriedTD, time.Now().Sub(startTimeTD), time.Now().Sub(startTimeAD), est)
			}
		}
		// Verbose Special "DBDS" Ends Here - No effect on operation

		if stratWinsTD > 0 {
			deckWinsAD += 1
		} else {
			deckLossesAD += 1
		}
		stratWinsAD += stratWinsTD
		stratWinsTD = 0
		stratLossesAD += stratLossesTD
		stratLossesTD = 0
		stratLossesGML_AD += stratLossesGML_TD
		stratLossesGML_TD = 0
		stratLossesNMA_AD += stratLossesNMA_TD
		stratLossesNMA_TD = 0
		stratLossesRB_AD += stratLossesRB_TD
		stratLossesRB_TD = 0
		stratLossesSE_AD += stratLossesSE_TD
		stratLossesSE_TD = 0
		stratNumAD += (stratNumTD + 1) // Because we start at strategy 0 which is all best moves
		stratNumTD = 0
		mvsTriedAD += mvsTriedTD + 1
		mvsTriedTD = 0
		clear(priorBoards)
	}

	endTime := time.Now()
	elapsedTime := endTime.Sub(startTimeAD)
	fmt.Printf("\n******************\n\n" + "Decks:")
	fmt.Printf("\n   Played: %d", numberOfDecksToBePlayed)
	fmt.Printf("\n      Won: %d", deckWinsAD)
	fmt.Printf("\n     Lost: %d", deckLossesAD)
	averageElapsedTimePerDeck := float64(elapsedTime.Microseconds()) / float64(numberOfDecksToBePlayed)
	fmt.Printf("\nElapsed Time is %v.", elapsedTime)
	fmt.Printf("\nAverage Elapsed Time per Deck is %v    us.", averageElapsedTimePerDeck)
	fmt.Printf("\n\nStrategies:")
	fmt.Printf("\n   Tried: %d", stratNumAD)
	fmt.Printf("\n     Won: %d", stratLossesAD)
	fmt.Printf("\n    Lost: %d", stratWinsAD)
	fmt.Printf("\n\nStrategies Lost Detail:")
	fmt.Printf("\n   NMA: %d   (No Moves Available)", stratLossesNMA_AD)
	fmt.Printf("\n    RB: %d   (Repetitive Board)", stratLossesRB_AD)
	fmt.Printf("\n    SE: %d   (Strategy Exhausted)", stratLossesSE_AD)
	fmt.Printf("\nStrategy Losses at Game Length Limit is: %d", stratLossesGML_AD)
	if stratLossesNMA_AD+stratLossesRB_AD+stratLossesSE_AD+stratLossesGML_AD != stratLossesAD {
		fmt.Printf("\n     *********** Total Strategy Losses != Sum of strategy detail")
	}
	if stratLossesAD+stratWinsAD != stratNumAD {
		fmt.Printf("\n     *********** Strategies Tried != Strategies Lost + Strategies Won")
	}
	if findAllSuccessfulStrategies {
		fmt.Printf("\n\n Multiple Successful Startegies were found in some winng decks.")
		fmt.Printf("   Decks Won: %d\n", deckWinsAD)
		fmt.Printf("   Total winning strategies found: %d\n", stratWinsAD)
		fmt.Printf("   Average winning strategies found: %d\n", stratWinsAD/deckWinsAD)
	}
	// Verbose Special "WL" Starts Here - No effect on operation
	if strings.Contains(verboseSpecial, "/WL/") { // Deck Win Loss Summary Statistics
		fmt.Printf("\n\n\n Deck-by Deck Win/Loss Detail   (Copy to Excel to get headings to line up with the columns)")
		fmt.Printf("\n\n Deck\tW/L\tMoveNum 1ST-Win\tStratNum At 1st-Win Or At-Loss\tMvsTried At 1st-Win Or At-Loss\tMoveNum Min-Win If-Find-All\tMoveNum Max-Win If-Find-All\tElapsed Time At 1st-Win Or At Loss\n")
		for dN, detail := range deckWinLossDetail {
			fmt.Printf("\n  %5v\t  %v\t%4v\t%8v\t%8v\t%4v\t%4v", dN, detail.winLoss, detail.moveNumAt1stWinOrAtLoss, detail.stratNumAt1stWinOrAtLoss, detail.mvsTriedAt1stWinOrAtLoss, detail.moveNumMinWinIfFindAll, detail.moveNumMaxWinIfFindAll, detail.ElapsedTimeAt1stWinOrAtLoss)
		}
	}
	// Verbose Special "WL" Ends Here - No effect on operation
}
