package test

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/xiahua/ifonly/internal"
	"github.com/xiahua/ifonly/pkg"
)

func InitializeDeck() (deck []pkg.Card) {
	for range 2 {
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
				fmt.Printf("{Suit: pkg.%s, Value: %d},\n", cc.Suit, cc.Value)
			}
			fmt.Println("jokerV", jokerV)
			//t.Fatalf("第 %v 次错误, num is not 13", i)
		}
	}
}

func GetCardsResult(cards []pkg.Card) []int {
	var myCards []int
	for _, c := range cards {
		if c.Suit == pkg.A {
			myCards = append(myCards, c.Value+48)
		} else if c.Suit == pkg.B {
			myCards = append(myCards, c.Value+32)
		} else if c.Suit == pkg.C {
			myCards = append(myCards, c.Value+16)
		} else if c.Suit == pkg.D {
			myCards = append(myCards, c.Value)
		} else if c.Suit == pkg.JokerSuit {
			myCards = append(myCards, 79)
		} else if c.Suit == pkg.JokerSuit {
			myCards = append(myCards, 78)
		}
	}

	if len(myCards) == 0 {
		return []int{0}
	}
	return myCards
}

func GetResponse(cards ...[][]pkg.Card) [][]int {
	var res [][]int

	for _, cardDim := range cards {
		for _, card := range cardDim {
			if len(card) > 0 {
				res = append(res, GetCardsResult(card))
			}
		}
	}

	for _, card := range res {
		sort.Ints(card)
	}
	return res
}

func DisSorting1DimTo2Dim(cardsRaw []pkg.Card) [][]pkg.Card {
	var response [][]pkg.Card

	response = append(response, cardsRaw)

	return response
}

func Test1(t *testing.T) {
	fmt.Println("执行了吗？1")
	cards := []pkg.Card{
		{Suit: pkg.D, Value: 1},
		{Suit: pkg.D, Value: 12},
		{Suit: pkg.D, Value: 13},
		{Suit: pkg.D, Value: 2},
		{Suit: pkg.D, Value: 6},
		{Suit: pkg.D, Value: 9},
		{Suit: pkg.D, Value: 7},
		{Suit: pkg.D, Value: 10},

		{Suit: pkg.C, Value: 2},
		{Suit: pkg.C, Value: 4},

		{Suit: pkg.A, Value: 7},
		{Suit: pkg.A, Value: 9},
		{Suit: pkg.A, Value: 12},
	}

	jokerV := 7

	result := internal.NewPlanner(cards, jokerV).Run()

	totalNum := 0
	for _, r := range result {
		totalNum += len(r)
	}

	if totalNum != 13 {
		t.Error("牌数不等于13", totalNum)
	}

	var data = []int{2, 6, 57, 60}
	if !areArraysEqual(GetResponse(result)[len(result)-1], data) {
		t.Error("数据不正确", result[len(result)-1], data)
	}
}

func Test2(t *testing.T) {
	fmt.Println("执行了吗？2")
	cards := []pkg.Card{
		{Suit: pkg.A, Value: 2},
		{Suit: pkg.A, Value: 3},
		{Suit: pkg.A, Value: 4},
		{Suit: pkg.A, Value: 5},
		{Suit: pkg.A, Value: 7},
		{Suit: pkg.A, Value: 9},

		{Suit: pkg.B, Value: 2},
		{Suit: pkg.B, Value: 5},

		{Suit: pkg.C, Value: 1},
		{Suit: pkg.C, Value: 13},
		{Suit: pkg.C, Value: 9},
		{Suit: pkg.C, Value: 2},

		{Suit: pkg.D, Value: 5},
	}

	jokerV := 5

	result := internal.NewPlanner(cards, jokerV).Run()

	totalNum := 0
	for _, r := range result {
		totalNum += len(r)
	}

	if totalNum != 13 {
		t.Error("牌数不等于13", totalNum)
	}

	var data = []int{55}
	if !areArraysEqual(GetResponse(result)[len(result)-1], data) {
		t.Error("数据不正确", result[len(result)-1], data)
	}
}

