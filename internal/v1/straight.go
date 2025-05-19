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

	straight, leftover = getStraightByGroup(joker, jokers, groupedBySuit)

	// var descCards = make([][]byte, len(groupedBySuit))
	// copy(descCards, groupedBySuit)
	// slices.Reverse(descCards)

	// straight2, leftover2 := getStraightByGroup(joker, jokers, descCards)
	// score1 := ifonlyutils.CalcScore(leftover, joker)
	// score2 := ifonlyutils.CalcScore(leftover2, joker)
	// if score2 < score1 {
	// 	straight = straight2
	// 	leftover = leftover2
	// }

	return
}

func getStraightByGroup(joker byte, jokers []byte, datas [][]byte) (straight [][]byte, leftover []byte) {
	for _, cards := range datas {
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
			tempSuitStraight, tempSuitLeftOver := GetGapStraight(data, tempJokers, joker, true)
			tempScore := ifonlyutils.CalcScore(tempSuitLeftOver, joker)

			descData := make([]byte, len(data))
			copy(descData, data)
			slices.Reverse(descData)

			tempSuitStraight2, tempSuitLeftOver2 := GetGapStraight(descData, tempJokers, joker, false)
			tempScore2 := ifonlyutils.CalcScore(tempSuitLeftOver2, joker)
			if tempScore2 < tempScore {
				tempScore = tempScore2
				tempSuitStraight = tempSuitStraight2
				tempSuitLeftOver = tempSuitLeftOver2
			}

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

func GetGapStraight(cards []byte, jokers []byte, joker byte, asc bool) (result [][]byte, leftover []byte) {
	if len(cards) < 2 || len(jokers) < 1 {
		return nil, append(cards, jokers...)
	}

	tempCards := make([]byte, len(cards))
	tempJokers := make([]byte, len(jokers))
	copy(tempCards, cards)
	copy(tempJokers, jokers)

	diff := func(a, b byte) int {
		if asc {
			return int(b) - int(a)
		} else {
			return int(a) - int(b)
		}
	}
	isGapOverRange := func(a, b byte) bool { return diff(a, b) > 2 }
	isSameCard := func(a, b byte) bool { return diff(a, b) == 0 }
	isConsecutive := func(a, b byte) bool { return diff(a, b) == 1 }
	isGapCard := func(a, b byte) bool { return diff(a, b) == 2 }

	straight := []byte{tempCards[0]}
	isUsed := false
	for i := 1; i < len(cards); i++ {
		last := straight[len(straight)-1]
		next := cards[i]

		slen := len(straight)

		isStart := slen == 1
		isTooShort := slen < 2
		isValid := slen > 1

		switch {
		case isGapOverRange(last, next) && isStart:
			straight = straight[1:]
			straight = append(straight, next)
		case isGapOverRange(last, next) && isValid:
			// too far to continue, ignore
			continue
		case isGapOverRange(last, next) && isTooShort:
			straight = straight[1:]
			straight = append(straight, next)
		case isSameCard(last, next):
			continue
		case isConsecutive(last, next):
			straight = append(straight, next)
		case last == joker || last > 0x4e:
			straight = append(straight, tempCards[i-1])
		case !isUsed && isGapCard(last, next) && len(jokers) > 0:
			isUsed = true
			straight = append(straight, jokers[0], next)
			jokers = jokers[1:]
		}

	}
	if len(straight) == 2 && !isUsed && len(jokers) > 0 {
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
