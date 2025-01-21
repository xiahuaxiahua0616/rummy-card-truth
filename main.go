package main

import (
	"fmt"
	"rummy-card-truth/internal"
	"rummy-card-truth/pkg"
)

func main() {

	cards := []pkg.Card{
		{Suit: pkg.JokerSuit, Value: 0},
		{Suit: pkg.A, Value: 6},
		{Suit: pkg.B, Value: 8},
		{Suit: pkg.B, Value: 9},
		{Suit: pkg.A, Value: 3},
		{Suit: pkg.A, Value: 9},
		{Suit: pkg.B, Value: 12},
		{Suit: pkg.A, Value: 1},
		{Suit: pkg.D, Value: 3},
		{Suit: pkg.D, Value: 7},
		{Suit: pkg.D, Value: 8},
		{Suit: pkg.A, Value: 5},
		{Suit: pkg.A, Value: 4},
	}
	result := internal.NewPlanner(cards, 6).Run()
	num := 0
	var res []pkg.Card
	for _, r := range result {
		res = append(res, r...)
		num += len(r)
	}
	fmt.Println("差距", pkg.SliceDifferent(cards, res))
	fmt.Println("结果", result, "长度: ", num)
}