func Test3(t *testing.T) {
	fmt.Println("执行了吗？3")
	cards := []pkg.Card{
		{Suit: pkg.A, Value: 4},
		{Suit: pkg.A, Value: 9},
		{Suit: pkg.A, Value: 11},
		{Suit: pkg.A, Value: 12},

		{Suit: pkg.B, Value: 4},
		{Suit: pkg.B, Value: 5},

		{Suit: pkg.C, Value: 1},
		{Suit: pkg.C, Value: 2},
		{Suit: pkg.C, Value: 9},
		{Suit: pkg.C, Value: 11},
		{Suit: pkg.C, Value: 12},
		{Suit: pkg.C, Value: 13},

		{Suit: pkg.D, Value: 9},
	}

	jokerV := 5

	result := internal.NewPlanner(cards, jokerV).Run()

	totalNum := 0
	for _, r := range result {
		totalNum += len(r)
	}

	if totalNum != 13 {
		t.Error("牌数不等于13", totalNum)
	}

	var data = []int{18, 36, 52}
	if !areArraysEqual(GetResponse(result)[len(result)-1], data) {
		t.Error("数据不正确", result[len(result)-1], data)
	}
}

func Test4(t *testing.T) {
	fmt.Println("执行了吗？4")
	cards := []pkg.Card{
		{Suit: pkg.JokerSuit, Value: 0},
		{Suit: pkg.A, Value: 7},

		{Suit: pkg.B, Value: 7},
		{Suit: pkg.B, Value: 8},
		{Suit: pkg.B, Value: 9},
		{Suit: pkg.B, Value: 10},
		{Suit: pkg.B, Value: 12},

		{Suit: pkg.C, Value: 5},

		{Suit: pkg.D, Value: 4},
		{Suit: pkg.D, Value: 5},
		{Suit: pkg.D, Value: 11},
		{Suit: pkg.D, Value: 12},
		{Suit: pkg.D, Value: 13},
	}

	jokerV := 8

	result := internal.NewPlanner(cards, jokerV).Run()

	totalNum := 0
	for _, r := range result {
		totalNum += len(r)
	}

	if totalNum != 13 {
		t.Error("牌数不等于13", totalNum)
	}

	var data = []int{4, 5, 21}
	if !areArraysEqual(GetResponse(result)[len(result)-1], data) {
		t.Error("数据不正确", result[len(result)-1], data)
	}
}

func Test5(t *testing.T) {
	fmt.Println("执行了吗？5")
	cards := []pkg.Card{
		{Suit: pkg.A, Value: 8},
		{Suit: pkg.A, Value: 9},
		{Suit: pkg.A, Value: 10},
		{Suit: pkg.A, Value: 11},
		{Suit: pkg.A, Value: 12},

		{Suit: pkg.B, Value: 1},
		{Suit: pkg.B, Value: 8},
		{Suit: pkg.B, Value: 10},

		{Suit: pkg.C, Value: 4},
		{Suit: pkg.C, Value: 5},
		{Suit: pkg.C, Value: 10},

		{Suit: pkg.D, Value: 7},
		{Suit: pkg.D, Value: 10},
	}

	jokerV := 1

	result := internal.NewPlanner(cards, jokerV).Run()

	totalNum := 0
	for _, r := range result {
		totalNum += len(r)
	}

	if totalNum != 13 {
		t.Error("牌数不等于13", totalNum)
	}

	var data = []int{7, 40}
	if !areArraysEqual(GetResponse(result)[len(result)-1], data) {
		t.Error("数据不正确", result[len(result)-1], data)
	}
}

