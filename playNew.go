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

var aMwinCounterThisDeck = 0
var aMearlyWinCounterThisDeck = 0
var aMstandardWinCounterThisDeck = 0
var aMlossCounterThisDeck = 0

// var aMStratlossesAtGLLThisDeck = 0
var aMStratlossesAtNoMovesThisDeck = 0
var aMStratlossesAtRepMveThisDeck = 0
var aMStratlossesExhaustedThisDeck = 0
var aMStratNumThisDeck = 0
var aMmvsTriedThisDeck = 0

type boardInfo struct {
	mN           int
	aMmvsTriedTD int
	// May add linked list and stats later
}

var priorBoards = make(map[string]boardInfo)

func playNew(reader csv.Reader) {
	if verbose > 1 {
		fmt.Printf("firstDeckNum: %v, numberOfDecksToBePlayed: %v, verbose: %v, verboseSpecial: %v, findAllSuccessfulStrategies: %v, printTree: %v, reader: %v", firstDeckNum, numberOfDecksToBePlayed, verbose, findAllSuccessfulStrategies, printTree, reader)
	}

	var aMwinCounterAllDecks = 0
	var aMearlyWinCounterAllDecks = 0
	var aMstandardWinCounterAllDecks = 0
	var aMlossCounterAllDecks = 0
	// var aMStratlossesAtGLLAllDecks = 0
	var aMStratlossesAtNoMovesAllDecks = 0
	var aMStratlossesAtRepMveAllDecks = 0
	var aMStratlossesExhaustedAllDecks = 0
	var aMStratTriedAllDecks int = 0
	var MovesTriedAllDecks int = 0

	aMstartTimeAllDecks := time.Now()
	aMstartTimeThisDeck := time.Now()

	for deckNum := firstDeckNum; deckNum < (firstDeckNum + numberOfDecksToBePlayed); deckNum++ {

		// Verbose Special Starts Here - No effect on operation
		if strings.Contains(verboseSpecial, "D") {
			aMstartTimeThisDeck = time.Now()
		}
		// Verbose Special Ends Here - No effect on operation

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

		// Verbose Special Starts Here - No effect on operation
		if strings.Contains(verboseSpecial, "D") { // Deck Statistics
			fmt.Printf("\nDeck: #%d  Result: %v\n", deckNum, result1, result2)
			//lossCounter := 1 - aMwinCounterAllDecks
			endTime := time.Now()
			elapsedTime := endTime.Sub(aMstartTimeThisDeck)
			//fmt.Printf("Strategy Losses at Game Length Limit is: %d\n", aMStratlossesAtGLLThisDeck)
			fmt.Printf("Strategies Played %d\n", aMStratNumThisDeck+1)
			fmt.Printf("     Strategy Losses at No Moves Available is %d\n", aMStratlossesAtNoMovesThisDeck)
			fmt.Printf("     Strategy Losses at Repetitive Move is %d\n", aMStratlossesAtRepMveThisDeck)
			fmt.Printf("     Strategy Losses at Moves Exhausted is %d\n", aMStratlossesExhaustedThisDeck)
			fmt.Printf("     Total Strategy Losses %d + Games Won %d = %d",
				aMStratlossesAtNoMovesThisDeck+aMStratlossesAtRepMveThisDeck+aMStratlossesExhaustedThisDeck,
				aMwinCounterThisDeck,
				aMStratlossesAtNoMovesThisDeck+aMStratlossesAtRepMveThisDeck+aMStratlossesExhaustedThisDeck+aMwinCounterThisDeck)
			fmt.Printf("          Should equal Strategies Played %d\n", aMStratNumThisDeck)
			fmt.Printf("Elapsed Time is %v.\n", elapsedTime)
		}
		// Verbose Special Ends Here - No effect on operation

		aMwinCounterAllDecks += aMwinCounterThisDeck
		aMwinCounterThisDeck = 0
		aMearlyWinCounterAllDecks += aMearlyWinCounterThisDeck
		aMearlyWinCounterThisDeck = 0
		aMstandardWinCounterAllDecks += aMstandardWinCounterThisDeck
		aMstandardWinCounterThisDeck = 0
		aMlossCounterAllDecks += aMlossCounterThisDeck
		aMlossCounterThisDeck = 0
		// aMStratlossesAtGLLAllDecks += aMStratlossesAtGLLThisDeck
		//aMStratlossesAtGLLThisDeck = 0
		aMStratlossesAtNoMovesAllDecks += aMStratlossesAtNoMovesThisDeck
		aMStratlossesAtNoMovesThisDeck = 0
		aMStratlossesAtRepMveAllDecks += aMStratlossesAtRepMveThisDeck
		aMStratlossesAtRepMveThisDeck = 0
		aMStratlossesExhaustedAllDecks += aMStratlossesExhaustedThisDeck
		aMStratlossesExhaustedThisDeck = 0
		aMStratTriedAllDecks += (aMStratNumThisDeck + 1) // Because we start at strategy 0 which is all best moves
		aMStratNumThisDeck = 0
		MovesTriedAllDecks += aMmvsTriedThisDeck
		aMmvsTriedThisDeck = 0
		clear(priorBoards)
	}
	lossCounter := numberOfDecksToBePlayed - aMwinCounterAllDecks
	endTime := time.Now()
	elapsedTime := endTime.Sub(aMstartTimeAllDecks)
	fmt.Printf("\n\n******************\n\nNumber of Decks Played is: %d.\n", numberOfDecksToBePlayed)
	fmt.Printf("Total Decks Won is: %d of which: %d were Early Wins and %d were Standard Wins\n", aMwinCounterAllDecks, aMearlyWinCounterAllDecks, aMstandardWinCounterAllDecks)
	fmt.Printf("Total Decks Lost is: %d which should equal Counted losses: %d\n", lossCounter, aMlossCounterThisDeck)
	//fmt.Printf("Strategy Losses at Game Length Limit is: %d\n", aMStratlossesAtGLL)
	fmt.Printf("Strategies Played %d\n", aMStratTriedAllDecks)
	fmt.Printf("     Strategy Losses at No Moves Available is %d\n", aMStratlossesAtNoMovesAllDecks)
	fmt.Printf("     Strategy Losses at Repetitive Move is %d\n", aMStratlossesAtRepMveAllDecks)
	fmt.Printf("     Strategy Losses at Moves Exhausted is %d\n", aMStratlossesExhaustedAllDecks)
	fmt.Printf("     Total Strategy Losses %d + Games Won %d = %d",
		aMStratlossesAtNoMovesAllDecks+aMStratlossesAtRepMveAllDecks+aMStratlossesExhaustedAllDecks,
		aMwinCounterThisDeck,
		aMStratlossesAtNoMovesAllDecks+aMStratlossesAtRepMveAllDecks+aMStratlossesExhaustedAllDecks+aMwinCounterThisDeck)
	fmt.Printf("          Should equal Strategies Played %d\n", aMStratTriedAllDecks)
	averageElapsedTimePerDeck := float64(elapsedTime.Microseconds()) / float64(numberOfDecksToBePlayed)
	fmt.Printf("Elapsed Time is %v.\n", elapsedTime)
	fmt.Printf("Average Elapsed Time per Deck is %vus.\n", averageElapsedTimePerDeck)
}
