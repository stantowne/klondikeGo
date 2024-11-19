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
		GameLengthLimit int `yaml:"game length limit in moves"`
	} `yaml:"play original"`
	PlayNew struct {
		GameLengthLimit     int  `yaml:"game length limit in moves tried"`
		FindAllWinStrats    bool `yaml:"find all winning strategies?"`
		ReportingDeckByDeck bool //not part of yaml file, derived after yaml file is unmarshalled & validated
		ReportingMoveByMove bool //not part of yaml file, derived after yaml file is unmarshalled & validated
		ReportingType       struct {
			DeckByDeck  bool `yaml:"deck by deck"` // referred to as "DbD_R", "DbD_S" or "DbD_VS", in calls to prntMDet and calls thereto
			MoveByMove  bool `yaml:"move by move"` // referred to as "MbM_R", "MbM_S" or "MbM_VS", in calls to prntMDet and calls thereto
			Tree        bool `yaml:"tree"`         // referred to as "Tree_R", "Tree_N" or "Tree_VN", in calls to prntMDet and calls thereto
			NoReporting bool //not part of yaml file, derived after yaml file is unmarshalled & validated   CONSIDER DELETING
		} `yaml:"reporting"`
		DeckByDeckReportingOptions struct {
			Type string `yaml:"typeDbD"`
		} `yaml:"deck by deck reporting options"`
		MoveByMoveReportingOptions struct {
			Type string `yaml:"typeMbM"`
		} `yaml:"move by move reporting options"`
		TreeReportingOptions struct {
			Type                     string        `yaml:"typeTree"`
			TreeSleepBetwnMoves      time.Duration `yaml:"sleep between moves"`
			TreeSleepBetwnStrategies time.Duration `yaml:"sleep between strategies"`
		}
		RestrictReporting   bool //not part of yaml file, derived after yaml file is unmarshalled & validated
		RestrictReportingTo struct {
			DeckStartVal          int `yaml:"starting deck number"`
			DeckContinueFor       int `yaml:"continue for how many decks"`
			MovesTriedStartVal    int `yaml:"starting move number"`
			MovesTriedContinueFor int `yaml:"continue for how many moves"`
		} `yaml:"restrict reporting to"`
		WinLossReport bool `yaml:"report final deck by deck win loss record"`
	} `yaml:"play all moves"`
}

