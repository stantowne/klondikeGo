package main

func detectMultiFlip(bIn board, moveCounter int, singleGame bool, logicVersion string) []move {
	var aMoves []move //available Moves
	var bMultiFlip board
	fullLengthStockAndWaste := len(bIn.stock) + len(bIn.waste)
	if len(bIn.stock)+len(bIn.waste) > 0 {
		bMultiFlip = bIn.copyBoard() // make copy so that bIn is never changed - all references should be to bMultiFlip

		/*		// test only delete when sure correct used to force a starting len of wast by changing the k = 24
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
			if (topWasteCard%3 == 2) || // 3rd, 6th, 9th etc. flip all the way around and back again = "*" in ModAnalysis.txt
				(topWasteCard == fullLengthStockAndWaste-1) || // the last card in stock + waste    = "L" in ModAnalysis.txt
				(topWasteCard >= startingLenOfWaste-1 &&
					topWasteCard%3 == modOfStartingTopWasteCard) { // Normal Flips of 3 (or fewer)  = "3" in ModAnalysis.txt
				aMoves = append(aMoves, adjustMoves(logicVersion, topWasteCard, detectAcrossMoves(bMultiFlip, moveCounter))...)
				aMoves = append(aMoves, adjustMoves(logicVersion, topWasteCard, detectDownMoves(bMultiFlip, moveCounter))...)
				/*				// NOTE: detectDownMoves will only ever return 1 move
									downMove := detectDownMoves(bMultiFlip, moveCounter)
									if downMove != nil {
										downMove[0].name = "mFlMDown"
										switch logicVersion {
										case "multiflip1":
											downMove[0].priority = moveBasePriority["mFlMDown"] + 50*(13-downMove[0].cardToMove.Rank) + (topWasteCard)
										case "multiflip2":
											//						downMove[0].priority = moveBasePriority["mFlMDown"] + 50*(topWasteCard) + (13 - downMove[0].cardToMove.Rank)
											downMove[0].priority = moveBasePriority["mFlMDown"] + 50*(topWasteCard) + (13 - downMove[0].cardToMove.Rank)
										}
										downMove[0].multiFlipCardsToFlip = topWasteCard
										aMoves = append(aMoves, downMove...)
									}
								// NOTE: detectAcrossMoves will only ever return 1 move
									acrossMove := detectAcrossMoves(bMultiFlip, moveCounter)
									if acrossMove != nil {
										switch logicVersion {
										case "multiflip1":
											switch acrossMove[0].cardToMove.Rank {
											// TODO IMPLEMENT NEW PRIORITY CALCULATIONS
											case 1: // Ace
												acrossMove[0].name = "mFlMAceAcross"
												acrossMove[0].priority = moveBasePriority["mFlMAceAcross"] + 50*(13-acrossMove[0].cardToMove.Rank) + (topWasteCard)
											case 2: // Deuce
												acrossMove[0].name = "mFlMDeuceAcross"
												acrossMove[0].priority = moveBasePriority["mFlMDeuceAcross"] + 50*(13-acrossMove[0].cardToMove.Rank) + (topWasteCard)
											default: // 3+
												acrossMove[0].name = "mFlM3PlusAcross"
												acrossMove[0].priority = moveBasePriority["mFlM3PlusAcross"] + 50*(13-acrossMove[0].cardToMove.Rank) + (topWasteCard)
											}
										case "multiflip2":
											switch acrossMove[0].cardToMove.Rank {
											case 1: // Ace
												acrossMove[0].name = "mFlMAceAcross"
												acrossMove[0].priority = moveBasePriority["mFlMAceAcross"] + 50*(topWasteCard) + (13 - acrossMove[0].cardToMove.Rank)
											case 2: // Deuce
												acrossMove[0].name = "mFlMDeuceAcross"
												acrossMove[0].priority = moveBasePriority["mFlMDeuceAcross"] + 50*(topWasteCard) + (13 - acrossMove[0].cardToMove.Rank)
											default: // 3+
												acrossMove[0].name = "mFlM3PlusAcross"
												acrossMove[0].priority = moveBasePriority["mFlM3PlusAcross"] + 50*(topWasteCard) + (13 - acrossMove[0].cardToMove.Rank)
											}
										}
										acrossMove[0].multiFlipCardsToFlip = topWasteCard
										aMoves = append(aMoves, acrossMove...)
									}
				*/
			}
		}

	}
	return aMoves
}

func adjustMoves(logicVersion string, topWasteCard int, moves []move) []move {
	for i := range moves {
		moves[i].multiFlipCardsToFlip = topWasteCard
		if moves[i].name != "badMove" {
			moves[i].name = "mFlM" + moves[i].name[4:]
		}
		switch logicVersion {
		case "multiflip1":
			moves[i].priority = moves[i].priority*moveBasePriority["origPrityMultiple"] +
				moveBasePriority[moves[i].name] +
				(moveBasePriority["posStkWaBase"]+moveBasePriority["posStkWaSign"]*topWasteCard)*moveBasePriority["posStkWaMultiple"]
			switch moves[i].name {
			case "mFlMAceAcross", "mFlMDeuceAcross", "mFlM3PlusAcross":
				moves[i].priority += (moves[i].cardToMove.Rank - 1) * moveBasePriority["rankMultiple"] // Lowest rank first
			case "mFlMDown":
				moves[i].priority += (12 - moves[i].cardToMove.Rank) * moveBasePriority["rankMultiple"] // Highest rank first
			}
		case "multiflip2":
		case "multiflip3":
		}
	}
	return moves
}
