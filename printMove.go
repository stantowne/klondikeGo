package main

import (
	"strconv"
)

func printMove(m move) string {

	outS := moveShortName[m.name] + "  "
	switch m.name {
	case "moveAceAcross", "moveDeuceAcross", "move3PlusAcross":
		outS += "Move: " + m.cardToMove.pStr() + " from waste to the: " + string(m.cardToMove.suitSymbol()) + "Pile  Priority: " + strconv.Itoa(m.priority)
	case "moveAceUp", "moveDeuceUp", "move3PlusUp":
		outS += "Move: " + m.cardToMove.pStr() + " up from column: " + strconv.Itoa(m.fromCol) + " To: " + string(m.cardToMove.suitSymbol()) + "Pile  Priority: " + strconv.Itoa(m.priority)
	case "moveDown":
		outS += "Move: " + m.cardToMove.pStr() + " down from waste to column: " + strconv.Itoa(m.toCol) + "  Priority: " + strconv.Itoa(m.priority)
	case "moveEntireColumn":
		outS += "Move: The cards starting with: " + m.MovePortion[0].pStr() + " From column: " + strconv.Itoa(m.fromCol) + " To column: " + strconv.Itoa(m.toCol) + "  Priority: " + strconv.Itoa(m.priority)
	case "movePartialColumn":
		/*	outS += "Move: the cards starting with: " + m.MovePortion[0].pStr() + " From column: " + strconv.Itoa(m.fromCol) + " To column: " + strconv.Itoa(m.toCol) + "\n" +
			"                         then:  move the card above " + m.MovePortion[0].pStr() + "in column: " + strconv.Itoa(m.fromCol) + " to the appropriate pile based on its suit.  Priority: " + strconv.Itoa(m.priority)*/
		outS += "   movePartialColumn"
	case "flipStockToWaste":
		outS += "Move: Flip just the 3 (or fewer) top cards from stock to waste"
	case "flipWasteToStock":
		outS += "Move: Flip the entire waste pile to stock"
	case "":
		outS += "No Prior Move "
	default:
		outS += "Unknown move name " + m.name
	}

	return outS
}