func main() {
	/*
		NOTE: Certain options of VerboseSpecial and/or pMD are incompatible:

			!!!!!! ADD CHECK TO SAY DBD AND DBDS can not be BOTH included in verbosespecial
			!!!!!!                 and that neither CAN be selected if the sixth argument pMD.pType is anything other than "X"
			!!!!!!  PROGRESSdddd can not be selected with argument 5 = TW, TS, or TSS

	*/
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

	// completing cfg

	/*****************************************

		Temporary till Yaml fixed start here

	*****************************************/

	cfg.PlayNew.DeckByDeckReportingOptions.Type = "very short"
	cfg.PlayNew.MoveByMoveReportingOptions.Type = "regular"
	cfg.PlayNew.TreeReportingOptions.Type = "regular"

	/*****************************************

		Temporary till Yaml fixed end here

	*****************************************/

	// make all strings in cfg EXCEPT cfg.General.TOutputTo lower case
	//Try to replace with a for range loop of cfg in future
	cfg.General.TypeOfPlay = strings.ToLower(cfg.General.TypeOfPlay)
	cfg.PlayNew.DeckByDeckReportingOptions.Type = strings.ToLower(cfg.PlayNew.DeckByDeckReportingOptions.Type)
	cfg.PlayNew.MoveByMoveReportingOptions.Type = strings.ToLower(cfg.PlayNew.MoveByMoveReportingOptions.Type)
	cfg.PlayNew.TreeReportingOptions.Type = strings.ToLower(cfg.PlayNew.TreeReportingOptions.Type)

	// ReportingMoveByMove zero value is false
	if cfg.PlayNew.ReportingType.DeckByDeck {
		cfg.PlayNew.ReportingDeckByDeck = true
	}
	// COMMENT
	/*if cfg.PlayNew.ReportingType.MoveByMove || cfg.PlayNew.ReportingType.Tree {
		cfg.PlayNew.ReportingMoveByMove = true
	}*/
	// RestrictReporting zero value is false
	if cfg.PlayNew.RestrictReportingTo.DeckStartVal != 0 ||
		cfg.PlayNew.RestrictReportingTo.DeckContinueFor != 0 ||
		cfg.PlayNew.RestrictReportingTo.MovesTriedStartVal != 0 ||
		cfg.PlayNew.RestrictReportingTo.MovesTriedContinueFor != 0 {
		cfg.PlayNew.RestrictReporting = true
	}
	cfg.PlayNew.ReportingType.NoReporting = !(cfg.PlayNew.ReportingType.DeckByDeck || cfg.PlayNew.ReportingType.MoveByMove || cfg.PlayNew.ReportingType.Tree)

	/*// Sixth
	pMdArgs := strings.Split(args[6], ",")
	l := len(pMdArgs)

	if l >= 1 {
		if pMdArgs[0] == "BB" || pMdArgs[0] == "BBS" || pMdArgs[0] == "BBSS" || pMdArgs[0] == "TW" || pMdArgs[0] == "TS" || pMdArgs[0] == "TSS" || pMdArgs[0] == "X" {
			pMD.pType = pMdArgs[0]
		} else {
			println("Sixth argument part 1 invalid - args[6];  arg[6] parts must be separated by commas: *1*,2,3,4,5,6")
			println("  Must start with BB, BBS, BBSS, TW, TS, TSS or X")
			os.Exit(1)
		}
	}


	// Arguments 5 & 6 above apply only to playNew			****************************************************

	/* Check for incompatible options among argument 5 or 6:
	          DBD AND DBDS can not be BOTH included in verbosespecial
			      and that neither CAN be selected if the sixth argument pMD.pType is anything other than "X"
			  PROGRESSdddd can not be selected with argument 5 = TW, TS, or TS
	*/
	/*if strings.Contains(cLArgs.verboseSpecial, ";DBD;") || strings.Contains(cLArgs.verboseSpecial, ";DBDS;") {
		if strings.Contains(cLArgs.verboseSpecial, ";DBD;") && strings.Contains(cLArgs.verboseSpecial, ";DBDS;") {
			println("Fifth argument cannot specify BOTH DBD and DBDS")
			os.Exit(1)
		} else {
			if pMD.pType != "X" {
				println("Fifth argument of \"DBD\" and \"DBDS\" incompatible with sixth argument equal to anything other than \"X\"")
				os.Exit(1)
			}
		}
		if pMD.pType == "TW" || pMD.pType == "TS" || pMD.pType == "TSS" {
			println("Sixth argument of \"TW\", \"TS\", or \"TSS\" incompatible with fifth argument (verboseSpecial of type \"PROGRESSdddd\"")
			os.Exit(1)
		}
	}
	*/
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
	}
	if cfg.General.TypeOfPlay == "playall" && cfg.PlayNew.ReportingType.DeckByDeck {
		_, err = pfmt.Printf("Deck By Deck Reporting: \n"+
			"                                           Type: %v\n"+
			"    Move Progress Reporting Cycles, in Millions: %v\n",
			cfg.PlayNew.DeckByDeckReportingOptions.Type,
			cfg.General.ProgressCounter)
	}
	if cfg.General.TypeOfPlay == "playall" && cfg.PlayNew.ReportingType.MoveByMove {
		_, err = pfmt.Printf("Move By Move Reporting: \n"+
			"                                           Type: %v\n"+
			"    Move Progress Reporting Cycles, in Millions: %v\n",
			cfg.PlayNew.MoveByMoveReportingOptions.Type,
			cfg.General.ProgressCounter)
		// add code here to turn progress reporting off if output to file unless figure out how to print some stuff to file and progress to console
		// add code for incompatible with ????
	}
	if cfg.General.TypeOfPlay == "playall" && cfg.PlayNew.ReportingType.Tree {
		_, err = pfmt.Printf("Tree Reporting: \n"+
			"                        Type: %v\n"+
			"         TreeSleepBetwnMoves: %v\n"+
			"    TreeSleepBetwnStrategies: %v\n",
			cfg.PlayNew.MoveByMoveReportingOptions.Type,
			cfg.PlayNew.TreeReportingOptions.TreeSleepBetwnMoves,
			cfg.PlayNew.TreeReportingOptions.TreeSleepBetwnStrategies)
	}
	if cfg.General.TypeOfPlay == "playall" && cfg.PlayNew.RestrictReporting {
		_, err = pfmt.Printf("\nReporting Restricted To\n"+
			"                       Staring with Deck: %v\n"+
			"                            Continue for: %v decks (0 = all the rest)\n"+
			"             Starting with Moves Tried #: %v\n"+
			"                            Continue for: %v moves tried (0 = all the rest)\n",
			cfg.PlayNew.RestrictReportingTo.DeckStartVal,
			cfg.PlayNew.RestrictReportingTo.DeckContinueFor,
			cfg.PlayNew.RestrictReportingTo.MovesTriedStartVal,
			cfg.PlayNew.RestrictReportingTo.MovesTriedContinueFor)
	}
	if cfg.General.TypeOfPlay == "playall" {
		_, err = pfmt.Printf("\nGame Length Limit, in millions: %v\n",
			cfg.PlayNew.GameLengthLimit)
	}
	if cfg.General.OutputTo == "console" {
		cfg.PlayNew.TreeReportingOptions.TreeSleepBetwnMoves *= 100_000_000
		cfg.PlayNew.TreeReportingOptions.TreeSleepBetwnStrategies *= 100_000_000
		if cfg.General.TypeOfPlay == "playorig" {
			cfg.General.ProgressCounter = 0
		} else {
			if cfg.PlayNew.ReportingType.MoveByMove || cfg.PlayNew.ReportingType.Tree {
				cfg.General.ProgressCounter = 0
			}
		}
	} else {
		cfg.PlayNew.TreeReportingOptions.TreeSleepBetwnMoves = 0 // Change this if figure out how to print to file AND console
		cfg.PlayNew.TreeReportingOptions.TreeSleepBetwnStrategies = 0
		if cfg.General.TypeOfPlay == "playorig" {
			cfg.General.ProgressCounter = 1
		}
	}
	cfg.General.ProgressCounterLastPrintTime = time.Now()

	/*	for _, v = range cfg {
		}*/

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
		playNew(*reader, cfg)
	}
	if cfg.General.TypeOfPlay == "playorig" {
		moveBasePriority = moveBasePriorityOrig
		playOrig(*reader, cfg)
	}
}
