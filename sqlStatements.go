package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
)

var WinningMoves_ID_inserted int64

func sqlExec(verb string, table string, cfg *Configuration, vPA *variablesSpecificToPlayAll, vPO *variablesSpecificToPlayOrig) string {
	returnResult := ""
	verb = strings.ToLower(verb)
	table = strings.ToLower(table)
	var err error
	var err2 error
	var stmt *sql.Stmt
	var rows *sql.Rows
	//var q string
	switch verb {
	case "insert":
		switch table {
		case "priority":
			stmt, err = db.Prepare("INSERT INTO [dbo].[Priority] ([Priority_SHA256], [moveAceAcross], [moveDeuceAcross], [move3PlusAcross], [moveDown], [moveEntireColumn], [flipWasteToStock], [flipStockToWaste], [movePartialColumn], [moveAceUp], [moveDeuceUp], [move3PlusUp], [badMove], [flipSt->W Max-0], [flipSt->W Max-1], [flipSt->W Max-2], [flipSt->W Max-3], [flipSt->W Max-4], [flipSt->W Max-5], [flipSt->W Max-6], [flipSt->W Max-7]) VALUES (@Priority_SHA256, @moveAceAcross, @moveDeuceAcross, @move3PlusAcross, @moveDown, @moveEntireColumn, @flipWasteToStock, @flipStockToWaste, @movePartialColumn, @moveAceUp, @moveDeuceUp, @move3PlusUp, @badMove, @flipStToW_Max_0, @flipStToW_Max_1, @flipStToW_Max_2, @flipStToW_Max_3, @flipStToW_Max_4, @flipStToW_Max_5, @flipStToW_Max_6, @flipStToW_Max_7); ")
			defer stmt.Close()
			if err != nil {
				if cfg.General.OutputTo != "console" {
					oW = os.Stdout
				}
				fmt.Printf("Table: %v   Verb: %v   Error preparing   Error: %v ", table, verb, err)
				os.Exit(1)
			}
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
		case "runcfg":
			stmt, err = db.Prepare("INSERT INTO [dbo].[RunCfg] ([Priority_SHA256], [RunStartTime], [GitVersion], [HostName], [DeckFileName], [Decks], [FirstDeckNum], [NumberOfDecksToBePlayed], [List], [TypeOfPlay], [Verbose], [OutputTo], [outWriterFileName]) OUTPUT inserted.Run_ID VALUES (@Priority_SHA256, @RunStartTime, @GitVersion, @HostName, @DeckFileName, @Decks, @FirstDeckNum, @NumberOfDecksToBePlayed, @List, @TypeOfPlay, @Verbose, @OutputTo, @outWriterFileName);")
			defer stmt.Close()
			if err != nil {
				if cfg.General.OutputTo != "console" {
					oW = os.Stdout
				}
				fmt.Printf("Table: %v   Verb: %v   Error preparing Error: %v ", table, verb, err)
				os.Exit(1)
			}
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
			for rows.Next() {
				err2 = rows.Scan(&cfg.General.RunID)
				if err2 != nil {
					if cfg.General.OutputTo != "console" {
						oW = os.Stdout
					}
					fmt.Printf("Table: %v   Verb: %v   Error getting Run_ID   Error: %v ", table, verb, err)
					os.Exit(1)
				}
			}
		case "cfg_playall":
			stmt, err = db.Prepare("INSERT INTO [dbo].[Cfg_PlayAll] (Run_ID, GameLengthLimit, DeckByDeck, MoveByMove, Tree, NoReporting, DbD_Type, MbM_Type, Tree_Type, TreeSleepBetwnMoves, TreeSleepBetwnMovesDur, TreeSleepBetwnStrategies, TreeSleepBetwnStrategiesDur, RestrictReporting, RestrictRept_DeckStartVal, RestrictRept_DeckContinueFor, RestrictRept_MovesTriedStartVal, RestrictRept_MovesTriedContinueFor, ProgressCounter, SaveResultsToSQL, SQLConnectionString) VALUES (@Run_ID, @GameLengthLimit, @DeckByDeck, @MoveByMove, @Tree, @NoReporting, @DbD_Type, @MbM_Type, @Tree_Type, @TreeSleepBetwnMoves, @TreeSleepBetwnMovesDur, @TreeSleepBetwnStrategies, @TreeSleepBetwnStrategiesDur, @RestrictReporting, @RestrictRept_DeckStartVal, @RestrictRept_DeckContinueFor, @RestrictRept_MovesTriedStartVal, @RestrictRept_MovesTriedContinueFor, @ProgressCounter, @SaveResultsToSQL, @SQLConnectionString); ")
			defer stmt.Close()
			if err != nil {
				if cfg.General.OutputTo != "console" {
					oW = os.Stdout
				}
				fmt.Printf("Table: %v   Verb: %v   Error preparing   Error: %v ", table, verb, err)
				os.Exit(1)
			}
			// Execute the prepared statement
			_, err = stmt.Exec(
				sql.Named("Run_ID", cfg.General.RunID),
				sql.Named("GameLengthLimit", cfg.PlayAll.GameLengthLimit),
				sql.Named("DeckByDeck", cfg.PlayAll.ReportingType.DeckByDeck),
				sql.Named("MoveByMove", cfg.PlayAll.ReportingType.MoveByMove),
				sql.Named("Tree", cfg.PlayAll.ReportingType.Tree),
				sql.Named("NoReporting", cfg.PlayAll.ReportingType.NoReporting),
				sql.Named("DbD_Type", cfg.PlayAll.DeckByDeckReportingOptions.Type),
				sql.Named("MbM_Type", cfg.PlayAll.MoveByMoveReportingOptions.Type),
				sql.Named("Tree_Type", cfg.PlayAll.TreeReportingOptions.Type),
				sql.Named("TreeSleepBetwnMoves", cfg.PlayAll.TreeReportingOptions.TreeSleepBetwnMoves),
				sql.Named("TreeSleepBetwnMovesDur", cfg.PlayAll.TreeReportingOptions.TreeSleepBetwnMovesDur),
				sql.Named("TreeSleepBetwnStrategies", cfg.PlayAll.TreeReportingOptions.TreeSleepBetwnStrategies),
				sql.Named("TreeSleepBetwnStrategiesDur", cfg.PlayAll.TreeReportingOptions.TreeSleepBetwnStrategiesDur),
				sql.Named("RestrictReporting", cfg.PlayAll.RestrictReporting),
				sql.Named("RestrictRept_DeckStartVal", cfg.PlayAll.RestrictReportingTo.DeckStartVal),
				sql.Named("RestrictRept_DeckContinueFor", cfg.PlayAll.RestrictReportingTo.DeckContinueFor),
				sql.Named("RestrictRept_MovesTriedStartVal", cfg.PlayAll.RestrictReportingTo.MovesTriedStartVal),
				sql.Named("RestrictRept_MovesTriedContinueFor", cfg.PlayAll.RestrictReportingTo.MovesTriedContinueFor),
				sql.Named("ProgressCounter", cfg.PlayAll.ProgressCounter),
				sql.Named("SaveResultsToSQL", cfg.PlayAll.SaveResultsToSQL),
				sql.Named("SQLConnectionString", cfg.PlayAll.SQLConnectionString),
			)
		case "cfg_playorig":
			if vPO != nil { // Added to eliminate error about unused variable vPO
				stmt, err = db.Prepare("INSERT INTO [dbo].[Cfg_PlayOrig] (Length, GameLengthLimit) VALUES (@Length, @GameLengthLimit); ")
				defer stmt.Close()
				if err != nil {
					if cfg.General.OutputTo != "console" {
						oW = os.Stdout
					}
					fmt.Printf("Table: %v   Verb: %v   Error preparing   Error: %v ", table, verb, err)
					os.Exit(1)
				}
				// Execute the prepared statement
				_, err = stmt.Exec(
					sql.Named("Run_ID", cfg.General.RunID),
					sql.Named("Length", cfg.PlayOrig.Length),
					sql.Named("GameLengthLimit", cfg.PlayOrig.GameLengthLimit),
				)
			}
		case "playall_statistics":
			stmt, err = db.Prepare("INSERT INTO [dbo].[PlayAll_Statistics] (Deck_ID, Run_ID, winning_MovesSHA256, mvsTried, stratNum, stratTried, stratWins, stratLosses, stratLossesGLE, stratLossesGLEAb, stratLossesNMA, stratLossesRB, stratLossesMajSE, stratLossesMinSE, stratLossesEL, winningMovesCnt, unqBoards, elapsedTime, moveNumMax, moveNumAtWin) VALUES (@Deck_ID, @Run_ID, @winningMoves_SHA256, @mvsTried, @stratNum, @stratTried, @stratWins, @stratLosses, @stratLossesGLE, @stratLossesGLEAb, @stratLossesNMA, @stratLossesRB, @stratLossesMajSE, @stratLossesMinSE, @stratLossesEL, @winningMovesCnt, @unqBoards, @elapsedTime, @moveNumMax, @moveNumAtWin); ")
			defer stmt.Close()
			if err != nil {
				if cfg.General.OutputTo != "console" {
					oW = os.Stdout
				}
				fmt.Printf("Table: %v   Verb: %v   Error preparing   Error: %v ", table, verb, err)
				os.Exit(1)
			}
			// Execute the prepared statement
			_, err = stmt.Exec(
				sql.Named("Deck_ID", vPA.TDotherSQL.deckNum),
				sql.Named("Run_ID", cfg.General.RunID),
				sql.Named("WinningMoves_SHA256", NewNullString(vPA.TDotherSQL.winningMovesSHA256)),
				sql.Named("mvsTried", vPA.TD.mvsTried),
				sql.Named("stratNum", vPA.TD.stratNum),
				sql.Named("stratTried", vPA.TD.stratTried),
				sql.Named("stratWins", vPA.TD.stratWins),
				sql.Named("stratLosses", vPA.TD.stratLosses),
				sql.Named("stratLossesGLE", vPA.TD.stratLossesGLE),
				sql.Named("stratLossesGLEAb", vPA.TD.stratLossesGLEAb),
				sql.Named("stratLossesNMA", vPA.TD.stratLossesNMA),
				sql.Named("stratLossesRB", vPA.TD.stratLossesRB),
				sql.Named("stratLossesMajSE", vPA.TD.stratLossesMajSE),
				sql.Named("stratLossesMinSE", vPA.TD.stratLossesMinSE),
				sql.Named("stratLossesEL", vPA.TD.stratLossesEL),
				sql.Named("winningMovesCnt", vPA.TD.winningMovesCnt),
				sql.Named("unqBoards", vPA.TD.unqBoards),
				sql.Named("elapsedTime", vPA.TD.elapsedTime),
				sql.Named("moveNumMax", vPA.TDotherSQL.moveNumMax),
				sql.Named("moveNumAtWin", vPA.TDotherSQL.moveNumAtWin),
			)
		case "winningmoves":
			returnResult = sqlExec("Query", "WinningMoves", cfg, vPA, nil)
			if returnResult == "New Set of Winning Moves" {
				stmt, err = db.Prepare("INSERT INTO [dbo].[WinningMoves] ([WinningMoves_SHA256]) OUTPUT inserted.WinningMoves_ID VALUES (@WinningMoves_SHA256); ")
				defer stmt.Close()
				if err != nil {
					if cfg.General.OutputTo != "console" {
						oW = os.Stdout
					}
					fmt.Printf("Table: %v   Verb: %v   Error preparing   Error: %v ", table, verb, err)
					os.Exit(1)
				}
				// Execute the prepared statement
				rows, err = stmt.Query(
					sql.Named("WinningMoves_SHA256", NewNullString(vPA.TDotherSQL.winningMovesSHA256)),
				)
				if rows != nil {
					for rows.Next() {
						err2 = rows.Scan(&WinningMoves_ID_inserted)
						if err2 != nil {
							if cfg.General.OutputTo != "console" {
								oW = os.Stdout
							}
							fmt.Printf("Table: %v   Verb: %v   Error getting Run_ID   Error: %v ", table, verb, err)
							os.Exit(1)
						}
					}
					sqlExec("Insert", "WinningMoves_Detail", cfg, vPA, nil)
				}
			}
		case "winningmoves_detail":
			for k := range vPA.TDotherSQL.winningMoves {
				stmt, err = db.Prepare("INSERT INTO [dbo].[WinningMoves_Detail] (WinningMoves_ID, MoveNum, name, priority, toPile, toCol, fromCol, MovePortionStartIdx, cardToMoveRank, cardToMoveSuit, cardToMoveFaceUp, colCardFlip) VALUES (@WinningMoves_ID, @MoveNum, @name, @priority, @toPile, @toCol, @fromCol, @MovePortionStartIdx, @cardToMoveRank, @cardToMoveSuit, @cardToMoveFaceUp, @colCardFlip); ")
				defer stmt.Close()
				// Execute the prepared statement
				_, err = stmt.Exec(
					sql.Named("WinningMoves_ID", WinningMoves_ID_inserted),
					sql.Named("MoveNum", k),
					sql.Named("name", vPA.TDotherSQL.winningMoves[k].name),
					sql.Named("priority", vPA.TDotherSQL.winningMoves[k].priority),
					sql.Named("toPile", vPA.TDotherSQL.winningMoves[k].toPile),
					sql.Named("toCol", vPA.TDotherSQL.winningMoves[k].toCol),
					sql.Named("fromCol", vPA.TDotherSQL.winningMoves[k].fromCol),
					sql.Named("MovePortionStartIdx", vPA.TDotherSQL.winningMoves[k].MovePortionStartIdx),
					sql.Named("cardToMoveRank", vPA.TDotherSQL.winningMoves[k].cardToMove.Rank),
					sql.Named("cardToMoveSuit", vPA.TDotherSQL.winningMoves[k].cardToMove.Suit),
					sql.Named("cardToMoveFaceUp", vPA.TDotherSQL.winningMoves[k].cardToMove.FaceUp),
					sql.Named("colCardFlip", vPA.TDotherSQL.winningMoves[k].colCardFlip),
				)
			}
		case "boardcode":
			stmt, err = db.Prepare("INSERT INTO [dbo].[boardCode] (boardCode_str, Deck_ID) VALUES (@boardCode_str, @Deck_ID); ")
			if err != nil {
				errHandler(cfg.General.OutputTo, table, verb, "preparing", err) /*	if cfg.General.OutputTo != "console" {
						oW = os.Stdout
					}
					fmt.Printf("Table: %v   Verb: %v   Error preparing   Error: %v ", table, verb, err)
					os.Exit(1)*/
			}
			defer stmt.Close()
			// Execute the prepared statement
			_, err = stmt.Exec(
				sql.Named("boardCode_str", vPA.TDotherSQL.boardCodeOfDeckAsString),
				sql.Named("Deck_ID", vPA.TDotherSQL.deckNum),
			)
		}
	case "query":
		switch table {
		case "winningmoves":
			var x string
			var row *sql.Row
			row = db.QueryRow("SELECT 'x' x FROM [dbo].[WinningMoves] WHERE WinningMoves_SHA256 = @w; ", sql.Named("w", vPA.TDotherSQL.winningMovesSHA256))
			if row.Scan(&x) == sql.ErrNoRows {
				return "New Set of Winning Moves"
			} else {
				return "Old Set of Winning Moves"
			}
		case "playall_statistics_gle":
			var max_mvsTried int
			var row *sql.Row
			row = db.QueryRow("SELECT MAX([mvsTried]) max_mvsTried FROM [dbo].[PlayAll_Statistics] WHERE Deck_ID = @Deck_ID AND [stratLossesGLE] > 0 AND NOT EXISTS (SELECT 'x' x FROM [dbo].[PlayAll_Statistics] WHERE Deck_ID = @Deck_ID AND ([stratWins] > 0 OR [stratLosses] > 0 )); ", sql.Named("Deck_ID", vPA.TDotherSQL.deckNum))
			if row.Scan(&max_mvsTried) == sql.ErrNoRows || max_mvsTried == 0 {
				return "Skip"
			} else {
				if max_mvsTried > cfg.PlayAll.GameLengthLimit*1000000 {
					_, _ = pfmt.Printf("\nDeck: %v is not yet solved.  However it will be skipped as GLL for this run is: %v < GLL of a previous run: %v\n", vPA.TDotherSQL.deckNum, cfg.PlayAll.GameLengthLimit*1000000, max_mvsTried)
					if cfg.General.OutputTo != "console" {
						_, _ = pfmt.Fprintf(oW, "Deck: %v is not yet solved.  However it will be skipped as GLL for this run is: %v < GLL of a previous run: %v\n", vPA.TDotherSQL.deckNum, cfg.PlayAll.GameLengthLimit, max_mvsTried)
					}
					return "Skip"
				} else {
					return "NoSkip"
				}
			}
		}
	}
	if err != nil {
		errHandler(cfg.General.OutputTo, table, verb, "SQL Error", err)
		/*		if cfg.General.OutputTo != "console" {
					oW = os.Stdout
				}
				fmt.Printf("SQL Error %v %v: %v", verb, table, err)
				panic(fmt.Sprintf("SQL Error %v %v: %v", verb, table, err))*/
	}
	return returnResult
}

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func errHandler(OutputTo string, table string, verb string, doing string, err error) {
	if OutputTo != "console" {
		oW = os.Stdout
	}
	fmt.Printf("Table: %v   Verb: %v   Error %v   Error: %v ", table, verb, doing, err)
	os.Exit(1)
}
