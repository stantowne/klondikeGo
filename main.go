package main

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"time"
)

var moveBasePriority = map[string]int{
	"moveAceAcross":     300,
	"moveDeuceAcross":   400,
	"move3PlusAcross":   900,
	"moveDown":          600,
	"moveEntireColumn":  500,
	"flipWasteToStock":  1000,
	"flipStockToWaste":  1100,
	"movePartialColumn": 700,
	"moveAceUp":         100,
	"moveDeuceUp":       200,
	"move3PlusUp":       800,
	"badMove":           1200,
}

func main() {

	//Used to record how many of each move type is executed during an attempt.
	moveTypes := map[string]int{ //11 moveTypes
		"moveAceAcross":     0,
		"moveDeuceAcross":   0,
		"move3PlusAcross":   0,
		"moveDown":          0,
		"moveEntireColumn":  0,
		"flipWasteToStock":  0,
		"flipStockToWaste":  0,
		"movePartialColumn": 0,
		"moveAceUp":         0,
		"moveDeuceUp":       0,
		"move3PlusUp":       0,
	}

	const gameLengthLimit = 150 //increasing to 200 does not increase win rate

	firstDeckNum, _ := strconv.Atoi(os.Args[1])
	numberOfDecksToBePlayed, _ := strconv.Atoi(os.Args[2])
	length, _ := strconv.Atoi(os.Args[3])  //length of each strategy (which also determines the # of strategies - 2^n)
	verbose, _ := strconv.Atoi(os.Args[4]) //the greater the number the more verbose

	numberOfStrategies := 1 //number of initial strategies
	for i := 0; i < length; i++ {
		numberOfStrategies = numberOfStrategies * 2
	}

	var veryVerbose = false
	if verbose > 2 {
		veryVerbose = true
	}

	var decks = deckReader("decks-made-2022-01-15_count_10000-dict.json") //contains decks 0-999 from Python version
	startTime := time.Now()
	winCounter := 0
	attemptsAvoidedCounter := 0
	lossesAtGLL := 0
	lossesAtNoMoves := 0
newDeck:
	for deckNum := firstDeckNum; deckNum < (firstDeckNum + numberOfDecksToBePlayed); deckNum++ {
		if verbose > 1 {
			fmt.Printf("\nDeck #%d:\n", deckNum)
		}
	newInitialOverrideStrategy:
		for iOS := 0; iOS < numberOfStrategies; iOS++ {
			//deal deck onto board
			var b = dealDeck(decks[deckNum])
			var priorBoardNullWaste board //used in Loss Detector
			if verbose > 1 {
				fmt.Printf("Start play of deck %v using initial override strategy %v.\n", deckNum, iOS)
			}

			//make this slice of int with length = 0 and capacity = gameLengthLimit
			aMovesNumberOf := make([]int, 0, gameLengthLimit) //number of available Moves

			for moveCounter := 1; moveCounter < gameLengthLimit+2; moveCounter++ { //start with 1 to line up with Python version
				var aMoves []move //available Moves
				aMoves = append(aMoves, detectUpMoves(b, moveCounter, veryVerbose)...)
				aMoves = append(aMoves, detectAcrossMoves(b, moveCounter, veryVerbose)...)
				aMoves = append(aMoves, detectEntireColumnMoves(b, moveCounter, veryVerbose)...)
				aMoves = append(aMoves, detectDownMoves(b, moveCounter, veryVerbose)...)
				aMoves = append(aMoves, detectPartialColumnMoves(b, moveCounter, veryVerbose)...)
				aMoves = append(aMoves, detectFlipStockToWaste(b, moveCounter, veryVerbose)...)
				aMoves = append(aMoves, detectFlipWasteToStock(b, moveCounter, veryVerbose)...)

				//detects Loss
				if len(aMoves) == 0 { //No available moves; game lost.
					if verbose > 1 {
						fmt.Printf("Initial Override Strategy: %v\n", iOS)
						fmt.Printf("****Deck %v: XXXXGame lost after %v moves\n", deckNum, moveCounter)
					}
					if verbose > 2 {
						fmt.Printf("GameLost: Frequency of each moveType:\n%v\n", moveTypes)
						fmt.Printf("GameLost: aMovesNumberOf:\n%v\n", aMovesNumberOf)
					}
					lossesAtNoMoves++
					continue newInitialOverrideStrategy
				}

				// if more than one move is available, sort them
				if len(aMoves) > 1 { //sort them by priority if necessary
					sort.SliceStable(aMoves, func(i, j int) bool {
						return aMoves[i].priority < aMoves[j].priority
					})
				}

				selectedMove := aMoves[0]

				//Initial Override Strategy logic
				mC := moveCounter - 1 // for this part of the program a zero-based move counter is needed
				if mC < length {
					if iOS&(1<<mC) != 0 {
						selectedMove = aMoves[len(aMoves)-1]
					}
				}

				b = moveMaker(b, selectedMove) //***Main Program Statement

				//Detects Win
				if len(b.piles[0])+len(b.piles[1])+len(b.piles[2])+len(b.piles[3]) == 52 {
					winCounter++
					attemptsAvoidedCounter = attemptsAvoidedCounter + numberOfStrategies - iOS

					if verbose > 0 {
						fmt.Printf("Deck %v, played using initialOverrideStrategy %v: Game won after %v moves. \n", deckNum, iOS, mC)
					}
					if verbose > 1 {
						fmt.Printf("GameWon: aMovesNumberOf:\n%v\n", aMovesNumberOf)
					}
					if verbose > 1 {
						fmt.Printf("GameWon: Frequency of each moveType:\n%v\n", moveTypes)
					}
					continue newDeck
				}
				//Detects Loss
				if aMoves[0].name == "flipWasteToStock" {
					if moveCounter < 20 {
						priorBoardNullWaste = b
					} else if reflect.DeepEqual(b, priorBoardNullWaste) {
						if verbose > 1 {
							fmt.Printf("*****Loss detected after %v moves\n", moveCounter)
						}
						continue newInitialOverrideStrategy
					} else {
						priorBoardNullWaste = b
					}
				}
			}
			lossesAtGLL++
			if verbose > 0 {
				fmt.Printf("Deck %v, played using Initial Override Strategy %v: Game not won\n", deckNum, iOS)
			}
			if verbose > 1 {
				fmt.Printf("Game Not Won:  Frequency of each moveType:\n%v\n", moveTypes)
				fmt.Printf("Game Not Won: aMovesNumberOf:\n%v\n", aMovesNumberOf)
			}
		}

	}
	possibleAttempts := numberOfStrategies * numberOfDecksToBePlayed
	lossCounter := numberOfDecksToBePlayed - winCounter
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	fmt.Printf("\nNumber of Decks PLayed is %v.\n", numberOfDecksToBePlayed)
	fmt.Printf("Length of Initial Override Strategies is %v.\n", length)
	fmt.Printf("Number of Initial Override Strategies is %v.\n", numberOfStrategies)
	fmt.Printf("Number of Possible Attempts is %v.\n", possibleAttempts)
	averageElapsedTimePerDeck := float64(elapsedTime.Milliseconds()) / float64(numberOfDecksToBePlayed)
	fmt.Printf("Elapsed Time is %v.\n", elapsedTime)
	fmt.Printf("Total Decks Won is %v\n", winCounter)
	fmt.Printf("Total Decks Lost is %v\n", lossCounter)
	fmt.Printf("Losses at Game Length Limit is %v\n", lossesAtGLL)
	fmt.Printf("Losses at No Moves Available is %v\n", lossesAtNoMoves)
	fmt.Printf("Number of Attempts Avoided ia %v\n", attemptsAvoidedCounter)
	fmt.Printf("Average Elapsed Time per Deck is %fms.\n", averageElapsedTimePerDeck)

}
