package main

func detectAvailableMoves(b board, moveCounter int, singleGame bool, logicVersion string) []move {
	var aMoves []move //available Moves
	if logicVersion == "orig" {
		aMoves = append(aMoves, detectUpMoves(b, moveCounter)...)
		aMoves = append(aMoves, detectAcrossMoves(b, moveCounter)...)
		aMoves = append(aMoves, detectMecNotThoughtful(b, moveCounter, singleGame)...)
		aMoves = append(aMoves, detectDownMoves(b, moveCounter)...)
		aMoves = append(aMoves, detectPartialColumnMoves(b, moveCounter, singleGame)...)
		aMoves = append(aMoves, detectFlipStockToWaste(b, moveCounter)...)
		aMoves = append(aMoves, detectFlipWasteToStock(b, moveCounter)...)
	} else {
		aMoves = append(aMoves, detectUpMoves(b, moveCounter)...)
		aMoves = append(aMoves, detectAcrossMoves(b, moveCounter)...) // TODO consider commenting out
		aMoves = append(aMoves, detectMecNotThoughtful(b, moveCounter, singleGame)...)
		aMoves = append(aMoves, detectDownMoves(b, moveCounter)...) // TODO consider commenting out
		aMoves = append(aMoves, detectPartialColumnMoves(b, moveCounter, singleGame)...)
		aMoves = append(aMoves, detectMultiFlip(b, moveCounter)...) // Does flip stk->Waste, Waste->Stk and moveDown and move Across
	}
	return aMoves
}

var moveBasePriority = map[string]int{}

var moveBasePriorityPlayOrigOriginal = map[string]int{
	"moveAceUp":         100,
	"moveDeuceUp":       200,
	"moveAceAcross":     300,
	"moveDeuceAcross":   400,
	"moveDown":          500,
	"moveEntireColumn":  600,
	"movePartialColumn": 700,
	"move3PlusUp":       800,
	"move3PlusAcross":   900,
	"flpMMveAceAcross":  9600, // NOT USED WHEN LogicVersion = "original
	"flpMMve2Across":    9700, // NOT USED WHEN LogicVersion = "original
	"flpMMveDown":       9800, // NOT USED WHEN LogicVersion = "original
	"flpMMve3UpAcross":  9900, // NOT USED WHEN LogicVersion = "original
	"flipWasteToStock":  1000, //flip moves have the lowest priority
	"flipStockToWaste":  1100, //flip moves have the lowest priority
	"badMove":           1200, // a legal move which is worse than a mere flip

}

var moveBasePriorityPlayAllOriginal = map[string]int{
	"moveAceUp":         100,
	"moveDeuceUp":       200,
	"moveAceAcross":     300,
	"moveDeuceAcross":   400,
	"moveDown":          500,
	"moveEntireColumn":  600,
	"movePartialColumn": 700,
	"move3PlusUp":       800,
	"move3PlusAcross":   900,
	"flpMMveAceAcross":  9600, // NOT USED WHEN LogicVersion = "original
	"flpMMve2Across":    9700, // NOT USED WHEN LogicVersion = "original
	"flpMMveDown":       9800, // NOT USED WHEN LogicVersion = "original
	"flpMMve3UpAcross":  9900, // NOT USED WHEN LogicVersion = "original
	"flipWasteToStock":  1000, //flip moves have the lowest priority
	"flipStockToWaste":  1100, //flip moves have the lowest priority
	"badMove":           1200, // a legal move which is worse than a mere flip
}

