package main

import (
	"fmt"
)

func sqlExec(verb string, table string, cfg *Configuration, vPA *variablesSpecificToPlayAll, vPO *variablesSpecificToPlayOrig) string {
	switch verb {
	case "Insert":
		switch table {
		case "Priority":
		case "RunCfg":
			result, err := db.Exec("INSERT INTO [dbo].[RunCfg] ([Priority_SHA256], [RunStartTime], [GitVersion], [HostName], [DeckFileName], [Decks], [FirstDeckNum], [NumberOfDecksToBePlayed], [List], [TypeOfPlay], [Verbose], [OutputTo], [outWriterFileName]) VALUES (RunStartTime, PrioritySHA256, cfg.General.GitVersion, cfg.General.HostName, cfg.General.DeckFileName, cfg.General.Decks, cfg.General.FirstDeckNum, cfg.General.NumberOfDecksToBePlayed, cfg.General.List, cfg.General.TypeOfPlay, cfg.General.Verbose, cfg.General.OutputTo, cfg.General.outWriter, cfg.General.outWriterFileName")
			if err != nil {
				fmt.Printf("SQL Error %v %v: %v", verb, table, err)
				panic(fmt.Sprintf("SQL Error %v %v: %v", verb, table, err))
			}
			id, err := result.LastInsertId()
			if err != nil {
				fmt.Printf("SQL Error LastInsertID %v %v: %v", verb, table, err)
				panic(fmt.Sprintf("SQL Error LastInsertID %v %v: %v", verb, table, err))
			}
			cfg.General.RunID = id
			return ""
		case "Cfg_PlayAll":
		case "Cfg_PlayOrig":
		case "PlayAllStatistics":
		case "WinningMoves":
		case "WinningMoves_Detail":
		case "boardCode":
		}
	}
	return ""
}
