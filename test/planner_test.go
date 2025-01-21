package test

import (
	"fmt"
	"math/rand"
	"rummy-card-truth/internal"
	"rummy-card-truth/pkg"
	"testing"
	"time"
)

func InitializeDeck() (deck []pkg.Card) {
	for i := 0; i < 2; i++ {
		for _, suit := range []pkg.SuitVal{pkg.A, pkg.B, pkg.C, pkg.D} {
			for value := 1; value <= 13; value++ {
				deck = append(deck, pkg.Card{Suit: suit, Value: value})
			}
		}

		// 添加大小王
		deck = append(deck, pkg.Card{Suit: pkg.JokerSuit, Value: 0})
		deck = append(deck, pkg.Card{Suit: pkg.JokerSuit, Value: 0})
	}

	return
}

func ShuffleDeck(deck []pkg.Card) []pkg.Card {
	rand.NewSource(time.Now().UnixNano()) // 设置随机种子
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
	return deck
}

func DealCards(deck *[]pkg.Card, numCards int) []pkg.Card {
	// numCards不能超过排堆大小
	if numCards > len(*deck) {
		panic("too many cards requested")
	}
	hand := (*deck)[:numCards]
	*deck = (*deck)[numCards:]
	return hand
}

func TestLenIsTrue(t *testing.T) {

	var totalNum time.Duration = 0

	// 计算平均耗时
	defer func() {
		fmt.Println("平均耗时:", totalNum/100000000)
	}()

	for i := range 1000000 {
		desk := InitializeDeck()
		ShuffleDeck(desk)
		cards := DealCards(&desk, 13)

		jokerV := rand.Intn(13) + 1

		CalcStart := time.Now()

		result := internal.NewPlanner(cards, jokerV).Run()

		totalNum += time.Since(CalcStart)

		num := 0
		for _, r := range result {
			num += len(r)
		}

		if num != 13 {
			fmt.Println(result)
			fmt.Printf("第 %v 次错误, num is not 13 \n", i)

			for _, cc := range cards {
				res := ""
				if cc.Suit == pkg.A {
					res = "pkg.A"
				} else if cc.Suit == pkg.B {
					res = "pkg.B"
				} else if cc.Suit == pkg.C {
					res = "pkg.C"
				} else if cc.Suit == pkg.D {
					res = "pkg.D"
				} else if cc.Suit == pkg.JokerSuit {
					res = "pkg.JokerSuit"
				}
				fmt.Printf("{Suit: %s, Value: %d},\n", res, cc.Value)
			}
			fmt.Println("jokerV", jokerV)
			//t.Fatalf("第 %v 次错误, num is not 13", i)
		}
	}
}
