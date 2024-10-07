package main

type move struct {
	name                string
	priority            int //lower is better
	toPile              int
	toCol               int
	fromCol             int
	MovePortionStartIdx int
	MovePortion         []Card //used in mec, mpc,
	cardToMove          Card   //used in Up, Down, Across and mpc
	colCardFlip         bool   //does the move result in a column Card flip

}
