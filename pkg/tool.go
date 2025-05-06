package pkg

import (
	"fmt"
	"sort"
)

// SliceDifferent 切片差值
func SliceDifferent(a, b []Card) []Card {
	// 用 map 记录 b 中每张卡片的数量
	bCount := make(map[Card]int)
	for _, card := range b {
		bCount[card]++ // 记录每张卡片出现的次数
	}

	var cards []Card
	// 遍历 a，检查每张卡片是否在 b 中以及出现的次数
	for _, card := range a {
		if count, found := bCount[card]; found && count > 0 {
			bCount[card]-- // b 中减少一次计数
		} else {
			cards = append(cards, card) // 如果 b 中没有或计数为 0，则加入差值
		}
	}

	return cards
}

func CalculateScore(rawCards []Card, jokerVal int) int {
	score := 0
	for _, card := range rawCards {
		if card.Value == jokerVal {
			continue
		}

		if card.Value == 1 || card.Value > 10 {
			score += 10
		} else {
			score += card.Value
		}
	}
	return score
}

func CardValue1To14(rawCards []Card) []Card {
	for i, card := range rawCards {
		if card.Value == 1 {
			rawCards[i].Value = 14
		}
	}
	return rawCards
}

func CardValue14To1(rawCards []Card) []Card {
	for i, card := range rawCards {
		if card.Value == 14 {
			rawCards[i].Value = 1
		}
	}
	return rawCards
}

// GetSeqAllPossible 获取顺子全部的可能性
func GetSeqAllPossible(rawCards []Card) [][]Card {
	var result [][]Card

	if len(rawCards) < 3 {
		fmt.Println(rawCards, "牌数不足")
		return result
	}

	if len(rawCards) == 3 {
		return [][]Card{
			rawCards,
		}
	}

	if !JudgeIsSeq(rawCards) {
		fmt.Println(rawCards, "不是顺子")
		return result
	}

	for start := range rawCards {
		for end := start + 2; end <= len(rawCards); end++ {
			if JudgeIsSeq(rawCards[start:end]) {
				result = append(result, rawCards[start:end])
			}
		}
	}

	return result
}

func JudgeIsSeq(rawCards []Card) bool {
	if len(rawCards) < 3 {
		return false
	}

	var cards = make([]Card, len(rawCards))
	copy(cards, rawCards)

	// 判断升序顺子
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Value < cards[j].Value
	})

	// 检查升序顺子
	for i := 1; i < len(cards); i++ {
		if cards[i].Value-cards[i-1].Value != 1 {
			// 如果不是升序顺子，检查倒序
			// 先重新排序为倒序

			// 把A变成14
			CardValue1To14(cards)

			sort.Slice(cards, func(i, j int) bool {
				return cards[i].Value > cards[j].Value
			})

			// 检查倒序顺子
			for i := 1; i < len(cards); i++ {
				if cards[i-1].Value-cards[i].Value != 1 {
					return false
				}
			}
			// 如果倒序顺子也不成立，返回 false
			return true
		}
	}
	return true
}
