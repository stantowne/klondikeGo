package main

import (
	"fmt"
	"strconv"
)

type card struct {
	Rank   int  `json:"rank"` // 1 ace; 11 jack, etc.
	Suit   int  `json:"suit"` // 0 clubs; 1 diamonds; 2 spades; 3 hearts
	FaceUp bool `json:"faceUp"`
}

func (c *card) flipCardUp() {
	(*c).FaceUp = true
}

func (c *card) flipCardUp2() card {
	return card{c.Rank, c.Suit, true}
}

func (c card) packCard() byte {
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

	/*var fU int
	switch c.FaceUp {
	case true:
		fU = 1
	default:
		fU = 0
	}*/

	// Use bitwise OR to overlay s (Suit) and FaceUp on the shifted r (which is an integer)
	r = r | s | fU

	// return the rightmost byte of the integer
	return byte(r)
}

func unPackByte2Card(y byte) card {

	// See method packCard() for structure of the byte
	// & is binary AND operator, returns for any bit, & results in 1 if that bit in both operands is 1
	r := int(y&0b_01111000) >> 3
	s := int(y&0b_00000110) >> 1
	f := int(y&0b_00000001) >> 0

	fU := true
	if f != 1 {
		fU = false
	}

	return card{
		Rank:   r,
		Suit:   s,
		FaceUp: fU,
	}
}

func (c *card) color() string {
	var clr string
	switch c.Suit {
	case 0, 2:
		clr = "Black"
	case 1, 3:
		clr = "Red"
	}
	return clr
}

func (c *card) suitSymbol() rune {
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
	}
	return symbol
}

func (c *card) rankSymbol() string {
	var symbol string
	switch {
	case c.Rank < 10:
		symbol = "0" + strconv.Itoa(c.Rank)
	default:
		symbol = strconv.Itoa(c.Rank)
	}
	return symbol
}

func (c *card) faceSymbol() string {
	var symbol string
	switch c.FaceUp {
	case true:
		symbol = "UP"
	case false:
		symbol = "DN"
	}
	return symbol
}

func (c *card) pStr() string {
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

func test(c card) bool {
	return unPackByte2Card(c.packCard()) == c
}

func testCardPackUnPack(args []string) {
	var testCard, rebuiltCard card
	r, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("bad input -- first argument")
	}
	s, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("bad input -- second argument")
	}
	testCard.Rank = r
	testCard.Suit = s
	testCard.FaceUp = true
	fmt.Printf("TestCard: %+v\n", testCard)
	packed := testCard.packCard()
	fmt.Printf("PackCard: %08b\n", packed)
	rebuiltCard = unPackByte2Card(packed)
	fmt.Printf("RebuiltCard: %+v\n", rebuiltCard)
	fmt.Printf("Does the round trip work: %v\n", testCard == rebuiltCard)
}
