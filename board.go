package main

type board struct {
	columns [7]column //column[x][0] rests on the table for all x
	piles   [4]pile   //  pile[y][0] rests on the table for all y
	stock   []card    //    stock[0] rests on the table
	waste   []card    //    waste[0] rests on the table
}
