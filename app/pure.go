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
			return cards[i].Value < cards[j].Value
		} else {
			return cards[i].Value > cards[j].Value
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
			if cards[i].Value-seq[len(seq)-1].Value == 1*factors {
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

func GetPureWithJoker(rawCards []Card, jokerVal int, isAsc bool) (pureWithJoker [][]Card, overCard []Card) {
	// 我们接收一组牌，在这组牌当中找到顺子，返回结果和剩余牌
	// 为什么要复制一份？因为切片是指针类型，如果直接操作会影响外面的数据
	// 这里的职责是找到顺子，然后返回结果和剩余其他的不是该函数的要点。
	if len(rawCards) < 3 {
		return nil, rawCards
	}

	var cards []Card
	_ = copier.Copy(&cards, &rawCards)

	// 找出手牌中所有的Joker
	var jokers []Card

	// 获取joker
	cards, jokers = GetJokers(cards, jokerVal)
	if len(jokers) < 1 {
		return nil, rawCards
	}

	sort.Slice(cards, func(i, j int) bool {
		if isAsc {
			return cards[i].Value < cards[j].Value
		} else {
			return cards[i].Value > cards[j].Value
		}
	})

	// 计算因子，为了兼容降序
	factors := 1
	if !isAsc {
		factors = -1
	}

	seq := []Card{cards[0]}
	isUsed := false

	for i := 1; i < len(cards); i++ {
		seqNextVal := seq[len(seq)-1].Value
		currentVal := cards[i].Value

		if currentVal-seqNextVal == 0 {
			// 相同的牌不进行处理
			overCard = append(overCard, cards[i])
			continue
		}

		if currentVal-seqNextVal == 1*factors {
			// 连续的牌
			seq = append(seq, cards[i])
			continue
		}

		if seqNextVal == jokerVal || seq[len(seq)-1].Suit == JokerSuit {
			// 间隙 == 2的牌
			seq = append(seq, cards[i-1])
			continue
		}

		if !isUsed && currentVal-seqNextVal == 2*factors {
			isUsed = true
			seq = append(seq, jokers[0], cards[i])
			jokers = jokers[1:]
			continue
		}
	}

	if len(seq) >= 3 {
		pureWithJoker = append(pureWithJoker, seq)
		// 从 cards 中移除 seq
		overCard = difference(cards, seq)
	}

	return pureWithJoker, overCard
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

func GetSet(cards []Card) (setCards [][]Card, overCards []Card) {
	result := make(map[int][]Card)
	// 按值分组
	for _, card := range cards {
		if result[card.Value] == nil {
			result[card.Value] = append(result[card.Value], card)
			continue
		}

		isExist := false
		for _, v := range result[card.Value] {
			if v.Suit == card.Suit {
				isExist = true
				break
			}
		}
		if isExist {
			overCards = append(overCards, card)
		} else {
			result[card.Value] = append(result[card.Value], card)
		}
	}

	for i, r := range result {
		if len(r) >= 3 {
			setCards = append(setCards, r)
			delete(result, i)
			continue
		}
		overCards = append(overCards, r...)
	}

	return setCards, overCards
}

func GetSetWithJoker(cards []Card, jokerVal int) (setWithJoker [][]Card, overCards []Card) {
	result := make(map[int][]Card)

	var jokers []Card
	// 获取joker
	cards, jokers = GetJokers(cards, jokerVal)
	if len(jokers) < 1 {
		return nil, cards
	}

	// 按值分组
	for _, card := range cards {
		if result[card.Value] == nil {
			result[card.Value] = append(result[card.Value], card)
			continue
		}

		isExist := false
		for _, v := range result[card.Value] {
			if v.Suit == card.Suit {
				isExist = true
				break
			}
		}
		if isExist {
			overCards = append(overCards, card)
		} else {
			result[card.Value] = append(result[card.Value], card)
		}
	}

	var keys []int
	for k := range result {
		keys = append(keys, k)
	}

	// 对key切片进行排序，从大到小
	sort.Sort(sort.Reverse(sort.IntSlice(keys)))

	// 消耗1张Joker牌
	for _, k := range keys {
		if len(jokers) < 1 {
			break
		}
		if len(result[k]) == 2 && len(jokers) >= 1 {
			temp := append(result[k], jokers[0])

			setWithJoker = append(setWithJoker, temp)
			jokers = jokers[1:]
			delete(result, k)
			continue
		}
	}

	// 消耗2张Joker牌
	for _, k := range keys {
		if len(jokers) < 2 {
			break
		}
		if len(result[k]) == 1 && len(jokers) >= 2 {
			temp := append(result[k], jokers[0], jokers[1])

			setWithJoker = append(setWithJoker, temp)
			jokers = jokers[2:]
			delete(result, k)
			continue
		}
	}

	for _, r := range result {
		overCards = append(overCards, r...)
	}
	return setWithJoker, overCards
}

func GetJokers(rawCards []Card, jokerVal int) (cards []Card, jokers []Card) {
	for i := 0; i < len(rawCards); i++ {
		if rawCards[i].Value == jokerVal || rawCards[i].Suit == JokerSuit {
			jokers = append(jokers, rawCards[i])
		} else {
			cards = append(cards, rawCards[i])
		}
	}
	return cards, jokers
}
