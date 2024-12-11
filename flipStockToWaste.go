package main

import (
	"fmt"
	"os"
)

func flipStockToWaste(b board) board {
	l := len(b.stock)
	if l > 2 {
		b.waste = append(b.waste, b.stock[l-1].flipCardUp2(), b.stock[l-2].flipCardUp2(), b.stock[l-3].flipCardUp2())
		b.stock = b.stock[:l-3]
	} else if l == 2 {
		b.waste = append(b.waste, b.stock[l-1].flipCardUp2(), b.stock[l-2].flipCardUp2())
		b.stock = b.stock[:l-2]
	} else if l == 1 {
		b.waste = append(b.waste, b.stock[l-1].flipCardUp2())
		b.stock = b.stock[:l-1]
	} else {
		fmt.Printf("Error: attempted to flip from empty stock")
		os.Exit(1)
	}
	return b
}

/*func flipStockToWasteV(b board, v int) board {
	l := len(b.stock)
	if l <= 0 {
		fmt.Printf("Error: attempted to flip from empty stock")
		os.Exit(1)
	}
	b.waste = append(b.waste, b.stock[l-1].flipCardUp2())
	b.stock = b.stock[:l-1]
	if l > 1 {
		flipStockToWasteV(b, v-1)
	}
	return b
}
*/
