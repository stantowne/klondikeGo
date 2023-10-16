package main

import (
	"fmt"
	"os"
)

func detectAcrossMoves(b board, mc int, _ bool) []move {
	var moves []move
	if mc < 0 {
		return moves
	}
	var m move
	lastWasteCard, _, err := last(b.waste)
	if err != nil {
		fmt.Printf("detectAcrossMoves: error calling last ob b.waste %v\n", err)
		os.Exit(1)
	}

	if lastWasteCard.Rank == 1 {
		m = move{
			name:       "moveAceAcross",
			priority:   moveBasePriority["moveAceAcross"],
			toPile:     lastWasteCard.Suit,
			cardToMove: lastWasteCard,
		}
		moves = append(moves, m)
		return moves
	}
	if lastWasteCard.Rank == 2 && len(b.piles[lastWasteCard.Suit]) == lastWasteCard.Rank-1 {
		m = move{
			name:       "moveDeuceAcross",
			priority:   moveBasePriority["moveDeuceAcross"],
			toPile:     lastWasteCard.Suit,
			cardToMove: lastWasteCard,
		}
		moves = append(moves, m)
		return moves
	}
	if len(b.piles[lastWasteCard.Suit]) == lastWasteCard.Rank-1 {
		m = move{
			name:       "move3PlusAcross",
			priority:   moveBasePriority["move3PlusAcross"],
			toPile:     lastWasteCard.Suit,
			cardToMove: lastWasteCard,
		}
		moves = append(moves, m)

	}
	return moves
}
