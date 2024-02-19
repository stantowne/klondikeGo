package main

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"time"
)

// a point in an attempt
type point struct {
	board  board  //board before move
	aMoves []move //available moves
	move   move   //move to be made
}

// record of an attempt is a slice of point
var record []point

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

	const initialFlipsMax = 8
	const gameLengthLimit = 150 //increasing to 200 does not increase win rate
	firstDeckNum, _ := strconv.Atoi(os.Args[1])
	numberOfDecksToBePlayed, _ := strconv.Atoi(os.Args[2])

	singleGame := false
	if numberOfDecksToBePlayed == 1 {
		singleGame = true
	}

	var decks = deckReader("decks-made-2022-01-15_count_10000-dict.json") //contains decks 0-999 from Python version
	startTime := time.Now()
	winCounter := 0
	var winsByInitialFlips [initialFlipsMax]int
newDeck:
	for deckNum := firstDeckNum; deckNum < (firstDeckNum + numberOfDecksToBePlayed); deckNum++ {
	newInitialFlips:
		for initialFlips := 0; initialFlips < initialFlipsMax; initialFlips++ {
			//deal deck onto board
			var b = dealDeck(decks[deckNum])
			var priorBoardNullWaste board //used in Loss Detector
			if singleGame {
				fmt.Printf("Start play of deck %v after %v initial flips from stock to waste.\n", deckNum, initialFlips)
			}
			//make the initial flips
			for flip := 0; flip <= initialFlips; flip++ {
				b = flipStockToWaste(b)
			}
			//make this slice of int with length = 0 and capacity = gameLengthLimit
			aMovesNumberOf := make([]int, 0, gameLengthLimit) //number of available Moves

			for movecounter := 1; movecounter < gameLengthLimit+2; movecounter++ { //start with 1 to line up with Python version
				if singleGame {
					//fmt.Println("\n\n**********************************************************")
					//fmt.Printf("Looking for Move %v\n", movecounter)
				}
				var aMoves []move //available Moves
				aMoves = append(aMoves, detectUpMoves(b, movecounter, singleGame)...)
				aMoves = append(aMoves, detectAcrossMoves(b, movecounter, singleGame)...)
				aMoves = append(aMoves, detectEntireColumnMoves(b, movecounter, singleGame)...)
				aMoves = append(aMoves, detectDownMoves(b, movecounter, singleGame)...)
				aMoves = append(aMoves, detectPartialColumnMoves(b, movecounter, singleGame)...)
				aMoves = append(aMoves, detectFlipStockToWaste(b, movecounter, singleGame)...)
				aMoves = append(aMoves, detectFlipWasteToStock(b, movecounter, singleGame)...)

				//detects Loss
				if len(aMoves) == 0 { //No available moves; game lost.
					if singleGame {
						fmt.Printf("InitialFlips: %v\n", initialFlips)
						fmt.Printf("Deck %v: Game lost after %v moves\n", deckNum, movecounter)
						fmt.Printf("GameLost: Frequency of each moveType:\n%v\n", moveTypes)
						fmt.Printf("GameLost: aMovesNumberOf:\n%v\n", aMovesNumberOf)
					}
					continue newInitialFlips
				}
				if len(aMoves) > 1 { //sort them by priority if necessary
					sort.SliceStable(aMoves, func(i, j int) bool {
						return aMoves[i].priority < aMoves[j].priority
					})
				}
				aMovesNumberOf = append(aMovesNumberOf, len(aMoves)) //record their number
				if singleGame {
					fmt.Printf("Move %v to be made is %+v\n", movecounter, aMoves[0])
				}
				selectedMove := aMoves[0]
				//additional logic goes here
				point := point{
					board:  b,
					aMoves: aMoves,
					move:   selectedMove,
				}
				record = append(record, point)
				b = moveMaker(b, selectedMove) //***Main Program Statement
				moveTypes[selectedMove.name]++
				if singleGame {
					//fmt.Printf("After move %v the board is as follows:\n", movecounter)
					printBoard(b)
				}
				//Detects Win
				if len(b.piles[0])+len(b.piles[1])+len(b.piles[2])+len(b.piles[3]) == 52 {
					winCounter++
					winsByInitialFlips[initialFlips]++
					longest := longestRunOfOne(aMovesNumberOf)
					if numberOfDecksToBePlayed < 10000 {
						fmt.Printf("Deck %v, played with %v initial flips from stock to waste: Game won after %v moves. Longest run of one is: %v\n", deckNum, initialFlips, movecounter, longest)
					}
					if singleGame {
						fmt.Printf("GameWon: aMovesNumberOf:\n%v\n", aMovesNumberOf)
					}
					if singleGame {
						fmt.Printf("GameWon: Frequency of each moveType:\n%v\n", moveTypes)
					}
					continue newDeck
				}
				///Detects Loss
				if aMoves[0].name == "flipWasteToStock" {
					if movecounter < 50 {
						priorBoardNullWaste = b
					} else if reflect.DeepEqual(b, priorBoardNullWaste) {
						if singleGame {
							fmt.Printf("Loss detected after %v moves\n", movecounter)
						}
						continue newInitialFlips
					} else {
						priorBoardNullWaste = b
					}
				}
			}
			if singleGame {
				fmt.Printf("Deck %v, played with %v initial flips from stock to waste: Game not won\n", deckNum, initialFlips)
				fmt.Printf("Game Not Won:  Frequency of each moveType:\n%v\n", moveTypes)
				fmt.Printf("Game Not Won: aMovesNumberOf:\n%v\n", aMovesNumberOf)
			}
		}

	}
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	averageElapsedTimePerDeck := float64(elapsedTime.Milliseconds()) / float64(numberOfDecksToBePlayed)
	attempts := initialFlipsMax * numberOfDecksToBePlayed
	for i, v := range winsByInitialFlips {
		attempts = attempts - (v * (initialFlipsMax - (i + 1)))
	}
	averageElapsedTimePerAttempt := float64(elapsedTime.Milliseconds()) / float64(attempts)
	averageAttemptsPerDeck := float64(attempts) / float64(numberOfDecksToBePlayed)
	fmt.Printf("\nElapsed Time is %v.\n", elapsedTime)
	fmt.Printf("Total Decks PLayed: %v. Total Decks Won: %v\n", numberOfDecksToBePlayed, winCounter)
	fmt.Printf("Wins by Number of Initial Flips is %v\n", winsByInitialFlips)
	fmt.Printf("Total Attempts Made: %v. Attempts per Deck: %f\n", attempts, averageAttemptsPerDeck)
	fmt.Printf("Average Elapsed Time per Deck is %fms.\n", averageElapsedTimePerDeck)
	fmt.Printf("Average Elapsed Time per Attempt is %fms.\n", averageElapsedTimePerAttempt)
}

func longestRunOfOne(aML []int) int { //availableMoveList
	var runOfOnes = 0
	var longestRun = 0
	for _, num := range aML {
		if num == 1 {
			runOfOnes++
			if runOfOnes >= longestRun {
				longestRun = runOfOnes
			}
		} else {
			runOfOnes = 0
		}
	}
	return longestRun
}
