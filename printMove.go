package main

import (
	"strconv"
)

func printMove(m move) (string, string) {
	outS := moveShortName[m.name] + "  "
	outS2 := ""
	switch m.name {
	case "moveAceAcross", "moveDeuceAcross", "move3PlusAcross":
		outS += "Move: " + m.cardToMove.pStr() + "from waste to the " + string(m.cardToMove.suitSymbolColored()) + "Pile  Priority: " + strconv.Itoa(m.priority)
	case "moveAceUp", "moveDeuceUp", "move3PlusUp":
		outS += "Move: " + m.cardToMove.pStr() + "up from column " + strconv.Itoa(m.fromCol) + " to the " + string(m.cardToMove.suitSymbolColored()) + "Pile  Priority: " + strconv.Itoa(m.priority)
	case "moveDown":
		outS += "Move: " + m.cardToMove.pStr() + "down from waste to column " + strconv.Itoa(m.toCol) + "  Priority: " + strconv.Itoa(m.priority)
	case "moveEntireColumn":
		outS += "Move: The cards starting with " + m.MovePortion[0].pStr() + "from column: " + strconv.Itoa(m.fromCol) + " to column: " + strconv.Itoa(m.toCol) + "  Priority: " + strconv.Itoa(m.priority)
	case "movePartialColumn":
		outS += "Move: The cards starting with: " + m.MovePortion[0].pStr() + "from column " + strconv.Itoa(m.fromCol) + " to column: " + strconv.Itoa(m.toCol) + "  Priority: " + strconv.Itoa(m.priority)
		outS2 += "\n                                 then: move the card above " + m.MovePortion[0].pStr() + "in column " + strconv.Itoa(m.fromCol) + " to the appropriate pile based on its suit."
	case "flipStockToWaste":
		outS += "Move: Flip just the 3 (or fewer) top cards from stock to waste"
	case "flipWasteToStock":
		outS += "Move: Flip the entire waste pile to stock"
	case "":
		outS += "No Prior Move "
	default:
		outS += "Unknown move name " + m.name
	}

	return outS, outS2
}
