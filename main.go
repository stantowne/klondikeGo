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
	"strings"
	"time"
)

var moveNumMax int

// Setup pfmt to print thousands with commas
var pfmt = message.NewPrinter(language.English)
var err error
var singleGame bool // = true

type Configuration struct {
	General struct {
		Decks                        string `yaml:"decks"`                        // must be "consecutive" or "list"
		FirstDeckNum                 int    `yaml:"first deck number"`            // must be non-negative integer
		NumberOfDecksToBePlayed      int    `yaml:"number of decks to be played"` //must be non-negative integer
		List                         []int
		TypeOfPlay                   string `yaml:"type of play"` // must be "playOrig" or "playAll"
		Verbose                      int    `yaml:"verbose"`
		OutputTo                     string `yaml:"outputTo"`
		ProgressCounter              int    `yaml:"progress counter in millions"`
		ProgressCounterLastPrintTime time.Time
	} `yaml:"general"`
	PlayOrig struct {
		Length          int `yaml:"length of initial override strategy"`
		GameLengthLimit int `yaml:"game length limit moveCounter"`
	} `yaml:"play original"`
	PlayAll struct {
		GameLengthLimit  int  `yaml:"game length limit in moves tried"`
		FindAllWinStrats bool `yaml:"find all winning strategies?"`
		ReportingType    struct {
			DeckByDeck  bool `yaml:"deck by deck"` // referred to as "DbD_R", "DbD_S" or "DbD_VS", in calls to prntMDet and calls thereto
			MoveByMove  bool `yaml:"move by move"` // referred to as "MbM_R", "MbM_S" or "MbM_VS", in calls to prntMDet and calls thereto
			Tree        bool `yaml:"tree"`         // referred to as "Tree_R", "Tree_N" or "Tree_VN", in calls to prntMDet and calls thereto
			NoReporting bool //not part of yaml file, derived after yaml file is unmarshalled & validated   CONSIDER DELETING
		} `yaml:"reporting"`
		DeckByDeckReportingOptions struct {
			Type string `yaml:"type"`
		} `yaml:"deck by deck reporting options"`
		MoveByMoveReportingOptions struct {
			Type string `yaml:"type"`
		} `yaml:"move by move reporting options"`
		TreeReportingOptions struct {
			Type                        string `yaml:"type"`
			TreeSleepBetwnMoves         int    `yaml:"sleep between moves"`
			TreeSleepBetwnMovesDur      time.Duration
			TreeSleepBetwnStrategies    int `yaml:"sleep between strategies"`
			TreeSleepBetwnStrategiesDur time.Duration
		} `yaml:"tree reporting options"`
		RestrictReporting   bool //not part of yaml file, derived after yaml file is unmarshalled & validated
		RestrictReportingTo struct {
			DeckStartVal          int `yaml:"starting deck number"`
			DeckContinueFor       int `yaml:"continue for how many decks"`
			MovesTriedStartVal    int `yaml:"starting move number"`
			MovesTriedContinueFor int `yaml:"continue for how many moves"`
		} `yaml:"restrict reporting to"`
		WinLossReport       bool   `yaml:"print final deck by deck win loss record"`
		SaveResultsToSQL    bool   `yaml:"save results to SQL"`
		SQLConnectionString string `yaml:"sql connection string"`
	} `yaml:"play all moves"`
}