func Test6(t *testing.T) {
	fmt.Println("执行了吗？6")
	cards := []pkg.Card{
		{Suit: pkg.A, Value: 1},
		{Suit: pkg.A, Value: 2},
		{Suit: pkg.A, Value: 3},
		{Suit: pkg.A, Value: 6},
		{Suit: pkg.A, Value: 7},
		{Suit: pkg.A, Value: 11},

		{Suit: pkg.B, Value: 6},
		{Suit: pkg.B, Value: 12},

		{Suit: pkg.C, Value: 1},
		{Suit: pkg.C, Value: 3},
		{Suit: pkg.C, Value: 13},

		{Suit: pkg.D, Value: 3},
		{Suit: pkg.JokerSuit, Value: 0},
	}

	jokerV := 5

	result := internal.NewPlanner(cards, jokerV).Run()

	totalNum := 0
	for _, r := range result {
		totalNum += len(r)
	}

	if totalNum != 13 {
		t.Error("牌数不等于13", totalNum)
	}

	var data = []int{3, 19, 38, 44, 54, 55, 59}
	if !areArraysEqual(GetResponse(result)[len(result)-1], data) {
		t.Error("数据不正确", result[len(result)-1], data)
	}
}

func Test7(t *testing.T) {
	fmt.Println("执行了吗？7")
	cards := []pkg.Card{
		{Suit: pkg.D, Value: 3},
		{Suit: pkg.D, Value: 5},
		{Suit: pkg.JokerSuit, Value: 0},

		{Suit: pkg.C, Value: 11},
		{Suit: pkg.C, Value: 12},
		{Suit: pkg.C, Value: 13},

		{Suit: pkg.B, Value: 2},
		{Suit: pkg.B, Value: 3},
		{Suit: pkg.B, Value: 4},
		{Suit: pkg.B, Value: 8},
		{Suit: pkg.B, Value: 9},
		{Suit: pkg.B, Value: 9},
		{Suit: pkg.B, Value: 10},

		{Suit: pkg.A, Value: 10},
	}

	jokerV := 9

	result := internal.NewPlanner(cards, jokerV).Run()

	totalNum := 0
	for _, r := range result {
		totalNum += len(r)
	}

	if totalNum != 14 {
		// todo:: 本身就有14张
		t.Error("牌数不等于13", totalNum)
	}

	var data = []int{3, 5}
	var data2 = []int{40, 78}
	if !areArraysEqual(GetResponse(result)[len(result)-1], data) && !areArraysEqual(GetResponse(result)[len(result)-1], data2) {
		t.Error("数据不正确", result[len(result)-1], data, data2)
	}
}

func Test8(t *testing.T) {
	fmt.Println("执行了吗？8")
	cards := []pkg.Card{
		{Suit: pkg.D, Value: 5},
		{Suit: pkg.D, Value: 6},
		{Suit: pkg.D, Value: 12},

		{Suit: pkg.C, Value: 4},
		{Suit: pkg.C, Value: 5},
		{Suit: pkg.C, Value: 6},
		{Suit: pkg.C, Value: 6},

		{Suit: pkg.B, Value: 2},
		{Suit: pkg.B, Value: 3},
		{Suit: pkg.B, Value: 6},
		{Suit: pkg.B, Value: 12},

		{Suit: pkg.A, Value: 10},
		{Suit: pkg.A, Value: 11},
	}

	jokerV := 11

	result := internal.NewPlanner(cards, jokerV).Run()

	totalNum := 0
	for _, r := range result {
		totalNum += len(r)
	}

	if totalNum != 13 {
		t.Error("牌数不等于13", totalNum)
	}

	var data = []int{5, 12, 44, 58}
	if !areArraysEqual(GetResponse(result)[len(result)-1], data) {
		t.Error("数据不正确", result[len(result)-1], data)
	}
}

func Test9(t *testing.T) {
	fmt.Println("执行了吗？9")
	cards := []pkg.Card{
		{Suit: pkg.D, Value: 5},
		{Suit: pkg.D, Value: 6},
		{Suit: pkg.D, Value: 12},

		{Suit: pkg.C, Value: 4},
		{Suit: pkg.C, Value: 5},
		{Suit: pkg.C, Value: 6},
		{Suit: pkg.C, Value: 6},

		{Suit: pkg.B, Value: 2},
		{Suit: pkg.B, Value: 3},
		{Suit: pkg.B, Value: 6},
		{Suit: pkg.B, Value: 12},

		{Suit: pkg.A, Value: 10},
		{Suit: pkg.A, Value: 11},
	}

	jokerV := 11

	result := internal.NewPlanner(cards, jokerV).Run()

	totalNum := 0
	for _, r := range result {
		totalNum += len(r)
	}

	if totalNum != 13 {
		t.Error("牌数不等于13", totalNum)
	}

	var data = []int{5, 12, 44, 58}
	if !areArraysEqual(GetResponse(result)[len(result)-1], data) {
		t.Error("数据不正确", result[len(result)-1], data)
	}
}

