package main

import (
	"strconv"
	"strings"
)

func printMove(m move, useLongName bool) (string, string) {
	var outS, outS2 string
	if useLongName {
		outS = m.name + strings.Repeat(" ", 20-len(m.name)) + "  "
	} else {
		outS = moveShortName[m.name] + "  "
	}
	outS2 = ""
	switch m.name {
	case "moveAceAcross", "moveDeuceAcross", "move3PlusAcross":
		outS += "Move the " + m.cardToMove.pStr() + "from waste to the " + string(m.cardToMove.suitSymbolColored()) + "Pile"
	case "moveAceUp", "moveDeuceUp", "move3PlusUp":
		outS += "Move the " + m.cardToMove.pStr() + "up from column " + strconv.Itoa(m.fromCol) + " to the " + string(m.cardToMove.suitSymbolColored()) + "Pile"
	case "moveDown":
		outS += "Move the " + m.cardToMove.pStr() + "down from waste to column " + strconv.Itoa(m.toCol)
	case "moveEntireColumn":
		outS += "Move the cards starting with " + m.MovePortion[0].pStr() + "from column: " + strconv.Itoa(m.fromCol) + " to column: " + strconv.Itoa(m.toCol)
	case "movePartialColumn":
		outS += "Move the cards starting with: " + m.MovePortion[0].pStr() + "from column " + strconv.Itoa(m.fromCol) + " to column: " + strconv.Itoa(m.toCol)
		outS2 += "\n                           then move the card above " + m.MovePortion[0].pStr() + "in column " + strconv.Itoa(m.fromCol) + " to the appropriate pile based on its suit."

		/*outS = "movePartialColumn"*/
	case "flipStockToWaste":
		outS += "Flip just the 3 (or fewer) top cards from stock to waste"
	case "flipWasteToStock":
		outS += "Flip the entire waste pile to stock"
	case "":
		outS += "No Prior Move "
	default:
		outS += "Unknown move name " + m.name
	}
	if useLongName {
		if len(outS2) != 0 {
			outS2 = outS2[2:]
		}
	} else {
		if outS != "No Prior Move " && outS != "Unknown move name "+m.name {
			outS = outS + "  Priority: " + strconv.Itoa(m.priority)
		}
	}
	return outS, outS2
}
