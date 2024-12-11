package main

func detectAvailableMoves(b board, moveCounter int, singleGame bool) []move {
	var aMoves []move //available Moves
	aMoves = append(aMoves, detectUpMoves(b, moveCounter)...)
	aMoves = append(aMoves, detectAcrossMoves(b, moveCounter)...)
	aMoves = append(aMoves, detectMecNotThoughtful(b, moveCounter, singleGame)...)
	aMoves = append(aMoves, detectDownMoves(b, moveCounter)...)
	aMoves = append(aMoves, detectPartialColumnMoves(b, moveCounter, singleGame)...)
	aMoves = append(aMoves, detectFlipStockToWaste(b, moveCounter)...)
	aMoves = append(aMoves, detectFlipWasteToStock(b, moveCounter)...)
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
	"flipSt->W Max-0":   9990, //flip MAX - (0 * 3)    i.e. Up to 24 cards    where MAX = len(stock) + len(waste)   Suggest flip waste=>stock then flip cards
	"flipSt->W Max-1":   9991, //flip MAX - (1 * 3)    i.e. Up to 21 cards
	"flipSt->W Max-2":   9992, //flip MAX - (2 * 3)    i.e. Up to 18 cards    where 999 > priority of any moves that are between columns
	"flipSt->W Max-3":   9993, //flip MAX - (3 * 3)    i.e. Up to 15 cards                ????or that expose a new column card ????
	"flipSt->W Max-4":   9994, //flip MAX - (4 * 3)    i.e. Up to 12 cards              < priority of moves not in >
	"flipSt->W Max-5":   9995, //flip MAX - (5 * 3)    i.e. Up to  9 cards    ONLY When mod(len(stock)) == 0
	"flipSt->W Max-6":   9996, //flip MAX - (6 * 3)    i.e. Up to  6 cards
	"flipSt->W Max-7":   9997, //flip MAX - (7 * 3)    i.e. Up to  3 cards

}

var moveBasePriorityAll = map[string]int{
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
	"flipSt->W Max-0":   9990, //flip MAX - (0 * 3)    i.e. Up to 24 cards    where MAX = len(stock) + len(waste)   Suggest flip waste=>stock then flip cards
	"flipSt->W Max-1":   9991, //flip MAX - (1 * 3)    i.e. Up to 21 cards
	"flipSt->W Max-2":   9992, //flip MAX - (2 * 3)    i.e. Up to 18 cards    where 999 > priority of any moves that are beteen columns
	"flipSt->W Max-3":   9993, //flip MAX - (3 * 3)    i.e. Up to 15 cards                ????or that expose a new column card ????
	"flipSt->W Max-4":   9994, //flip MAX - (4 * 3)    i.e. Up to 12 cards              < priority of moves not in >
	"flipSt->W Max-5":   9995, //flip MAX - (5 * 3)    i.e. Up to  9 cards    ONLY When mod(len(stock)) == 0
	"flipSt->W Max-6":   9996, //flip MAX - (6 * 3)    i.e. Up to  6 cards
	"flipSt->W Max-7":   9997, //flip MAX - (7 * 3)    i.e. Up to  3 cards

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
	"flipSt->W Max-0":   "S>W M0",
	"flipSt->W Max-1":   "S>W M1",
	"flipSt->W Max-2":   "S>W M2",
	"flipSt->W Max-3":   "S>W M3",
	"flipSt->W Max-4":   "S>W M4",
	"flipSt->W Max-5":   "S>W M5",
	"flipSt->W Max-6":   "S>W M6",
	"flipSt->W Max-7":   "S>W M7",
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
	"flipSt->W Max-0":   " S>W M0 ",
	"flipSt->W Max-1":   " S>W M1 ",
	"flipSt->W Max-2":   " S>W M2 ",
	"flipSt->W Max-3":   " S>W M3 ",
	"flipSt->W Max-4":   " S>W M4 ",
	"flipSt->W Max-5":   " S>W M5 ",
	"flipSt->W Max-6":   " S>W M6 ",
	"flipSt->W Max-7":   " S>W M7 ",
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
	"flipSt->W Max-7":   0}
