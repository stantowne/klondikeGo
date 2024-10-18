package main

import "strconv"

func printMove(m move) string {

	outS := moveShortName[m.name] + "  "
	switch m.name {
	case "moveAceAcross", "moveDeuceAcross", "move3PlusAcross":
		outS += "Move: " + m.cardToMove.pStr() + " To the: " + string(m.cardToMove.suitSymbol()) + "Pile  Priority: " + strconv.Itoa(m.priority)
	case "moveAceUp", "moveDeuceUp", "move3PlusUp":
		outS += "Move: " + m.cardToMove.pStr() + " From column: " + strconv.Itoa(m.fromCol) + " To: " + string(m.cardToMove.suitSymbol()) + "Pile  Priority: " + strconv.Itoa(m.priority)
	case "moveDown":
		outS += "Move: " //+ m.cardToMove.pStr() + " From column: " + strconv.Itoa(m.fromCol) + " To: " + string(m.cardToMove.suitSymbol()) + "Pile  Priority: " + strconv.Itoa(m.priority)
	case "moveEntireColumn":
		outS += "Move the cards starting with: " + m.MovePortion[0].pStr() + " From column: " + strconv.Itoa(m.fromCol) + " To column: " + strconv.Itoa(m.toCol) + "  Priority: " + strconv.Itoa(m.priority)
	case "movePartialColumn":
		outS += "Move: "
	case "flipStockToWaste":
		outS += "Move: "
	case "flipWasteToStock":
		outS += "Move: "
	}

	return outS
}
