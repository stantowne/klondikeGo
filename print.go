package main

import (
	"fmt"
	"strconv"
)

func printBoardSum(b board) {
	for i := 0; i < 7; i++ {
		fmt.Fprintf(oW, "length of column %v is %v\n", i, len(b.columns[i]))
	}
	for i := 0; i < 4; i++ {
		fmt.Fprintf(oW, "length of pile %v is %v\n", i, len(b.piles[i]))
	}
	fmt.Fprintf(oW, "length of stock is %v\n", len(b.stock))
	fmt.Fprintf(oW, "length of waste is %v\n", len(b.waste))
}
func printBoard(b board) {

	sStock := "stock(" + strconv.Itoa(len(b.stock)) + "):"
	for j := 0; j < len(b.stock); j++ {
		sStock = sStock + b.stock[j].pStr()
	}
	fmt.Fprintf(oW, "\n%v\n", sStock) //print stock

	sWaste := "waste(" + strconv.Itoa(len(b.waste)) + "):"
	for j := 0; j < len(b.waste); j++ {
		sWaste = sWaste + b.waste[j].pStr()
	}
	fmt.Fprintf(oW, "%v\n\n", sWaste) //print waste

	for i := 0; i < 4; i++ {
		s := "pile " + strconv.Itoa(i) + ": "
		for j := 0; j < len(b.piles[i]); j++ {
			s = s + b.piles[i][j].pStr()
		}
		fmt.Fprintf(oW, "%v\n", s) //print the piles

	}

	fmt.Fprintln(oW, "\nColumns:")
	cardWidth := "      " //6 spaces
	spacer := "    "      //4 space
	//determine how many rows must be shown
	numberOfRows := 0
	for cc := 0; cc < 7; cc++ {
		if len(b.columns[cc]) > numberOfRows {
			numberOfRows = len(b.columns[cc])
		}
	}
	fmt.Fprintf(oW, "  0         1         2         3         4         5         6\n")
	for r := 0; r < numberOfRows; r++ {
		s := ""
		for ccc := 0; ccc < 7; ccc++ {
			if r > len(b.columns[ccc])-1 {
				s = s + cardWidth + spacer
			} else {
				s = s + b.columns[ccc][r].pStr() + spacer
			}
		}
		fmt.Fprintf(oW, "%v\n", s)
	}

}
