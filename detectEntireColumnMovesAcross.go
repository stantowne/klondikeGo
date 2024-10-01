package main

import (
	"fmt"
	"os"
)

func detectMecInner(b board, _ int, _ bool) []move {
	var moves []move
	for frmColNum := 0; frmColNum < 7; frmColNum++ {
		firstFaceUpIndex, FaceUpPortion, err := faceUpPortion(b.columns[frmColNum])
		if err != nil {
			fmt.Printf("detectEntireColumnMovesAcross: error calling faceUpPortion on b.columns[%v] %v\n", frmColNum, err)
			os.Exit(1)
		}
		if len(FaceUpPortion) == 0 {
			continue
		}
		for step := 1; step < 7; step++ {
			toColNum := (frmColNum + step) % 7
			lastCard, _, err := last(b.columns[toColNum])
			if err != nil {
				fmt.Printf("detectEntireColumnMovesAcross : error calling last on b.columns[%v] %v\n", toColNum, err)
				os.Exit(1)
			}
			if ((FaceUpPortion[0].Rank == 13) && //if the FaceUpPortion of the fromCol begins with a King AND
				b.columns[frmColNum][0].FaceUp == false && //fromCol begins with a face down card AND
				(len(b.columns[toColNum]) == 0)) || //the toCol is empty OR
				(FaceUpPortion[0].Rank == lastCard.Rank-1 && //three needs four AND
					FaceUpPortion[0].color() != lastCard.color()) { //red needs black
				m := move{
					name:                "moveEntireColumn",
					priority:            moveBasePriority["badMove"], //This is an initial assignment.  See below
					toCol:               toColNum,
					fromCol:             frmColNum,
					MovePortionStartIdx: firstFaceUpIndex,
					MovePortion:         FaceUpPortion,
				}
				if firstFaceUpIndex > 0 { //if there is a face down portion of the from column
					m.colCardFlip = true
				}

				// recall that lower priority is better
				// if any of the three conditions which make a mec worth making exist
				// then lower the priority to the moveBasePriority minus the length of the face down portion of the from column
				// so a longer face down portion results lower priority
				// I have confirmed that minus wins more decks than plus

				if firstFaceUpIndex > 0 || //if there is a face down portion of the from column OR
					(firstFaceUpIndex == 0 && kingReady(b, frmColNum)) || //an empty column will result and there is a king ready to move there OR
					sisterCardInUpPortion(b, lastCard, toColNum, frmColNum) { //sister card in up portion of the other columns combined
					m.priority = moveBasePriority["moveEntireColumn"] - firstFaceUpIndex //Reassignment of moveBasePriority

					moves = append(moves, m)
				}
				/*if sisterCardInUpPortion(b, lastCard, toColNum, frmColNum) {
					m.priority = moveBasePriority["moveEntireColumn"] - firstFaceUpIndex
					moves = append(moves, m)
				}
				if firstFaceUpIndex > 0 {
					m.priority = moveBasePriority["moveEntireColumn"] - firstFaceUpIndex
					moves = append(moves, m)
				}
				if firstFaceUpIndex == 0 && kingReady(b, frmColNum) {
					m.priority = moveBasePriority["moveEntireColumn"] - firstFaceUpIndex
					moves = append(moves, m)
				}*/
			}
			continue
		}
	}
	return moves
}

func detectMecNotThoughtful(b board, mc int, singleGame bool) []move {
	var movesFirstLevel []move
	movesFirstLevel = detectMecInner(b, mc, singleGame)
	if singleGame && len(movesFirstLevel) > 1 {
		fmt.Printf("mc is %d -- MEC detected are %v\n", mc, movesFirstLevel)
		printBoard(b)
	}
	return movesFirstLevel
}

func detectMecThoughtful(b board, mc int, singleGame bool) []move {
	var movesFirstLevel []move
	movesFirstLevel = detectMecInner(b, mc, singleGame)
	if singleGame && len(movesFirstLevel) > 1 {
		fmt.Printf("mc is %d -- first level MEC(s) detected is(are) %v\n", mc, movesFirstLevel)
		printBoard(b)
	}
	var replacementMovesFirstLevel []move
	if len(movesFirstLevel) > 1 { // if there are two or more first level MEC
		for i := 0; i < len(movesFirstLevel); i++ { // for each first level MEC
			var movesNextLevel []move
			nextB := moveMaker(b, movesFirstLevel[i])              // play it
			movesNextLevel = detectMecInner(nextB, mc, singleGame) // detect whether the resulting board offers a MEC
			if len(movesNextLevel) > 0 {                           // if movesNextLevel[i] is not empty, it contains an MEC which means the ith move of movesFirstLevel is preferable
				replacementMovesFirstLevel = append(replacementMovesFirstLevel, movesFirstLevel[i]) //put the ith move from the First Level into it
				if singleGame {
					fmt.Printf("checking element %d of first level MEC after which the board would be\n", i)
					printBoard(nextB)
					fmt.Printf("next level MEC detected are %v\n", movesNextLevel)
				}
			}
		}
		if len(replacementMovesFirstLevel) != 0 {
			if singleGame {
				fmt.Printf("return replacementMovesFirstLevel: %v\n", replacementMovesFirstLevel)
			}
			return replacementMovesFirstLevel
		}
	}
	if singleGame && len(movesFirstLevel) > 0 {
		fmt.Printf("return movesFirstLevel: %v\n", movesFirstLevel)
	}
	return movesFirstLevel
}
