package main

import (
	"fmt"
	"os"
)

func detectEntireColumnMoves(b board, mc int, singleGame bool) []move {
	//mc expects movecounter
	specialMove := 200
	var moves []move
	if mc < 0 {
		return moves
	}
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
			if singleGame && mc == specialMove {
				fmt.Printf("within detectEntireColumnMovesAcross: frmColNum is %v, step is %v, toColNum is %v, \nFaceUpPortion[0].Rank is %v, \nlen(b.columns[toColNum] is %v\n", frmColNum, step, toColNum, FaceUpPortion[0].Rank, len(b.columns[toColNum]))
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

				if singleGame && mc == specialMove {
					fmt.Printf("within MEC: initially m is %+v\n", m)
				}
				// recall that lower priority is better
				// if any of the three conditions which make a mec worth making exist
				// then lower the priority to the moveBasePriority minus the length of the face down portion of the from column
				// so a longer face down portion results lower priority
				// i have confirmed that minus wins more decks than plus

				if firstFaceUpIndex > 0 || //if there is a face down portion of the from column OR
					(firstFaceUpIndex == 0 && kingReady(b, frmColNum)) || //an empty column will result and there is a king ready to move there OR
					sisterCardInUpPortion(b, lastCard, toColNum, frmColNum) { //sister card in up portion of the other columns combined
					m.priority = moveBasePriority["moveEntireColumn"] - firstFaceUpIndex //Reassignment of moveBasePriority
				}

				if singleGame && mc == specialMove {
					fmt.Printf("within MEC: m is then adjusted to %+v\n", m)
				}

				moves = append(moves, m)
			}
			continue
		}
	}
	return moves
}
