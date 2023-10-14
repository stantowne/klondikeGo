package main

type deck []card

func (d deck) firstRest() (card, deck) {
	first := d[0]
	rest := d[1:]
	return first, rest
}
func dealDeck(d deck) board {
	var b board
	for j := 0; j < 7; j++ {
		for i := 0; i+j < 7; i++ {
			first, rest := d.firstRest()
			b.columns[i+j] = append(b.columns[i+j], first)
			d = rest
		}
	}

	for j := 0; j < 7; j++ {
		b.columns[j][j].flipCardUp()
	}

	b.stock = reverseSlice(d) //to conform to initialization in my python klondikeGo

	return b
}
func reverseSlice(input []card) []card {
	var output []card
	for i := len(input) - 1; i > -1; i-- {
		output = append(output, input[i])
	}
	return output
}
