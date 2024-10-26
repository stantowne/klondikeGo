package main

func detectWinEarly(b board) bool {
	// If stock and waste are empty and the top card (i.e. closest to the table) in all of the columns is face up,
	//    that implies that the rest of the column is face up and in order ready to be moved up to the piles
	// Then since we already determined that stock and waste are empty all cards not in the columns must be neatly
	//    stacked in the piles
	// Therefore we WIN

	if len(b.waste) != 0 {
		return false
	}
	if len(b.stock) != 0 {
		return false
	}
	colNum := 0
	for colNum = 0; colNum < 7; colNum++ {
		if len(b.columns[colNum]) > 0 && b.columns[colNum][0].FaceUp == false {
			return false
		}
	}
	return true
}
