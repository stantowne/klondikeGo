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

/*
func flipStockToWasteAll(b board) board {
	for i := len(b.stock) - 1; i >= 0; i-- {
		b.waste = append(b.waste, b.stock[i].flipCardUp2())
	}
	b.stock = make([]Card, 0, len(b.waste))
	return b
}
*/

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
