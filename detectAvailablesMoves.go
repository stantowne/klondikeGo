package main

func detectAvailableMoves(b board, moveCounter int, singleGame bool, passedLogicVersion string) []move {

	moveBasePriorityPlayOrigMultiflip2 = moveBasePriorityPlayOrigMultiflip1
	moveBasePriorityPlayOrigMultiflip3 = moveBasePriorityPlayOrigMultiflip1
	moveBasePriorityPlayAllMultiflip2 = moveBasePriorityPlayAllMultiflip1
	moveBasePriorityPlayAllMultiflip3 = moveBasePriorityPlayAllMultiflip1

	var aMoves []move                     //available Moves
	if passedLogicVersion == "original" { // TODO Figure out why is needed ????
		aMoves = append(aMoves, detectUpMoves(b, moveCounter)...)
		aMoves = append(aMoves, detectAcrossMoves(b, moveCounter)...)
		aMoves = append(aMoves, detectMecNotThoughtful(b, moveCounter, singleGame)...)
		aMoves = append(aMoves, detectDownMoves(b, moveCounter)...)
		aMoves = append(aMoves, detectPartialColumnMoves(b, moveCounter, singleGame)...)
		aMoves = append(aMoves, detectFlipStockToWaste(b, moveCounter)...)
		aMoves = append(aMoves, detectFlipWasteToStock(b, moveCounter)...)
	} else {
		aMoves = append(aMoves, detectMecNotThoughtful(b, moveCounter, singleGame)...)
		aMoves = append(aMoves, detectPartialColumnMoves(b, moveCounter, singleGame)...)
		aMoves = append(aMoves, detectMultiFlip(b, moveCounter, singleGame, passedLogicVersion)...) // Does flip stk->Waste, Waste->Stk and moveDown and move Across
	}
	return aMoves
}

var moveBasePriority = map[string]int{}

var moveBasePriorityPlayOrigOriginal = map[string]int{
	"moveAceUp":         10,
	"moveDeuceUp":       20,
	"moveAceAcross":     30,
	"moveDeuceAcross":   40,
	"moveDown":          50,
	"moveEntireColumn":  60,
	"movePartialColumn": 70,
	"move3PlusUp":       80,
	"move3PlusAcross":   90,
	"flipWasteToStock":  95, //flip moves have the lowest priority     NOT USED WHEN LogicVersion = "multiflip1" or "multiflip2" or "multiflip3"
	"flipStockToWaste":  99, //flip moves have the lowest priority     NOT USED WHEN LogicVersion = "multiflip1" or "multiflip2" or "multiflip3"
	//
	//                                "mFlM" means multi flip move
	//
	//                                posStkWa = Position of card to be moved in combined Stock + Waste
	//
	"mFlMAceUp":         11000,
	"mFlMDeuceUp":       12000,
	"mFlMAceAcross":     13000,
	"mFlMDeuceAcross":   14000,
	"mFlMDown":          15000,
	"mFlMEntireColumn":  16000,
	"mFlMPartialColumn": 17000,
	"mFlM3PlusUp":       18000,
	"mFlM3PlusAcross":   19000,
	"posStkWaBase":      24,
	"posStkWaSign":      -1,
	"posStkWaMultiple":  10_000_000,
	"origPrityMultiple": 1,
	"rankMultiple":      100_000,
	"badMove":           2_000_000_001, // a legal move which is worse than a mere flip (Max of int32 = 2,147,483,647)
}

var moveBasePriorityPlayOrigMultiflip1 = map[string]int{
	"moveAceUp":         10,
	"moveDeuceUp":       20,
	"moveAceAcross":     30,
	"moveDeuceAcross":   40,
	"moveDown":          50,
	"moveEntireColumn":  60,
	"movePartialColumn": 70,
	"move3PlusUp":       80,
	"move3PlusAcross":   90,
	"flipWasteToStock":  95, //flip moves have the lowest priority     NOT USED WHEN LogicVersion = "multiflip1" or "multiflip2" or "multiflip3"
	"flipStockToWaste":  99, //flip moves have the lowest priority     NOT USED WHEN LogicVersion = "multiflip1" or "multiflip2" or "multiflip3"
	//
	//                                "mFlM" means multi flip move
	//
	//                                posStkWa = Position of card to be moved in combined Stock + Waste
	//
	"mFlMAceUp":         11000,
	"mFlMDeuceUp":       12000,
	"mFlMAceAcross":     13000,
	"mFlMDeuceAcross":   14000,
	"mFlMDown":          15000,
	"mFlMEntireColumn":  16000,
	"mFlMPartialColumn": 17000,
	"mFlM3PlusUp":       18000,
	"mFlM3PlusAcross":   19000,
	"posStkWaBase":      24,
	"posStkWaSign":      -1,
	"posStkWaMultiple":  10_000_000,
	"origPrityMultiple": 1,
	"rankMultiple":      100_000,
	"badMove":           2_000_000_001, // a legal move which is worse than a mere flip (Max of int32 = 2,147,483,647)
}

