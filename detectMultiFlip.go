package main

func detectMultiFlip(bIn board, moveCounter int, singleGame bool, logicVersion string) []move {
	var aMoves []move //available Moves
	fullLengthStockAndWaste := len(bIn.stock) + len(bIn.waste)
	if len(bIn.stock)+len(bIn.waste) > 0 {
		bMultiFlip := bIn.copyBoard() // make copy so that bIn is never changed - all references should be to bMultiFlip

		/*		// test only
				for k := 1; k <= 24; k++ {
					l := len(bMultiFlip.stock)
					bMultiFlip.waste = append(bMultiFlip.waste, bMultiFlip.stock[l-1].flipCardUp2())
					bMultiFlip.stock = bMultiFlip.stock[:l-1]
				}
		*/
		startingLenOfStock := len(bMultiFlip.stock) // index is within combined stockAndWaste i.e. the bIn "available to move card"

		// We will be using moveMaker with a move of "flipWasteToStock" to do the work of getting a starting
		//     state where all stock and waste cards are in stock
		// create the move
		var flipWasteToStockMove move
		flipWasteToStockMove.name = "flipWasteToStock"
		bMultiFlip = moveMaker(bMultiFlip, flipWasteToStockMove)
		for i := fullLengthStockAndWaste - 1; i >= 0; i-- {
			if (i%3 == 2) || // 3rd, 6th, 9th etc
				(i == fullLengthStockAndWaste-1 || // the last card in stock + waste
					(startingLenOfStock != 0 &&
						(i-startingLenOfStock+1) <= fullLengthStockAndWaste) &&
						(i-startingLenOfStock+1)%3 == 0) {
				//fmt.Printf("\n\n i: %v   startingIndxOfTopWaste: %v", i, startingIndxOfTopWaste)
				aMoves = append(aMoves /* TODO Adjust Priority */, detectDownMoves(bMultiFlip, moveCounter)...)
				// TODO ADJUST MovePortionStartIdx
				//fmt.Printf("\n down: %v", aMoves)
				aMoves = append(aMoves /* Adjust Priority */, detectAcrossMoves(bMultiFlip, moveCounter)...)
				//fmt.Printf("\n down and across: %v", aMoves)
			}
			// flip one card
			if i != 0 {
				l := len(bMultiFlip.stock)
				bMultiFlip.waste = append(bMultiFlip.waste, bMultiFlip.stock[l-1].flipCardUp2())
				bMultiFlip.stock = bMultiFlip.stock[:l-1]
			}
		}

	}
	return aMoves
	/*
		bOut.waste = make([]Card, len(bIn.waste))
		copy(bOut.waste, bIn.waste)

		// Now copy this starting position for future use in variable fullStockAndWaste
		fullStockAndWaste := make([]Card, len(bMultiFlip.stock)+len(bMultiFlip.waste))
		copy(fullStockAndWaste, bMultiFlip.stock)

	*/
}
