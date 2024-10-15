package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

func main() {

	firstDeckNum, _ := strconv.Atoi(args[1])
	numberOfDecksToBePlayed, _ := strconv.Atoi(args[2])
	length, _ := strconv.Atoi(args[3]) //length of each strategy (which also determines the # of strategies - 2^n)
	// open the deck file here in main - will be read in both playOrig and PlayNew
	inputFileName := "decks-made-2022-01-15_count_10000-dict.csv"
	file, err := os.Open(inputFileName)
	if err != nil {
		log.Println("Cannot open inputFileName:", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			println("could not close file:", err)
		}
	}(file)
	reader := csv.NewReader(file)

	// playOrig will execute the original code designed to play either the "Best" move or the best move modified by the IOS strategy
	// of substituting FlipToWaste as described in more detail below.  This was formally known as the "playBestOrIOS" strategy
	// and developed under a function of that name in project branch "tree".
	// To avoid issues with old "tree" branch code the function playOrig has been created by refactoring and adding passed arguments.
	//
	playOrig(firstDeckNum, numberOfDecksToBePlayed, length, reader)
	//testCardPackUnPack(os.Args)
	//testBoardCodeDeCode(os.Args)
}
