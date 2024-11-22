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
var singleGame bool // = true

type ConfigurationSubsetForSQLWriting struct { // STAN not sure we even need to create this it is simply here for me to communicatewhat needs to be written
	RunStartTime time.Time
	GitVersion   string // Stan we need to figure out how to get this
	General      struct {
		Verbose  int    `yaml:"verbose"`
		OutputTo string `yaml:"outputTo"`
	}
	PlayAll struct {
		GameLengthLimit  int  `yaml:"game length limit in moves tried"`
		FindAllWinStrats bool `yaml:"find all winning strategies?"`
		ReportingType    struct {
			DeckByDeck bool `yaml:"deck by deck"` // referred to as "DbD_R", "DbD_S" or "DbD_VS", in calls to prntMDet and calls thereto
			MoveByMove bool `yaml:"move by move"` // referred to as "MbM_R", "MbM_S" or "MbM_VS", in calls to prntMDet and calls thereto
			Tree       bool `yaml:"tree"`         // referred to as "Tree_R", "Tree_N" or "Tree_VN", in calls to prntMDet and calls thereto
		} `yaml:"reporting"`
		DeckByDeckReportingOptions struct {
			Type string `yaml:"type"`
		} `yaml:"deck by deck reporting options"`
		MoveByMoveReportingOptions struct {
			Type string `yaml:"type"`
		} `yaml:"move by move reporting options"`
		TreeReportingOptions struct {
			Type                     string `yaml:"type"`
			TreeSleepBetwnMoves      int    `yaml:"sleep between moves"`
			TreeSleepBetwnStrategies int    `yaml:"sleep between strategies"`
		} `yaml:"tree reporting options"`
		RestrictReporting   bool //not part of yaml file, derived after yaml file is unmarshalled & validated
		RestrictReportingTo struct {
			DeckStartVal          int `yaml:"starting deck number"`
			DeckContinueFor       int `yaml:"continue for how many decks"`
			MovesTriedStartVal    int `yaml:"starting move number"`
			MovesTriedContinueFor int `yaml:"continue for how many moves"`
		} `yaml:"restrict reporting to"`
		PrintWinningMoves bool `yaml:"print winning moves"`
		ProgressCounter   int  `yaml:"progress counter in millions"`
	} `yaml:"play all moves"`
}

func main() {
	// Temp **********************
	PrintWinningMoves = true
	// Temp **********************

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

	validateConfig(cfg)

	if cfg.General.NumberOfDecksToBePlayed == 1 {
		singleGame = true
	}

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

	// ******************************************
	//
	// Print out the configuration:
	//
	// Always a good idea to print out the program source of the output.
	//    Will be especially useful when we figure out how to include the versioning

	printConfig(cfg)
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