func main() {

	// unmarshal YAML file
	cfg := Configuration{}
	data, err3 := os.ReadFile("./config.yml") // err3 used to avoid shadowing err
	if err3 != nil {
		panic(err3)
	}
	if err4 := yaml.Unmarshal(data, &cfg); err4 != nil {
		panic(err4)
	} // err4 used to avoid shadowing err

	// validate cfg after unmarshal
	// BEFORE validation, make all strings in cfg EXCEPT cfg.PlayAll.SQLConnectionString lower case
	cfg.General.TypeOfPlay = strings.ToLower(cfg.General.TypeOfPlay)
	cfg.PlayAll.DeckByDeckReportingOptions.Type = strings.ToLower(cfg.PlayAll.DeckByDeckReportingOptions.Type)
	cfg.PlayAll.MoveByMoveReportingOptions.Type = strings.ToLower(cfg.PlayAll.MoveByMoveReportingOptions.Type)
	cfg.PlayAll.TreeReportingOptions.Type = strings.ToLower(cfg.PlayAll.TreeReportingOptions.Type)
	if strings.ToLower(cfg.General.OutputTo) == "console" {
		cfg.General.OutputTo = strings.ToLower(cfg.General.OutputTo)
	}

	// validate cfg after unmarshal and conversion to lower case
	if cfg.General.FirstDeckNum < 0 || cfg.General.FirstDeckNum > 9999 {
		println("FirstDeckNum invalid")
		println("FirstDeckNum must be non-negative integer less than 10,000")
		os.Exit(1)
	}

	//the line below should be changed if the input file contains more than 10,000 decks
	if cfg.General.NumberOfDecksToBePlayed < 1 || cfg.General.NumberOfDecksToBePlayed > (10000-cfg.General.FirstDeckNum) {
		println("numberOfDecksToBePlayed invalid")
		println("numberOfDecksToBePlayed must be 1 or more, but not more than 10,000 minus firstDeckNum")
		os.Exit(1)
	}

	if cfg.General.NumberOfDecksToBePlayed == 1 {
		singleGame = true
	}

	// 		In the line below, 24 is arbitrarily set; 24 would result in 16,777,216 attempts per deck
	// 		I have never run the program with length greater than 16
	// 		Depending upon the size of an int, length could be 32 or greater, but the program may never finish
	if cfg.PlayOrig.Length < -1 || cfg.PlayOrig.Length > 24 {
		println("length invalid")
		println("Length must be an integer > -1 and <= 24")
		os.Exit(1)
	}

	if cfg.General.Verbose >= 10 || cfg.General.Verbose < 0 {
		println("verbose invalid")
		println("verbose must be a non-negative integer no greater than 10")
		os.Exit(1)
	}

	/*

				Additional cfg validations needed:

				1. cfg.PlayAll.ReportingType.xxx no more than 1 true
				2. cfg.PlayAll.xxxReportingOptions.Type  are valid
				3. TreeSleepBetwnMoves and TreeSleepBetwnStrategies non negative integer

		        I checked all of the commented out validations from the former command line arguments (which I have now deleted)
		        and they are all now included in the above or section below (Search for "start ProgressCounterOverRides")

			    Can you catch non-numeric entry into integer in yaml without a panic or who cares?
			    Can you catch fractional entry into integer in yaml without a panic or who cares?
		        Can you catch non boolean entry into integer in yaml without a panic or who cares?
	*/

	// completing cfg

	cfg.PlayAll.TreeReportingOptions.TreeSleepBetwnMovesDur = time.Duration(cfg.PlayAll.TreeReportingOptions.TreeSleepBetwnMoves*100) * time.Millisecond
	cfg.PlayAll.TreeReportingOptions.TreeSleepBetwnStrategiesDur = time.Duration(cfg.PlayAll.TreeReportingOptions.TreeSleepBetwnStrategies*100) * time.Millisecond

	// RestrictReporting zero value is false
	if cfg.PlayAll.RestrictReportingTo.DeckStartVal != 0 ||
		cfg.PlayAll.RestrictReportingTo.DeckContinueFor != 0 ||
		cfg.PlayAll.RestrictReportingTo.MovesTriedStartVal != 0 ||
		cfg.PlayAll.RestrictReportingTo.MovesTriedContinueFor != 0 {
		cfg.PlayAll.RestrictReporting = true
	}
	cfg.PlayAll.ReportingType.NoReporting = !(cfg.PlayAll.ReportingType.DeckByDeck || cfg.PlayAll.ReportingType.MoveByMove || cfg.PlayAll.ReportingType.Tree)

	// ******************************************
	//
	// Print out the configuration:
	//
	// Always a good idea to print out the program source of the output.
	//    Will be especially useful when we figure out how to include the versioning

	fmt.Printf("\nCalling Program: %v\n\n", os.Args[0])
	fmt.Printf("\nRun Start Time: %15s\n\n", time.Now().Format("2006.01.02  3:04:05 pm"))
	_, err = pfmt.Printf("General:\n"+
		"            Number Of Decks To Be Played: %v\n"+
		"                      Starting with deck: %v\n"+
		"                            Type of Play: %v\n"+
		"                                 Verbose: %v\n\n",
		cfg.General.NumberOfDecksToBePlayed,
		cfg.General.FirstDeckNum,
		cfg.General.TypeOfPlay,
		cfg.General.Verbose)

	// DO BELOW edits to ProgressCounter so that it will be set to 0 when inappropriate and based on tests in prints will not print out
	//               start ProgressCounterOverRides
	if cfg.General.OutputTo == "console" {
		if cfg.General.TypeOfPlay == "playorig" {
			cfg.General.ProgressCounter = 0
		} else {
			if cfg.PlayAll.ReportingType.DeckByDeck {
				cfg.General.ProgressCounter *= 1_000_000
			} else {
				cfg.General.ProgressCounter = 0
			}
		}
	} else {
		cfg.PlayAll.TreeReportingOptions.TreeSleepBetwnMoves = 0      // Delete this if figure out how to print to file AND console
		cfg.PlayAll.TreeReportingOptions.TreeSleepBetwnStrategies = 0 // Delete this if figure out how to print to file AND console
		if cfg.General.TypeOfPlay == "playorig" {
			cfg.General.ProgressCounter = 0 // Change this to 1 if figure out how to print to file AND console
		} else {
			cfg.General.ProgressCounter = 0 // replace this with cfg.General.ProgressCounter *= 1_000_000 if figure out how to print to file AND console
		}
	}
	//               end ProgressCounterOverRides
	//DO ABOVE edits to ProgressCounter so that it will be set to 0 when inappropriate and based on tests in prints will not print out

	cfg.General.ProgressCounterLastPrintTime = time.Now()

	if cfg.General.TypeOfPlay == "playorig" {
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
		if cfg.General.ProgressCounter != 0 {
			// Note: will not ever print due to the overrides (Search for "start ProgressCounterOverRides")
			//       until figure out how to print to console and file and until code added to playOrig
			_, err = pfmt.Printf("    Move Progress Reporting Cycles, in Decks: %v\n", cfg.General.ProgressCounter)
		}
	} else {
		if cfg.PlayAll.ReportingType.NoReporting {
			_, err = pfmt.Printf("No Deck-by-Deck, Move-by-Move or Tree Reporting\n")
		} else {
			if cfg.PlayAll.ReportingType.DeckByDeck {
				_, err = pfmt.Printf("Deck By Deck Reporting: \n"+
					"                                           Type: %v\n",
					cfg.PlayAll.DeckByDeckReportingOptions.Type)
				if cfg.General.ProgressCounter != 0 {
					_, err = pfmt.Printf("    Move Progress Reporting Cycles, in Millions: %v\n", cfg.General.ProgressCounter)
				}
			}
			if cfg.PlayAll.ReportingType.MoveByMove {
				_, err = pfmt.Printf("Move By Move Reporting: \n"+
					"                                           Type: %v\n",
					cfg.PlayAll.MoveByMoveReportingOptions.Type)
				if cfg.General.ProgressCounter != 0 {
					_, err = pfmt.Printf("    Move Progress Reporting Cycles, in Millions: %v\n", cfg.General.ProgressCounter)
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
				if cfg.General.ProgressCounter != 0 {
					_, err = pfmt.Printf("    Move Progress Reporting Cycles, in Millions: %v\n", cfg.General.ProgressCounter)
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
	if cfg.General.TypeOfPlay == "playall" {
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
	file, err1 := os.Open(inputFileName) // err1 used to avoid shadowing err
	if err1 != nil {
		log.Println("Cannot open inputFileName:", err1)
	}
	defer func(file *os.File) {
		err2 := file.Close() // err2 used to avoid shadowing err
		if err2 != nil {
			println("could not close file:", err2)
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

	if cfg.General.TypeOfPlay == "playall" {
		moveBasePriority = moveBasePriorityNew
		playAll(*reader, &cfg)
	}
	if cfg.General.TypeOfPlay == "playorig" {
		moveBasePriority = moveBasePriorityOrig
		playOrig(*reader, &cfg)
	}
}
