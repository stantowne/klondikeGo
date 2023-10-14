package main

import (
	"github.com/TwiN/go-color"
	"strconv"
)

type card struct {
	Rank   int  `json:"rank"` // 1 ace; 11 jack, etc.
	Suit   int  `json:"suit"` // 0 clubs; 1 diamonds; 2 spades; 3 hearts
	FaceUp bool `json:"faceUp"`
}

type stock []card
type waste []card
type column []card
type pile []card

func (c *card) flipCardUp() {
	(*c).FaceUp = true
}

func (c card) flipCardUp2() card {
	return card{c.Rank, c.Suit, true}
}

func (c card) color() string {
	var clr string
	switch c.Suit {
	case 0, 2:
		clr = "black"
	case 1, 3:
		clr = "red"
	}
	return clr
}

func (c card) suitSymbol() rune {
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

func (c card) rankSymbol() string {
	var symbol string
	switch {
	case c.Rank < 10:
		symbol = "0" + strconv.Itoa(c.Rank)
	default:
		symbol = strconv.Itoa(c.Rank)
	}
	return symbol
}

func (c card) faceSymbol() string {
	var symbol string
	switch c.FaceUp {
	case true:
		symbol = "UP"
	case false:
		symbol = "DN"
	}
	return symbol
}

func (c card) pStr() string {
	var sSuit string
	if c.Suit == 0 || c.Suit == 2 {
		sSuit = string(c.suitSymbol())
	} else {
		sSuit = color.Ize(color.Red, string(c.suitSymbol()))
	}

	var sFace string
	if c.FaceUp {
		sFace = color.Ize(color.Green, c.faceSymbol())
	} else {
		sFace = c.faceSymbol()
	}

	return c.rankSymbol() + sSuit + sFace + " "

}