package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const gameLengthLimitOrig = 150 // max moveCounter; increasing to 200 does not increase win rate
const gameLengthLimitNew = 150  // max moveNum
var gameLengthLimit int

var firstDeckNum int
var numberOfDecksToBePlayed int
var length int
var verbose int
var findAllSuccessfulStrategies bool
var printTree string

var err error
var singleGame = true
var AllMvStratNum int

func main() {
	/*

		Command line arguments

		args[0] = program name
		args[1] = # of the first deck to be used from within the pre-stored decks used to standardize testing
		args[2] = # of decks to be played
		args[3] = length of IOS strategy - see comments below	(applicable to playOrig only)
		args[4] = verbosity switch for messages
		args[5] = findAllSuccessfulStrategies 					(applicable playNew only)
		args[6] = printTree                   					(applicable playNew only)
	*/

	// Always a good idea to print out the program source of the output.
	//    Will be especially useful when we figure out how to include the versioning

	args := os.Args

	println("Calling Program:  ", args[0])

	firstDeckNum, err = strconv.Atoi(args[1])
	if err != nil || firstDeckNum < 0 || firstDeckNum > 9999 {
		println("first argument invalid")
		println("firstDeckNum must be non-negative integer less than 10,000")
		os.Exit(1)
	}

	numberOfDecksToBePlayed, err = strconv.Atoi(args[2])

	//the line below should be changed if the input file contains more than 10,000 decks
	if err != nil || numberOfDecksToBePlayed < 1 || numberOfDecksToBePlayed > (10000-firstDeckNum) {
		println("second argument invalid")
		println("numberOfDecksToBePlayed must be 1 or more, but not more than 10,000 minus firstDeckNum")
		os.Exit(1)
	}

	if numberOfDecksToBePlayed > 1 {
		singleGame = false
	}

	//
	// If length = -1 then execute playNew
	//
	// If length >= 0 then execute playOrig
	//    length  = 1 playBestOrIOS will just play the best move
	//    length >  1 playBestOrIOS will just play the best move or force a flip
	//

	length, err = strconv.Atoi(args[3]) //length of each strategy (which also determines the # of strategies - 2^n)

	// 		In the line below, 24 is arbitrarily set; 24 would result in 16,777,216 attempts per deck
	// 		I have never run the program with length greater than 16
	// 		Depending upon the size of an int, length could be 32 or greater, but the program may never finish
	if err != nil || length < -1 || length > 24 {
		println("Third argument invalid")
		println("Length must be an integer >= -1 and <= 24")
		os.Exit(1)
	}

	verbose, err = strconv.Atoi(args[4]) //the greater the number the more printing to standard output (terminal)
	// in the line below, 10 is arbitrarily set; at present all values greater than 1 result in the same
	if err != nil || verbose < 0 || verbose > 10 {
		println("fourth argument invalid")
		println("verbose must be a non-negative integer no greater than 10")
		os.Exit(1)
	}

	// Arguments 5 & 6 below applies only to playNew			****************************************************
	// But they must be on command line anyway
	switch strings.TrimSpace(args[5])[0:1] {
	case "A", "a":
		findAllSuccessfulStrategies = true
	case "F", "f":
		findAllSuccessfulStrategies = false
	default:
		println("Fifth argument invalid")
		println("  findAllSuccessfulStrategies must be either:")
		println("     'F' or 'f' - Normal case stop after finding a successful set of moves")
		println("     'A' or 'a' - See how many paths to success you can find")
		os.Exit(1)
	}

	switch strings.TrimSpace(args[6])[0:1] {
	case "X", "x":
		printTree = "X" // Do Not Print Tree
	case "C", "c":
		printTree = "C" // Print Tree to Console
	case "F", "f":
		printTree = "F" // Print Tree to File -- Note not yet implemented
	default:
		println("Sixth argument invalid")
		println("  printTree must be either:")
		println("     'X' or 'x' - See how many paths to success you can find")
		println("     'C' or 'c' - See how many paths to success you can find")
		println("     'F' or 'f' - Normal case stop after finding a successful set of moves")
		os.Exit(1)
	}
	// Argument above applies only to playNew			****************************************************

	// Open the deck file here in main - will be read in both playOrig and PlayNew
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

	if firstDeckNum > 0 {
		for i := 0; i < firstDeckNum; i++ {
			_, err = reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println("Cannot read from inputFileName", err)
			}
		}
	}

	if length != -1 {
		// playOrig will execute the original code designed to play either the "Best" move or the best move modified by the IOS strategy
		// of substituting FlipToWaste as described in more detail below.  This was formally known as the "playBestOrIOS" strategy
		// and developed under a function of that name in project branch "tree".
		// To avoid issues with old "tree" branch code the function playOrig has been created by refactoring and adding passed arguments.

		gameLengthLimit = gameLengthLimitOrig
		moveBasePriority = moveBasePriorityOrig
		playOrig(*reader)
	} else {
		gameLengthLimit = gameLengthLimitNew
		moveBasePriority = moveBasePriorityNew
		playNew(*reader)
	}

}
