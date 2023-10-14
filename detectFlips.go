package main

func detectFlipStockToWaste(b board, mc int, singleGame bool) []move {
	var moves []move
	if mc < 0 {
		return moves
	}
	if len(b.stock) > 0 {
		m := move{
			name:     "flipStockToWaste",
			priority: 1000,
		}
		moves = append(moves, m)
		return moves
	}
	return moves
}

func detectFlipWasteToStock(b board, mc int, singleGame bool) []move {
	var moves []move
	if mc < 0 {
		return moves
	}
	if len(b.stock) == 0 && len(b.waste) > 0 {
		m := move{
			name:     "flipWasteToStock",
			priority: 1100,
		}
		moves = append(moves, m)
		return moves
	}
	return moves
}
