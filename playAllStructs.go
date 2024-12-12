package main

import "time"

type Statistics struct {
	mvsTried         int
	stratNum         int // stratNum starts at 0             Both are not really needed but it was written that way and not worth changing
	stratTried       int // Strategies TRIED = stratNum + 1  Both are not really needed but it was written that way and not worth changing
	stratWins        int
	stratLosses      int
	stratLossesGLE   int // Strategy Game Length Exceeded
	stratLossesGLEAb int // Strategy Game Length Exceeded Aborted moves
	stratLossesNMA   int // Strategy No Moves Available
	stratLossesRB    int // Strategy Repetitive Board
	stratLossesMajSE int // Strategy Exhausted Major
	stratLossesMinSE int // Strategy Exhausted Minor
	stratLossesEL    int // Strategy Early Loss
	winningMovesCnt  int // = length of winning moves
	unqBoards        int
	elapsedTime      time.Duration
}

type TDotherSQL struct { // Variables in addition to TD that should be output to SQL
	moveNumMax   int
	moveNumAtWin int
}

type variablesSpecificToPlayAll struct {
	TD         Statistics
	TDotherSQL TDotherSQL
	TDother    struct { // Variables NOT common to TD and AD
		startTime     time.Time
		treePrevMoves string // Used to retain values between calls to prntMDetTree for a single deck - Needed for when the strategy "Backs Uo"
		winningMoves  []move
		priorBoards   map[bCode]bool // NOTE: bcode is an array of 65 ints as defined in board.go
	}
	AD      Statistics
	ADother struct {
		startTime       time.Time
		moveNumMax      int
		moveNumAtWinMin int
		moveNumAtWinMax int
		decksPlayed     int
		decksWon        int
		decksLost       int
		decksLostGLE    int
	}
}
