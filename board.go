package main

type board struct {
	columns [7]column //column[x][0] rests on the table for all x
	piles   [4]pile   //  pile[y][0] rests on the table for all y
	stock   []card    //    stock[0] rests on the table
	waste   []card    //    waste[0] rests on the table
}
type (
	bCode [65]byte // byte 0 code + 1 stock + 1 waste + 4 piles + 7 columns + 52 cards= 66
	// note array index of course goes from 0 to 65 !!!
)

func (b board) bCodeFromBoard() bCode {
	//
	// This method takes a struct of type board which contains four fields:
	//		columns:	an array of 7 column (each column being a slice of card)
	//		piles:		an array of 4 pile (each pile being a slice of card)
	//		stock:		a slice of card, and waste, a slice of card.
	//
	//	Thus there are a total of 13 slices (7+4+1+1).  Holding the 52 cards of a standard deck.
	//
	//  Note that the cards in board are originally stored in the structure card using two integers representing the rank and suit
	//   	and a bool for faceup(1) or down(0)(see struc in card.go).
	//
	//	This method uses packCard (a method of Card see card.go) to pack each cards rank, suit and FaceUp
	//	   into the 7 rightmost bits of a single byte.  Or 52 bytes.
	//
	//	So, in order to completely describe a board we need 52 bytes for the cards, and 13 bytes ("Flags") to mark which each slice begins.
	//
	//  The 13 "Flags" are each one byte long and are as follows:
	//
	//	      0b_10000000 = start of Stock      decimal = 128
	//	      0b_10000001 = start of Waste      decimal = 129
	//	      0b_10000010 = start of Pile 0     decimal = 130      leftmost pile
	//	      0b_10000011 = start of Pile 1     decimal = 131
	//	      0b_10000100 = start of Pile 2     decimal = 132
	//	      0b_10000101 = start of Pile 3     decimal = 133
	//	      0b_10000110 = start of Column 0   decimal = 134      leftmost column
	//	      0b_10000111 = start of Column 1   decimal = 135
	//	      0b_10001000 = start of Column 2   decimal = 136
	//	      0b_10001001 = start of Column 3   decimal = 137
	//	      0b_10001010 = start of Column 4   decimal = 138
	//	      0b_10001011 = start of Column 5   decimal = 139
	//	      0b_10001100 = start of Column 6   decimal = 140
	//
	//
	//
	//  Notes: 1. 52 cards plus 13 flags = 65 bytes
	//
	//         2. Not all flags are required and the number of packed cards following each flag can be 0
	//                However this method and the inverse method ( boardFromBCode ) both work with any or all of the flags present.
	//                Note: The method bCodeFromBoard DOES create all flags even if some groups are empty
	//
	//         3. The number of packed cards following each flag can be 0
	//
	//         3. Bytes not used when/if all flags are NOT present must be 0.
	//
	//                3a. All zero bytes MUST be at the end of the array.
	//                3b. The first 0 byte encountered will cause the method boardFrombCode to ignore the remainder of the bytes in bCode
	//                3c. Note: No packed card can ever have a value of 0 as the ranks start with ace = 1
	//
	//         4. There is no set order to the flag/groups.
	//
	//         5. The cards in each group appear in order with the first card after the flag being the card on the surface of the table.
	//
	//         6. The flags all have a 1 in the leftmost bit so all are >= 127.
	//

	bC := bCode{}
	i := 0

	bC[i] = 0b_10000000
	i++
	for _, c := range b.stock {
		bC[i] = c.packCard()
		i++
	}

	bC[i] = 0b_10000001
	i++
	for _, c := range b.waste {
		bC[i] = c.packCard()
		i++
	}

	var p, col uint8
	for p = 0; p < 4; p++ {
		bC[i] = 0b_10000010 + p
		i++
		for _, c := range b.piles[p] {
			bC[i] = c.packCard()
			i++
		}
	}

	for col = 0; col < 7; col++ {
		bC[i] = 0b_10000110 + col
		i++
		for _, c := range b.columns[col] {
			bC[i] = c.packCard()
			i++
		}
	}

	return bC
}

func (bC bCode) boardFrombCode() board {
	//
	// See method bCodeFromBoard (a method of board) for comments which will explain how this method (boardFrombCode) , which is the inverse of bCodeFromBoard
	//
	b := board{}
	i := 0
	for i <= 64 {
		switch bC[i] {
		case 0:
			panic("or first byte of bCode not >= 128 and <= 140")
		case 128:
			i++
			for i <= 64 && bC[i] < 128 {
				b.stock = append(b.stock, unPackByte2Card(bC[i]))
				i++
			}
		case 129:
			i++
			for i <= 64 && bC[i] < 128 {
				b.waste = append(b.waste, unPackByte2Card(bC[i]))
				i++
			}
		case 130, 131, 132, 133:
			pileNum := bC[i] - 130
			i++
			for i <= 64 && bC[i] < 128 {
				b.piles[pileNum] = append(b.piles[pileNum], unPackByte2Card(bC[i]))
				i++
			}
		case 134, 135, 136, 137, 138, 139, 140:
			colNum := bC[i] - 134
			i++
			//if i <= 64 {
			for i <= 64 && bC[i] < 128 {
				b.columns[colNum] = append(b.columns[colNum], unPackByte2Card(bC[i]))
				i++
			}
		//	}
		default:
			panic("Flag > 140")
		}
	}

	return b
}
