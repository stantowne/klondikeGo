package main

import (
	"encoding/csv"
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"io"
	"log"
	"strconv"
	"time"
)

var aMwinCounter = 0
var aMearlyWinCounter = 0
var aMstandardWinCounter = 0
var aMlossCounter = 0

// var aMStratlossesAtGLL = 0
var aMStratlossesAtNoMovesThisDeck = 0
var aMStratlossesAtRepMveThisDeck = 0
var aMStratlossesExhaustedThisDeck = 0
var aMStratNumThisDeck int
var aMmvsTriedThisDeck int

var priorBoards = make(map[string]bool)

func playNew(reader csv.Reader) {
	if verbose > 1 {
		fmt.Printf("firstDeckNum: %v, numberOfDecksToBePlayed: %v, verbose: %v, verboseSpecial: %v, findAllSuccessfulStrategies: %v, printTree: %v, reader: %v", firstDeckNum, numberOfDecksToBePlayed, verbose, findAllSuccessfulStrategies, printTree, reader)
	}

	var aMStratlossesAtNoMovesAllDecks = 0
	var aMStratlossesAtRepMveAllDecks = 0
	var aMStratlossesExhaustedAllDecks = 0
	var aMStratNumAllDecks int = 0
	var MovesTriedAllDecks int = 0

	aMstartTime := time.Now()

	for deckNum := firstDeckNum; deckNum < (firstDeckNum + numberOfDecksToBePlayed); deckNum++ {
		aMStratNumThisDeck = 1
		aMmvsTriedThisDeck = 0
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

		firstMoveNull := move{}

		playAllMoveS(b, firstMoveNull, 0, deckNum)

		aMStratlossesAtNoMovesAllDecks += aMStratlossesAtNoMovesThisDeck
		aMStratlossesAtNoMovesThisDeck = 0
		aMStratlossesAtRepMveAllDecks += aMStratlossesAtRepMveThisDeck
		aMStratlossesAtRepMveThisDeck = 0
		aMStratlossesExhaustedAllDecks += aMStratlossesExhaustedThisDeck
		aMStratlossesExhaustedThisDeck = 0
		aMStratNumAllDecks += aMStratNumThisDeck
		aMStratNumThisDeck = 0
		MovesTriedAllDecks += aMmvsTriedThisDeck
		aMmvsTriedThisDeck = 0
		clear(priorBoards)
	}
	lossCounter := numberOfDecksToBePlayed - aMwinCounter
	endTime := time.Now()
	elapsedTime := endTime.Sub(aMstartTime)
	var p = message.NewPrinter(language.English)
	_, err = p.Printf("\nNumber of Decks Played is %d.\n", numberOfDecksToBePlayed)
	if err != nil {
		fmt.Println("Number of Decks Played cannot print")
	}
	_, err = p.Printf("Total Decks Won is %d of which %d were Early Wins and %d were Standard Wins\n", aMwinCounter, aMearlyWinCounter, aMstandardWinCounter)
	if err != nil {
		fmt.Println("Total Decks Won cannot print")
	}
	_, err = p.Printf("Total Decks Lost is %d   Counted losses are %d should equal.\n", lossCounter, aMlossCounter)
	if err != nil {
		fmt.Println("Total Decks Lost cannot print")
	}
	//	_, err = p.Printf("Strategy Losses at Game Length Limit is %d\n", aMStratlossesAtGLL)
	//	if err != nil {
	//		fmt.Println("Strategy Losses at Game Length Limit cannot print")
	//	}
	_, err = p.Printf("Strategy Losses at No Moves Available is %d\n", aMStratlossesAtNoMovesAllDecks)
	if err != nil {
		fmt.Println("Strategy Losses at No Moves Available cannot print")
	}
	_, err = p.Printf("Strategy Losses at Repetitive Move is %d\n", aMStratlossesAtRepMveAllDecks)
	if err != nil {
		fmt.Println("Strategy Losses at Repetitive Move cannot print")
	}
	_, err = p.Printf("Strategy Losses at Repetitive Move is %d\n", aMStratlossesExhaustedAllDecks)
	if err != nil {
		fmt.Println("Strategy Losses at Repetitive Move cannot print")
	}
	averageElapsedTimePerDeck := float64(elapsedTime.Milliseconds()) / float64(numberOfDecksToBePlayed)
	fmt.Printf("Elapsed Time is %v.\n", elapsedTime)
	fmt.Printf("Average Elapsed Time per Deck is %fms.\n", averageElapsedTimePerDeck)
}
