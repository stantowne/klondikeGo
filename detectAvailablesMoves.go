package main

func detectAvailableMoves(b board, moveCounter int, singleGame bool, passedLogicVersion string) []move {
	var aMoves []move                                                         //available Moves
	if passedLogicVersion == "original" || passedLogicVersion == "playOrig" { // TODO Figure out why is needed ????
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
		aMoves = append(aMoves, detectMultiFlip(b, moveCounter, passedLogicVersion)...) // Does flip stk->Waste, Waste->Stk and moveDown and move Across
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
	// "flpMMveAceAcross":  9600, //                                            NOT USED WHEN LogicVersion = "original   See TODO in detectAvailableMoves! as of now needs ' or cfg.General.Type = "playOrig" '
	// "flpMMve2Across":    9700, //                                            NOT USED WHEN LogicVersion = "original
	// "flpMMveDown":       9800, //                                            NOT USED WHEN LogicVersion = "original
	// "flpMMve3UpAcross":  9900, //                                            NOT USED WHEN LogicVersion = "original
	"flipWasteToStock": 1000,          //flip moves have the lowest priority            NOT USED WHEN LogicVersion = any "multiflip"
	"flipStockToWaste": 1100,          //flip moves have the lowest priority            NOT USED WHEN LogicVersion = any "multiflip"
	"badMove":          2_147_483_647, // a legal move which is worse than a mere flip (Max of int32 = 2,147,483,647)
}

// Need these in future when TODO in detectAvailableMoves is fixed
var moveBasePriorityPlayOrigMultiflip = map[string]int{
	"badMove": 2_147_483_647, // a legal move which is worse than a mere flip (Max of int32 = 2,147,483,647)
}

var moveBasePriorityPlayOrigMultiflipTest = map[string]int{
	"badMove": 2_147_483_647, // a legal move which is worse than a mere flip (Max of int32 = 2,147,483,647)
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
	// "flpMMveAceAcross":  9600, //                                	        NOT USED WHEN LogicVersion = "original   See TODO in detectAvailableMoves! as of now needs ' or cfg.General.Type = "playOrig" '
	// "flpMMve2Across":    9700, //                                	        NOT USED WHEN LogicVersion = "original
	// "flpMMveDown":       9800, //                                	        NOT USED WHEN LogicVersion = "original
	// "flpMMve3UpAcross":  9900, //                                	        NOT USED WHEN LogicVersion = "original
	"flipWasteToStock": 1000, //flip moves have the lowest priority            NOT USED WHEN LogicVersion = any "multiflip"
	"flipStockToWaste": 1100, //flip moves have the lowest priority            NOT USED WHEN LogicVersion = any "multiflip"
	"badMove":          1200, // a legal move which is worse than a mere flip
}

/*	"flpMMveAceAcross":  4000,  // + 50 * (24 - posStkWa) + (rank - 1 )   min, max = 4000, 5212
	"flpMMve2Across":    4001,  // + 50 * (24 - posStkWa) + (rank - 1 )   min, max = 4001, 5213
	"flpMMveDown":       4002,  // + 50 * (24 - posStkWa) + (13 - rank)   min, max = 4002, 5214 Move highest rank first
	"flpMMve3UpAcross":  4003,  // + 50 * (24 - posStkWa) + (rank - 1 )   min, max = 4003, 5215 Move lowest rank first
	"flipWasteToStock":  10000, //flip moves have the lowest priority     NOT USED WHEN LogicVersion = "multiflip"
	"flipStockToWaste":  11000, //flip moves have the lowest priority     NOT USED WHEN LogicVersion = "multiflip"
	"badMove":   2_147_483_647, // a legal move which is worse than a mere flip (Max of int32 = 2,147,483,647)
*/

var moveBasePriorityPlayAllMultiflip = map[string]int{
	"moveAceUp":         100,
	"moveDeuceUp":       200,
	"moveAceAcross":     300,
	"moveDeuceAcross":   450,
	"moveDown":          500,
	"moveEntireColumn":  600,  // TODO TEST IF interlacing flpMMveDown and flpMMve3UpAcross is better
	"movePartialColumn": 700,  // TODO TEST IF FLIPPING RANK AND posStkWas is better
	"move3PlusUp":       800,  //                                                                posStkWas = Position of card to be
	"move3PlusAcross":   900,  //                                                                moved in combined Stock + Waste
	"flpMMveAceAcross":  4000, // + posStkWaMultiple * (posStkWaBase + posStkWaSign * posStkWa) + (rank - 1 )   min, max = 4000, 5212 Move lowest rank first
	"flpMMve2Across":    4001, // + posStkWaMultiple * (posStkWaBase + posStkWaSign * posStkWa) + (rank - 1 )   min, max = 4001, 5213 Move lowest rank first
	"flpMMveDown":       4002, // + posStkWaMultiple * (posStkWaBase + posStkWaSign * posStkWa) + (13 - rank)   min, max = 4002, 5214 Move highest rank first
	"flpMMve3UpAcross":  4003, // + posStkWaMultiple * (posStkWaBase + posStkWaSign * posStkWa) + (rank - 1 )   min, max = 4003, 5215 Move lowest rank first
	//"flipWasteToStock":  10000, //flip moves have the lowest priority     NOT USED WHEN LogicVersion = "multiflip"
	//"flipStockToWaste":  11000, //flip moves have the lowest priority     NOT USED WHEN LogicVersion = "multiflip"
	"posStkWaBase":     24,            //
	"posStkWaSign":     -1,            //
	"posStkWaMultiple": 50,            //
	"badMove":          2_147_483_647, // a legal move which is worse than a mere flip (Max of int32 = 2,147,483,647)
}

var moveBasePriorityPlayAllMultiflipTest = map[string]int{
	"moveAceUp":         100,
	"moveDeuceUp":       200,
	"moveAceAcross":     300,
	"moveDeuceAcross":   450,
	"moveDown":          500,
	"moveEntireColumn":  600,  // TODO TEST IF interlacing flpMMveDown and flpMMve3UpAcross is better
	"movePartialColumn": 700,  // TODO TEST IF FLIPPING RANK AND posStkWas is better
	"move3PlusUp":       800,  //                                                                posStkWas = Position of card to be
	"move3PlusAcross":   900,  //                                                                moved in combined Stock + Waste
	"flpMMveAceAcross":  4000, // + posStkWaMultiple * (posStkWaBase + posStkWaSign * posStkWa) + (rank - 1 )   min, max = 4000, 5212 Move lowest rank first
	"flpMMve2Across":    4001, // + posStkWaMultiple * (posStkWaBase + posStkWaSign * posStkWa) + (rank - 1 )   min, max = 4001, 5213 Move lowest rank first
	"flpMMveDown":       4002, // + posStkWaMultiple * (posStkWaBase + posStkWaSign * posStkWa) + (13 - rank)   min, max = 4002, 5214 Move highest rank first
	"flpMMve3UpAcross":  4003, // + posStkWaMultiple * (posStkWaBase + posStkWaSign * posStkWa) + (rank - 1 )   min, max = 4003, 5215 Move lowest rank first
	//"flipWasteToStock":  10000, //flip moves have the lowest priority     NOT USED WHEN LogicVersion = "multiflip"
	//"flipStockToWaste":  11000, //flip moves have the lowest priority     NOT USED WHEN LogicVersion = "multiflip"
	"posStkWaBase":     24,            //
	"posStkWaSign":     -1,            //
	"posStkWaMultiple": 50,            //
	"badMove":          2_147_483_647, // a legal move which is worse than a mere flip (Max of int32 = 2,147,483,647)
}

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
