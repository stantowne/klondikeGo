package main

import (
	"encoding/csv"
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"io"
	"log"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

// because aMoves sometimes includes moves with priority badMove, which would sort to the very end
// this function is needed to that, when iOS needs NOT to call the best move
// it instead calls a flip, rather than a bad move.
// oddly, although this function would seem to be necessary, introducing it had no effect on the win rate
func findFlip(moves []move) move {
	for i := 0; i < len(moves); i++ {
		if moves[i].name == "flipStockToWaste" {
			return moves[i]
		}
		if moves[i].name == "flipWasteToStock" {
			return moves[i]
		}
	}
	fmt.Println("This line from findFlip should never execute")
	return moves[len(moves)-1]
}

func playOrig(reader csv.Reader, cfg *Configuration) {

	// comment probably obsolete
	// Need to define variable err type error here.  Originally it was implicitly created by the following statement and then reused many times
	//   inputFileName := "decks-made-2022-01-15_count_10000-dict.csv"
	// That statement has been moved up into main so we need to explicitly create it here.

	numberOfStrategies := 1 << cfg.PlayOrig.Length //number of initial strategies

	startTime := time.Now()
	winCounter := 0
	earlyWinCounter := 0
	attemptsAvoidedCounter := 0
	lossesAtGLL := 0
	lossesAtNoMoves := 0
	regularLosses := 0
	var losses [][]string
	firstDeckNum := cfg.General.FirstDeckNum
	numberOfDecksToBePlayed := cfg.General.NumberOfDecksToBePlayed
	length := cfg.PlayOrig.Length
	verbose := cfg.General.Verbose
newDeck:
	for deckNum := firstDeckNum; deckNum < (firstDeckNum + numberOfDecksToBePlayed); deckNum++ {
		if deckNum%1000 == 0 {
			fmt.Printf("\nStarting Deck Number %v at %v", deckNum, time.Now())
		}
		protoDeck, err := reader.Read() // protoDeck is a slice of strings: rank, suit, rank, suit, etc.
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Cannot read from inputFileName:", err)
		}

		if cfg.General.Verbose > 1 {
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

	newInitialOverrideStrategy:
		for iOS := 0; iOS < numberOfStrategies; iOS++ {
			//deal Deck onto board
			var b = dealDeck(d)
			var priorBoardNullWaste board //used in Loss Detector
			if cfg.General.Verbose > 1 {
				fmt.Printf("Start play of Deck %v using initial override strategy %v.\n", deckNum, iOS)
			}

			//make this slice of int with length = 0 and capacity = gameLengthLimitOrig
			aMovesNumberOf := make([]int, 0, cfg.PlayOrig.GameLengthLimit) //number of available Moves

			for moveCounter := 1; moveCounter < cfg.PlayOrig.GameLengthLimit+2; moveCounter++ { //start with 1 to line up with Python version
				aMoves := detectAvailableMoves(b, moveCounter, cfg.General.NumberOfDecksToBePlayed == 1)

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
					if iOS == numberOfStrategies-1 {
						loss := []string{strconv.Itoa(deckNum), "lossAtNoMoves"}
						losses = append(losses, loss)
					}
					continue newInitialOverrideStrategy
				}

				// if more than one move is available, sort them
				if len(aMoves) > 1 { //sort them by priority if necessary
					sort.SliceStable(aMoves, func(i, j int) bool {
						return aMoves[i].priority < aMoves[j].priority
					})
				}
				// DanTest remove below
				for z := range aMoves {
					if aMoves[z].cardToMove.Rank == 0 && strings.Contains("moveDown/", aMoves[z].name+"/") || aMoves[z].name == "" {
						fmt.Printf("\n****************\nBad Move generated %v\n\n", aMoves[z])
						fmt.Printf("  Within these generated moves: ")
						for w := range aMoves {
							fmt.Printf("\n%v", aMoves[w])
						}
						fmt.Printf("\n\n")
						printBoard(b)
						fmt.Printf("\n\n%d", b.boardCode(deckNum))
					}
				}
				// DanTest remove above

				selectedMove := aMoves[0]

				//Initial Override Strategy logic
				mC := moveCounter - 1 // for this part of the program a zero-based move counter is needed
				//below:  example -> if length is 8, then this IF is satisfied for mc = 0, 1, 2, 3, 4, 5, 6 & 7
				if mC > -1 && mC < length {
					// below: & is bitwise AND which means look, bit by bit, at each operand result is 0 unless both bits are 1
					// below: first operand is the strategy number which also expresses the strategy
					// below: second operand is all zeros except the mC bit from the right.
					// below: so result is 0 unless mC(th) bit of the strategy, from right, is 1
					if iOS&(1<<mC) != 0 {
						selectedMove = findFlip(aMoves)
						//selectedMove = aMoves[len(aMoves)-1]  //see explanation of findFlip for why this line has been replaced
					}
				}

				b = moveMaker(b, selectedMove) //***Main Program Statement

				//quickTestBoardCodeDeCode(b, deckNum, length, iOS, moveCounter)

				//Detect Early Win
				if detectWinEarly(b) {
					earlyWinCounter++
					winCounter++
					attemptsAvoidedCounter = attemptsAvoidedCounter + numberOfStrategies - iOS

					if verbose > 0 {
						fmt.Printf("Deck %v, played using initialOverrideStrategy %v: Game won early after %v moves. \n", deckNum, iOS, mC)
					}
					if verbose > 1 {
						fmt.Printf("GameWon: aMovesNumberOf:\n%v\n", aMovesNumberOf)
					}
					if verbose > 1 {
						fmt.Printf("GameWon: Frequency of each moveType:\n%v\n", moveTypes)
					}
					continue newDeck
				}

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
					if moveCounter < 20 { // changed from < 20
						priorBoardNullWaste = b
					} else if reflect.DeepEqual(b, priorBoardNullWaste) {
						if verbose > 1 {
							fmt.Printf("*****Loss detected after %v moves\n", moveCounter)
						}
						regularLosses++
						if iOS == numberOfStrategies-1 {
							loss := []string{strconv.Itoa(deckNum), "regularLoss"}
							losses = append(losses, loss)
						}
						continue newInitialOverrideStrategy
					} else {
						priorBoardNullWaste = b
					}
				}
			}
			lossesAtGLL++
			loss := []string{strconv.Itoa(deckNum), "lossAtGameLengthLimit"}
			losses = append(losses, loss)

			if verbose > 0 {
				fmt.Printf("Deck %v, played using Initial Override Strategy %v: Game not won\n", deckNum, iOS)
			}
			if verbose > 1 {
				fmt.Printf("Game Not Won:  Frequency of each moveType:\n%v\n", moveTypes)
				fmt.Printf("Game Not Won: aMovesNumberOf:\n%v\n", aMovesNumberOf)
			}
		}

	}
	fileName := "playOrigLosses-firstDeck-" +
		strconv.Itoa(firstDeckNum) +
		"-strategyLength-" +
		strconv.Itoa(length) +
		"-numberOfDecks-" +
		strconv.Itoa(numberOfDecksToBePlayed) +
		".csv"
	file, err := os.Create(fileName)
	if err != nil {
		log.Println("Cannot create csv file:", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	err = writer.WriteAll(losses)
	if err != nil {
		log.Println("Cannot write csv file:", err)
	}
	possibleAttempts := numberOfStrategies * numberOfDecksToBePlayed
	lossCounter := numberOfDecksToBePlayed - winCounter
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	percentageAttemptsAvoided := 100.0 * float64(attemptsAvoidedCounter) / float64(possibleAttempts)
	var p = message.NewPrinter(language.English)
	fmt.Printf("\nDate & Time Completed is %v\n", endTime)
	_, err = p.Printf("Number of Decks Played is %d, starting with Deck %d.\n", numberOfDecksToBePlayed, firstDeckNum)
	if err != nil {
		fmt.Println("Number of Decks Played cannot print")
	}
	fmt.Printf("Length of Initial Override Strategies is %d.\n", length)
	fmt.Printf("Number of Initial Override Strategies Per Deck is %d.\n", numberOfStrategies)
	_, err = p.Printf("Number of Possible Attempts is %d.\n", possibleAttempts)
	if err != nil {
		fmt.Println("Number of Possible Attempts cannot print")
	}
	averageElapsedTimePerDeck := float64(elapsedTime.Milliseconds()) / float64(numberOfDecksToBePlayed)
	fmt.Printf("Elapsed Time is %v.\n", elapsedTime)
	_, err = p.Printf("Total Decks Won is %d of which %d were Early Wins\n", winCounter, earlyWinCounter)
	if err != nil {
		fmt.Println("Total Decks Won cannot print")
	}
	_, err = p.Printf("Total Decks Lost is %d\n", lossCounter)
	if err != nil {
		fmt.Println("Total Decks Lost cannot print")
	}
	_, err = p.Printf("Losses at Game Length Limit is %d\n", lossesAtGLL)
	if err != nil {
		fmt.Println("Losses at Game Length Limit cannot print")
	}
	_, err = p.Printf("Losses at No Moves Available is %d\n", lossesAtNoMoves)
	if err != nil {
		fmt.Println("Losses at No Moves Available cannot print")
	}
	_, err = p.Printf("Regular Losses is %d\n", regularLosses)
	if err != nil {
		fmt.Println("Regular Losses cannot print")
	}
	_, err = p.Printf("Number of Attempts Avoided ia %d\n", attemptsAvoidedCounter)
	if err != nil {
		fmt.Println("Number of Attempts Avoided cannot print")
	}
	_, err = p.Printf("Percentage of Possible Attempts Avoided is %v\n", percentageAttemptsAvoided)
	if err != nil {
		fmt.Println("Percentage of Possible Attempts Avoided cannot print")
	}
	fmt.Printf("Average Elapsed Time per Deck is %vms.\n", averageElapsedTimePerDeck)

}