func Test10(t *testing.T) {
	fmt.Println("执行了吗？10")
	cards := []pkg.Card{
		{Suit: pkg.A, Value: 1},
		{Suit: pkg.A, Value: 2},
		{Suit: pkg.A, Value: 3},

		{Suit: pkg.B, Value: 11},
		{Suit: pkg.B, Value: 12},
		{Suit: pkg.B, Value: 2},
		{Suit: pkg.B, Value: 3},
		{Suit: pkg.B, Value: 6},
		{Suit: pkg.B, Value: 7},
		{Suit: pkg.B, Value: 9},

		{Suit: pkg.C, Value: 4},

		{Suit: pkg.D, Value: 6},
		{Suit: pkg.JokerSuit, Value: 0},
	}

	jokerV := 4

	result := internal.NewPlanner(cards, jokerV).Run()

	totalNum := 0
	for _, r := range result {
		totalNum += len(r)
	}

	if totalNum != 13 {
		t.Error("牌数不等于13", totalNum)
	}

	var data = []int{6, 34, 35}
	if !areArraysEqual(GetResponse(result)[len(result)-1], data) {
		t.Error("数据不正确", result[len(result)-1], data)
	}
}

func Test11(t *testing.T) {
	fmt.Println("执行了吗？11")
	cards := []pkg.Card{
		{Suit: pkg.A, Value: 3},
		{Suit: pkg.A, Value: 4},
		{Suit: pkg.A, Value: 5},
		{Suit: pkg.A, Value: 6},
		{Suit: pkg.A, Value: 10},

		{Suit: pkg.B, Value: 1},
		{Suit: pkg.B, Value: 2},
		{Suit: pkg.B, Value: 9},
		{Suit: pkg.B, Value: 13},

		{Suit: pkg.C, Value: 6},
		{Suit: pkg.C, Value: 10},
		{Suit: pkg.C, Value: 13},

		{Suit: pkg.D, Value: 1},
	}

	jokerV := 2
	result := internal.NewPlanner(cards, jokerV).Run()

	totalNum := 0
	for _, r := range result {
		totalNum += len(r)
	}

	if totalNum != 13 {
		t.Error("牌数不等于13", totalNum)
	}

	var data = []int{1, 22, 26, 29, 41, 58}
	if !areArraysEqual(GetResponse(result)[len(result)-1], data) {
		t.Error("数据不正确", result[len(result)-1], data)
	}
}

func Test12(t *testing.T) {
	fmt.Println("执行了吗？12")
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

	jokerV := 7

	for range 10000 {
		result := internal.NewPlanner(cards, jokerV).Run()

		totalNum := 0
		for _, r := range result {
			totalNum += len(r)
		}

		if totalNum != 13 {
			t.Error("牌数不等于13", totalNum)
		}

		var data = []int{6, 18, 20, 37, 52, 60}
		if !areArraysEqual(GetResponse(result)[len(result)-1], data) {
			t.Error("数据不正确", result[len(result)-1], data)
		}
	}
}

func Test13(t *testing.T) {
	fmt.Println("执行了吗？13")
	cards := []pkg.Card{
		{Suit: pkg.D, Value: 1},
		{Suit: pkg.D, Value: 2},
		{Suit: pkg.D, Value: 12},
		{Suit: pkg.D, Value: 13},
		{Suit: pkg.D, Value: 9},
		{Suit: pkg.D, Value: 10},
		{Suit: pkg.D, Value: 11},

		{Suit: pkg.C, Value: 9},

		{Suit: pkg.B, Value: 2},
		{Suit: pkg.B, Value: 3},
		{Suit: pkg.B, Value: 4},
		{Suit: pkg.B, Value: 5},
		{Suit: pkg.B, Value: 9},
	}

	jokerV := 8
	result := internal.NewPlanner(cards, jokerV).Run()

	totalNum := 0
	for _, r := range result {
		totalNum += len(r)
	}

	if totalNum != 13 {
		t.Error("牌数不等于13", totalNum)
	}

	var data = []int{2}

	if !areArraysEqual(GetResponse(result)[len(result)-1], data) {
		t.Error("数据不正确", result[len(result)-1], data)
	}
}

