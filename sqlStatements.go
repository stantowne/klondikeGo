package main

import (
	"database/sql"
	"fmt"
	"os"
)

func sqlExec(verb string, table string, cfg *Configuration, vPA *variablesSpecificToPlayAll, vPO *variablesSpecificToPlayOrig) {
	//var result sql.Result
	var err error
	var stmt *sql.Stmt
	var rows sql.Rows
	switch verb {
	case "Insert":
		switch table {
		case "Priority":
			/*			q = "INSERT INTO "
						q += " ([Priority_SHA256], [moveAceAcross], [moveDeuceAcross], [move3PlusAcross], [moveDown], [moveEntireColumn], [flipWasteToStock], [flipStockToWaste], [movePartialColumn], [moveAceUp], [moveDeuceUp], [move3PlusUp], [badMove], [flipSt->W Max-0], [flipSt->W Max-1], [flipSt->W Max-2], [flipSt->W Max-3], [flipSt->W Max-4], [flipSt->W Max-5], [flipSt->W Max-6], [flipSt->W Max-7])"
						q += " VALUES "
						q += " (@Priority_SHA256,  @moveAceAcross,  @moveDeuceAcross,  @move3PlusAcross,  @moveDown,  @moveEntireColumn,  @flipWasteToStock,  @flipStockToWaste,  @movePartialColumn,  @moveAceUp,  @moveDeuceUp,  @move3PlusUp,  @badMove,  @flipStToW_Max_0,  @flipStToW_Max_1,  @flipStToW_Max_2,  @flipStToW_Max_3,  @flipStToW_Max_4,  @flipStToW_Max_5,  @flipStToW_Max_6,  @flipStToW_Max_7);"
			*/
			stmt, err = db.Prepare("INSERT INTO [dbo].[Priority] ([Priority_SHA256], [moveAceAcross], [moveDeuceAcross], [move3PlusAcross], [moveDown], [moveEntireColumn], [flipWasteToStock], [flipStockToWaste], [movePartialColumn], [moveAceUp], [moveDeuceUp], [move3PlusUp], [badMove], [flipSt->W Max-0], [flipSt->W Max-1], [flipSt->W Max-2], [flipSt->W Max-3], [flipSt->W Max-4], [flipSt->W Max-5], [flipSt->W Max-6], [flipSt->W Max-7]) VALUES (@Priority_SHA256, @moveAceAcross, @moveDeuceAcross, @move3PlusAcross, @moveDown, @moveEntireColumn, @flipWasteToStock, @flipStockToWaste, @movePartialColumn, @moveAceUp, @moveDeuceUp, @move3PlusUp, @badMove, @flipStToW_Max_0, @flipStToW_Max_1, @flipStToW_Max_2, @flipStToW_Max_3, @flipStToW_Max_4, @flipStToW_Max_5, @flipStToW_Max_6, @flipStToW_Max_7); ")
			//stmt, err = db.Prepare(q)
			defer stmt.Close()
			// Execute the prepared statement
			_, err = stmt.Exec(
				sql.Named("Priority_SHA256", string(cfg.General.PrioritySHA256[:])),
				sql.Named("moveAceAcross", moveBasePriority["moveAceAcross"]),
				sql.Named("moveDeuceAcross", moveBasePriority["moveDeuceAcross"]),
				sql.Named("move3PlusAcross", moveBasePriority["move3PlusAcross"]),
				sql.Named("moveDown", moveBasePriority["moveDown"]),
				sql.Named("moveEntireColumn", moveBasePriority["moveEntireColumn"]),
				sql.Named("flipWasteToStock", moveBasePriority["flipWasteToStock"]),
				sql.Named("flipStockToWaste", moveBasePriority["flipStockToWaste"]),
				sql.Named("movePartialColumn", moveBasePriority["movePartialColumn"]),
				sql.Named("moveAceUp", moveBasePriority["moveAceUp"]),
				sql.Named("moveDeuceUp", moveBasePriority["moveDeuceUp"]),
				sql.Named("move3PlusUp", moveBasePriority["move3PlusUp"]),
				sql.Named("badMove", moveBasePriority["badMove"]),
				sql.Named("flipStToW_Max_0", moveBasePriority["flipSt->W Max-0"]),
				sql.Named("flipStToW_Max_1", moveBasePriority["flipSt->W Max-1"]),
				sql.Named("flipStToW_Max_2", moveBasePriority["flipSt->W Max-2"]),
				sql.Named("flipStToW_Max_3", moveBasePriority["flipSt->W Max-3"]),
				sql.Named("flipStToW_Max_4", moveBasePriority["flipSt->W Max-4"]),
				sql.Named("flipStToW_Max_5", moveBasePriority["flipSt->W Max-5"]),
				sql.Named("flipStToW_Max_6", moveBasePriority["flipSt->W Max-6"]),
				sql.Named("flipStToW_Max_7", moveBasePriority["flipSt->W Max-7"]),
			)
			if err != nil {
				fmt.Printf("Table: %v   Verb: %v   Error: %v getting Hostname", table, verb, err)
				os.Exit(1)
			}
		case "RunCfg":
			//fmt.Printf("PrioritySha256: %v\n", string(cfg.General.PrioritySHA256[:]))
			stmt, err = db.Prepare("INSERT INTO [dbo].[RunCfg] ([Priority_SHA256], [RunStartTime], [GitVersion], [HostName], [DeckFileName], [Decks], [FirstDeckNum], [NumberOfDecksToBePlayed], [List], [TypeOfPlay], [Verbose], [OutputTo], [outWriterFileName]) OUTPUT inserted.Run_ID VALUES (@Priority_SHA256, @RunStartTime, @GitVersion, @HostName, @DeckFileName, @Decks, @FirstDeckNum, @NumberOfDecksToBePlayed, @List, @TypeOfPlay, @Verbose, @OutputTo, @outWriterFileName);")
			defer stmt.Close()
			// Execute the prepared statement
			rows, err = stmt.Query(
				sql.Named("Priority_SHA256", string(cfg.General.PrioritySHA256[:])),
				sql.Named("RunStartTime", cfg.General.RunStartTime),
				sql.Named("GitVersion", cfg.General.GitVersion),
				sql.Named("HostName", cfg.General.HostName),
				sql.Named("DeckFileName", cfg.General.DeckFileName),
				sql.Named("Decks", cfg.General.Decks),
				sql.Named("FirstDeckNum", cfg.General.FirstDeckNum),
				sql.Named("NumberOfDecksToBePlayed", cfg.General.NumberOfDecksToBePlayed),
				sql.Named("List", cfg.General.List),
				sql.Named("TypeOfPlay", cfg.General.TypeOfPlay),
				sql.Named("Verbose", cfg.General.Verbose),
				sql.Named("OutputTo", cfg.General.OutputTo),
				sql.Named("outWriterFileName", cfg.General.outWriterFileName),
			)
			if err != nil {
				fmt.Printf("Table: %v   Verb: %v   Error: %v getting Hostname", table, verb, err)
				os.Exit(1)
			}
			for rows.Next() {
				err = rows.Scan(&cfg.General.RunID)
				if err != nil {
					fmt.Printf("Table: %v   Verb: %v   Error getting Run_ID   Error: %v getting Hostname", table, verb, err)
					os.Exit(1)
				}
			}
			return

		case "Cfg_PlayAll":
		case "Cfg_PlayOrig":
		case "PlayAllStatistics":
		case "WinningMoves":
		case "WinningMoves_Detail":
		case "boardCode":
		}
		if err != nil {
			fmt.Printf("SQL Error %v %v: %v", verb, table, err)
			panic(fmt.Sprintf("SQL Error %v %v: %v", verb, table, err))
		}
	}
	return
}
