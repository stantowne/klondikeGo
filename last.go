package main

func last(sl []card) (lastCard card, residue []card, err error) {
	if len(sl) > 0 {
		lastIndex := len(sl) - 1
		lastCard := sl[lastIndex]
		residue := sl[:lastIndex]
		return lastCard, residue, nil
	}

	var y card
	var z []card //empty slice
	return y, z, nil

}
