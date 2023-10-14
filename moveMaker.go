package main

import (
	"fmt"
	"os"
)

func moveMaker(b board, m move) board {
	if m.name == "moveAceAcross" || m.name == "moveDeuceAcross" || m.name == "move3PlusAcross" {
		_, residue, err := last(b.waste)
		if err != nil {
			fmt.Printf("from Move...Across: %v/n", err)
			os.Exit(1)
		}
		b.piles[m.toPile] = append(b.piles[m.toPile], m.cardToMove)
		b.waste = residue
		return b
	}
	if m.name == "moveDown" {
		_, residue, err := last(b.waste)
		if err != nil {
			fmt.Printf("from MoveDown: %v/n", err)
			os.Exit(1)
		}
		b.columns[m.toCol] = append(b.columns[m.toCol], m.cardToMove)
		b.waste = residue
		return b
	}
	if m.name == "moveAceUp" || m.name == "moveDeuceUp" || m.name == "move3PlusUp" {
		_, residue, err := last(b.columns[m.fromCol])
		if err != nil {
			fmt.Printf("from Move...Up: %v/n", err)
			os.Exit(1)
		}
		b.piles[m.toPile] = append(b.piles[m.toPile], m.cardToMove)
		b.columns[m.fromCol] = residue
		l := len(residue)
		if l > 0 {
			b.columns[m.fromCol][l-1].flipCardUp()
			return b
		}
	}
	if m.name == "moveEntireColumn" {
		b.columns[m.toCol] = append(b.columns[m.toCol], m.MovePortion...)
		b.columns[m.fromCol] = b.columns[m.fromCol][:m.MovePortionStartIdx]
		l := len(b.columns[m.fromCol])
		if l > 0 {
			b.columns[m.fromCol][l-1].flipCardUp()
		}
		return b
	}
	if m.name == "movePartialColumn" {
		b.columns[m.toCol] = append(b.columns[m.toCol], m.MovePortion...)
		b.piles[m.toPile] = append(b.piles[m.toPile], b.columns[m.fromCol][m.MovePortionStartIdx-1])
		b.columns[m.fromCol] = b.columns[m.fromCol][:m.MovePortionStartIdx-1]
		l := len(b.columns[m.fromCol])
		if l > 0 {
			b.columns[m.fromCol][l-1].flipCardUp()
		}
		return b
	}
	if m.name == "flipStockToWaste" {
		l := len(b.stock)
		if l > 2 {
			b.waste = append(b.waste, b.stock[l-1].flipCardUp2(), b.stock[l-2].flipCardUp2(), b.stock[l-3].flipCardUp2())
			b.stock = b.stock[:l-3]
			return b
		}
		if len(b.stock) == 2 {
			b.waste = append(b.waste, b.stock[l-1].flipCardUp2(), b.stock[l-2].flipCardUp2())
			b.stock = b.stock[:l-2]
			return b
		}
		if len(b.stock) == 1 {
			b.waste = append(b.waste, b.stock[l-1].flipCardUp2())
			b.stock = b.stock[:l-1]
			return b
		}
		fmt.Printf("Error: attempted to flip from empty stock")
		os.Exit(1)
	}
	if m.name == "flipWasteToStock" {
		var s card
		for i := len(b.waste) - 1; i > -1; i-- {
			s = card{b.waste[i].Rank, b.waste[i].Suit, false}
			b.stock = append(b.stock, s)
		}
		w := make([]card, 0, 24)
		b.waste = w
	}
	return b
}
