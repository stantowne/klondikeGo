package main

import (
	"strconv"
)

type Card struct {
	Rank   int  `json:"rank"` // 1 ace; 11 jack, etc.
	Suit   int  `json:"suit"` // 0 clubs; 1 diamonds; 2 spades; 3 hearts
	FaceUp bool `json:"faceUp"`
}

func (c *Card) flipCardUp() {
	(*c).FaceUp = true
}

func (c *Card) flipCardUp2() Card {
	return Card{c.Rank, c.Suit, true}
}

func (c Card) packCard() byte {
	// We only care about the rightmost 4 bits of Rank and, the rightmost 2 bits of Suit and, the 1 bit of FaceUp
	// since by definition Rank is never greater than 13 and suit is never greater than 4
	//
	// Rank shift bits 3-0 to 6-3
	r := c.Rank << 3
	// Suit shift bits 1-0 to 2-1
	s := c.Suit << 1
	// FaceUp 's one bit stays at position 0 - convert it to an integer
	fU := 0
	if c.FaceUp {
		fU = 1
	}

	// Use bitwise OR to overlay s (Suit) and FaceUp on the shifted r (which is an integer)
	r = r | s | fU

	// return the rightmost byte of the integer
	return byte(r)
}

func unPackByte2Card(y byte) Card {

	// See method packCard() for structure of the byte
	// & is binary AND operator, returns for any bit, & results in 1 if that bit in both operands is 1
	r := int(y&0b_01111000) >> 3
	s := int(y&0b_00000110) >> 1
	f := int(y&0b_00000001) >> 0

	fU := true
	if f != 1 {
		fU = false
	}

	return Card{
		Rank:   r,
		Suit:   s,
		FaceUp: fU,
	}
}

func (c *Card) color() string {
	var clr string
	switch c.Suit {
	case 0, 2:
		clr = "Black"
	case 1, 3:
		clr = "Red"
	}
	return clr
}

func (c *Card) suitSymbol() rune {
	var symbol rune
	switch c.Suit {
	case 0: //clubs
		symbol = '\u2663'
	case 1: //diamonds
		symbol = '\u2666'
	case 2: //spades
		symbol = '\u2660'
	case 3: //hearts
		symbol = '\u2665'

		return symbol
	}
	return symbol
}

func (c *Card) rankSymbol() string {
	var symbol string
	switch {
	case c.Rank < 10:
		symbol = "0" + strconv.Itoa(c.Rank)
	default:
		symbol = strconv.Itoa(c.Rank)
	}
	return symbol
}

func (c *Card) faceSymbol() string {
	var symbol string
	switch c.FaceUp {
	case true:
		symbol = "UP"
	case false:
		symbol = "DN"
	}
	return symbol
}

func (c *Card) pStr() string {
	var sSuit string
	var Reset = "\033[m" //These are ANSI escape codes for colors
	var Red = "\033[31m"
	var Green = "\033[32m"
	if c.Suit == 0 || c.Suit == 2 {
		sSuit = string(c.suitSymbol())
	} else {
		sSuit = Red + string(c.suitSymbol()) + Reset
	}

	var sFace string
	if c.FaceUp {
		sFace = Green + c.faceSymbol() + Reset
	} else {
		sFace = c.faceSymbol()
	}
	return c.rankSymbol() + sSuit + sFace + " "
}

func (c *Card) suitSymbolColored() string {
	var sSuit string
	var Reset = "\033[m" //These are ANSI escape codes for colors
	var Red = "\033[31m"
	if c.Suit == 0 || c.Suit == 2 {
		sSuit = string(c.suitSymbol())
	} else {
		sSuit = Red + string(c.suitSymbol()) + Reset
	}

	return sSuit
}
