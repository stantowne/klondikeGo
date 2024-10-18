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

var aMstartTime = time.Now()
var aMwinCounter = 0
var aMearlyWinCounter = 0
var aMlossesAtGLL = 0
var aMlossesAtNoMoves = 0
var aMregularLosses = 0
var aMlossesAtLoop = 0
var AmStratNumThisDeck int
var aMmvsTriedThisDeck int
var priorBoards = make(map[string]bool)

func playNew(reader csv.Reader) {
	if verbose > 1 {
		fmt.Printf("firstDeckNum: %v, numberOfDecksToBePlayed: %v, verbose: %v, verboseSpecial: %v, findAllSuccessfulStrategies: %v, printTree: %v, reader: %v", firstDeckNum, numberOfDecksToBePlayed, verbose, findAllSuccessfulStrategies, printTree, reader)
	}

	var AllMvStratNumAllDecks int = 0
	var MovesTriedAllDecks int = 0
	aMstartTime := time.Now()

	for deckNum := firstDeckNum; deckNum < (firstDeckNum + numberOfDecksToBePlayed); deckNum++ {
		AmStratNumThisDeck = 1
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
		AllMvStratNumAllDecks += AmStratNumThisDeck
		AmStratNumThisDeck = 0
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
	_, err = p.Printf("Total Decks Won is %d of which %d were Early Wins\n", aMwinCounter, aMearlyWinCounter)
	if err != nil {
		fmt.Println("Total Decks Won cannot print")
	}
	_, err = p.Printf("Total Decks Lost is %d\n", lossCounter)
	if err != nil {
		fmt.Println("Total Decks Lost cannot print")
	}
	_, err = p.Printf("Losses at Game Length Limit is %d\n", aMlossesAtGLL)
	if err != nil {
		fmt.Println("Losses at Game Length Limit cannot print")
	}
	_, err = p.Printf("Losses at No Moves Available is %d\n", aMlossesAtNoMoves)
	if err != nil {
		fmt.Println("Losses at No Moves Available cannot print")
	}
	_, err = p.Printf("Regular Losses is %d\n", aMregularLosses)
	if err != nil {
		fmt.Println("Regular Losses cannot print")
	}
	averageElapsedTimePerDeck := float64(elapsedTime.Milliseconds()) / float64(numberOfDecksToBePlayed)
	fmt.Printf("Elapsed Time is %v.\n", elapsedTime)
	fmt.Printf("Average Elapsed Time per Deck is %fms.\n", averageElapsedTimePerDeck)
}
