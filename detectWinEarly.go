package main

func detectWinEarly(b board) bool {
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
