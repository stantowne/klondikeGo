package main

import (
	"fmt"
	"strconv"
)

func printMove(m move) string {

	outS := moveShortName[m.name] + "  "
	switch m.name {
	case "moveAceAcross", "moveDeuceAcross", "move3PlusAcross":
		outS += "Move: " + m.cardToMove.pStr() + " from waste to the: " + string(m.cardToMove.suitSymbol()) + "Pile  Priority: " + strconv.Itoa(m.priority)
	case "moveAceUp", "moveDeuceUp", "move3PlusUp":
		outS += "Move: " + m.cardToMove.pStr() + " from column: " + strconv.Itoa(m.fromCol) + " To: " + string(m.cardToMove.suitSymbol()) + "Pile  Priority: " + strconv.Itoa(m.priority)
	case "moveDown":
		outS += "Move: " + m.cardToMove.pStr() + " from waste to column: " + strconv.Itoa(m.toCol) + "  Priority: " + strconv.Itoa(m.priority)
	case "moveEntireColumn":
		outS += "Move the cards starting with: " + m.MovePortion[0].pStr() + " From column: " + strconv.Itoa(m.fromCol) + " To column: " + strconv.Itoa(m.toCol) + "  Priority: " + strconv.Itoa(m.priority)
	case "movePartialColumn":
		//outS += "Move the cards starting with: " + m.MovePortion[0].pStr() + " From column: " + strconv.Itoa(m.fromCol) + " To column: " + strconv.Itoa(m.toCol) + "THEN move card " + m.cardToMove.pStr() + " to the: " + string(m.cardToMove.suitSymbol()) + "Pile  Priority: " + strconv.Itoa(m.priority)
		fmt.Printf("\n")
	case "flipStockToWaste":
		outS += "Move: Flip just the 3 top cards from stock to waste"
	case "flipWasteToStock":
		outS += "Move: Flip the entire waste pile to stock"
	case "":
		outS += "No Prior Move "
	default:
		outS += "Unknown move name " + m.name
	}

	return outS
}
