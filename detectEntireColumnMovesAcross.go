package main

import (
	"fmt"
	"os"
)

func detectMecInner(b board) []move {
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
			}
			continue
		}
	}
	return moves
}

func detectMecNotThoughtful(b board) []move {
	return detectMecInner(b)
}

func detectMecThoughtful(b board) []move {
	//mc expects movecounter
	var movesFirstLevel []move
	movesFirstLevel = detectMecInner(b)
	if len(movesFirstLevel) == 0 { // no MEC
		return movesFirstLevel
	}
	if len(movesFirstLevel) == 1 { // only one MEC
		return movesFirstLevel
	}
	//two or more MEC
	var boards []board // slice of boards
	for i := 0; i < len(movesFirstLevel); i++ {
		boards[i] = moveMaker(b, movesFirstLevel[i]) //play each MEC

	}

	return nil
}
