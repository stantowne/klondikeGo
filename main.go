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

const gameLengthLimitOrig = 150     // max moveCounter; increasing to 200 does not increase win rate
const gameLengthLimitNew = 50000000 // max mvsTriedTD
var gameLengthLimit int
var moveNumMax int

type Configuration struct {
	General struct {
		Decks                   string `yaml:"decks"`                        // must be "consecutive" or "list"
		FirstDeckNum            int    `yaml:"first deck number"`            // must be non-negative integer
		NumberOfDecksToBePlayed int    `yaml:"number of decks to be played"` //must be non-negative integer
		List                    []int
		SuitSymbol              string `yaml:"suit symbol"`  // must be "icon" or "alpha"
		RankSymbol              string `yaml:"rank symbol"`  // must be "number" or "alpha"
		TypeOfPlay              string `yaml:"type of play"` // must be "playOrig" or "playAll"
		Verbose                 int    `yaml:"verbose"`
		OutputTo                string `yaml:"outputTo"`
	} `yaml:"general"`
	PlayOrig struct {
		Length int `yaml:"length of initial override strategy"`
	} `yaml:"play original"`
	PlayNew struct {
		DeckByDeckReport       string `yaml:"deck by deck result report"`
		ProgressCounter        int    `yaml:"move count reporting points"`
		WinLossReport          bool   `yaml:"win loss report"`
		FindAllWinStrats       bool   `yaml:"find all winning strategies?"`
		PrintMoveDetails       bool   `yaml:"print move details?"`
		PrintMoveDetailOptions struct {
			Type                  string `yaml:"type of PMD"`
			DeckStartVal          int    `yaml:"starting deck number"`
			DeckContinueFor       int    `yaml:"continue for how many decks"`
			MovesTriedStartVal    int    `yaml:"starting move number"`
			MovesTriedContinueFor int    `yaml:"continue for how many moves"`
		} `yaml:"print move detail options"`
	} `yaml:"play all moves"`
}

var err error
var singleGame bool // = true

