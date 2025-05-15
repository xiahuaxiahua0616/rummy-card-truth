package v1

import (
	"slices"

	"github.com/xiahua/ifonly/pkg"
	"github.com/xiahuaxiahua0616/ifonlyutils/ifonlyutils"
)

func GetStraight(cards []byte, joker byte) (straight [][]byte, leftover []byte) {
	groupedBySuit := ifonlyutils.GroupBySuit(cards)

	for suit, cards := range groupedBySuit {
		if len(cards) < 3 || suit == pkg.JokerSuitV2 {
			leftover = append(leftover, cards...)
			continue
		}

		slices.Sort(cards)
		// 去重
		cards, duplicates := ifonlyutils.UniqueAndDuplicates(cards)

		datas := [][]byte{
			cards,
			ifonlyutils.Conv1to14(cards),
		}

		var tempStraight [][]byte
		var tempLeftover []byte

		var score int
		for i, data := range datas {
			tempSuitStraight, tempSuitLeftover := findAllStraights(data)
			tempScore := ifonlyutils.CalcScore(tempSuitLeftover, joker)
			if i == 0 || score > tempScore {
				tempStraight = tempSuitStraight
				tempLeftover = tempSuitLeftover
				score = tempScore
				continue
			}
		}
		straight = append(straight, tempStraight...)
		leftover = append(leftover, duplicates...)
		leftover = append(leftover, tempLeftover...)
	}
	return
}

func findAllStraights(cards []byte) (straights [][]byte, leftover []byte) {
	for {
		prevLen := len(cards)
		slices.Sort(cards)

		straight, rest := ifonlyutils.GetStraightOnyByOne(cards)
		if len(straight) >= 3 {
			straights = append(straights, ifonlyutils.Conv14to1(straight))
		}
		cards = rest

		if len(cards) < 3 || len(cards) == prevLen {
			break
		}
	}
	leftover = append(leftover, cards...)
	return straights, leftover
}

func GetStraightWithJoker(cards []byte, joker byte) (straight [][]byte, leftover []byte) {
	var jokers []byte
	cards, jokers = getJokerV2(cards, joker)
	groupedBySuit := ifonlyutils.GroupBySuit(cards)

	for _, cards := range groupedBySuit {
		// 提取joker牌
		if cards == nil || len(jokers) == 0 || len(cards)+len(jokers) < 3 {
			leftover = append(leftover, cards...)
			continue
		}

		// 去重
		cards, duplicates := ifonlyutils.UniqueAndDuplicates(cards)
		leftover = append(leftover, duplicates...)

		slices.Sort(cards)

		datas := [][]byte{
			cards,
			ifonlyutils.Conv1to14(cards),
		}

		var tempStraight [][]byte
		var tempLeftover []byte

		var score int
		for i, data := range datas {
			if len(data) < 2 && len(jokers) < 1 {
				leftover = append(leftover, data...)
				continue
			}
			// 找到当前可以带joker的合法顺子
			tempSuitStraight, tempSuitLeftOver := GetGapStraight(data, jokers)
			tempScore := ifonlyutils.CalcScore(tempSuitLeftOver, joker)
			if i == 0 || score > tempScore {
				tempStraight = tempSuitStraight
				tempLeftover = tempSuitLeftOver
				score = tempScore
			}
		}
		for _, s := range tempStraight {
			jokers = SliceDiffWithDup(jokers, s)
		}
		jokers = SliceDiffWithDup(jokers, tempLeftover)
		straight = append(straight, tempStraight...)
		leftover = append(leftover, duplicates...)
		leftover = append(leftover, tempLeftover...)
		// leftover = append(leftover, joker)
	}
	leftover = append(leftover, jokers...)
	return
}

func GetGapStraight(cards []byte, jokers []byte) (result [][]byte, leftover []byte) {
	tempCards := []byte{}
	for len(cards) > 0 {
		// 初始化起始
		if len(tempCards) == 0 {
			tempCards = append(tempCards, cards[0])
			cards = cards[1:]
			continue
		}

		last := tempCards[len(tempCards)-1]
		next := cards[0]

		if next == last+1 {
			// 连续牌
			tempCards = append(tempCards, next)
			cards = cards[1:]

			if len(tempCards) == 3 {
				result = append(result, tempCards)
				tempCards = []byte{}
			}
		} else if next == last+2 && len(jokers) > 0 {
			// 中间断一个，用癞子补
			tempCards = append(tempCards, jokers[0], next)
			jokers = jokers[1:]
			cards = cards[1:]
			result = append(result, tempCards)
			tempCards = []byte{}
		} else {
			// 不连，扔掉当前 temp
			leftover = append(leftover, tempCards...)
			tempCards = []byte{}
		}
	}

	// 最后处理剩余的temp和jokers
	if len(tempCards) == 2 && len(jokers) > 0 {
		tempCards = append(tempCards, jokers[0])
		jokers = jokers[1:]
		result = append(result, tempCards)
	} else if len(tempCards) == 1 && len(jokers) >= 2 {
		// 特殊处理：1张原牌 + 2个癞子 → 成为顺子
		tempCards = append(tempCards, jokers[0], jokers[1])
		jokers = jokers[2:]
		result = append(result, tempCards)
	} else {
		leftover = append(leftover, tempCards...)
	}

	leftover = append(leftover, cards...)
	leftover = append(leftover, jokers...)

	for i, r := range result {
		result[i] = ifonlyutils.Conv14to1(r)
	}

	return result, ifonlyutils.Conv14to1(leftover)
}
