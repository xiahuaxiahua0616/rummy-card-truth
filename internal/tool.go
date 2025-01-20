package internal

import "rummy-card-truth/pkg"

func difference(a, b []pkg.Card) []pkg.Card {
	// 用 map 记录 b 中每张卡片的数量
	bCount := make(map[pkg.Card]int)
	for _, card := range b {
		bCount[card]++ // 记录每张卡片出现的次数
	}

	var cards []pkg.Card
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