func main() {
	/*


				  NOTE: No appreciable time penalty                        DBD  = print Deck-by-deck detail info after each deck 									playNew Only (in playAllMovesS)
					       for any option other than                       DBDS = print Deck-by-deck SHORT detail info after each deck								playNew Only (in playAllMovesS)WL  = print deck summary Win/Loss stats after all decks to see which decks won and which lost    playNew Only (in playNew)
				           PROGRESSdddd                         		   SUITSYMBOL = print S, D, C, H instead of runes - defaults to runes
				                                                  		   RANKSYMBOL = print Ac, Ki, Qu, Jk instead of 01, 11, 12, 13 - defaults to numeric
				                                                  		   WL = Win/Loss record for each deck printed at end
				  NOTE Time penalty at: GML = 800,000,000      		   PROGRESSdddd = Print the deckNum, mvsTriedTD, moveNum, stratNumTD, unqBoards every dddd movesTriedTD tried overwriting the previous printing
				       PROGRESS500000 = X.X%                       		                        dddd = 0 will be treated as if the /PROGRESS0000/ was NOT in the verbose special string !!!!!!!
				       PROGRESS500000 = X.X%                                           		    if dddd is left out then a default of every 10,000 movesTriedTD will be used
		               PROGRESS50000  = X.X%
				                                                  								PROGRESSdddd is preprocessed below (soon to move to playNew)

				                                                  		   BELL = Ring bell after any deck taking more than 000 minutes (Not yet Implemented)

														args[5] = findAllWinStrats 	      - (applicable playNew only)
								       printMoveDetail  args[6] = pMD		              - struct type: printMoveDetail - (applicable playNew only)
													              	as a string to be parsed of the form:
													           			pType,startType,deckStartVal,deckContinueFor,outputTo

																		where:
																			pType = empty or X - do not print NOTE: Default if args[6] is not on command line
																				  = BB         - Board by Board detail
												                                  = BBS        - Board by Board Short detail
												                                  = BBSS       - Board by Board Super Short detail
													                              = TW         - print Tree in Wide mode     8 char per move
																				  = TS         - print Tree in Skinny mode   5 char per move
											                                      = TSS        - print Tree in Super Skinny mode   3 char per move
										                               These next four limit at what point and for how long move detail should actually be printed.
																			deckStartVal    	  = Non-negative integer (Default 0)
																			deckContinueFor  	  = Non-negative integer (Default 0 which indicates forever)
																			movesTriedTDStartVal    = Non-negative integer (Default 0)
																			movesTriedTDContinueFor = Non-negative integer (Default 0 which indicates forever)

										                                    outputTo = C = Console (default)
																                     = file name and path (if applicable)
													                                   Note: if file name is present then startType. deckStartVal and ContinueFor
													                                         must be present or delineated with ":"

														           	and placed into a package level struct pMD of type pMd which can be seen above:

								NOTE: Certain options of VerboseSpecial and/or pMD are incompatible:

									!!!!!! ADD CHECK TO SAY DBD AND DBDS can not be BOTH included in verbosespecial
									!!!!!!                 and that neither CAN be selected if the sixth argument pMD.pType is anything other than "X"
									!!!!!!  PROGRESSdddd can not be selected with argument 5 = TW, TS, or TSS

	*/

	// load YAML file
	cfg := Configuration{}

	data, err := os.ReadFile("./config.yml")
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}

	//cfg.verboseSpecialProgressCounterLastPrintTime = time.Now()

	// Setup pfmt to print thousands with commas
	var pfmt = message.NewPrinter(language.English)

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
		println("Length must be an integer >= -1 and <= 24")
		os.Exit(1)
	}

	if cfg.General.Verbose >= 10 || cfg.General.Verbose < 0 {
		println("verbose invalid")
		println("verbose must be a non-negative integer no greater than 10")
		os.Exit(1)
	}

	//// Arguments 5 & 6 below applies only to playNew			****************************************************
	//// But they must be on command line anyway
	//switch strings.TrimSpace(args[5])[0:1] {
	//case "A", "a":
	//	cLArgs.findAllWinStrats = true
	//case "F", "f":
	//	cLArgs.findAllWinStrats = false
	//default:
	//	println("Fifth argument invalid - args[5]")
	//	println("  findAllWinStrats must be either:")
	//	println("     'F' or 'f' - Normal case stop after finding a successful set of moves")
	//	println("     'A' or 'a' - See how many paths to success you can find")
	//	os.Exit(1)
	//}

	/*// Sixth
	pMdArgs := strings.Split(args[6], ",")
	l := len(pMdArgs)

	if l >= 1 {
		if pMdArgs[0] == "BB" || pMdArgs[0] == "BBS" || pMdArgs[0] == "BBSS" || pMdArgs[0] == "TW" || pMdArgs[0] == "TS" || pMdArgs[0] == "TSS" || pMdArgs[0] == "X" {
			pMD.pType = pMdArgs[0]
		} else {
			println("Sixth argument part 1 invalid - args[6];  arg[6] parts are separated by commas: *1*,2,3,4,5,6")
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

	if cfg.General.TypeOfPlay == "playOrig" {
		nOfS := 1 << cfg.PlayOrig.Length //number of initial strategies
		_, err = pfmt.Printf(" Style: Original iOS (Initial Override Strategies)\n\n"+
			"                     iOS strategy length: %v\n"+
			"          Max possible attempts per deck: %v\n"+
			"       Total possible attempts all decks: %v\n\n",
			cfg.PlayOrig.Length,
			nOfS,
			nOfS*cfg.General.NumberOfDecksToBePlayed)
	}
	if cfg.General.TypeOfPlay == "playAll" {
		_, err = pfmt.Printf("Deck by Deck Reporting Form: %v\n"+
			"Progress Move Counter: %v\n"+
			"Win Loss Report: %v\n"+
			"Find All Successful Strategies : %v\n"+
			"", cfg.PlayNew.DeckByDeckReport, cfg.PlayNew.ProgressCounter, cfg.PlayNew.WinLossReport, cfg.PlayNew.FindAllWinStrats)
	}
	if cfg.General.TypeOfPlay == "playAll" && !cfg.PlayNew.PrintMoveDetails {
		_, err = pfmt.Printf("No Print Move Details\n")
	}
	if cfg.General.TypeOfPlay == "PlayAll" && cfg.PlayNew.PrintMoveDetails {
		_, err = pfmt.Printf("Print Move Details Options\n"+
			"                              Type: %v\n"+
			"                       Staring with Deck: %v\n"+
			"                            Continue for: %v decks (0 = all the rest)\n"+
			"             Starting with Moves Tried #: %v\n"+
			"                            Continue for: %v moves tried (0 = all the rest)\n"+
			cfg.PlayNew.PrintMoveDetailOptions.Type,
			cfg.PlayNew.PrintMoveDetailOptions.DeckStartVal,
			cfg.PlayNew.PrintMoveDetailOptions.DeckContinueFor,
			cfg.PlayNew.PrintMoveDetailOptions.MovesTriedStartVal,
			cfg.PlayNew.PrintMoveDetailOptions.MovesTriedContinueFor)
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
		gameLengthLimit = gameLengthLimitNew
		moveBasePriority = moveBasePriorityNew
		_, err = pfmt.Printf("\n                         GameLengthLimit: %v Move Counter\n\n\n", gameLengthLimit)
		playNew(*reader, cfg)
	}
	if cfg.General.TypeOfPlay == "playOrig" {
		gameLengthLimit = gameLengthLimitOrig
		_, err = pfmt.Printf("\n                         GameLengthLimit: %v Moves Tried\n\n\n", gameLengthLimit)
		moveBasePriority = moveBasePriorityOrig
		playOrig(*reader, cfg)
	}
}
