package main

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

type column []Card
type pile []Card

type board struct {
	columns [7]column //column[x][0] rests on the table for all x
	piles   [4]pile   //  pile[y][0] rests on the table for all y
	stock   []Card    //    stock[0] rests on the table
	waste   []Card    //    waste[0] rests on the table
}
type (
	bCode [65]byte // 1 stock + 1 waste + 4 piles + 7 columns + 52 cards = 65
	// note array index of course goes from 0 to 64
)

func (b board) boardCode() bCode {
	//
	// This method takes a struct of type board which contains four fields:
	//		columns:	an array of 7 column (each column being a slice of Card)
	//		piles:		an array of 4 pile (each pile being a slice of Card)
	//		stock:		a slice of Card, and
	//		waste:		a slice of Card.
	//
	//	Thus, there are a total of 13 slices (7+4+1+1) holding the 52 cards of a standard deck.
	//
	//  Note that the cards in board are originally stored in a Card struct using two integers representing the rank and suit
	//   	and a bool for faceUp (true is Up, false is Down) (See card.go).
	//
	//	This method uses packCard (a method of Card see card.go) to pack each Card's rank, suit and FaceUp
	//	   into the 7 rightmost bits of a single byte.
	//
	//	So, in order to completely describe a board we need 52 bytes for the cards, and 13 bytes ("Flags") to mark which each slice begins.
	//
	//  The 13 "Flags" are each one byte long and are as follows:
	//
	//	      0b_10000000 = start of Stock      decimal = 128
	//	      0b_10000001 = start of Waste      decimal = 129
	//	      0b_10000010 = start of Pile 0     decimal = 130      leftmost pile (corresponding to suit 0 which is clubs)
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
	//                3b. The first 0 byte encountered will cause the method boardDeCode to ignore the remainder of the bytes in bCode
	//                3c. Note: No packed Card can ever have a value of 0 as the ranks start with ace = 1
	//
	//         4. There is no set order to the flag/groups.  NOTE TO DAN:  Is this so?
	//
	//         5. The cards in each group appear in order with the first Card after the flag being the Card on the surface of the table.
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

func (bC bCode) boardDeCode() board {
	//
	// See method boardCode (a method of board) for comments which will explain how this method (boardDeCode), which is the inverse of bCodeFromBoard, works
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

func testBoardCodeDeCode(args []string) {

	firstDeckNum, _ := strconv.Atoi(args[1])
	numberOfDecksToBePlayed, _ := strconv.Atoi(args[2])
	verbose, _ := strconv.Atoi(args[3]) //the greater the number the more verbose

	var decks = deckReader("decks-made-2022-01-15_count_10000-dict.json") //contains decks 0-999 from Python version

	for deckNum := firstDeckNum; deckNum < (firstDeckNum + numberOfDecksToBePlayed); deckNum++ {
		if verbose > 1 {
			fmt.Printf("\nDeck #%d:\n", deckNum)
		}
		//TempTest
		var b = dealDeck(decks[deckNum])
		fmt.Println("Original Board")
		printBoard(b)

		var bC = b.boardCode()
		fmt.Printf("%v \t    %08b \n\n", deckNum, bC)

		var bRoundTrip = bC.boardDeCode()
		fmt.Println("RoundTrip Board")
		printBoard(bRoundTrip) //TempTest end
	}
}

func quickTestBoardCodeDeCode(b board, deckNum int, length int, iOS int, mC int) {
	bCode := b.boardCode()
	roundTripResult := bCode.boardDeCode()
	if !reflect.DeepEqual(b, roundTripResult) {
		fmt.Println("quickTestBoardCodeDeCode() failed")
		fmt.Println("deckNum: ", deckNum)
		fmt.Println("length: ", length)
		fmt.Println("iOS: ", iOS)
		fmt.Println("moveCounter: ", mC)
		printBoard(b)
		fmt.Printf("\t   %08b \n\n", bCode)
		printBoard(roundTripResult)
		os.Exit(1)
	}
}
