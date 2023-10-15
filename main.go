package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"
)

// type gameResult struct {
// 	result string
// }

func main() {
	// Used to control the order in which the detection functions are executed.
	// This is not important since the slice of type move is sorted.
	detectionFunctions := []func(board, int, bool) []move{
		detectUpMoves,
		detectAcrossMoves,
		detectEntireColumnMoves,
		detectDownMoves,
		detectPartialColumnMovesAcross,
		detectFlipStockToWaste,
		detectFlipWasteToStock,
	}
	/*
		moveBasePriority := map[string]int{ //11 moveTypes + bad Move
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
		}*/

	//Used to record how many of each move type is executed during a game.
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

	gameLengthLimit := 150
	firstDeckNum, _ := strconv.Atoi(os.Args[1])
	numberOfDecksToBePlayed, _ := strconv.Atoi(os.Args[2])

	singleGame := false
	if numberOfDecksToBePlayed == 1 {
		singleGame = true
	}

	var decks decks = deckReader("decks-made-2022-01-15_count_1000-dict.json") //contains decks 0-999 from Python version
	startTime := time.Now()
	winCounter := 0
outer:
	for deckNum := firstDeckNum; deckNum < (firstDeckNum + numberOfDecksToBePlayed); deckNum++ {
		var b board = dealDeck(decks[deckNum])
		if singleGame {
			fmt.Printf("\nStart play of deck %v\n", deckNum)
			fmt.Printf("Starting board is as follows:\n")
			printBoard(b)
		}
		availableMoveLengthRecord := make([]int, 0, gameLengthLimit)
		for movecounter := 1; movecounter < gameLengthLimit+2; movecounter++ { //start with 1 to line up with Python version
			if singleGame {
				fmt.Println("\n\n**********************************************************")
				fmt.Printf("Looking for Move %v\n", movecounter)
			}
			var availableMoves []move
			for _, m := range detectionFunctions { //find all availableMoves
				availableMoves = append(availableMoves, m(b, movecounter, singleGame)...)
			}
			if len(availableMoves) == 0 {
				if singleGame {
					fmt.Printf("Deck %v: Game lost after %v moves\n", deckNum, movecounter)
					fmt.Printf("GameLost: Frequency of each moveType:\n%v\n", moveTypes)
					fmt.Printf("GameLost: availableMovesLengthRecord:\n%v\n", availableMoveLengthRecord)
				}
				continue outer
			}
			if len(availableMoves) > 1 { //sort them by priority if necessary
				sort.SliceStable(availableMoves, func(i, j int) bool {
					return availableMoves[i].priority < availableMoves[j].priority
				})
			}
			availableMoveLengthRecord = append(availableMoveLengthRecord, len(availableMoves)) //record their number
			if singleGame {
				fmt.Printf("\nMove %v to be made is %+v\n", movecounter, availableMoves[0])
			}
			b = moveMaker(b, availableMoves[0]) //***Main Program Statement
			moveTypes[availableMoves[0].name]++
			if singleGame {
				fmt.Printf("After move %v the board is as follows:\n", movecounter)
				printBoard(b)
			}
			if len(b.piles[0])+len(b.piles[1])+len(b.piles[2])+len(b.piles[3]) == 52 {
				winCounter++
				if singleGame {
					fmt.Println("")
				}
				fmt.Printf("Deck %v: Game won after %v moves\n", deckNum, movecounter)
				if singleGame {
					fmt.Printf("GameWon: Frequency of each moveType:\n%v\n", moveTypes)
					fmt.Printf("GameWon: availableMovesLengthRecord:\n%v\n", availableMoveLengthRecord)
				}
				continue outer
			}
		}
		if singleGame {
			fmt.Printf("Deck %v: Game not won\n", deckNum)
			fmt.Printf("Game Not Won:  Frequency of each moveType:\n%v\n", moveTypes)
			fmt.Printf("Game Not Won: availableMovesLengthRecord:\n%v\n", availableMoveLengthRecord)
		}

	}
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	averageElapsedTime := float64(elapsedTime.Nanoseconds()) / float64(numberOfDecksToBePlayed)
	fmt.Printf("\nTotal Games PLayed: %v: Total Games Won: %v\n", numberOfDecksToBePlayed, winCounter)
	fmt.Printf("Elapsed Time is %v; Average Elapsed Time per Game is %fns.\n", elapsedTime, averageElapsedTime)
}
