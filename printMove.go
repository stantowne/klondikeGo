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
	case "mFlMAceAcross", "mFlMDeuceAcross", "flpMMve3UpAcross", "flpMMveDown": // MultiFlip
		outS += "Flip the stock/waste until the " + strconv.Itoa(m.MovePortionStartIdx+1)
		switch (m.MovePortionStartIdx + 1) % 10 {
		case 1:
			outS += "st "
		case 2:
			outS += "nd "
		case 3:
			outS += "rd "
		default:
			outS += "th "
		}
		outS += "card " + m.cardToMove.pStrC() + "is on top of waste pile."
		if m.name == "flpMMveDown" {
			outS2 += "                           then move the " + m.cardToMove.pStrC() + "down from waste to column " + strconv.Itoa(m.toCol) + "\n"
		} else {
			outS2 += "                           then move the " + m.cardToMove.pStrC() + "across from waste to the " + string(m.cardToMove.suitSymbolColored()) + "Pile\n"
		}
	case "moveAceAcross", "moveDeuceAcross", "move3PlusAcross":
		outS += "Move the " + m.cardToMove.pStrC() + "from waste to the " + string(m.cardToMove.suitSymbolColored()) + "Pile"
	case "moveAceUp", "moveDeuceUp", "move3PlusUp":
		outS += "Move the " + m.cardToMove.pStrC() + "up from column " + strconv.Itoa(m.fromCol) + " to the " + string(m.cardToMove.suitSymbolColored()) + "Pile"
	case "moveDown":
		outS += "Move the " + m.cardToMove.pStrC() + "down from waste to column " + strconv.Itoa(m.toCol)
	case "moveEntireColumn":
		outS += "Move the cards starting with " + m.MovePortion[0].pStrC() + "from column: " + strconv.Itoa(m.fromCol) + " to column: " + strconv.Itoa(m.toCol)
	case "movePartialColumn":
		outS += "Move the cards starting with: " + m.MovePortion[0].pStrC() + "from column " + strconv.Itoa(m.fromCol) + " to column: " + strconv.Itoa(m.toCol)
		outS2 += "\n                           then move the card above " + m.MovePortion[0].pStrC() + "in column " + strconv.Itoa(m.fromCol) + " to the appropriate pile based on its suit.\n"

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
