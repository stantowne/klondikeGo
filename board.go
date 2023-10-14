package main

type board struct {
	columns [7]column //column[x][0] rests on the table for all x
	piles   [4]pile   //  pile[y][0] rests on the table for all y
	stock   []card    //    stock[0] rests on the table
	waste   []card    //    waste[0] rests on the table
}

// detects if there is a king ready to move to empty column
func kingReady(b board, fromColNum int, toColNum int) bool {
	if len(b.waste) > 0 && b.waste[len(b.waste)-1].Rank == 13 {
		return true
	}
	for i, col := range b.columns {
		if i == fromColNum || len(col) == 0 {
			continue
		}
		firstFUIndex, _, _ := faceUpPortion(col)
		if firstFUIndex == 0 {
			continue
		}
		if col[firstFUIndex].Rank == 13 {
			return true
		}
	}
	return false
}

func sisterCardInUpPortion(b board, c card, toColNum int, fromColNum int) bool {
	sisterCard := card{
		Rank:   c.Rank,
		Suit:   (c.Suit + 2) % 4,
		FaceUp: true,
	}
	for i, col := range b.columns {
		if i == fromColNum || i == toColNum || len(col) == 0 {
			continue
		}
		for _, crd := range b.columns[i] {
			if crd == sisterCard {
				return true
			}
		}
	}
	return false
}
