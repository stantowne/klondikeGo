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

// Temp **********************
var PrintWinningMoves bool

// Temp **********************

type ConfigurationSubsetOnlyForSQLWriting struct { // STAN not sure we even need to create this it is simply here for me to communicatewhat needs to be written
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
		ProgressCounter int `yaml:"progress counter in millions"`
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

	fmt.Printf("\nCalling Program: %v\n\n", os.Args[0])
	fmt.Printf("\nRun Start Time: %15s\n\n", cfg.RunStartTime.Format("2006.01.02  3:04:05 pm"))
	_, err = pfmt.Printf("General:\n"+
		"            Number Of Decks To Be Played: %v\n"+
		"                      Starting with deck: %v\n"+
		"                            Type of Play: %v\n"+
		"                                 Verbose: %v\n\n",
		cfg.General.NumberOfDecksToBePlayed,
		cfg.General.FirstDeckNum,
		cfg.General.TypeOfPlay,
		cfg.General.Verbose)

	cfg.PlayAll.ProgressCounter *= 1_000_000

	cfg.PlayAll.ProgressCounterLastPrintTime = time.Now()

	if cfg.General.TypeOfPlay == "playOrig" {
		nOfS := 1 << cfg.PlayOrig.Length //number of initial strategies
		_, err = pfmt.Printf(" Style: Original iOS (Initial Override Strategies)\n\n"+
			"                     iOS strategy length: %v\n"+
			"          Max possible attempts per deck: %v\n"+
			"       Total possible attempts all decks: %v\n"+
			"                       Game Length Limit: %v\n\n",
			cfg.PlayOrig.Length,
			nOfS,
			nOfS*cfg.General.NumberOfDecksToBePlayed,
			cfg.PlayOrig.GameLengthLimit)
	} else {
		if cfg.PlayAll.ReportingType.NoReporting {
			_, err = pfmt.Printf("No Deck-by-Deck, Move-by-Move or Tree Reporting\n")
		} else {
			if cfg.PlayAll.ReportingType.DeckByDeck {
				_, err = pfmt.Printf("Deck By Deck Reporting: \n"+
					"                                           Type: %v\n",
					cfg.PlayAll.DeckByDeckReportingOptions.Type)
				if cfg.PlayAll.ProgressCounter != 0 {
					_, err = pfmt.Printf("    Move Progress Reporting Cycles, in Millions: %v\n", cfg.PlayAll.ProgressCounter)
				}
			}
			if cfg.PlayAll.ReportingType.MoveByMove {
				_, err = pfmt.Printf("Move By Move Reporting: \n"+
					"                                           Type: %v\n",
					cfg.PlayAll.MoveByMoveReportingOptions.Type)
				if cfg.PlayAll.ProgressCounter != 0 {
					_, err = pfmt.Printf("    Move Progress Reporting Cycles, in Millions: %v\n", cfg.PlayAll.ProgressCounter)
				}
			}
			if cfg.PlayAll.ReportingType.Tree {
				_, err = pfmt.Printf("Tree Reporting: \n"+
					"                        Type: %v\n"+
					"         TreeSleepBetwnMoves: %v\n"+
					"    TreeSleepBetwnStrategies: %v\n",

					cfg.PlayAll.MoveByMoveReportingOptions.Type,
					cfg.PlayAll.TreeReportingOptions.TreeSleepBetwnMoves,
					cfg.PlayAll.TreeReportingOptions.TreeSleepBetwnStrategies)
				if cfg.PlayAll.ProgressCounter != 0 {
					_, err = pfmt.Printf("    Move Progress Reporting Cycles, in Millions: %v\n", cfg.PlayAll.ProgressCounter)
				}
			}
			if cfg.PlayAll.RestrictReporting {
				_, err = pfmt.Printf("\nReporting Restricted To\n"+
					"                       Staring with Deck: %v\n"+
					"                            Continue for: %v decks (0 = all the rest)\n"+
					"             Starting with Moves Tried #: %v\n"+
					"                            Continue for: %v moves tried (0 = all the rest)\n",
					cfg.PlayAll.RestrictReportingTo.DeckStartVal,
					cfg.PlayAll.RestrictReportingTo.DeckContinueFor,
					cfg.PlayAll.RestrictReportingTo.MovesTriedStartVal,
					cfg.PlayAll.RestrictReportingTo.MovesTriedContinueFor)
			}
		}
	}
	if cfg.General.TypeOfPlay == "playAll" {
		_, err = pfmt.Printf("\nGame Length Limit, in millions: %v\n",
			cfg.PlayAll.GameLengthLimit)
	}

	// ******************************************
	//
	// Done printing out the configuration:
	//
	// Set up the code for reading the decks skipping to the firstDeckNum if not 0
	//
	// ******************************************

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
