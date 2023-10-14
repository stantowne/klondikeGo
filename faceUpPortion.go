package main

import "errors"

// returns the index of the first FaceUpCard and the entire FaceUp portion
func faceUpPortion(col column) (int, []card, error) {
	for i, crd := range col {
		if crd.FaceUp {
			firstFaceUpIndex := i
			faceUpPortion := col[i:]
			return firstFaceUpIndex, faceUpPortion, nil
		}
	}
	return -1, nil, errors.New("call to faceUpPortion resulted in no FaceUp card")
}
