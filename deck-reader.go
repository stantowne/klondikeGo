package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Decks []Deck

func check(e error) {
	if e != nil {
		panic(e)
	}
}

/*
func printRawDeck(d Deck) {
	for i, c := range d {
		fmt.Printf("Index: %v : Card: %+v\n", i, c)
	}
} */

func DeckReader(s string) Decks {

	dat, err := os.ReadFile(s)
	check(err)

	var o Decks
	errJ := json.Unmarshal(dat, &o)
	if errJ != nil {
		fmt.Println("unmarshall error:", errJ)
	}
	return o

	//fmt.Printf("%+v\n", o[0])
	//fmt.Println(len(o))
	//printRawDeck(o[0])

}
