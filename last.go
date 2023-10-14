package main

import "errors"

func last(sl []card) (lastCard card, residue []card, err error) {
	if len(sl) > 0 {
		lastIndex := len(sl) - 1
		lastCard := sl[lastIndex]
		residue := sl[:lastIndex]
		return lastCard, residue, nil
	}

	var x card

	if len(sl) == 0 {
		return x, nil, nil
	}
	return x, nil, errors.New("error in call to last")
}
