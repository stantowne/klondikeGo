package main

func detectFlipStockToWaste(b board, mc int) []move {
	var moves []move
	if mc < 0 {
		return moves
	}
	if len(b.stock) > 0 {
		m := move{
			name:     "flipStockToWaste",
			priority: moveBasePriority["flipStockToWaste"],
		}
		moves = append(moves, m)
		return moves
	}
	return moves
}

func detectFlipWasteToStock(b board, mc int) []move {
	var moves []move
	if mc < 0 {
		return moves
	}
	if len(b.stock) == 0 && len(b.waste) > 0 {
		m := move{
			name:     "flipWasteToStock",
			priority: moveBasePriority["flipWasteToStock"],
		}
		moves = append(moves, m)
		return moves
	}
	return moves
}
