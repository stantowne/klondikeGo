package main

import (
	"encoding/csv"
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"runtime/debug"
	"time"
)

// Setup pfmt to print thousands with commas
var pfmt = message.NewPrinter(language.English)

// Create the Short named package variable "oW" for cfg.General.outWriter
var oW *os.File

func main() {
	// unmarshal YAML file
	cfg := Configuration{}
	cfg.General.RunStartTime = time.Now()
	cfg.General.GitVersion = "" // Stan we need to figure out how to get this
	data, err := os.ReadFile("./config.yml")
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}
	//I need to confirm that, by virtue of no error being returned, we know that all bools have valid values.

	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				cfg.General.GitVersion = setting.Value // Stan we need to figure out how to get this
				break
			}
		}
	}
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Printf("Error: %v getting Hostname", err)
		os.Exit(1)
	}
	cfg.General.GitSystem = hostname

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

	cfg.General.outWriter = os.Stdout
	oW = cfg.General.outWriter
	configPrint(cfg) // Print FIRST time to stout

	if cfg.General.OutputTo != "console" {
		// create file
		cfg.General.outWriter, err = os.Create(cfg.General.OutputTo + "__" + cfg.General.RunStartTime.Format("2006.01.02_15.04.05_-0700") + ".txt")
		if err != nil {
			fmt.Printf("Error: %v  Error creating output file: %s", err, cfg.General.OutputTo+"__"+cfg.General.RunStartTime.Format("__2006.01.02_15.04.05_-0700")+".txt")
			os.Exit(1)
		}
		oW = cfg.General.outWriter
		// remember to close the file
		defer cfg.General.outWriter.Close()
		configPrint(cfg) // Print SECOND time to file
	}
	// Fill the Short named package variable "oW" for cfg.General.outWriter

	cfg.PlayAll.ProgressCounter *= 1_000_000
	cfg.PlayAll.ProgressCounterLastPrintTime = time.Now()

	inputFileName := cfg.General.DeckFileName
	file, err := os.Open(inputFileName)
	if err != nil {
		fmt.Println("Error: %v  Cannot open Deck inputFileName: %s", err, cfg.General.DeckFileName)
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
				fmt.Println("Error: %v  Cannot read from Deck inputFileName: %s", err, cfg.General.DeckFileName)
			}
		}
	}

	if cfg.General.TypeOfPlay == "playAll" {
		moveBasePriority = moveBasePriorityAll
		playAll(*reader, &cfg)
	}
	if cfg.General.TypeOfPlay == "playOrig" {
		moveBasePriority = moveBasePriorityOrig
		playOrig(*reader, &cfg)
	}
}
