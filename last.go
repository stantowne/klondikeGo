package main

import "errors"

func last(sl []Card) (lastCard Card, residue []Card, err error) {
	if len(sl) > 0 {
		lastIndex := len(sl) - 1
		lastCard := sl[lastIndex]
		residue := sl[:lastIndex]
		return lastCard, residue, nil
	}

	var x Card

	if len(sl) == 0 {
		return x, nil, nil
	}
	return x, nil, errors.New("error in call to last")
}
