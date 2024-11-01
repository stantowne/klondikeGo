package main

import (
	"fmt"
	"os"
)

func detectDownMoves(b board, mc int, _ bool) []move {
	var moves []move
	if mc < 0 {
		return moves
	}
	var m move
	if len(b.waste) == 0 { // waste is empty
		return moves
	}
	lastWasteCard, _, err := last(b.waste)
	if err != nil {
		fmt.Printf("detectDownMoves: error calling last on b.waste %v\n", err)
		os.Exit(1)
	}
	if lastWasteCard.Rank == 1 {
		m = move{
			name:       "moveAceAcross", //always better to move an ace across
			priority:   moveBasePriority["moveAceAcross"],
			toPile:     lastWasteCard.Suit,
			cardToMove: lastWasteCard,
		}
		moves = append(moves, m)
		return moves
	}
	if lastWasteCard.Rank == 2 && len(b.piles[lastWasteCard.Suit]) == 1 {
		m = move{
			name:       "moveDeuceAcross", //always better to move a deuce across, if possible
			priority:   moveBasePriority["moveDeuceAcross"],
			toPile:     lastWasteCard.Suit,
			cardToMove: lastWasteCard,
		}
		moves = append(moves, m)
		return moves
	}
	for i := 0; i < 7; i++ {
		lastColumnCard, _, err := last(b.columns[i])
		if err != nil {
			fmt.Printf("detectDownMoves: error calling last on b.columns[%vi] %v\n", i, err)
			continue
		}
		if (lastWasteCard.Rank == lastColumnCard.Rank-1 && lastWasteCard.color() != lastColumnCard.color()) ||
			(lastWasteCard.Rank == 13 && len(b.columns[i]) == 0) {
			m = move{
				name:       "moveDown",
				priority:   moveBasePriority["moveDown"],
				toCol:      i,
				cardToMove: lastWasteCard,
			}
			moves = append(moves, m)
		}
	}
	return moves
}
