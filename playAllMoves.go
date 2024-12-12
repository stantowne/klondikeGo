package main

import (
	"math"
	"sort"
	"strconv"
	"time"
)

func playAllMoves(bIn board, moveNum int, deckNum int, cfg *Configuration, vPA *variablesSpecificToPlayAll) (string, int) {

	/* Return Codes: NMA = Strategy Loss: No Moves Available

	   A "Strategy" is any unique set of moves on a path leading back to move 0 aka a "Branch"
	   Note: When available moves = 1 a new "Strategy"/"Branch" is NOT created
	         When available moves > 1 a new "Strategy"/"Branch" IS     created

	   RB     = Strategy Loss: Repetitive Board
	   MajSE  = Major Strategy Loss: all possible moves Exhausted
	   MinSE  = Minor Strategy Loss: all possible moves Exhausted Not Really Considered a strategy!
	   GLE    = Strategy Loss: GameLength Limit Exceeded
	   EL     = Strategy Loss: Early Loss
	   SW     = Strategy Win

	*/
	// add code for findAllWinStrats

	var aMoves []move //available Moves
	var recurReturnV1 string
	var recurReturnNum int

	if moveNum > vPA.TDotherSQL.moveNumMax {
		vPA.TDotherSQL.moveNumMax = moveNum
	}

	// Check to see if the gameLengthLimit has been exceeded.
	//If, treats this as a loss and returns with loss codes.
	if vPA.TD.mvsTried >= cfg.PlayAll.GameLengthLimit*1_000_000 {
		prntMDet(bIn, aMoves, 0, deckNum, moveNum, "MbM_ANY", 2, "GLE: Game Length of: %v exceeds limit: %v\n", strconv.Itoa(vPA.TD.mvsTried), strconv.Itoa(cfg.PlayAll.GameLengthLimit*1_000_000), cfg, vPA)
		vPA.TD.stratLossesGLE++
		prntMDetTreeReturnComment(" ==> GLE", deckNum, recurReturnNum, cfg, vPA)
		return "GLE", recurReturnNum + 1
	}

	// Find Next Moves
	aMoves = detectAvailableMoves(bIn, moveNum, cfg.General.NumberOfDecksToBePlayed == 1)

	if len(aMoves) == 0 {
		m := move{name: "No Moves Available"} // This is a pseudo move not created by detectAvailable Moves it exists to remember
		aMoves = append(aMoves, m)            //      this state and for the various printing routines that come below.  If a return
		//      was made from this point no history of it would exist.  See code at "// Check if No Moves Available"
	} else {
		// if more than one move is available, sort them
		if len(aMoves) > 1 { //sort them by priority if necessary
			sort.SliceStable(aMoves, func(i, j int) bool {
				return aMoves[i].priority < aMoves[j].priority
			})
		}
	}

	// Try all moves
	for i := range aMoves {

		/* Actually before we actually try all moves let's first: print (optionally based on cfg.PlayAll.RestrictReportingTo.pType) the incoming board
		      and check the incoming board for various end-of-strategy conditions
		   Note: This was done this way, so as to ensure that when returns backed up the moveNum, the board would reprint.
		*/

		// Possible increment to vPA.TD.stratNum
		if i != 0 {
			// Started at 0 in playAll for each deck.  Increment each time a second through nth move is made
			//     That way strategy 0 is the "All Best Moves" strategy.
			vPA.TD.stratNum++
			vPA.TD.stratTried++
		}

		// Print the incoming board EVEN IF we are returning to it to try the next available move
		//       This had to be done after possible increment to vPA.TD.stratNum so that each time a board is reprinted it shows the NEW strategy number
		//       Before when it was above the possible increment the board was printing out with the stratNum of the last failed strategy
		prntMDet(bIn, aMoves, i, deckNum, moveNum, "MbM_ANY", 1, "", "", "", cfg, vPA)

		if i == 0 {
			// Check for repetitive board
			// We only have to check if i == 0 because if i >= 0 then we are just returning to this board at which point
			//   the board will already have been checked (when i == 0 {past tense})
			bNewBcode := bIn.boardCode(deckNum) //  Just use the array itself
			// Have we seen this board before?
			if vPA.TDother.priorBoards[bNewBcode] {
				// OK we did see it before so return to try next available move (if any) in aMoves[] aka strategy
				vPA.TD.stratLossesRB++
				prntMDet(bIn, aMoves, i, deckNum, moveNum, "MbM_ANY", 2, "RB: Repetitive Board - \"Next Move\" yielded a repeat of a board.", "", "", cfg, vPA)
				prntMDetTreeReturnComment(" ==> RB", deckNum, recurReturnNum, cfg, vPA)
				return "RB", recurReturnNum + 1 // Repetitive Board
			} else {
				// Remember the board state by putting it into the map "vPA.TDother.priorBoards"
				vPA.TDother.priorBoards[bNewBcode] = true
			}
		}

		// Check if No Moves Available
		if i == 0 && aMoves[0].name == "No Moves Available" {
			vPA.TD.stratLossesNMA++
			prntMDet(bIn, aMoves, i, deckNum, moveNum, "MbM_ANY", 2, "NMA: No Moves Available: Strategy Lost %v%v\n", "", "", cfg, vPA)
			prntMDetTreeReturnComment(" ==> NMA", deckNum, recurReturnNum, cfg, vPA)
			return "NMA", recurReturnNum + 1
		}

		//Detect Win (formerly Early Win)
		if detectWinEarly(bIn) {
			vPA.TD.stratWins++

			prntMDetTreeReturnComment(" ==> DECK WON", deckNum, recurReturnNum, cfg, vPA)
			vPA.TDotherSQL.moveNumAtWin = moveNum
			return "SW", recurReturnNum + 1 //  Strategy Win
		}

		// OK, done with the various end-of-strategy conditions
		// let's print out the list of available moves and make the next available move
		// The board state was already printed above
		prntMDet(bIn, aMoves, i, deckNum, moveNum, "MbM_R", 3, "", "", "", cfg, vPA) // Print Available Moves
		prntMDetTree(bIn, aMoves, i, deckNum, moveNum, cfg, vPA)

		bNew := bIn.copyBoard() // Critical Must use copyBoard!!!

		// ********** 1st of the 2 MOST IMPORTANT statements in this function:  ******************************
		bNew = moveMaker(bNew, aMoves[i])

		vPA.TD.mvsTried++

		prntMDet(bIn, aMoves, i, deckNum, moveNum, "MbM_S", 2, "\n      bIn: %v\n", "", "", cfg, vPA)
		prntMDet(bNew, aMoves, i, deckNum, moveNum, "MbM_S", 2, "     bNew: %v", "", "", cfg, vPA)

		if cfg.PlayAll.ProgressCounter > 0 && math.Mod(float64(vPA.TD.mvsTried+vPA.AD.mvsTried), float64(cfg.PlayAll.ProgressCounter)) <= 0.1 {
			decksPlayed := float64(deckNum-cfg.General.FirstDeckNum) + .5
			decksCompleted := float64(vPA.ADother.decksWon + vPA.ADother.decksLost + vPA.ADother.decksLostGLE)
			estRemTimeAD := time.Duration(float64(time.Since(vPA.ADother.startTime)) * (float64(cfg.General.NumberOfDecksToBePlayed) - decksPlayed) / decksPlayed)
			// NOTE THE FOLLOWING PRINT STATEMENT NEVER GOES TO FILE - ALWAYS TO CONSOLE
			_, _ = pfmt.Printf("\rDk: %d  Mvs: %vmm  Strats: %vmm  UnqBoards: %vmm  MaxMoveNum: %v  Elapsed: %s  estRem: %s  W/L/GLE: %v/%v/%v  W/L/GLE %%: %3.1f/%3.1f/%3.1f\r",
				deckNum, (vPA.TD.mvsTried+vPA.AD.mvsTried)/1000000, vPA.AD.stratNum/1000000, vPA.AD.unqBoards/1000000, vPA.ADother.moveNumMax, time.Since(vPA.ADother.startTime).Round(6*time.Second).String(), estRemTimeAD.Round(6*time.Second), vPA.ADother.decksWon, vPA.ADother.decksLost, vPA.ADother.decksLostGLE, float64(vPA.ADother.decksWon)/decksCompleted*100.0, float64(vPA.ADother.decksLost)/decksCompleted*100.0, float64(vPA.ADother.decksLostGLE)/decksCompleted*100.0)
		}

		// ********** 2nd of the 2 MOST IMPORTANT statements in this function:  ******************************
		recurReturnV1, recurReturnNum = playAllMoves(bNew, moveNum+1, deckNum, cfg, vPA)

		if recurReturnV1 == "SW" {
			// save winning moves into a slice in reverse
			vPA.TDother.winningMoves = append(vPA.TDother.winningMoves, aMoves[i])
			return recurReturnV1, recurReturnNum + 1 // return up the call stack to end strategies search  if findAllWinStrats false, and we had a win
		}
		if recurReturnV1 == "GLE" {
			return recurReturnV1, recurReturnNum + 1 //
		}
	}

	if moveNum != 0 { // The initial call to playAllMoves from playALL (when MoveNum = 0) is only there to display the initial board
		if len(aMoves) == 1 { // and does not represent a move tried - Therefore returning from the initial call is not either a Maj or MinSE
			vPA.TD.stratLossesMinSE++
		} else {
			vPA.TD.stratLossesMajSE++
		}
	}
	cmt := "  MajSE   Major Strategy loss: all possible moves Exhausted%v%v"
	if recurReturnNum == 0 {
		prntMDet(bIn, aMoves, len(aMoves), deckNum, moveNum, "DbDorMbM", 2, cmt, "", "", cfg, vPA)
	}
	prntMDetTreeReturnComment("^Maj SE^", deckNum, recurReturnNum, cfg, vPA)
	return "SE", recurReturnNum + 1 //  Strategy Exhausted
}
