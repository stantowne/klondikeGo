package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func playAll(reader csv.Reader, cfg *Configuration) {
	firstDeckNum := cfg.General.FirstDeckNum                       // Shorthand name but really is a copy - OK since never changed (but would Pointer or address be better?)
	numberOfDecksToBePlayed := cfg.General.NumberOfDecksToBePlayed // Shorthand name but really is a copy - OK since never changed (but would Pointer or address be better?)
	var vPA variablesSpecificToPlayAll
	vPA.TD.stratTried = 1
	vPA.TDother.startTime = time.Now()
	vPA.ADother.startTime = time.Now()

	for deckNum := firstDeckNum; deckNum < (firstDeckNum + numberOfDecksToBePlayed); deckNum++ {
		if !cfg.PlayAll.ReportingType.NoReporting && !(cfg.PlayAll.ReportingType.DeckByDeck && cfg.PlayAll.DeckByDeckReportingOptions.Type != "regular") {
			_, _ = fmt.Fprintf(oW, "\n\n******************************************************************************************************\n")
		}
		protoDeck, err := reader.Read() // protoDeck is a slice of strings: rank, suit, rank, suit, etc.
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Cannot read from inputFileName:", err)
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
		//temp		AllMvStratNum := 0
		var b = dealDeck(d)

		/* *************************************************************************************

		OK lets try to play this deck!

		This statement is executed once per deck and transfers program execution.
		When this statement returns the deck has been played.

		****************************************************************************************
		*/

		result1, result2 := playAllMoves(b, 0, deckNum, cfg, &vPA)

		/* *************************************************************************************

		The deck has been played

		All the rest is printing and preparation for the next deck

		****************************************************************************************
		*/

		if vPA.TD.stratLossesGLE > 0 {
			vPA.TD.stratLossesGLEAb = result2 - 2
		}
		vPA.TD.elapsedTime = time.Since(vPA.TDother.startTime)
		vPA.TD.winningMovesCnt = len(vPA.TDother.winningMoves)
		var dummy []move

		vPA.TD.unqBoards = len(vPA.TDother.priorBoards)
		var s1 string
		s2 := "   DECK"
		s1 = strconv.Itoa(deckNum)
		if result1 == "SW" {
			vPA.ADother.decksWon += 1
			s2 += " WON "
		} else {
			if vPA.TD.stratLossesGLE > 0 {
				s2 += " GLE "
				vPA.ADother.decksLostGLE += 1
			} else {
				vPA.TD.stratLosses++
				vPA.ADother.decksLost += 1
				s2 += " LOST"
			}
		}
		prntMDet(b, dummy, 1, deckNum, 1, "MbM_ANY", 2, "\n   "+s2+"\n", "", "", cfg, &vPA)
		prntMDetTreeReturnComment("\n   "+s2+"\n", deckNum, 0, cfg, &vPA)

		// End of Deck Statistics Reporting
		prntMDet(b, dummy, 0, deckNum, 1, "ANY", 0, s1, s2, result1, cfg, &vPA)

		// Reverse the collected winning moves and print them if needed
		if (cfg.PlayAll.SaveResultsToSQL || cfg.PlayAll.PrintWinningMoves) && vPA.TD.winningMovesCnt != 0 {
			PrintWinningMoves(cfg, &vPA)
		}

		if cfg.PlayAll.SaveResultsToSQL {
			// write ConfigurationSubsetOnlyForSQLWriting and vPA.TD out to sql/csv here
		}

		// Done printing and SQL writing for this deck
		// collect statistics into all deck variables (AD) and clear this deck (TD) variables for the next deck
		if vPA.ADother.moveNumMax == 0 || vPA.ADother.moveNumMax < vPA.TDotherSQL.moveNumMax {
			vPA.ADother.moveNumMax = vPA.TDotherSQL.moveNumMax
		}
		if vPA.TD.stratWins > 0 && (vPA.ADother.moveNumAtWinMin == 0 || vPA.ADother.moveNumAtWinMin > vPA.TDotherSQL.moveNumAtWin) {
			vPA.ADother.moveNumAtWinMin = vPA.TDotherSQL.moveNumAtWin
		}
		if vPA.TD.stratWins > 0 && (vPA.ADother.moveNumAtWinMax == 0 || vPA.ADother.moveNumAtWinMax < vPA.TDotherSQL.moveNumAtWin) {
			vPA.ADother.moveNumAtWinMax = vPA.TDotherSQL.moveNumAtWin
		}

		vPA.TDotherSQL.moveNumAtWin = 0
		vPA.TDotherSQL.moveNumMax = 0 //to keep track of length of the longest strategy so far

		vPA.AD.mvsTried += vPA.TD.mvsTried
		vPA.TD.mvsTried = 0
		vPA.AD.stratNum += vPA.TD.stratNum
		vPA.TD.stratNum = 0
		vPA.AD.stratTried += vPA.TD.stratTried
		vPA.TD.stratTried = 1 // NOTE: Starts at 1 not 0
		vPA.AD.stratWins += vPA.TD.stratWins
		vPA.TD.stratWins = 0
		vPA.AD.stratLosses += vPA.TD.stratLosses
		vPA.TD.stratLosses = 0
		vPA.AD.stratLossesGLE += vPA.TD.stratLossesGLE
		vPA.TD.stratLossesGLE = 0
		vPA.AD.stratLossesGLEAb += vPA.TD.stratLossesGLEAb
		vPA.TD.stratLossesGLEAb = 0
		vPA.AD.stratLossesNMA += vPA.TD.stratLossesNMA
		vPA.TD.stratLossesNMA = 0
		vPA.AD.stratLossesRB += vPA.TD.stratLossesRB
		vPA.TD.stratLossesRB = 0
		vPA.AD.stratLossesMajSE += vPA.TD.stratLossesMajSE
		vPA.TD.stratLossesMajSE = 0
		vPA.AD.stratLossesMinSE += vPA.TD.stratLossesMinSE
		vPA.TD.stratLossesMinSE = 0
		vPA.AD.stratLossesEL += vPA.TD.stratLossesEL
		vPA.TD.stratLossesEL = 0
		vPA.AD.unqBoards += vPA.TD.unqBoards
		vPA.TD.unqBoards = 0
		vPA.AD.winningMovesCnt += vPA.TD.winningMovesCnt
		vPA.TD.winningMovesCnt = 0
		vPA.AD.elapsedTime += vPA.TD.elapsedTime

		vPA.ADother.decksPlayed++
		vPA.TDother.treePrevMoves = ""
		vPA.TDother.winningMoves = nil
		vPA.TDother.startTime = time.Now()
		clear(vPA.TDother.priorBoards)
		vPA.TD.stratNum = 0
		vPA.TD.stratTried = 1 // NOTE: Starts at 1 not 0
	}

	// At this point, all decks to be played have been played.  Time to report aggregate won loss.
	// From this point on, the program only prints.

	printSummaryStats(cfg, &vPA)           // Print End of Run Stuff to either file or console
	if cfg.General.OutputTo != "console" { // Print End of Run Stuff again forcing to the console
		oW = os.Stdout
		printSummaryStats(cfg, &vPA)
	}

	if cfg.PlayAll.WinLossReport { // Deck Win Loss Summary Statistics
		// Close sql/csv file for writing and open it for reading and report it here
	}
}
