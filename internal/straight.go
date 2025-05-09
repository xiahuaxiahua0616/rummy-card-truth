package internal

import (
	"slices"

	"github.com/xiahua/ifonly/pkg"
	"github.com/xiahuaxiahua0616/ifonlyutils/ifonlyutils"
)

func GetStraight(cards []byte, joker byte) ([][]byte, []byte) {
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

		var straight [][]byte
		var leftover []byte
		var score int
		for i, data := range datas {
			tempStraight, tempLeftover := findAllStraights(data)
			tempScore := ifonlyutils.CalcScore(tempLeftover, joker)
			if i == 0 || score > tempScore {
				straight = tempStraight
				leftover = tempLeftover
				score = tempScore
				continue
			}
		}
		leftover = append(leftover, duplicates...)

		return straight, leftover
	}
	return nil, nil
}

func findAllStraights(cards []byte) (straights [][]byte, remaining []byte) {
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
	remaining = append(remaining, cards...)
	return straights, remaining
}
