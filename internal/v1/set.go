package v1

import (
	"sort"

	"github.com/xiahuaxiahua0616/ifonlyutils/ifonlyutils"
)

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
