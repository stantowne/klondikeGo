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

	var decks decks = deckReader("decks-made-2022-01-15_count_10000-dict.json") //contains decks 0-999 from Python version
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
			var aMoves []move //available Moves
			aMoves = append(aMoves, detectUpMoves(b, movecounter, singleGame)...)
			aMoves = append(aMoves, detectAcrossMoves(b, movecounter, singleGame)...)
			aMoves = append(aMoves, detectEntireColumnMoves(b, movecounter, singleGame)...)
			aMoves = append(aMoves, detectDownMoves(b, movecounter, singleGame)...)
			aMoves = append(aMoves, detectPartialColumnMoves(b, movecounter, singleGame)...)
			aMoves = append(aMoves, detectFlipStockToWaste(b, movecounter, singleGame)...)
			aMoves = append(aMoves, detectFlipWasteToStock(b, movecounter, singleGame)...)

			if len(aMoves) == 0 {
				if singleGame {
					fmt.Printf("Deck %v: Game lost after %v moves\n", deckNum, movecounter)
					fmt.Printf("GameLost: Frequency of each moveType:\n%v\n", moveTypes)
					fmt.Printf("GameLost: availableMovesLengthRecord:\n%v\n", availableMoveLengthRecord)
				}
				continue outer
			}
			if len(aMoves) > 1 { //sort them by priority if necessary
				sort.SliceStable(aMoves, func(i, j int) bool {
					return aMoves[i].priority < aMoves[j].priority
				})
			}
			availableMoveLengthRecord = append(availableMoveLengthRecord, len(aMoves)) //record their number
			if singleGame {
				fmt.Printf("\nMove %v to be made is %+v\n", movecounter, aMoves[0])
			}
			b = moveMaker(b, aMoves[0]) //***Main Program Statement
			moveTypes[aMoves[0].name]++
			if singleGame {
				fmt.Printf("After move %v the board is as follows:\n", movecounter)
				printBoard(b)
			}
			if len(b.piles[0])+len(b.piles[1])+len(b.piles[2])+len(b.piles[3]) == 52 {
				winCounter++
				longest := longestRunOfOne(availableMoveLengthRecord)
				fmt.Printf("Deck %v: Game won after %v moves. Longest run of one is: %v\n", deckNum, movecounter, longest)
				if singleGame {
					fmt.Printf("GameWon: availableMovesLengthRecord:\n%v\n", availableMoveLengthRecord)
				}
				if singleGame {
					fmt.Printf("GameWon: Frequency of each moveType:\n%v\n", moveTypes)
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

func longestRunOfOne(aML []int) int { //availableMoveList
	var runOfOnes int = 0
	var longestRun int = 0
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