var moveBasePriorityPlayOrigMultiflip = map[string]int{
	"moveAceUp":         100,
	"moveDeuceUp":       200,
	"moveAceAcross":     300,
	"moveDeuceAcross":   450,
	"moveDown":          500,
	"moveEntireColumn":  600,   // TODO TEST IF interlacing flpMMveDown and flpMMve3UpAcross is better
	"movePartialColumn": 700,   // TODO TEST IF FLIPPING RANK AND posStkWas is better
	"move3PlusUp":       800,   //                                                                posStkWas = Position of card to be
	"move3PlusAcross":   900,   //                                                                moved in combined Stock + Waste
	"flpMMveAceAcross":  2000,  // + 100 - posStkWas                      min, max = 1076, 1100
	"flpMMve2Across":    3000,  // + 100 - posStkWas                      min, max = 2076, 2100
	"flpMMveDown":       4000,  // 50 * (13 - Rank - 1) + 24 - PosStkWa   min, max = 4000, 4624 Move highest Rank first
	"flpMMve3UpAcross":  5000,  // 50 * (Rank - 1) + 24 - PosStkWa        min, max = 5000, 5624 Move lowest Rank first
	"flipWasteToStock":  10000, //flip moves have the lowest priority    NOT USED WHEN LogicVersion = "multiflip"
	"flipStockToWaste":  11000, //flip moves have the lowest priority    NOT USED WHEN LogicVersion = "multiflip"
	"badMove":           12000, // a legal move which is worse than a mere flip
}

var moveBasePriorityPlayAllMultiflip = map[string]int{
	"moveAceUp":         100,
	"moveDeuceUp":       200,
	"moveAceAcross":     300,
	"moveDeuceAcross":   450,
	"moveDown":          500,
	"moveEntireColumn":  600,   // TODO TEST IF interlacing flpMMveDown and flpMMve3UpAcross is better
	"movePartialColumn": 700,   // TODO TEST IF FLIPPING RANK AND posStkWas is better
	"move3PlusUp":       800,   //                                                                posStkWas = Position of card to be
	"move3PlusAcross":   900,   //                                                                moved in combined Stock + Waste
	"flpMMveAceAcross":  2000,  // + 100 - posStkWas                      min, max = 1076, 1100
	"flpMMve2Across":    3000,  // + 100 - posStkWas                      min, max = 2076, 2100
	"flpMMveDown":       4000,  // 50 * (13 - Rank - 1) + 24 - PosStkWa   min, max = 4000, 4624 Move highest Rank first
	"flpMMve3UpAcross":  5000,  // 50 * (Rank - 1) + 24 - PosStkWa        min, max = 5000, 5624 Move lowest Rank first
	"flipWasteToStock":  10000, //flip moves have the lowest priority    NOT USED WHEN LogicVersion = "multiflip"
	"flipStockToWaste":  11000, //flip moves have the lowest priority    NOT USED WHEN LogicVersion = "multiflip"
	"badMove":           12000, // a legal move which is worse than a mere flip
}

// ANY CHANGES IN THESE MUST BE MADE IN moveShortName8 BELOW!!!!!!!!!!!!!
var moveShortName = map[string]string{
	"moveAceUp":         " A Up ",
	"moveDeuceUp":       " 2 Up ",
	"moveAceAcross":     "A Acr ",
	"moveDeuceAcross":   "2 Acr ",
	"moveDown":          " Down ",
	"moveEntireColumn":  "EntCol",
	"movePartialColumn": "ParCol",
	"move3PlusUp":       " 3+Up ",
	"move3PlusAcross":   "3+Acr ",
	"flpMMveAceAcross":  "fAAccr",
	"flpMMve2Across":    "f2Accr",
	"flpMMveDown":       "f Down",
	"flpMMve3UpAcross":  "f3+Acr",
	"flipWasteToStock":  "W->Stk", //flip moves have the lowest priority
	"flipStockToWaste":  "Stk->W", //flip moves have the lowest priority
	"badMove":           "badMve", // a legal move which is worse than a mere flip
}

// Used to record how many of each move type is executed during an attempt.
var moveTypes = map[string]int{
	"moveAceUp":         0,
	"moveDeuceUp":       0,
	"moveAceAcross":     0,
	"moveDeuceAcross":   0,
	"moveDown":          0,
	"moveEntireColumn":  0,
	"movePartialColumn": 0,
	"move3PlusUp":       0,
	"move3PlusAcross":   0,
	"flpMMveAceAcross":  0,
	"flpMMve2Across":    0,
	"flpMMveDown":       0,
	"flpMMve3UpAcross":  0,
	"flipWasteToStock":  0,
	"flipStockToWaste":  0,
	"badMove":           0,
}
