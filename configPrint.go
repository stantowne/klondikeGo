package main

func printConfig(c Configuration) {
	_, _ = pfmt.Printf("General:\n"+
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
		_, _ = pfmt.Printf(" Style: Original iOS (Initial Override Strategies)\n\n"+
			"                     iOS strategy length: %v\n"+
			"          Max possible attempts per deck: %v\n"+
			"       Total possible attempts all decks: %v\n"+
			"                       Game Length Limit: %v\n\n",
			c.PlayOrig.Length,
			nOfS,
			nOfS*c.General.NumberOfDecksToBePlayed,
			c.PlayOrig.GameLengthLimit)
	} else {
		if c.PlayAll.ReportingType.NoReporting {
			_, _ = pfmt.Printf("No Deck-by-Deck, Move-by-Move or Tree Reporting\n")
		} else {
			if c.PlayAll.ReportingType.DeckByDeck {
				_, _ = pfmt.Printf("Deck By Deck Reporting: \n"+
					"                                           Type: %v\n",
					c.PlayAll.DeckByDeckReportingOptions.Type)
				if c.PlayAll.ProgressCounter != 0 {
					_, _ = pfmt.Printf("    Move Progress Reporting Cycles, in Millions: %v\n", c.PlayAll.ProgressCounter)
				}
			}
			if c.PlayAll.ReportingType.MoveByMove {
				_, _ = pfmt.Printf("Move By Move Reporting: \n"+
					"                                           Type: %v\n",
					c.PlayAll.MoveByMoveReportingOptions.Type)
				if c.PlayAll.ProgressCounter != 0 {
					_, _ = pfmt.Printf("    Move Progress Reporting Cycles, in Millions: %v\n", c.PlayAll.ProgressCounter)
				}
			}
			if c.PlayAll.ReportingType.Tree {
				_, _ = pfmt.Printf("Tree Reporting: \n"+
					"                        Type: %v\n"+
					"         TreeSleepBetwnMoves: %v\n"+
					"    TreeSleepBetwnStrategies: %v\n",

					c.PlayAll.MoveByMoveReportingOptions.Type,
					c.PlayAll.TreeReportingOptions.TreeSleepBetwnMoves,
					c.PlayAll.TreeReportingOptions.TreeSleepBetwnStrategies)
				if c.PlayAll.ProgressCounter != 0 {
					_, _ = pfmt.Printf("    Move Progress Reporting Cycles, in Millions: %v\n", c.PlayAll.ProgressCounter)
				}
			}
			if c.PlayAll.RestrictReporting {
				_, _ = pfmt.Printf("\nReporting Restricted To\n"+
					"                       Staring with Deck: %v\n"+
					"                            Continue for: %v decks (0 = all the rest)\n"+
					"             Starting with Moves Tried #: %v\n"+
					"                            Continue for: %v moves tried (0 = all the rest)\n",
					c.PlayAll.RestrictReportingTo.DeckStartVal,
					c.PlayAll.RestrictReportingTo.DeckContinueFor,
					c.PlayAll.RestrictReportingTo.MovesTriedStartVal,
					c.PlayAll.RestrictReportingTo.MovesTriedContinueFor)
			}
		}
	}
	if c.General.TypeOfPlay == "playAll" {
		_, _ = pfmt.Printf("\nGame Length Limit, in millions: %v\n",
			c.PlayAll.GameLengthLimit)
	}

}
