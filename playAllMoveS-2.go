package main

import (
	"sort"
)

func playAllMoveS2(bIn board, moveNum int, deckNum int) (string, string) {
	/* Return Codes: SL  = Strategy Lost	 NMA = No Moves Available
	                 						 RB  = Repetitive Board
	                                         SE  = Strategy Exhausted
	                                         GML = GameLength Limit exceeded
	                 SW  = Strategy Win      EW  = Early Win
											 SW  = Standard Win  Obsolete all wins are early
	*/
	var aMoves []move //available Moves
	if moveNum > moveNumMax {
		moveNumMax = moveNum
	}
	if mvsTriedTD >= gameLengthLimit {
		stratLossesGML_TD++
		return "SL", "GML"
	}
	aMoves = detectAvailableMoves(bIn, moveNum, singleGame)
	if len(aMoves) == 0 {
		m := move{name: "No Moves Available"}
		aMoves = append(aMoves, m)
	} else {
		if len(aMoves) > 1 { //sort them by priority if necessary
			sort.SliceStable(aMoves, func(i, j int) bool {
				return aMoves[i].priority < aMoves[j].priority
			})
		}
	}
	// Try all moves
	for i := range aMoves {
		if i != 0 {
			stratNumTD++
		} else {
			// Check for repetitive board
			bNewBcode := bIn.boardCode(deckNum) //  consider modifying the boardCode and boardDeCode methods to produce strings
			bNewBcodeS := string(bNewBcode[:])  //  consider modifying the boardCode and boardDeCode methods to produce strings
			// Have we seen this board before?
			if _, ok := priorBoards[bNewBcodeS]; ok {
				// OK we did see it before so return to try next available move (if any) in aMoves[] aka strategy
				stratLossesRB_TD++
				return "SL", "RB" // Repetitive Board
			} else {
				// Remember the board state by putting it into the map "priorBoards"
				bInf := boardInfo{
					exists: true,
				}
				priorBoards[bNewBcodeS] = bInf
			}
		}
		// Check if No Moves Available
		if i == 0 && aMoves[0].name == "No Moves Available" {
			stratLossesNMA_TD++
			return "SL", "NMA"
		}
		//Detect Early Win
		if detectWinEarly(bIn) {
			stratWinsTD++
			return "SW", "EW" //  Strategy Early Win
		}
		bNew := bIn.copyBoard() // Critical Must use copyBoard
		bNew = moveMaker(bNew, aMoves[i])
		mvsTriedTD++
		recurReturnV1, recurReturnV2 := playAllMoveS(bNew, moveNum+1, deckNum)
		if recurReturnV1 == "SL" && recurReturnV2 == "GML" {
			return recurReturnV1, recurReturnV2 // return up the call stack to end strategies search  if findAllWinStrats false, and we had a win
		}
	}
	stratLossesSE_TD++
	return "SL", "SE" //  Strategy Exhausted
}
