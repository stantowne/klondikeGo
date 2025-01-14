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
		/* ??? */ aMoves = append(aMoves, detectAcrossMoves(b, moveCounter)...)
		aMoves = append(aMoves, detectMecNotThoughtful(b, moveCounter, singleGame)...)
		/* ??? */ aMoves = append(aMoves, detectDownMoves(b, moveCounter)...)
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
	"mMveAceAcross":     9600,
	"moveDeuceAcross":   400,
	"mMveDeuceAcross":   9700,
	"moveDown":          500,
	"mMveDown":          9800,
	"moveEntireColumn":  600,
	"movePartialColumn": 700,
	"move3PlusUp":       800,
	"move3PlusAcross":   900,
	"mMve3PlusAcross":   9900,
	"flipWasteToStock":  1000, //flip moves have the lowest priority
	"flipStockToWaste":  1100, //flip moves have the lowest priority
	"badMove":           1200, // a legal move which is worse than a mere flip
}

var moveBasePriorityPlayAllOriginal = map[string]int{
	"moveAceUp":         100,
	"moveDeuceUp":       200,
	"moveAceAcross":     300,
	"mMveAceAcross":     9700,
	"moveDeuceAcross":   400,
	"mMveDeuceAcross":   9800,
	"moveDown":          500,
	"mMveDown":          9900,
	"moveEntireColumn":  600,
	"movePartialColumn": 700,
	"move3PlusUp":       800,
	"move3PlusAcross":   900,
	"mMve3PlusAcross":   9900,
	"flipWasteToStock":  1000, //flip moves have the lowest priority
	"flipStockToWaste":  1100, //flip moves have the lowest priority
	"badMove":           1200, // a legal move which is worse than a mere flip
}

var moveBasePriorityPlayOrigMultiflip = map[string]int{
	"moveAceUp":         50,
	"moveDeuceUp":       100,
	"moveAceAcross":     200, //                                       PosStkWas = Position of card to be
	"mMveAceAcross":     200, // + 23 - PosStkWas min, max = 200, 223              moved in combined Stock + Waste
	"moveDeuceAcross":   250,
	"mMveDeuceAcross":   250, // + 23 - PosStkWas min, max = 250, 273
	"moveDown":          300,
	"mMveDown":          300, // 24 * (12 - Rank - 1) + 23 - PosStkWa min, max = 300, 599 Move highest Rank first
	"moveEntireColumn":  700,
	"movePartialColumn": 800,
	"move3PlusUp":       800,
	"move3PlusAcross":   800,
	"mMve3PlusAcross":   800,  // 24 * (Rank - 3) + 23 - PosStkWa min, max = 800, 1063 Move lowest Rank first
	"flipWasteToStock":  1200, //flip moves have the lowest priority
	"flipStockToWaste":  1300, //flip moves have the lowest priority
	"badMove":           1400, // a legal move which is worse than a mere flip
}

var moveBasePriorityPlayAllMultiflip = map[string]int{
	"moveAceUp":         50,
	"moveDeuceUp":       100,
	"moveAceAcross":     200, //                                       PosStkWas = Position of card to be
	"mMveAceAcross":     200, // + 23 - PosStkWas min, max = 200, 223              moved in combined Stock + Waste
	"moveDeuceAcross":   250,
	"mMveDeuceAcross":   250, // + 23 - PosStkWas min, max = 250, 273
	"moveDown":          300,
	"mMveDown":          300, // 24 * (12 - Rank - 1) + 23 - PosStkWa min, max = 300, 599 Move highest Rank first
	"moveEntireColumn":  700,
	"movePartialColumn": 800,
	"move3PlusUp":       800,
	"move3PlusAcross":   800,
	"mMve3PlusAcross":   800,  // 24 * (Rank - 3) + 23 - PosStkWa min, max = 800, 1063 Move lowest Rank first
	"flipWasteToStock":  1200, //flip moves have the lowest priority
	"flipStockToWaste":  1300, //flip moves have the lowest priority
	"badMove":           1400, // a legal move which is worse than a mere flip
}

// ANY CHANGES IN THESE MUST BE MADE IN moveShortName8 BELOW!!!!!!!!!!!!!
var moveShortName = map[string]string{
	"moveAceUp":         " A Up ",
	"moveDeuceUp":       " 2 Up ",
	"moveAceAcross":     "AAccr ",
	"mMveAceAcross":     "mAAccr",
	"moveDeuceAcross":   "2Accr ",
	"mMveDeuceAcross":   "m2Accr",
	"moveDown":          " Down ",
	"mMveDown":          "m Down",
	"moveEntireColumn":  "EntCol",
	"movePartialColumn": "ParCol",
	"move3PlusUp":       " 3+Up ",
	"move3PlusAcross":   "3+Accr",
	"mMve3PlusAcross":   "m3+Acr",
	"flipWasteToStock":  "W->Stk", //flip moves have the lowest priority
	"flipStockToWaste":  "Stk->W", //flip moves have the lowest priority
	"badMove":           "badMve", // a legal move which is worse than a mere flip
}

// Used to record how many of each move type is executed during an attempt.
var moveTypes = map[string]int{
	"moveAceUp":         0,
	"moveDeuceUp":       0,
	"moveAceAcross":     0,
	"mMveAceAcross":     0,
	"moveDeuceAcross":   0,
	"mMveDeuceAcross":   0,
	"moveDown":          0,
	"mMveDown":          0,
	"moveEntireColumn":  0,
	"movePartialColumn": 0,
	"move3PlusUp":       0,
	"move3PlusAcross":   0,
	"mMve3PlusAcross":   0,
	"flipWasteToStock":  0,
	"flipStockToWaste":  0,
	"badMove":           0,
}
