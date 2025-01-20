package internal

import (
	"rummy-card-truth/pkg"
	"sort"
)

func getSet(cards []pkg.Card) (setCards [][]pkg.Card, overCards []pkg.Card) {
	result := make(map[int][]pkg.Card)
	// 按值分组
	for _, card := range cards {
		if result[card.Value] == nil {
			result[card.Value] = append(result[card.Value], card)
			continue
		}

		isExist := false
		for _, v := range result[card.Value] {
			if v.Suit == card.Suit {
				isExist = true
				break
			}
		}
		if isExist {
			overCards = append(overCards, card)
		} else {
			result[card.Value] = append(result[card.Value], card)
		}
	}

	for i, r := range result {
		if len(r) >= 3 {
			setCards = append(setCards, r)
			delete(result, i)
			continue
		}
		overCards = append(overCards, r...)
	}

	return setCards, overCards
}

func getSetWithJoker(cards []pkg.Card, jokerVal int) (setWithJoker [][]pkg.Card, overCards []pkg.Card) {
	result := make(map[int][]pkg.Card)

	var jokers []pkg.Card
	// 获取joker
	cards, jokers = getJokers(cards, jokerVal)
	if len(jokers) < 1 {
		return nil, cards
	}

	// 按值分组
	for _, card := range cards {
		if result[card.Value] == nil {
			result[card.Value] = append(result[card.Value], card)
			continue
		}

		isExist := false
		for _, v := range result[card.Value] {
			if v.Suit == card.Suit {
				isExist = true
				break
			}
		}
		if isExist {
			overCards = append(overCards, card)
		} else {
			result[card.Value] = append(result[card.Value], card)
		}
	}

	var keys []int
	for k := range result {
		keys = append(keys, k)
	}

	// 对key切片进行排序，从大到小
	sort.Sort(sort.Reverse(sort.IntSlice(keys)))

	// 消耗1张Joker牌
	for _, k := range keys {
		if len(jokers) < 1 {
			break
		}
		if len(result[k]) == 2 && len(jokers) >= 1 {
			temp := append(result[k], jokers[0])

			setWithJoker = append(setWithJoker, temp)
			jokers = jokers[1:]
			delete(result, k)
			continue
		}
	}

	// 消耗2张Joker牌
	for _, k := range keys {
		if len(jokers) < 2 {
			break
		}
		if len(result[k]) == 1 && len(jokers) >= 2 {
			temp := append(result[k], jokers[0], jokers[1])

			setWithJoker = append(setWithJoker, temp)
			jokers = jokers[2:]
			delete(result, k)
			continue
		}
	}

	for _, r := range result {
		overCards = append(overCards, r...)
	}
	return setWithJoker, overCards
}
