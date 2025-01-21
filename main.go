package main

import (
	"fmt"
	"rummy-card-truth/internal"
	"rummy-card-truth/pkg"
)

func main() {

	cards := []pkg.Card{
		{Suit: pkg.A, Value: 7},
		{Suit: pkg.A, Value: 8},
		{Suit: pkg.A, Value: 9},
		{Suit: pkg.A, Value: 4},
		{Suit: pkg.A, Value: 12},

		{Suit: pkg.B, Value: 5},

		{Suit: pkg.C, Value: 4},
		{Suit: pkg.C, Value: 2},

		{Suit: pkg.D, Value: 1},
		{Suit: pkg.D, Value: 11},
		{Suit: pkg.D, Value: 12},
		{Suit: pkg.D, Value: 6},
		{Suit: pkg.JokerSuit, Value: 0},
	}
	result := internal.NewPlanner(cards, 7).Run()
	num := 0
	var res []pkg.Card
	for _, r := range result {
		res = append(res, r...)
		num += len(r)
	}
	fmt.Println("少返回了: ", pkg.SliceDifferent(cards, res))
	fmt.Println("多返回了: ", pkg.SliceDifferent(res, cards))
	fmt.Println("结果", result, "长度: ", num)
}
