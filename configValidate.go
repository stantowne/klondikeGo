package main

import (
	"fmt"
	"os"
)

// validate cfg struct
// bool fields not validated because an attempt to assign a non-bool value to a bool variable causes panic
func configValidate(c Configuration) {
	//
	//General section
	if !(c.General.Decks == "consecutive" || c.General.Decks == "list") {
		println("General.Decks invalid; must be either 'consecutive' or 'list'")
		defer os.Exit(1)
	}
	// this if statement and the next should be changed in the input file of decks contains greater or fewer than 10,000 decks
	if c.General.FirstDeckNum < 0 || c.General.FirstDeckNum > 410000 {
		println("General.FirstDeckNum invalid; must be non-negative integer less than 410,000")
		defer os.Exit(1)
	}
	if c.General.NumberOfDecksToBePlayed < 1 || c.General.NumberOfDecksToBePlayed > (410000-c.General.FirstDeckNum) {
		println("General.numberOfDecksToBePlayed invalid; must be 1 or more, but not more than 410,000 minus firstDeckNum")
		defer os.Exit(1)
	}
	if !(c.General.TypeOfPlay == "playOrig" || c.General.TypeOfPlay == "playAll") {
		println("General.TypeOfPlay invalid; must be either 'playOrig' or 'playAll'")
		defer os.Exit(1)
	}
	if c.General.Verbose >= 10 || c.General.Verbose < 0 {
		println("General.Verbose invalid; must be a non-negative integer no greater than 10")
		defer os.Exit(1)
	}
	//
	//PlayOrig Section
	if c.PlayOrig.Length < -1 || c.PlayOrig.Length > 24 {
		// in the line above, 24 is set arbitrarily; 24 would result in 16,777,216 possible attempts per deck
		// length could be greater, but the run time nearly doubles between n and n+1
		// I have never set Length to be greater than 16
		println("PlayOrig.Length invalid; must be a non-negative integer no greater than 24")
		defer os.Exit(1)
	}
	if c.PlayOrig.GameLengthLimit < 1 {
		println("gameLengthLimit invalid; must be a positive integer")
		defer os.Exit(1)
	}
	//
	//PlayAll Section
	if c.PlayAll.GameLengthLimit < 1 {
		println("gameLengthLimit invalid; must be a positive integer")
		os.Exit(1)
	}
	r := 0
	if c.PlayAll.ReportingType.DeckByDeck {
		r++
	}
	if c.PlayAll.ReportingType.MoveByMove {
		r++
	}
	if c.PlayAll.ReportingType.Tree {
		r++
	}
	if r > 1 {
		println("No more than one PlayAll reporting type allowed")
		defer os.Exit(1)
	}
	if c.PlayAll.ReportingType.MoveByMove &&
		!(c.PlayAll.MoveByMoveReportingOptions.Type == "regular" ||
			c.PlayAll.MoveByMoveReportingOptions.Type == "short" ||
			c.PlayAll.MoveByMoveReportingOptions.Type == "very short") {
		fmt.Println("Move By Move Reporting Type invalid; must be 'regular' 'short' or 'very short'")
		defer os.Exit(1)
	}
	if c.PlayAll.ReportingType.DeckByDeck &&
		!(c.PlayAll.DeckByDeckReportingOptions.Type == "regular" ||
			c.PlayAll.DeckByDeckReportingOptions.Type == "short" ||
			c.PlayAll.DeckByDeckReportingOptions.Type == "very short") {
		fmt.Println("DeckByDeck Reporting Type invalid; must be 'regular' 'short' or 'very short'")
		defer os.Exit(1)
	}
	if c.PlayAll.ReportingType.Tree &&
		!(c.PlayAll.TreeReportingOptions.Type == "regular" ||
			c.PlayAll.TreeReportingOptions.Type == "narrow" ||
			c.PlayAll.TreeReportingOptions.Type == "very narrow") {
		fmt.Println("DeckByDeck Reporting Type invalid; must be 'regular' 'narrow' or 'very narrow'")
		defer os.Exit(1)
	}
	if c.PlayAll.ReportingType.Tree &&
		((c.PlayAll.TreeReportingOptions.TreeSleepBetwnMoves < 0) ||
			(c.PlayAll.TreeReportingOptions.TreeSleepBetwnStrategies < 0)) {
		fmt.Println("Time Between Moves or Time Between Strategies invalid; both must be non-negative  integers")
		defer os.Exit(1)

	}
	if c.PlayAll.ProgressCounter < 0 {
		println("PlayAll.ProgressCounter invalid; must be a non-negative integer")
	}
}
