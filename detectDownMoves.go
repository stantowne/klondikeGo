package main

import (
	"fmt"
	"os"
)

func detectDownMoves(b board, mc int) []move {
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
		return moves // there is never a reason to move an Ace down
	}
	if lastWasteCard.Rank == 2 && len(b.piles[lastWasteCard.Suit]) == 1 {
		return moves // there is never a reason to move a Deuce down if it can be moved over
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
	// if lastWasteCard.Rank < 13, len(moves) can be 0, 1 or 2
	// if lastWasteCard.Rank == 13, len(moves) can be as great as 7 (i.e., when all 7 columns are empty)
	if len(moves) > 1 {
		moves = moves[:1] // sub slice containing just 0th element. This may not be optimal
	}
	return moves
}
