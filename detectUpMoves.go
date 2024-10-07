package main

import (
	"fmt"
	"os"
)

func detectUpMoves(b board, mc int, _ bool) []move {
	var moves []move
	if mc < 0 {
		return moves
	}
	for i := 0; i < 7; i++ {
		lastCard, _, err := last(b.columns[i])
		if err != nil {
			fmt.Printf("detectUpMoves: error calling last on b.columns[%v] %v\n", i, err)
			os.Exit(1)
		}
		if lastCard.Rank == 1 {
			m := move{
				name:       "moveAceUp",
				priority:   moveBasePriority["moveAceUp"],
				fromCol:    i,
				toPile:     lastCard.Suit,
				cardToMove: lastCard,
			}
			moves = append(moves, m)
			continue
		}
		if lastCard.Rank == 2 && len(b.piles[lastCard.Suit]) == lastCard.Rank-1 {
			m := move{
				name:       "moveDeuceUp",
				priority:   moveBasePriority["moveDeuceUp"],
				fromCol:    i,
				toPile:     lastCard.Suit,
				cardToMove: lastCard,
			}
			moves = append(moves, m)
			continue
		}
		if len(b.piles[lastCard.Suit]) == lastCard.Rank-1 {
			m := move{
				name:       "move3PlusUp",
				priority:   moveBasePriority["badMove"],
				fromCol:    i,
				toPile:     lastCard.Suit,
				cardToMove: lastCard,
			}
			if (lastCard.Rank <= 8) ||
				(len(b.piles[(lastCard.Suit+1)%4]) >= (lastCard.Rank-2) &&
					containedInAggregateUpPortion(b, Card{Rank: lastCard.Rank - 1, Suit: (lastCard.Suit + 3) % 4, FaceUp: true})) ||
				(len(b.piles[(lastCard.Suit+3)%4]) >= (lastCard.Rank-2) &&
					containedInAggregateUpPortion(b, Card{Rank: lastCard.Rank - 1, Suit: (lastCard.Suit + 1) % 4, FaceUp: true})) {
				m.priority = moveBasePriority["move3PlusUp"]
			}
			if len(b.piles[(lastCard.Suit+1)%4]) >= (lastCard.Rank-2) &&
				len(b.piles[(lastCard.Suit+3)%4]) >= (lastCard.Rank-2) {
				m.priority = moveBasePriority["move3PlusUp"]
			}

			moves = append(moves, m)
		}
	}
	return moves
}

func containedInAggregateUpPortion(b board, c Card) bool {
	for _, col := range b.columns {
		for _, crd := range col {
			if crd == c {
				return true
			}
		}
	}
	return false
}
