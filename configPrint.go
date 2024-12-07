package main

import (
	"fmt"
	"os"
	"sort"
)

func configPrint(c Configuration) {

	_, _ = fmt.Fprintf(oW, "\nCalling Program: %v\n", os.Args[0])
	_, _ = fmt.Fprintf(oW, "\nRun Start Time: %15s\n\n", c.General.RunStartTime.Format("2006.01.02  3:04:05 pm"))

	_, _ = pfmt.Fprintf(oW, "General:\n"+
		"            Number Of Decks To Be Played: %v\n"+
		"                      Starting with deck: %v\n"+
		"                            Type of Play: %v\n"+
		"                                 Verbose: %v\n\n",
		c.General.NumberOfDecksToBePlayed,
		c.General.FirstDeckNum,
		c.General.TypeOfPlay,
		c.General.Verbose)

	if c.General.TypeOfPlay == "playOrig" {
		nOfS := 1 << c.PlayOrig.Length //number of initial strategies
		_, _ = pfmt.Fprintf(oW, " Style: Original iOS (Initial Override Strategies)\n\n"+
			"                     iOS strategy length: %v\n"+
			"          Max possible attempts per deck: %v\n"+
			"       Total possible attempts all decks: %v\n"+
			"                       Game Length Limit: %v\n\n",
			c.PlayOrig.Length,
			nOfS,
			nOfS*c.General.NumberOfDecksToBePlayed,
			c.PlayOrig.GameLengthLimit)

		_, _ = pfmt.Fprintf(oW, "\nMove Priority Settings:\n\n")
		moveTypes := make([]string, 0, len(moveBasePriorityOrig))

		for priority := range moveBasePriorityOrig {
			moveTypes = append(moveTypes, priority)
		}

		// sort by priority before printing
		sort.SliceStable(moveTypes, func(i, j int) bool {
			return moveBasePriorityAll[moveTypes[i]] < moveBasePriorityAll[moveTypes[j]]
		})
		for i, moveType := range moveTypes {
			_, _ = pfmt.Fprintf(oW, "   %2v   %17s: %5v\n", i, moveTypes[i], moveBasePriorityAll[moveType])
		}
		_, _ = pfmt.Fprintf(oW, "\n\n")
	} else {
		_, _ = pfmt.Fprintf(oW, "          Game Length Limit, in millions: %v\n\n", c.PlayAll.GameLengthLimit)
		if c.PlayAll.ReportingType.NoReporting {
			_, _ = pfmt.Fprintf(oW, "No Deck-by-Deck, Move-by-Move or Tree Reporting\n")
		} else {
			if c.PlayAll.ReportingType.DeckByDeck {
				_, _ = pfmt.Fprintf(oW, "Deck By Deck Reporting: \n"+
					"                  Type: %v\n",
					c.PlayAll.DeckByDeckReportingOptions.Type)
			}
			if c.PlayAll.ReportingType.MoveByMove {
				_, _ = pfmt.Fprintf(oW, "Move By Move Reporting: \n"+
					"                  Type: %v\n",
					c.PlayAll.MoveByMoveReportingOptions.Type)
			}
			if c.PlayAll.ReportingType.Tree {
				_, _ = pfmt.Fprintf(oW, "Tree Reporting: \n"+
					"                       Type: %v\n"+
					"        TreeSleepBetwnMoves: %v\n"+
					"   TreeSleepBetwnStrategies: %v\n",

					c.PlayAll.MoveByMoveReportingOptions.Type,
					c.PlayAll.TreeReportingOptions.TreeSleepBetwnMoves,
					c.PlayAll.TreeReportingOptions.TreeSleepBetwnStrategies)
			}
			if c.PlayAll.RestrictReporting {
				_, _ = pfmt.Fprintf(oW, "\nReporting Restricted To\n"+
					"               Staring with Deck: %v\n"+
					"                    Continue for: %v decks (0 = all the rest)\n"+
					"     Starting with Moves Tried #: %v\n"+
					"                    Continue for: %v moves tried (0 = all the rest)\n",
					c.PlayAll.RestrictReportingTo.DeckStartVal,
					c.PlayAll.RestrictReportingTo.DeckContinueFor,
					c.PlayAll.RestrictReportingTo.MovesTriedStartVal,
					c.PlayAll.RestrictReportingTo.MovesTriedContinueFor)
			}
		}
		_, _ = pfmt.Fprintf(oW, "\nPrint Winning Moves: %v\n", c.PlayAll.PrintWinningMoves)
		_, _ = pfmt.Fprintf(oW, "Move Progress Reporting Cycles, in Millions: %-5v\n", c.PlayAll.ProgressCounter)
		_, _ = pfmt.Fprintf(oW, "Print final DbD W/L record: %v\n", c.PlayAll.WinLossReport)
		_, _ = pfmt.Fprintf(oW, "Save results to SQL: %v\n", c.PlayAll.SaveResultsToSQL)
		if c.PlayAll.SaveResultsToSQL {
			_, _ = pfmt.Fprintf(oW, "   SQL connection string: %v\n", c.PlayAll.SQLConnectionString)
		}

		_, _ = pfmt.Fprintf(oW, "\n\nMove Priority Settings:\n\n")
		moveTypes := make([]string, 0, len(moveBasePriorityAll))
		for priority := range moveBasePriorityAll {
			moveTypes = append(moveTypes, priority)
		}

		// sort by priority before printing
		sort.SliceStable(moveTypes, func(i, j int) bool {
			return moveBasePriorityAll[moveTypes[i]] < moveBasePriorityAll[moveTypes[j]]
		})
		for i, moveType := range moveTypes {
			_, _ = pfmt.Fprintf(oW, "   %2v   %17s: %5v\n", i, moveTypes[i], moveBasePriorityAll[moveType])
		}
		_, _ = pfmt.Fprintf(oW, "\n\n")
	}

}
