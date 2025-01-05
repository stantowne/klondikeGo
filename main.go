package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/csv"
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gopkg.in/yaml.v3"
	"io"
	/*
		_ "net/http/pprof" //  DELETE ME MEMORY PROFILING
	*/
	"os"
	"runtime/debug"
	"sort"
	"strconv"
	"time"
)

// Setup pfmt to print thousands with commas
var pfmt = message.NewPrinter(language.English)

// Create the Short named package variable "oW" for cfg.General.outWriter
var oW *os.File

var db *sql.DB
var tx *sql.Tx

var SQL_Time_elapsed time.Duration
var SQL_Start_Time time.Time

func main() {
	//                                                            DELETE ME MEMORY PROFILING
	/*	go func() {
			http.ListenAndServe("localhost:8080", nil)
		}()
	*/
	//                                                            DELETE ME MEMORY PROFILING

	//*****************************************************************
	// unmarshal YAML file
	cfg := Configuration{}

	cfg.General.RunStartTime = time.Now()
	cfg.General.GitVersion = ""
	data, err := os.ReadFile("./config.yml")
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}
	//I need to confirm that, by virtue of no error being returned, we know that all bools have valid values.

	//	configPrint2(cfg, 0, "cfg")

	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				cfg.General.GitVersion = setting.Value
				break
			}
		}
	}
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Printf("Error: %v getting Hostname", err)
		os.Exit(1)
	}
	cfg.General.HostName = hostname

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
	// Fill the Short named package variable "oW" for cfg.General.outWriter
	oW = os.Stdout
	if cfg.General.OutputTo != "console" {
		cfg.General.outWriterFileName = cfg.General.OutputTo
		cfg.PlayAll.TreeReportingOptions.TreeSleepBetwnMoves = 0
		cfg.PlayAll.TreeReportingOptions.TreeSleepBetwnMovesDur = 0
		cfg.PlayAll.TreeReportingOptions.TreeSleepBetwnStrategies = 0
		cfg.PlayAll.TreeReportingOptions.TreeSleepBetwnStrategiesDur = 0
		if cfg.General.Decks == "consecutive" {
			cfg.General.outWriterFileName += "_" + strconv.Itoa(cfg.General.FirstDeckNum) + "-" + strconv.Itoa(cfg.General.FirstDeckNum+cfg.General.NumberOfDecksToBePlayed-1)
		} else {
			cfg.General.outWriterFileName += "_List"
		}
		if cfg.General.TypeOfPlay == "playAll" {
			cfg.General.outWriterFileName += "_GLE" + strconv.Itoa(cfg.PlayAll.GameLengthLimit)
		} else {
			cfg.General.outWriterFileName += "_GLE" + strconv.Itoa(cfg.PlayOrig.GameLengthLimit)
		}
		cfg.General.outWriterFileName += "__" + cfg.General.RunStartTime.Format("2006.01.02_15.04.05_-0700") + ".txt"
	}
	configPrint(cfg) // Print FIRST time to stout
	if cfg.General.OutputTo != "console" {
		// create file
		cfg.General.outWriter, err = os.Create(cfg.General.outWriterFileName)
		if err != nil {
			fmt.Printf("Error: %v  Error creating output file: %s", err, cfg.General.outWriterFileName)
			os.Exit(1)
		}
		// Fill the Short named package variable "oW" for cfg.General.outWriter
		oW = cfg.General.outWriter
		// remember to close the file
		defer func(oW *os.File) {
			err := oW.Close()
			if err != nil {
				fmt.Printf("Error: %v  Error closing output file: %s msg: %v", err, cfg.General.outWriterFileName, err.Error())
				os.Exit(1)
			}
		}(oW)

		configPrint(cfg) // Print SECOND time to file
	}

	cfg.PlayAll.ProgressCounter *= 1_000_000

	//inputFileName := cfg.General.DeckFileName
	file, err := os.Open(cfg.General.DeckFileName)
	if err != nil {
		fmt.Printf("Error: %v  Cannot open Deck inputFileName: %s", err, cfg.General.DeckFileName)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			println("could not close file: %v  err: %v  errmsg: %v", err)
		}
	}(file)
	reader := csv.NewReader(file)

	// skip forward to the first deck to be played
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
		cfg.General.PrioritySHA256 = sha256ofMap(moveBasePriority)
		if cfg.PlayAll.SaveResultsToSQL {

			SQL_Start_Time = time.Now()

			sqlExec("Transaction", "Begin", &cfg, nil, nil)
			sqlExec("Insert", "Priority", &cfg, nil, nil)
			// Create row in sql file RunCfg and get back the Run_ID that was created on the insert and assign it to cfg.General.runID !!
			sqlExec("Insert", "RunCfg", &cfg, nil, nil)
			// Create row in sql file Cfg_PlayAll
			sqlExec("Insert", "Cfg_PlayAll", &cfg, nil, nil)
			sqlExec("Transaction", "Commit", &cfg, nil, nil)

			SQL_Time_elapsed += time.Since(SQL_Start_Time)

		}
		playAll(*reader, &cfg)
	}
	if cfg.General.TypeOfPlay == "playOrig" {
		moveBasePriority = moveBasePriorityOrig
		/*
		   THIS COMMENTED OUT CODE IS INCLUDED HERE IN CASE IT IS DECIDED TO MOVE ADD SAVE RESULTS TO SQL OPTION TO INCLUDE PLAYoRIG
		   in that case the variables cfg.PlayAll.SaveResultsToSQL and cfg.PlayAll.SQLConnectionString should be moved into cfg.General

		   		cfg.General.PrioritySHA256 = sha256ofMap(moveBasePriority)
		   		if cfg.PlayAllxxx.SaveResultsToSQL {
		   			// Calculate and set cfg.General.prioritySHA256 of moveBasePriority
		   			// Check if a row Priority_SHA256 == cfg.General.prioritySHA256 exists if not Create it !!!!
		   			sqlExec("Insert", "Priority", &cfg, nil, nil)
		   			// Create row in sql file RunCfg and get back the Run_ID that was created on the insert and assign it to cfg.General.runID !!
		   			sqlExec("Insert", "RunCfg", &cfg, nil, nil)
		   			// Create row in sql file Cfg_PlayAll
		   			sqlExec("Insert", "Cfg_Orig", &cfg, nil, nil)
		   		}

		*/
		playOrig(*reader, &cfg)
	}
}

func sha256ofMap(m map[string]int) [32]byte {
	var bytesOfSortedMap []byte
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		bytesOfSortedMap = append(bytesOfSortedMap, []byte(k)...) //fmt.Println(k, m[k])
		bytesOfSortedMap = append(bytesOfSortedMap, i32ToBytes(m[k])...)
	}
	return sha256.Sum256(bytesOfSortedMap)
}
