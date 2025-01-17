package app

import (
	"github.com/jinzhu/copier"
	"sort"
)

func GetPure(rawCards []Card, isAsc bool) (pure [][]Card, overCard []Card) {
	// 我们接收一组牌，在这组牌当中找到顺子，返回结果和剩余牌
	// 为什么要复制一份？因为切片是指针类型，如果直接操作会影响外面的数据
	// 这里的职责是找到顺子，然后返回结果和剩余其他的不是该函数的要点。
	if len(rawCards) < 3 {
		return nil, rawCards
	}

	var cards []Card
	_ = copier.Copy(&cards, &rawCards)

	sort.Slice(cards, func(i, j int) bool {
		if isAsc {
			return cards[i].Value() < cards[j].Value()
		} else {
			return cards[i].Value() > cards[j].Value()
		}
	})

	// 计算因子，为了兼容降序
	factors := 1
	if !isAsc {
		factors = -1
	}

	for {
		if len(cards) < 3 {
			overCard = append(overCard, cards...)
			break
		}

		// 2. 比较，找到连续的值
		seq := []Card{cards[0]}
		for i := 1; i < len(cards); i++ {
			if cards[i].Value()-seq[len(seq)-1].Value() == 1*factors {
				seq = append(seq, cards[i])
			} else {
				break
			}
		}

		if len(seq) >= 3 {
			if len(seq) >= 6 && (len(seq)%6 == 0 || len(seq)%6 == 2 || len(seq)%6 == 4) {
				// 如果长度符合条件，将序列分为两部分
				middle := len(seq) / 2
				pure = append(pure, seq[:middle], seq[middle:])
			} else {
				// 其他情况下，直接加入 pure
				pure = append(pure, seq)
			}
			// 从 cards 中移除 seq
			cards = difference(cards, seq)
		} else {
			// 如果 seq 长度小于 3，直接将第一张卡片加入 overCard
			// 目的是处理可能顺子的第一张牌并不是顺子的情况
			overCard = append(overCard, cards[0])
			cards = cards[1:]
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

func changeTo1Dim(rawCards [][]Card) []Card {
	var cards []Card
	for _, v := range rawCards {
		cards = append(cards, v...)
	}
	return cards
}
