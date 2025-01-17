package app

import (
	"github.com/jinzhu/copier"
	"sort"
)

func GetPure(rawCards []Card) (pure, overCard []Card) {
	// 我们接收一组牌，在这组牌当中找到顺子，返回结果和剩余牌
	// 为什么要复制一份？因为切片是指针类型，如果直接操作会影响外面的数据
	// 这里的职责是找到顺子，然后返回结果和剩余其他的不是该函数的要点。
	if len(rawCards) < 3 {
		return nil, rawCards
	}

	var cards []Card
	_ = copier.Copy(&cards, &rawCards)

	for {
		if len(cards) < 3 {
			break
		}

		minCard := cards[0]
		maxCard := cards[0]

		for _, c := range cards {
			if c.Value() < minCard.Value() {
				minCard = c
			}

			if c.Value() > maxCard.Value() {
				maxCard = c
			}
		}

		isSeq := maxCard.Value()-minCard.Value() == len(cards)-1
		if !isSeq {
			overCard = append(overCard, cards[0])
			cards = cards[1:]
		} else {
			pure = append(pure, cards...)
			cards = difference(cards, pure)
		}
	}
	return pure, overCard
}

func AscCards(cards []Card) {
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Value() < cards[j].Value()
	})
}

func difference(a, b []Card) []Card {
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
