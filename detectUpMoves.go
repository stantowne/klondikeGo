package main

import (
	"fmt"
	"os"
)

func detectUpMoves(b board, mc int) []move {
	var moves []move
	if mc < 0 {
		return moves
	}
	for i := 0; i < 7; i++ {
		if len(b.columns[i]) == 0 { // this column is empty
			continue
		}
		lastCard, _, err := last(b.columns[i])
		if err != nil {
			fmt.Printf("detectUpMoves: error calling last on b.columns[%v] %v\n", i, err)
			os.Exit(1)
		}
		// Below:  Moving an ace up is always a good move because an ace never is needed to hold the next lower card.
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
		// Below:  Moving a deuce up is always a good move because a deuce is never needed to hold the next lower card
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
		// Below:  At this point we know that last card has rank of 3 or more.
		// Such a card sometimes has utility as a last card so initially I assign the move up a low priority.
		if len(b.piles[lastCard.Suit]) == lastCard.Rank-1 {
			m := move{
				name:       "move3PlusUp",
				priority:   moveBasePriority["badMove"],
				fromCol:    i,
				toPile:     lastCard.Suit,
				cardToMove: lastCard,
			}
			// below:  The number in this line (8) was determined experimentally to produce the highest win rate.  Different deck samples could result in different results
			// below:  Move the last card up if its rank is greater than 8
			if (lastCard.Rank <= 8) ||
				//below:  Suppose the lastCard is 9 of spades (suit 2), which (if left down) could accept 8 of hearts or diamonds.
				//below:  First condition is satisfied if the length of the pile of hearts(suit 3), which is also the rank of the top card of the heart pile, is equal to or greater than 7
				//below:  This makes sense because (A) if it is 7, the 8 of hearts can be placed on the heart pile, and (B) if it is 8 the 8 of hearts is already on the heart pile
				//below:  The logical AND means the second condition must also be satisfied
				//below:  The second condition is satisfied in the 8 of diamonds is already in the face up portion of the columns
				//below:  So, if both conditions are true, there is no reason to retain the 9 of spades as the last card
				(len(b.piles[(lastCard.Suit+1)%4]) >= (lastCard.Rank-2) && //the length of one of the two piles of opposite color
					containedInAggregateUpPortion(b, Card{Rank: lastCard.Rank - 1, Suit: (lastCard.Suit + 3) % 4, FaceUp: true})) ||
				//below:  This is just the reverse of the above.  First condition diamonds; second condition hearts.
				(len(b.piles[(lastCard.Suit+3)%4]) >= (lastCard.Rank-2) &&
					containedInAggregateUpPortion(b, Card{Rank: lastCard.Rank - 1, Suit: (lastCard.Suit + 1) % 4, FaceUp: true})) {
				//below if any of the three tests are satisfied, give the move in question the move3PlusUp priority
				m.priority = moveBasePriority["move3PlusUp"]
				//Author is very proud of the above bit of code.
			}
			//below:  if both first conditions above are satisfied, there is no reason to test for the second condition.
			if len(b.piles[(lastCard.Suit+1)%4]) >= (lastCard.Rank-2) &&
				len(b.piles[(lastCard.Suit+3)%4]) >= (lastCard.Rank-2) {
				m.priority = moveBasePriority["move3PlusUp"]
			}

			//for consideration:  Maybe I should add if both second conditions are satisfied.

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
