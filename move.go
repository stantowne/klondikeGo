package main

type move struct {
	name                string
	priority            int //lower is better
	toPile              int
	toCol               int
	fromCol             int
	MovePortionStartIdx int
	MovePortion         []card //used in mec, mpc,
	cardToMove          card   //used in Up, Down, Across and mpc

}