var moveBasePriorityPlayOrigMultiflip2 = map[string]int{}

// When this is defined you MUST delete the appropriate line near 5 above

var moveBasePriorityPlayOrigMultiflip3 = map[string]int{}

// When this is defined you MUST delete the appropriate line near 5 above

var moveBasePriorityPlayAllOriginal = map[string]int{
	"moveAceUp":         10,
	"moveDeuceUp":       20,
	"moveAceAcross":     30,
	"moveDeuceAcross":   40,
	"moveDown":          50,
	"moveEntireColumn":  60,
	"movePartialColumn": 70,
	"move3PlusUp":       80,
	"move3PlusAcross":   90,
	"flipWasteToStock":  95, //flip moves have the lowest priority     NOT USED WHEN LogicVersion = "multiflip1" or "multiflip2" or "multiflip3"
	"flipStockToWaste":  99, //flip moves have the lowest priority     NOT USED WHEN LogicVersion = "multiflip1" or "multiflip2" or "multiflip3"
	//
	//                                "mFlM" means multi flip move
	//
	//                                posStkWa = Position of card to be moved in combined Stock + Waste
	//
	"mFlMAceUp":         11000,
	"mFlMDeuceUp":       12000,
	"mFlMAceAcross":     13000,
	"mFlMDeuceAcross":   14000,
	"mFlMDown":          15000,
	"mFlMEntireColumn":  16000,
	"mFlMPartialColumn": 17000,
	"mFlM3PlusUp":       18000,
	"mFlM3PlusAcross":   19000,
	"posStkWaBase":      24,
	"posStkWaSign":      -1,
	"posStkWaMultiple":  10_000_000,
	"origPrityMultiple": 1,
	"rankMultiple":      100_000,
	"badMove":           2_000_000_001, // a legal move which is worse than a mere flip (Max of int32 = 2,147,483,647)
}

var moveBasePriorityPlayAllMultiflip1 = map[string]int{
	"moveAceUp":         10,
	"moveDeuceUp":       20,
	"moveAceAcross":     30,
	"moveDeuceAcross":   40,
	"moveDown":          50,
	"moveEntireColumn":  60,
	"movePartialColumn": 70,
	"move3PlusUp":       80,
	"move3PlusAcross":   90,
	"flipWasteToStock":  95, //flip moves have the lowest priority     NOT USED WHEN LogicVersion = "multiflip1" or "multiflip2" or "multiflip3"
	"flipStockToWaste":  99, //flip moves have the lowest priority     NOT USED WHEN LogicVersion = "multiflip1" or "multiflip2" or "multiflip3"
	//
	//                                "mFlM" means multi flip move
	//
	//                                posStkWa = Position of card to be moved in combined Stock + Waste
	//
	"mFlMAceUp":         11000,
	"mFlMDeuceUp":       12000,
	"mFlMAceAcross":     13000,
	"mFlMDeuceAcross":   14000,
	"mFlMDown":          15000,
	"mFlMEntireColumn":  16000,
	"mFlMPartialColumn": 17000,
	"mFlM3PlusUp":       18000,
	"mFlM3PlusAcross":   19000,
	"posStkWaBase":      24,
	"posStkWaSign":      -1,
	"posStkWaMultiple":  10_000_000,
	"origPrityMultiple": 1,
	"rankMultiple":      100_000,
	"badMove":           2_000_000_001, // a legal move which is worse than a mere flip (Max of int32 = 2,147,483,647)
}

var moveBasePriorityPlayAllMultiflip2 = map[string]int{}

// When this is defined you MUST delete the appropriate line near 5 above

var moveBasePriorityPlayAllMultiflip3 = map[string]int{}

// When this is defined you MUST delete the appropriate line near 5 above

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
	"flipWasteToStock":  "W->Stk", //flip moves have the lowest priority
	"flipStockToWaste":  "Stk->W", //flip moves have the lowest priority
	"mFlMAceUp":         "f A Up",
	"mFlMDeuceUp":       "f 2 Up",
	"mFlMAceAcross":     "fA Acr",
	"mFlMDeuceAcross":   "f2 Acr",
	"mFlMDown":          "f Down",
	"mFlMEntireColumn":  "f EntC",
	"mFlMPartialColumn": "f ParC",
	"mFlM3PlusUp":       "f 3+Up",
	"mFlM3PlusAcross":   "f3+Acr",
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
	"flipWasteToStock":  0,
	"flipStockToWaste":  0,
	"mFlMAceUp":         0,
	"mFlMDeuceUp":       0,
	"mFlMAceAcross":     0,
	"mFlMDeuceAcross":   0,
	"mFlMDown":          0,
	"mFlMEntireColumn":  0,
	"mFlMPartialColumn": 0,
	"mFlM3PlusUp":       0,
	"mFlM3PlusAcross":   0,
	"badMove":           0,
}
