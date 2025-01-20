package main

func detectMultiFlip(bIn board, moveCounter int /*, singleGame bool*/, logicVersion string) []move {
	var aMoves []move //available Moves
	var bMultiFlip board
	fullLengthStockAndWaste := len(bIn.stock) + len(bIn.waste)
	if len(bIn.stock)+len(bIn.waste) > 0 {
		bMultiFlip = bIn.copyBoard() // make copy so that bIn is never changed - all references should be to bMultiFlip

		/*		// test only
				for k := 1; k <= 24; k++ {
					l := len(bMultiFlip.stock)
					bMultiFlip.waste = append(bMultiFlip.waste, bMultiFlip.stock[l-1].flipCardUp2())
					bMultiFlip.stock = bMultiFlip.stock[:l-1]
				}
		*/
		startingLenOfWaste := len(bMultiFlip.waste) // index is within combined stockAndWaste i.e. the bIn "available to move card"
		modOfStartingTopWasteCard := (startingLenOfWaste - 1) % 3
		// We will be using moveMaker with a move of "flipWasteToStock" to do the work of getting a starting
		//     state where all stock and waste cards are in stock
		// create the move
		var flipWasteToStockMove move
		flipWasteToStockMove.name = "flipWasteToStock"
		bMultiFlip = moveMaker(bMultiFlip, flipWasteToStockMove)
		for topWasteCard := 0; topWasteCard <= fullLengthStockAndWaste-1; topWasteCard++ {
			// flip one card
			l := len(bMultiFlip.stock)
			bMultiFlip.waste = append(bMultiFlip.waste, bMultiFlip.stock[l-1].flipCardUp2())
			bMultiFlip.stock = bMultiFlip.stock[:l-1]
			if (topWasteCard%3 == 2) || // 3rd, 6th, 9th etc flip all the way around and back again = "*" in ModAnalysis.txt
				(topWasteCard == fullLengthStockAndWaste-1) || // the last card in stock + waste    = "L" in ModAnalysis.txt
				(topWasteCard >= startingLenOfWaste-1 &&
					topWasteCard%3 == modOfStartingTopWasteCard) { // Normal Flips of 3 (or fewer)  = "3" in ModAnalysis.txt
				// NOTE: detectDownMoves will only ever return 1 move
				downMove := detectDownMoves(bMultiFlip, moveCounter)
				if downMove != nil {
					downMove[0].name = "flpMMveDown"
					switch logicVersion {
					case "multiflip":
						downMove[0].priority = moveBasePriority["flpMMveDown"] + 50*(13-downMove[0].cardToMove.Rank) + (topWasteCard)
					case "multifliptest":
						downMove[0].priority = moveBasePriority["flpMMveDown"] + 50*(topWasteCard) + (13 - downMove[0].cardToMove.Rank)
					}
					downMove[0].MovePortionStartIdx = topWasteCard
					aMoves = append(aMoves, downMove...)
				}
				// NOTE: detectAcrossMoves will only ever return 1 move
				acrossMove := detectAcrossMoves(bMultiFlip, moveCounter)
				if acrossMove != nil {
					switch logicVersion {
					case "multiflip":
						switch acrossMove[0].cardToMove.Rank {
						case 1: // Ace
							acrossMove[0].name = "flpMMveAceAcross"
							acrossMove[0].priority = moveBasePriority["flpMMveAceAcross"] + 50*(13-acrossMove[0].cardToMove.Rank) + (topWasteCard)
						case 2: // Deuce
							acrossMove[0].name = "flpMMve2Across"
							acrossMove[0].priority = moveBasePriority["flpMMve2Across"] + 50*(13-acrossMove[0].cardToMove.Rank) + (topWasteCard)
						default: // 3+
							acrossMove[0].name = "flpMMve3UpAcross"
							acrossMove[0].priority = moveBasePriority["flpMMve3UpAcross"] + 50*(13-acrossMove[0].cardToMove.Rank) + (topWasteCard)
						}
					case "multifliptest":
						switch acrossMove[0].cardToMove.Rank {
						case 1: // Ace
							acrossMove[0].name = "flpMMveAceAcross"
							acrossMove[0].priority = moveBasePriority["flpMMveAceAcross"] + 50*(topWasteCard) + (13 - acrossMove[0].cardToMove.Rank)
						case 2: // Deuce
							acrossMove[0].name = "flpMMve2Across"
							acrossMove[0].priority = moveBasePriority["flpMMve2Across"] + 50*(topWasteCard) + (13 - acrossMove[0].cardToMove.Rank)
						default: // 3+
							acrossMove[0].name = "flpMMve3UpAcross"
							acrossMove[0].priority = moveBasePriority["flpMMve3UpAcross"] + 50*(topWasteCard) + (13 - acrossMove[0].cardToMove.Rank)
						}
					}
					acrossMove[0].MovePortionStartIdx = topWasteCard
					aMoves = append(aMoves, acrossMove...)
				}
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
