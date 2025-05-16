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
		if cards == nil || len(jokers) == 0 || len(cards)+len(jokers) < 3 {
			leftover = append(leftover, cards...)
			continue
		}

		var tempJokers []byte = make([]byte, len(jokers))
		copy(tempJokers, jokers)

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
			if len(data) < 2 && len(tempJokers) < 1 {
				leftover = append(leftover, data...)
				continue
			}
			// 找到当前可以带joker的合法顺子
			tempSuitStraight, tempSuitLeftOver := GetGapStraight(data, tempJokers, joker)
			tempScore := ifonlyutils.CalcScore(tempSuitLeftOver, joker)
			if i == 0 || score > tempScore {
				// 这个为了解决jokers重复使用的问题
				for _, tss := range tempSuitStraight {
					jokers = SliceDiffWithDup(jokers, tss)
				}

				tempStraight = tempSuitStraight
				tempLeftover = tempSuitLeftOver
				score = tempScore
			}
		}

		tempLeftover = SliceDiffWithDup(tempLeftover, tempJokers)

		straight = append(straight, tempStraight...)
		leftover = append(leftover, duplicates...)
		leftover = append(leftover, tempLeftover...)
	}
	return
}

func GetGapStraight(cards []byte, jokers []byte, joker byte) (result [][]byte, leftover []byte) {
	if len(cards) < 2 || len(jokers) < 1 {
		return nil, append(cards, jokers...)
	}

	tempCards := make([]byte, len(cards))
	tempJokers := make([]byte, len(jokers))
	copy(tempCards, cards)
	copy(tempJokers, jokers)

	straight := []byte{tempCards[0]}
	isUsed := false
	for i := 1; i < len(cards); i++ {
		last := straight[len(straight)-1]
		next := cards[i]

		if (next-last > 2) && len(straight) == 1 {
			straight = straight[1:]
			straight = append(straight, cards[i])
			continue
		}

		if (next-last > 2) && len(straight) > 1 {
			continue
		}

		if next-last > 2 && len(straight) < 2 {
			straight = straight[1:]
			straight = append(straight, cards[i])
			continue
		}

		if next-last == 0 {
			continue
		}

		if next-last == 1 {
			straight = append(straight, cards[i])
			continue
		}

		if last == joker || last > 0x4e {
			straight = append(straight, cards[i-1])
		}

		if !isUsed && next-last == 2 {
			isUsed = true
			straight = append(straight, jokers[0], cards[i])
			jokers = jokers[1:]
			continue
		}

	}
	if len(straight) == 2 && !isUsed {
		isUsed = true
		straight = append(straight, jokers[0])
		jokers = jokers[1:]
	}

	if len(straight) >= 3 {
		convStraight := ifonlyutils.Conv14to1(straight)
		result = append(result, convStraight)
		leftover = SliceDiffWithDup(cards, straight)
	} else {
		leftover = cards
	}
	leftover = append(leftover, jokers...)
	leftover = ifonlyutils.Conv14to1(leftover)
	return result, leftover
}
