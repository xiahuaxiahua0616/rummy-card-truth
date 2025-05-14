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

		var score int
		for i, data := range datas {
			tempSuitStraight, tempLeftover := findAllStraights(data)
			tempScore := ifonlyutils.CalcScore(tempLeftover, joker)
			if i == 0 || score > tempScore {
				tempStraight = tempSuitStraight
				leftover = tempLeftover
				score = tempScore
				continue
			}
		}
		straight = append(straight, tempStraight...)
		leftover = append(leftover, duplicates...)
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
	groupedBySuit := ifonlyutils.GroupBySuit(cards)

	for suit, cards := range groupedBySuit {
		// 提取joker牌
		var jokers []byte
		cards, jokers = getJokerV2(cards, joker)
		if len(cards) < 3 && len(jokers) < 1 || suit == pkg.JokerSuitV2 || len(jokers) == 0 {
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

		var score int
		for i, data := range datas {
			if len(data) < 2 && len(jokers) < 1 {
				leftover = append(leftover, data...)
				continue
			}
			// 找到当前可以带joker的合法顺子
			tempStraight, tempLeftOver := GetGapStraight(cards, jokers)
			tempScore := ifonlyutils.CalcScore(tempLeftOver, joker)
			if i == 0 || score > tempScore {
				straight = tempStraight
				leftover = tempLeftOver
				score = tempScore
			}
		}
	}
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
	} else {
		leftover = append(leftover, tempCards...)
	}

	leftover = append(leftover, cards...)
	leftover = append(leftover, jokers...)
	return result, leftover
}
