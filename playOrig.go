package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"sort"
	"strconv"
	"time"
)

// because aMoves sometimes includes moves with priority badMove, which would sort to the very end
// this function is needed to ensure that, when iOS needs NOT to call the best move
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

type variablesSpecificToPlayOrig struct {
	numberOfStrategies     int
	startTime              time.Time
	endTime                time.Time
	winCounter             int
	earlyWinCounter        int
	attemptsAvoidedCounter int
	lossesAtGLL            int
	lossesAtNoMoves        int
	regularLosses          int
	losses                 [][]string
}

func playOrig(reader csv.Reader, cfg *Configuration) {
	var vPO variablesSpecificToPlayOrig
	vPO.numberOfStrategies = 1 << cfg.PlayOrig.Length //number of initial strategies
	vPO.startTime = time.Now()

newDeck:
	for deckNum := cfg.General.FirstDeckNum; deckNum < (cfg.General.FirstDeckNum + cfg.General.NumberOfDecksToBePlayed); deckNum++ {
		if deckNum%1000 == 0 {
			_, _ = fmt.Fprintf(oW, "\nStarting Deck Number %v at %v", deckNum, time.Now())
		}
		protoDeck, err := reader.Read() // protoDeck is a slice of strings: rank, suit, rank, suit, etc.
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Cannot read from inputFileName:", err)
		}

		if cfg.General.Verbose > 1 {
			_, _ = fmt.Fprintf(oW, "\nDeck #%d:\n", deckNum)
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
		for iOS := 0; iOS < vPO.numberOfStrategies; iOS++ {
			//deal Deck onto board
			var b = dealDeck(d)
			var priorBoardNullWaste board //used in Loss Detector
			if cfg.General.Verbose > 1 {
				_, _ = fmt.Fprintf(oW, "Start play of Deck %v using initial override strategy %v.\n", deckNum, iOS)
			}

			//make this slice of int with length = 0 and capacity = gameLengthLimitOrig
			aMovesNumberOf := make([]int, 0, cfg.PlayOrig.GameLengthLimit) //number of available Moves

			for moveCounter := 1; moveCounter < cfg.PlayOrig.GameLengthLimit+2; moveCounter++ { //start with 1 to line up with Python version
				aMoves := detectAvailableMoves(b, moveCounter, cfg.General.NumberOfDecksToBePlayed == 1)

				//detects Loss
				if len(aMoves) == 0 { //No available moves; game lost.
					if cfg.General.Verbose > 1 {
						_, _ = fmt.Fprintf(oW, "Initial Override Strategy: %v\n", iOS)
						_, _ = fmt.Fprintf(oW, "****Deck %v: XXXXGame lost after %v moves\n", deckNum, moveCounter)
					}
					if cfg.General.Verbose > 2 {
						_, _ = fmt.Fprintf(oW, "GameLost: Frequency of each moveType:\n%v\n", moveTypes)
						_, _ = fmt.Fprintf(oW, "GameLost: aMovesNumberOf:\n%v\n", aMovesNumberOf)
					}
					vPO.lossesAtNoMoves++
					if iOS == vPO.numberOfStrategies-1 {
						loss := []string{strconv.Itoa(deckNum), "lossAtNoMoves"}
						vPO.losses = append(vPO.losses, loss)
					}
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
				//below:  example -> if length is 8, then this IF is satisfied for mc = 0, 1, 2, 3, 4, 5, 6 & 7
				if mC > -1 && mC < cfg.PlayOrig.Length {
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

				//quickTestBoardCodeDeCode(b, deckNum, cfg.PlayOrig.Length, iOS, moveCounter)

				//Detect Early Win
				if detectWinEarly(b) {
					vPO.earlyWinCounter++
					vPO.winCounter++
					vPO.attemptsAvoidedCounter = vPO.attemptsAvoidedCounter + vPO.numberOfStrategies - iOS

					if cfg.General.Verbose > 0 {
						_, _ = fmt.Fprintf(oW, "Deck %v, played using initialOverrideStrategy %v: Game won early after %v moves. \n", deckNum, iOS, mC)
					}
					if cfg.General.Verbose > 1 {
						_, _ = fmt.Fprintf(oW, "GameWon: aMovesNumberOf:\n%v\n", aMovesNumberOf)
					}
					if cfg.General.Verbose > 1 {
						_, _ = fmt.Fprintf(oW, "GameWon: Frequency of each moveType:\n%v\n", moveTypes)
					}
					continue newDeck
				}

				//Detects Win
				if len(b.piles[0])+len(b.piles[1])+len(b.piles[2])+len(b.piles[3]) == 52 {
					vPO.winCounter++
					vPO.attemptsAvoidedCounter = vPO.attemptsAvoidedCounter + vPO.numberOfStrategies - iOS

					if cfg.General.Verbose > 0 {
						_, _ = fmt.Fprintf(oW, "Deck %v, played using initialOverrideStrategy %v: Game won after %v moves. \n", deckNum, iOS, mC)
					}
					if cfg.General.Verbose > 1 {
						_, _ = fmt.Fprintf(oW, "GameWon: aMovesNumberOf:\n%v\n", aMovesNumberOf)
					}
					if cfg.General.Verbose > 1 {
						_, _ = fmt.Fprintf(oW, "GameWon: Frequency of each moveType:\n%v\n", moveTypes)
					}
					continue newDeck
				}
				//Detects Loss
				if aMoves[0].name == "flipWasteToStock" {
					if moveCounter < 20 { // changed from < 20
						priorBoardNullWaste = b
					} else if reflect.DeepEqual(b, priorBoardNullWaste) {
						if cfg.General.Verbose > 1 {
							_, _ = fmt.Fprintf(oW, "*****Loss detected after %v moves\n", moveCounter)
						}
						vPO.regularLosses++
						if iOS == vPO.numberOfStrategies-1 {
							loss := []string{strconv.Itoa(deckNum), "regularLoss"}
							vPO.losses = append(vPO.losses, loss)
						}
						continue newInitialOverrideStrategy
					} else {
						priorBoardNullWaste = b
					}
				}
			}
			vPO.lossesAtGLL++
			loss := []string{strconv.Itoa(deckNum), "lossAtGameLengthLimit"}
			vPO.losses = append(vPO.losses, loss)

			if cfg.General.Verbose > 0 {
				_, _ = fmt.Fprintf(oW, "Deck %v, played using Initial Override Strategy %v: Game not won\n", deckNum, iOS)
			}
			if cfg.General.Verbose > 1 {
				_, _ = fmt.Fprintf(oW, "Game Not Won:  Frequency of each moveType:\n%v\n", moveTypes)
				_, _ = fmt.Fprintf(oW, "Game Not Won: aMovesNumberOf:\n%v\n", aMovesNumberOf)
			}
		}

	}
	if cfg.General.Verbose > 2 {
		fileName := "../playOrigLossesOutput/playOrigLosses-firstDeck-" +
			strconv.Itoa(cfg.General.FirstDeckNum) +
			"-strategyLength-" +
			strconv.Itoa(cfg.PlayOrig.Length) +
			"-numberOfDecks-" +
			strconv.Itoa(cfg.General.NumberOfDecksToBePlayed) +
			".csv"
		file, err := os.Create(fileName)
		if err != nil {
			log.Println("Cannot create csv file:", err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				println("cannot close file")
				os.Exit(1)
			}
		}(file)
		writer := csv.NewWriter(file)
		err = writer.WriteAll(vPO.losses)
		if err != nil {
			log.Println("Cannot write csv file:", err)
		}
	}
	vPO.endTime = time.Now()

	playOrigReport(vPO, cfg)               // Print End of Run Stuff to either file or console
	if cfg.General.OutputTo != "console" { // Print End of Run Stuff again forcing to the console
		oW = os.Stdout
		playOrigReport(vPO, cfg)
	}
}
