package main

func detectAvailableMoves(b board, moveCounter int, singleGame bool) []move {
	var aMoves []move //available Moves
	aMoves = append(aMoves, detectUpMoves(b, moveCounter, singleGame)...)
	aMoves = append(aMoves, detectAcrossMoves(b, moveCounter, singleGame)...)
	aMoves = append(aMoves, detectMecNotThoughtful(b, moveCounter, singleGame)...)
	aMoves = append(aMoves, detectDownMoves(b, moveCounter, singleGame)...)
	aMoves = append(aMoves, detectPartialColumnMoves(b, moveCounter, singleGame)...)
	aMoves = append(aMoves, detectFlipStockToWaste(b, moveCounter, singleGame)...)
	aMoves = append(aMoves, detectFlipWasteToStock(b, moveCounter, singleGame)...)
	return aMoves
}
