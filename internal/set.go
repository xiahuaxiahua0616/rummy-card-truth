package internal

import (
	"sort"

	"github.com/xiahua/ifonly/pkg"
	"github.com/xiahuaxiahua0616/ifonlyutils/ifonlyutils"
)

func getSet(cards []pkg.Card) (setCards [][]pkg.Card, overCards []pkg.Card) {
	result := make(map[int]map[pkg.SuitVal]pkg.Card)
	for _, card := range cards {
		if result[card.Value] == nil {
			result[card.Value] = make(map[pkg.SuitVal]pkg.Card)
		}
		if _, exists := result[card.Value][card.Suit]; exists {
			overCards = append(overCards, card)
		} else {
			result[card.Value][card.Suit] = card
		}
	}

	for _, suits := range result {
		if len(suits) >= 3 {
			set := make([]pkg.Card, 0, len(suits))
			for _, card := range suits {
				set = append(set, card)
			}
			setCards = append(setCards, set)
		} else {
			for _, card := range suits {
				overCards = append(overCards, card)
			}
		}
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
	overCards = append(overCards, jokers...)
	return setWithJoker, overCards
}

func GetSetV2(cards []byte) (sets [][]byte, leftover []byte) {
	groupMap := make(map[byte][]byte)

	for _, card := range cards {
		val := card & 0x0F
		groupMap[val] = append(groupMap[val], card)
	}

	for _, group := range groupMap {
		if len(group) < 3 {
			leftover = append(leftover, group...)
			continue
		}

		uniq, dup := ifonlyutils.UniqueAndDuplicates(group)
		leftover = append(leftover, dup...)

		if len(uniq) >= 3 {
			sets = append(sets, uniq)
		} else {
			leftover = append(leftover, uniq...)
		}
	}

	return
}

func GetSetWithJokerV2(cards []byte, joker byte) (sets [][]byte, leftover []byte) {
	groupMap := make(map[byte][]byte)

	filterCards, jokers := getJokerV2(cards, joker)

	// 对三条进行分组
	for _, card := range filterCards {
		val := card & 0x0F
		groupMap[val] = append(groupMap[val], card)
	}

	// 排序要从大到小排序
	var keys []byte
	for k := range groupMap {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		getWeight := func(val byte) int {
			if val == 1 || val >= 10 {
				return 14
			}
			return int(val)
		}
		return getWeight(keys[i]) > getWeight(keys[j])
	})

	useJokerToFormSet := func(need int) {
		for _, k := range keys {
			cards := groupMap[k]
			if len(cards) == 0 || len(jokers) < need {
				continue
			}
			if len(cards)+need == 3 {
				set := make([]byte, 0, 3)
				set = append(set, cards...)
				set = append(set, jokers[:need]...)
				sets = append(sets, set)
				jokers = jokers[need:]
				delete(groupMap, k)
			}
		}
	}

	useJokerToFormSet(1)
	useJokerToFormSet(2)

	// 剩下未用的加入 leftover
	for _, cards := range groupMap {
		leftover = append(leftover, cards...)
	}
	leftover = append(leftover, jokers...)

	return
}
