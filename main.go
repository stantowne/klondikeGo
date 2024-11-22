package main

import (
	"encoding/csv"
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"time"
)

// Setup pfmt to print thousands with commas
var pfmt = message.NewPrinter(language.English)

func main() {
	// unmarshal YAML file
	cfg := Configuration{}
	cfg.RunStartTime = time.Now()
	cfg.GitVersion = "" // Stan we need to figure out how to get this
	data, err := os.ReadFile("./config.yml")
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}
	//I need to confirm that, by virtue of no error being returned, we know that all bools have valid values.
	cfg.PlayAll.ReportingType.NoReporting =
		!(cfg.PlayAll.ReportingType.DeckByDeck ||
			cfg.PlayAll.ReportingType.MoveByMove ||
			cfg.PlayAll.ReportingType.Tree)

	configValidate(cfg)

	cfg.PlayAll.TreeReportingOptions.TreeSleepBetwnMovesDur =
		time.Duration(cfg.PlayAll.TreeReportingOptions.TreeSleepBetwnMoves*100) * time.Millisecond
	cfg.PlayAll.TreeReportingOptions.TreeSleepBetwnStrategiesDur =
		time.Duration(cfg.PlayAll.TreeReportingOptions.TreeSleepBetwnStrategies*100) * time.Millisecond
	// Null value RestrictReportingTo is false
	if cfg.PlayAll.RestrictReportingTo.DeckStartVal != 0 ||
		cfg.PlayAll.RestrictReportingTo.DeckContinueFor != 0 ||
		cfg.PlayAll.RestrictReportingTo.MovesTriedStartVal != 0 ||
		cfg.PlayAll.RestrictReportingTo.MovesTriedContinueFor != 0 {
		cfg.PlayAll.RestrictReporting = true
	}
	configPrint(cfg)
	cfg.PlayAll.ProgressCounter *= 1_000_000
	cfg.PlayAll.ProgressCounterLastPrintTime = time.Now()
	fmt.Printf("\nCalling Program: %v\n\n", os.Args[0])
	fmt.Printf("\nRun Start Time: %15s\n\n", cfg.RunStartTime.Format("2006.01.02  3:04:05 pm"))

	inputFileName := "decks-made-2022-01-15_count_10000-dict.csv"
	file, err := os.Open(inputFileName)
	if err != nil {
		log.Println("Cannot open inputFileName:", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			println("could not close file:", err)
		}
	}(file)
	reader := csv.NewReader(file)

	if cfg.General.FirstDeckNum > 0 {
		for i := 0; i < cfg.General.FirstDeckNum; i++ {
			_, err = reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println("Cannot read from inputFileName", err)
			}
		}
	}

	if cfg.General.TypeOfPlay == "playAll" {
		moveBasePriority = moveBasePriorityNew
		playAll(*reader, &cfg)
	}
	if cfg.General.TypeOfPlay == "playOrig" {
		moveBasePriority = moveBasePriorityOrig
		playOrig(*reader, &cfg)
	}
}
