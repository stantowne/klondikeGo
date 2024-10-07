package main

import "errors"

// returns the index of the first FaceUpCard and the entire FaceUp portion
func faceUpPortion(col column) (int, []Card, error) {
	for i, crd := range col {
		if crd.FaceUp {
			firstFaceUpIndex := i
			faceUpPortion := col[i:]
			return firstFaceUpIndex, faceUpPortion, nil
		}
	}
	if len(col) == 0 {
		return -1, nil, nil
	}
	return -1, nil, errors.New("error in call to faceUpPortion")
}
