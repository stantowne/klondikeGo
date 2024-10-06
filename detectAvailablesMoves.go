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

var moveBasePriority = map[string]int{
	"moveAceAcross":     300,
	"moveDeuceAcross":   400,
	"move3PlusAcross":   900,
	"moveDown":          500,
	"moveEntireColumn":  600,
	"flipWasteToStock":  1000, //flip moves have the lowest priority
	"flipStockToWaste":  1100, //flip moves have the lowest priority
	"movePartialColumn": 700,
	"moveAceUp":         100,
	"moveDeuceUp":       200,
	"move3PlusUp":       800,
	"badMove":           1200, // a legal move which is worse than a mere flip
}
