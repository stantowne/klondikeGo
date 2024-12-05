package main

import (
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func playOrigReport(vPO variablesSpecificToPlayOrig, cfg *Configuration) {
	possibleAttempts := vPO.numberOfStrategies * cfg.General.NumberOfDecksToBePlayed
	lossCounter := cfg.General.NumberOfDecksToBePlayed - vPO.winCounter
	elapsedTime := vPO.endTime.Sub(vPO.startTime)
	percentageAttemptsAvoided := 100.0 * float64(vPO.attemptsAvoidedCounter) / float64(possibleAttempts)
	var p = message.NewPrinter(language.English)

	fmt.Fprintf(oW, "\n\nDate & Time Completed is %v\n", vPO.endTime)
	var err error

	_, err = p.Fprintf(oW, "Number of Decks Played is %d, starting with Deck %d.\n", cfg.General.NumberOfDecksToBePlayed, cfg.General.FirstDeckNum)
	if err != nil {
		fmt.Println("Number of Decks Played cannot print")
	}

	fmt.Fprintf(oW, "Length of Initial Override Strategies is %d.\n", cfg.PlayOrig.Length)

	fmt.Fprintf(oW, "Number of Initial Override Strategies Per Deck is %d.\n", vPO.numberOfStrategies)

	_, err = p.Fprintf(oW, "Number of Possible Attempts is %d.\n", possibleAttempts)
	if err != nil {
		fmt.Println("Number of Possible Attempts cannot print")
	}

	averageElapsedTimePerDeck := float64(elapsedTime.Milliseconds()) / float64(cfg.General.NumberOfDecksToBePlayed)

	fmt.Fprintf(oW, "Elapsed Time is %v.\n", elapsedTime)

	_, err = p.Fprintf(oW, "Total Decks Won is %d of which %d were Early Wins\n", vPO.winCounter, vPO.earlyWinCounter)
	if err != nil {
		fmt.Println("Total Decks Won cannot print")
	}

	_, err = p.Fprintf(oW, "Total Decks Lost is %d\n", lossCounter)
	if err != nil {
		fmt.Println("Total Decks Lost cannot print")
	}

	_, err = p.Fprintf(oW, "Losses at Game Length Limit is %d\n", vPO.lossesAtGLL)
	if err != nil {
		fmt.Println("Losses at Game Length Limit cannot print")
	}

	_, err = p.Fprintf(oW, "Losses at No Moves Available is %d\n", vPO.lossesAtNoMoves)
	if err != nil {
		fmt.Println("Losses at No Moves Available cannot print")
	}

	_, err = p.Fprintf(oW, "Regular Losses is %d\n", vPO.regularLosses)
	if err != nil {
		fmt.Println("Regular Losses cannot print")
	}

	_, err = p.Fprintf(oW, "Number of Attempts Avoided ia %d\n", vPO.attemptsAvoidedCounter)
	if err != nil {
		fmt.Println("Number of Attempts Avoided cannot print")
	}

	_, err = p.Fprintf(oW, "Percentage of Possible Attempts Avoided is %v\n", percentageAttemptsAvoided)
	if err != nil {
		fmt.Println("Percentage of Possible Attempts Avoided cannot print")
	}

	fmt.Fprintf(oW, "Average Elapsed Time per Deck is %vms.\n", averageElapsedTimePerDeck)
}