func Test14(t *testing.T) {
	fmt.Println("执行了吗？14")
	cards := []pkg.Card{
		{Suit: pkg.D, Value: 7},
		{Suit: pkg.D, Value: 8},
		{Suit: pkg.D, Value: 9},

		{Suit: pkg.C, Value: 3},
		{Suit: pkg.C, Value: 10},
		{Suit: pkg.C, Value: 13},

		{Suit: pkg.B, Value: 3},

		{Suit: pkg.A, Value: 2},
		{Suit: pkg.A, Value: 3},
		{Suit: pkg.A, Value: 5},
		{Suit: pkg.A, Value: 8},
		{Suit: pkg.A, Value: 13},
		{Suit: pkg.JokerSuit, Value: 0},
	}

	jokerV := 11
	result := internal.NewPlanner(cards, jokerV).Run()

	totalNum := 0
	for _, r := range result {
		totalNum += len(r)
	}

	if totalNum != 13 {
		t.Error("牌数不等于13", totalNum)
	}

	var data = []int{19, 26, 29, 35, 56, 61}
	if !areArraysEqual(GetResponse(result)[len(result)-1], data) {
		t.Error("数据不正确", result[len(result)-1], data)
	}
}

func Test15(t *testing.T) {
	fmt.Println("执行了吗？15")
	cards := []pkg.Card{
		{Suit: pkg.D, Value: 2},
		{Suit: pkg.D, Value: 2},
		{Suit: pkg.D, Value: 4},
		{Suit: pkg.D, Value: 1},

		{Suit: pkg.B, Value: 3},
		{Suit: pkg.B, Value: 5},
		{Suit: pkg.B, Value: 7},
		{Suit: pkg.B, Value: 8},
		{Suit: pkg.B, Value: 10},
		{Suit: pkg.B, Value: 11},
		{Suit: pkg.B, Value: 12},
		{Suit: pkg.B, Value: 13},

		{Suit: pkg.A, Value: 10},
	}

	jokerV := 1
	result := internal.NewPlanner(cards, jokerV).Run()

	totalNum := 0
	for _, r := range result {
		totalNum += len(r)
	}

	if totalNum != 13 {
		t.Error("牌数不等于13", totalNum)
	}

	var data = []int{2, 2, 4, 35, 58}
	if !areArraysEqual(GetResponse(result)[len(result)-1], data) {
		t.Error("数据不正确", result[len(result)-1], data)
	}
}

func Test16(t *testing.T) {
	fmt.Println("执行了吗？16")
	cards := []pkg.Card{
		{Suit: pkg.B, Value: 11},
		{Suit: pkg.B, Value: 12},
		{Suit: pkg.A, Value: 10},
		{Suit: pkg.C, Value: 11},
		{Suit: pkg.C, Value: 10},
		{Suit: pkg.C, Value: 13},
		{Suit: pkg.D, Value: 8},
		{Suit: pkg.D, Value: 9},
		{Suit: pkg.D, Value: 10},
		{Suit: pkg.C, Value: 6},
		{Suit: pkg.C, Value: 9},
		{Suit: pkg.A, Value: 10},
		{Suit: pkg.A, Value: 2},
	}

	jokerV := 10
	result := internal.NewPlanner(cards, jokerV).Run()

	totalNum := 0
	for _, r := range result {
		totalNum += len(r)
	}

	if totalNum != 13 {
		t.Error("牌数不等于13", totalNum)
	}

	var data = []int{22, 25, 50, 58}
	if !areArraysEqual(GetResponse(result)[len(result)-1], data) {
		t.Error("数据不正确", result[len(result)-1], data)
	}
}

func areArraysEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	sort.Ints(a)
	sort.Ints(b)
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
