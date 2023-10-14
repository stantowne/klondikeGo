package main

import (
	"fmt"
)

func detectUpMoves(b board, mc int, singleGame bool) []move {
	verbose := 0
	var moves []move
	if mc < 0 {
		return moves
	}
	for i := 0; i < 7; i++ {
		lastCard, _, err := last(b.columns[i])
		if err != nil {
			if verbose > 0 {
				fmt.Printf("detectUpMoves: examining cloumn %v result is %v\n", i, err)
			}
			continue
		}
		if lastCard.Rank == 1 {
			m := move{
				name:       "moveAceUp",
				priority:   100,
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
				priority:   200,
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
				priority:   1200,
				fromCol:    i,
				toPile:     lastCard.Suit,
				cardToMove: lastCard,
			}
			if (lastCard.Rank <= 8) ||
				(len(b.piles[(lastCard.Suit+1)%4]) >= (lastCard.Rank-2) &&
					containedInAggregateUpPortion(b, card{Rank: lastCard.Rank - 1, Suit: (lastCard.Suit + 3) % 4, FaceUp: true})) ||
				(len(b.piles[(lastCard.Suit+3)%4]) >= (lastCard.Rank-2) &&
					containedInAggregateUpPortion(b, card{Rank: lastCard.Rank - 1, Suit: (lastCard.Suit + 1) % 4, FaceUp: true})) {
				m.priority = 800
			}
			if len(b.piles[(lastCard.Suit+1)%4]) >= (lastCard.Rank-2) &&
				len(b.piles[(lastCard.Suit+3)%4]) >= (lastCard.Rank-2) {
				m.priority = 800
			}

			moves = append(moves, m)
		}
	}
	return moves
}

func containedInAggregateUpPortion(b board, c card) bool {
	for _, col := range b.columns {
		for _, crd := range col {
			if crd == c {
				return true
			}
		}
	}
	return false
}
