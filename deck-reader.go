package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type decks []deck

func check(e error) {
	if e != nil {
		panic(e)
	}
}

/*
func printRawDeck(d deck) {
	for i, c := range d {
		fmt.Printf("Index: %v : Card: %+v\n", i, c)
	}
} */

func deckReader(s string) decks {

	dat, err := os.ReadFile(s)
	check(err)

	var o decks
	errJ := json.Unmarshal(dat, &o)
	if errJ != nil {
		fmt.Println("unmarshall error:", errJ)
	}
	return o

	//fmt.Printf("%+v\n", o[0])
	//fmt.Println(len(o))
	//printRawDeck(o[0])

}
