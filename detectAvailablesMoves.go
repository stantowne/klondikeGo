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

var moveBasePriority = map[string]int{}

var moveBasePriorityOrig = map[string]int{
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

var moveBasePriorityNew = map[string]int{
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

var moveShortNameNew = map[string]string{
	"moveAceAcross":     "AAccr ",
	"moveDeuceAcross":   "2Accr ",
	"move3PlusAcross":   "3+Accr",
	"moveDown":          " Down ",
	"moveEntireColumn":  "EntCol",
	"flipWasteToStock":  "W->Stk", //flip moves have the lowest priority
	"flipStockToWaste":  "Stk->W", //flip moves have the lowest priority
	"movePartialColumn": "ParCol",
	"moveAceUp":         " A Up ",
	"moveDeuceUp":       " 2 Up ",
	"move3PlusUp":       " 3+Up ",
	"badMove":           "badMve", // a legal move which is worse than a mere flip
}

// Used to record how many of each move type is executed during an attempt.
var moveTypes = map[string]int{
	"moveAceAcross":     0,
	"moveDeuceAcross":   0,
	"move3PlusAcross":   0,
	"moveDown":          0,
	"moveEntireColumn":  0,
	"flipWasteToStock":  0,
	"flipStockToWaste":  0,
	"movePartialColumn": 0,
	"moveAceUp":         0,
	"moveDeuceUp":       0,
	"move3PlusUp":       0,
}
