package main

func detectAvailableMoves(b board, moveCounter int, veryVerbose bool) []move {
	var aMoves []move //available Moves
	aMoves = append(aMoves, detectUpMoves(b, moveCounter, veryVerbose)...)
	aMoves = append(aMoves, detectAcrossMoves(b, moveCounter, veryVerbose)...)
	aMoves = append(aMoves, detectEntireColumnMoves(b, moveCounter, veryVerbose)...)
	aMoves = append(aMoves, detectDownMoves(b, moveCounter, veryVerbose)...)
	aMoves = append(aMoves, detectPartialColumnMoves(b, moveCounter, veryVerbose)...)
	aMoves = append(aMoves, detectFlipStockToWaste(b, moveCounter, veryVerbose)...)
	aMoves = append(aMoves, detectFlipWasteToStock(b, moveCounter, veryVerbose)...)
	return aMoves
}
