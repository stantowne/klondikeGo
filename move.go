package main

type move struct {
	name                 string
	priority             int //lower is better
	toPile               int
	toCol                int
	fromCol              int
	MovePortionStartIdx  int
	MovePortion          []Card //used in mec, mpc,                          // See comments below
	cardToMove           Card   //used in Up, Down, Across and mpc
	colCardFlip          bool   //does the move result in a column Card flip
	multiFlipCardsToFlip int    // used by MultiFlip

	// When winning moves are saved to SQL, the slice MovePortion is NOT written.
	// If and when a routine to read And EXECUTE the winning moves is written the
	// slice MUST be written by reading the then current board from MovePortionStartIdx
	// on down.
}
