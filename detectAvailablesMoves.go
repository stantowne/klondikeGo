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
	"flipSt->W Max-0":   9999, //flip MAX - 0 * 3 i.e. Up to 21 cards
	"flipSt->W Max-1":   9999, //flip MAX - 0 * 3 i.e. Up to 18 cards
	"flipSt->W Max-2":   9999, //flip MAX - 0 * 3 i.e. Up to 15 cards
	"flipSt->W Max-3":   9999, //flip MAX - 0 * 3 i.e. Up to 12 cards
	"flipSt->W Max-4":   9999, //flip MAX - 0 * 3 i.e. Up to  9 cards
	"flipSt->W Max-5":   9999, //flip MAX - 0 * 3 i.e. Up to  6 cards
	"flipSt->W Max-6":   9999, //flip MAX - 0 * 3 i.e. Up to  3 cards
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
	"flipSt->W Max-0":   9999, //flip MAX - 0 * 3 i.e. Up to 21 cards
	"flipSt->W Max-1":   9999, //flip MAX - 0 * 3 i.e. Up to 18 cards
	"flipSt->W Max-2":   9999, //flip MAX - 0 * 3 i.e. Up to 15 cards
	"flipSt->W Max-3":   9999, //flip MAX - 0 * 3 i.e. Up to 12 cards
	"flipSt->W Max-4":   9999, //flip MAX - 0 * 3 i.e. Up to  9 cards
	"flipSt->W Max-5":   9999, //flip MAX - 0 * 3 i.e. Up to  6 cards
	"flipSt->W Max-6":   9999, //flip MAX - 0 * 3 i.e. Up to  3 cards
}

// ANY CHANGES IN THESE MUST BE MADE IN moveShortName8 BELOW!!!!!!!!!!!!!
var moveShortName = map[string]string{
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
	"flipSt->W Max-0":   "S>W M0", //flip MAX - 0 * 3 i.e. Up to 21 cards
	"flipSt->W Max-1":   "S>W M1", //flip MAX - 0 * 3 i.e. Up to 18 cards
	"flipSt->W Max-2":   "S>W M2", //flip MAX - 0 * 3 i.e. Up to 15 cards
	"flipSt->W Max-3":   "S>W M3", //flip MAX - 0 * 3 i.e. Up to 12 cards
	"flipSt->W Max-4":   "S>W M4", //flip MAX - 0 * 3 i.e. Up to  9 cards
	"flipSt->W Max-5":   "S>W M5", //flip MAX - 0 * 3 i.e. Up to  6 cards
	"flipSt->W Max-6":   "S>W M6", //flip MAX - 0 * 3 i.e. Up to  3 cards
}

// ANY CHANGES IN THESE MUST BE MADE IN moveShortName ABOVE!!!!!!!!!!!!
//
//	These are used in playAllMoves func pmd when printing Board-byBoard detail in "TW" format (see arg[6] in main)
var moveShortName8 = map[string]string{
	"moveAceAcross":     " AAccr  ",
	"moveDeuceAcross":   " 2Accr  ",
	"move3PlusAcross":   " 3+Accr ",
	"moveDown":          "  Down  ",
	"moveEntireColumn":  " EntCol ",
	"flipWasteToStock":  " W->Stk ", //flip moves have the lowest priority
	"flipStockToWaste":  " Stk->W ", //flip moves have the lowest priority
	"movePartialColumn": " ParCol ",
	"moveAceUp":         "  A Up  ",
	"moveDeuceUp":       "  2 Up  ",
	"move3PlusUp":       "  3+Up  ",
	"badMove":           " badMve ", // a legal move which is worse than a mere flip
	"flipSt->W Max-0":   " S>W M0 ", //flip MAX - 0 * 3 i.e. Up to 21 cards
	"flipSt->W Max-1":   " S>W M1 ", //flip MAX - 0 * 3 i.e. Up to 18 cards
	"flipSt->W Max-2":   " S>W M2 ", //flip MAX - 0 * 3 i.e. Up to 15 cards
	"flipSt->W Max-3":   " S>W M3 ", //flip MAX - 0 * 3 i.e. Up to 12 cards
	"flipSt->W Max-4":   " S>W M4 ", //flip MAX - 0 * 3 i.e. Up to  9 cards
	"flipSt->W Max-5":   " S>W M5 ", //flip MAX - 0 * 3 i.e. Up to  6 cards
	"flipSt->W Max-6":   " S>W M6 ", //flip MAX - 0 * 3 i.e. Up to  3 cards
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
	"badMove":           0,
	"flipSt->W Max-0":   0,
	"flipSt->W Max-1":   0,
	"flipSt->W Max-2":   0,
	"flipSt->W Max-3":   0,
	"flipSt->W Max-4":   0,
	"flipSt->W Max-5":   0,
	"flipSt->W Max-6":   0,
}
