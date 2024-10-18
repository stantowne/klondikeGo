package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strconv"
)

func playNew(reader csv.Reader) {
	if verbose > 1 {
		fmt.Printf("firstDeckNum: %v, numberOfDecksToBePlayed: %v, verbose: %v, verboseSpecial: %v, findAllSuccessfulStrategies: %v, printTree: %v, reader: %v", firstDeckNum, numberOfDecksToBePlayed, verbose, findAllSuccessfulStrategies, printTree, reader)
	}
	/*	startTime := time.Now()
		winCounter := 0
		earlyWinCounter := 0
		attemptsAvoidedCounter := 0
		lossesAtGLL := 0
		lossesAtNoMoves := 0
		regularLosses := 0

	*/
	var AllMvStratNumAllDecks int = 0
	var MovesTriedAllDecks int = 0

	for deckNum := firstDeckNum; deckNum < (firstDeckNum + numberOfDecksToBePlayed); deckNum++ {
		AllMvStratNumThisDeck = 0
		MovesTriedThisDeck = 0
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
		AllMvStratNumAllDecks += AllMvStratNumThisDeck
		AllMvStratNumThisDeck = 0
		MovesTriedAllDecks += MovesTriedThisDeck
		MovesTriedThisDeck = 0
	}
}
