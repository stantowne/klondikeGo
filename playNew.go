package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strconv"
)

func playNew(firstDeckNum int, numberOfDecksToBePlayed int, verbose int, findAllSuccessfulStrategies bool, printTree string, reader csv.Reader) {
	if verbose > 1 {
		fmt.Printf("firstDeckNum: %v, numberOfDecksToBePlayed: %v, verbose: %v, findAllSuccessfulStrategies: %v, printTree: %v, reader: %v", firstDeckNum, numberOfDecksToBePlayed, verbose, findAllSuccessfulStrategies, printTree, reader)
	}

	const gameLengthLimitNew = 150 // max moveCounter
	/*	AllMvStratNumAllDecks := 0
		startTime := time.Now()
		winCounter := 0
		earlyWinCounter := 0
		attemptsAvoidedCounter := 0
		lossesAtGLL := 0
		lossesAtNoMoves := 0
		regularLosses := 0

		AllMvStratNumAllDecks = AllMvStratNumAllDecks + AllMvStratNum
		AllMvStratNum = 0
	*/

	for deckNum := firstDeckNum; deckNum < (firstDeckNum + numberOfDecksToBePlayed); deckNum++ {
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
		AllMvStratNum := 0
		var b = dealDeck(d)
		firstMoveNull := move{}
		playAllMoveS(b, firstMoveNull, 0, deckNum, verbose, findAllSuccessfulStrategies, printTree, AllMvStratNum)

	}
}
