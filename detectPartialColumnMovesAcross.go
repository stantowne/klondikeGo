package main

import (
	"fmt"
	"os"
)

func detectPartialColumnMoves(b board, mc int, singleGame bool) []move {
	var moves []move
	if mc < 0 {
		return moves
	}
outer:
	for frmColNum := 0; frmColNum < 7; frmColNum++ { //the function must always complete this for loop
		firstFaceUpIndex, FaceUpPortion, err := faceUpPortion(b.columns[frmColNum])
		if singleGame && mc == 1 {
			fmt.Printf("detectPartialColumnMoves: frmColNum is %v, firstFaceUpIndex is %v, FaceUpPortion is %v\n", frmColNum, firstFaceUpIndex, FaceUpPortion)
		}
		if err != nil {
			fmt.Printf("detectPartialColumnMoves: error calling faceUpPortion on b.columns[%v] %v\n", frmColNum, err)
			os.Exit(1)
		}
		if len(FaceUpPortion) > 1 {
			for stepdown := 0; stepdown < len(FaceUpPortion); stepdown++ {
				candidateMoveUpCard := FaceUpPortion[stepdown]
				if singleGame && mc == 1 {
					fmt.Printf("within detectPartialColumnMoves: stepdown is %v, candidateMoveUpCard is %v\n", stepdown, candidateMoveUpCard)
				}
				if candidateMoveUpCard.Rank == len(b.piles[candidateMoveUpCard.Suit])+1 {
					sisterCard := card{Rank: candidateMoveUpCard.Rank, Suit: (candidateMoveUpCard.Suit + 2) % 4, FaceUp: true}
					// now, see if the sisterCard is the last card of another column
					for step := 1; step < 7; step++ {
						toColNum := (frmColNum + step) % 7
						lastCard, _, _ := last(b.columns[toColNum])
						if singleGame && mc == 1 {
							fmt.Printf("within detectPartialColumnMoves:\n"+
								"step is %v\n"+
								"toColNum is %v\n"+
								"candidateMoveUpCard is %v\n"+
								"sisterCard is %v\n"+
								"LastCard is %v\n", step, toColNum, candidateMoveUpCard, sisterCard, lastCard)
						}
						if lastCard == sisterCard {
							m := move{
								name:                "movePartialColumn",
								priority:            700,
								toPile:              candidateMoveUpCard.Suit,
								toCol:               toColNum,
								fromCol:             frmColNum,
								MovePortionStartIdx: firstFaceUpIndex + stepdown + 1,
								MovePortion:         FaceUpPortion[stepdown+1:],
							}
							moves = append(moves, m)
							if singleGame {
								fmt.Printf("detectPartialColumnMoves: moves is %v\n", moves)
							}
							continue outer //because each candidateMoveUpCard (e.g., 7 of Hearts) has only a single sister (7 of Diamonds)
						}
					}
				}
			}
		}
	}
	return moves
}
