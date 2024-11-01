package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const gameLengthLimitOrig = 150     // max moveCounter; increasing to 200 does not increase win rate
const gameLengthLimitNew = 25000000 // max mvsTriedTD
var gameLengthLimit int

var firstDeckNum int
var numberOfDecksToBePlayed int
var length int
var verbose int
var verboseSpecial string
var findAllSuccessfulStrategies bool

type pMds struct {
	pType                 string
	deckStartVal          int
	deckContinueFor       int
	aMvsThisDkStartVal    int
	aMvsThisDkContinueFor int
	outputTo              string
}

var printMoveDetail pMds

var err error
var singleGame bool // = true

func main() {
	/*

						Command line arguments

						args[0] = program name
						args[1] = # of the first deck to be used from within the pre-stored decks used to standardize testing
						args[2] = # of decks to be played
						args[3] = length of IOS strategy - see comments below	(applicable to playOrig only)
						args[4] = verbosity switch for messages
						args[5] = findAllSuccessfulStrategies 					(applicable playNew only)
						args[6] = printMoveDetail		                  		(applicable playNew only)
					              	as a string to be parsed of the form:
					           			pType,startType,deckStartVal,deckContinueFor,outputTo

										where:
											pType = empty or X - do not print NOTE: Default if args[6] is not on command line
												  = BB         - Board by Board detail
				                                  = BBS        - Board by Board Short detail
				                                  = BBSS       - Board by Board Super Short detail
					                              = TW         - print Tree in Wide mode     8 char per move
												  = TS         - print Tree in Skinny mode   5 char per move
			                                      = TSS        - print Tree in Super Skinny mode   3 char per move
		                               These next four limit at what point and for how long move detail should actually be printed.
											deckStartVal    	  = Non-negative integer (Default 0)
											deckContinueFor  	  = Non-negative integer (Default 0 which indicates forever)
											aMvsThisDkStartVal    = Non-negative integer (Default 0)
											aMvsThisDkContinueFor = Non-negative integer (Default 0 which indicates forever)

		                                    outputTo = C = Console (default)
								                     = file name and path (if applicable)
					                                   Note: if file name is present then startType. deckStartVal and ContinueFor
					                                         must be present or delineated with ":"

						           	and placed into a package level struct printMoveDetail of type pMd which can be seen above:
	*/

	// Always a good idea to print out the program source of the output.
	//    Will be especially useful when we figure out how to include the versioning

	printMoveDetail.pType = "X"
	printMoveDetail.deckStartVal = 0
	printMoveDetail.deckContinueFor = 0
	printMoveDetail.aMvsThisDkStartVal = 0
	printMoveDetail.aMvsThisDkContinueFor = 0
	printMoveDetail.outputTo = "C"

	args := os.Args

	// Convert all the arguments to upper case

	for i := range args {
		args[i] = strings.ToUpper(args[i])
	}

	firstDeckNum, err = strconv.Atoi(args[1])
	if err != nil || firstDeckNum < 0 || firstDeckNum > 9999 {
		println("first argument invalid - args[1]")
		println("firstDeckNum must be non-negative integer less than 10,000")
		os.Exit(1)
	}

	numberOfDecksToBePlayed, err = strconv.Atoi(args[2])

	//the line below should be changed if the input file contains more than 10,000 decks
	if err != nil || numberOfDecksToBePlayed < 1 || numberOfDecksToBePlayed > (10000-firstDeckNum) {
		println("second argument invalid - args[2]")
		println("numberOfDecksToBePlayed must be 1 or more, but not more than 10,000 minus firstDeckNum")
		os.Exit(1)
	}

	if numberOfDecksToBePlayed == 1 {
		singleGame = true

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
		println("Third argument invalid - args[3]")
		println("Length must be an integer >= -1 and <= 24")
		os.Exit(1)
	}

	// The first character of verbose must be a value from 0 - 9.  Higher numbers indicate more detailed msgs should be printed
	//  The remaining characters are used to form verboseSpecial.  Code in the program will look for specific values within
	//  verbose special to indicate that optional printing should be done.
	verboseSpecial = args[4]
	verbose, err = strconv.Atoi(verboseSpecial[0:1])
	if err != nil || verbose >= 10 || verbose < 0 {
		println("fourth argument invalid - args[4]")
		println("verbose must be a non-negative integer no greater than 10")
		os.Exit(1)
	}
	verboseSpecial = verboseSpecial[1:]
	/* Verbose Special codes implemented:  CASE IS IGNORED  Place some divider such as / when multiple specials are requested

	   DBD = print Deck-by-deck detail info after each Move 											playNew Only (in playAllMovesS)
	   DBDS = print Deck-by-deck SHORT detail info after each Move 											playNew Only (in playAllMovesS)WL  = print deck summary Win/Loss stats after all decks to see which decks won and which lost    playNew Only (in playNew)
	   SUITSYMBOL = print S, D, C, H instead of runes - defaults to runes
	   RANKSYMBOL = print Ac, Ki, Qu, Jk instead of 01, 11, 12, 13 - defaults to numeric
	*/

	// Arguments 5 & 6 below applies only to playNew			****************************************************
	// But they must be on command line anyway
	switch strings.TrimSpace(args[5])[0:1] {
	case "A", "a":
		findAllSuccessfulStrategies = true
	case "F", "f":
		findAllSuccessfulStrategies = false
	default:
		println("Fifth argument invalid - args[5]")
		println("  findAllSuccessfulStrategies must be either:")
		println("     'F' or 'f' - Normal case stop after finding a successful set of moves")
		println("     'A' or 'a' - See how many paths to success you can find")
		os.Exit(1)
	}
	pMdArgs := strings.Split(args[6], ",")
	l := len(pMdArgs)

	if l >= 1 {
		if pMdArgs[0] == "BB" || pMdArgs[0] == "BBS" || pMdArgs[0] == "BBSS" || pMdArgs[0] == "TW" || pMdArgs[0] == "TS" || pMdArgs[0] == "X" {
			printMoveDetail.pType = pMdArgs[0]
		} else {
			println("Sixth argument part 1 invalid - args[6] arg[6] parts are  separated by commas: *1*,2,3,4,5,6")
			println("  Must start with BB, BBS, BBSS, TW, TS, TSS or X")
			os.Exit(1)
		}
	}
	if l >= 2 {
		printMoveDetail.deckStartVal, err = strconv.Atoi(pMdArgs[1])
		if err != nil || printMoveDetail.deckStartVal < 0 {
			println("Sixth argument part 2 invalid - args[6] arg[6] parts are  separated by commas: 1,*2*,3,4,5,6")
			println("must be a non-negative integer")
			os.Exit(1)
		}
	}
	if l >= 3 {
		printMoveDetail.deckContinueFor, err = strconv.Atoi(pMdArgs[2])
		if err != nil || printMoveDetail.deckContinueFor < 0 {
			println("Sixth argument part 3 invalid - args[6] arg[6] parts are  separated by commas: 1,2,*3*,4,5,6")
			println("must be a non-negative integer")
			os.Exit(1)
		}
	}
	if l >= 4 {
		printMoveDetail.aMvsThisDkStartVal, err = strconv.Atoi(pMdArgs[3])
		if err != nil || printMoveDetail.aMvsThisDkStartVal < 0 {
			println("Sixth argument part 4 invalid - args[6] arg[6] parts are  separated by commas: 1,2,3,*4*,5,6")
			println("must be a non-negative integer")
			os.Exit(1)
		}
	}
	if l >= 5 {
		printMoveDetail.aMvsThisDkContinueFor, err = strconv.Atoi(pMdArgs[4])
		if err != nil || printMoveDetail.aMvsThisDkContinueFor < 0 {
			println("Sixth argument part 5 invalid - args[6] arg[6] parts are  separated by commas: 1,2,3,4,*5*,6")
			println("must be a non-negative integer")
			os.Exit(1)
		}
	}
	if l >= 6 {
		if pMdArgs[5] == "C" {
			printMoveDetail.outputTo = pMdArgs[5]
		} else {
			// add if-clause here to test if valid filename and it can be overwritten
			println("Sixth argument part 6 invalid - args[6] arg[6] parts are  separated by commas: 1,2,3,4,5,*6*")
			println("must be a C for Console or a valid file name")
			os.Exit(1)
		}
	}

	// Argument above applies only to playNew			****************************************************

	fmt.Printf("\nCalling Program: %v\n\n", args[0])
	fmt.Printf("Command Line Arguments:\n"+
		"            Number Of Decks To Be Played: %v\n"+
		"                      Starting with deck: %v\n\n", numberOfDecksToBePlayed, firstDeckNum)
	if length != -1 {
		nOfS := 1 << length //number of initial strategies
		fmt.Printf(" Style: Original IOS (Initial Override Strategies)\n\n"+
			"                     IOS strategy length: %v\n"+
			"      Max IOS startegy attempts per deck: %v\n"+
			"       Total strategy attempts all decks: %v\n\n", length, nOfS, nOfS*numberOfDecksToBePlayed)
	} else {
		fmt.Printf(" Style: New AllMvs (All Moves Possible)\n\n" +
			"   Max AllMvs startegy attempts per deck: Variable\n\n")
	}
	fmt.Printf("                           Verbose level: %v\n"+
		"                   Verbose special codes: %v\n",
		verbose, verboseSpecial)
	if length == -1 {
		fmt.Printf("\n Print Move Detail Options:\n"+
			"          Find All Successful Strategies: %v\n"+
			"                              Print Type: %v\n"+
			"                       Staring with Deck: %v\n"+
			"                            Continue for: %v decks (0 = all the rest)\n"+
			"             Starting with Moves Tried #: %v\n"+
			"                            Continue for: %v moves tried (0 = all the rest)\n",
			findAllSuccessfulStrategies,
			printMoveDetail.pType,
			printMoveDetail.deckStartVal,
			printMoveDetail.deckContinueFor,
			printMoveDetail.aMvsThisDkStartVal,
			printMoveDetail.aMvsThisDkContinueFor)
		if printMoveDetail.outputTo == "C" {
			fmt.Printf("                         Print Output to: Console\n")
		} else {
			fmt.Printf("                    Print Output to file: %v  (not yet implemented)\n", printMoveDetail.outputTo)
		}
	}

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
		fmt.Printf("\n                         GameLengthLimit: %v\n\n\n", gameLengthLimit)
		playOrig(*reader)
	} else {
		gameLengthLimit = gameLengthLimitNew
		fmt.Printf("\n                         GameLengthLimit: %v  (not yet implemented)\n\n\n", gameLengthLimit)
		moveBasePriority = moveBasePriorityNew
		playNew(*reader)
	}

}
