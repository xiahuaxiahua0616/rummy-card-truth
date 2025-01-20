package internal

import (
	"github.com/jinzhu/copier"
	"rummy-card-truth/pkg"
	"sort"
)

func GetPure(rawCards []pkg.Card, isAsc bool) (pure [][]pkg.Card, overCard []pkg.Card) {
	// 我们接收一组牌，在这组牌当中找到顺子，返回结果和剩余牌
	// 为什么要复制一份？因为切片是指针类型，如果直接操作会影响外面的数据
	// 这里的职责是找到顺子，然后返回结果和剩余其他的不是该函数的要点。
	if len(rawCards) < 3 {
		return nil, rawCards
	}

	var cards []pkg.Card
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
		seq := []pkg.Card{cards[0]}
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

func GetPureWithJoker(rawCards []pkg.Card, jokerVal int, isAsc bool) (pureWithJoker [][]pkg.Card, overCard []pkg.Card) {
	// 我们接收一组牌，在这组牌当中找到顺子，返回结果和剩余牌
	// 为什么要复制一份？因为切片是指针类型，如果直接操作会影响外面的数据
	// 这里的职责是找到顺子，然后返回结果和剩余其他的不是该函数的要点。
	if len(rawCards) < 3 {
		return nil, rawCards
	}

	var cards []pkg.Card
	_ = copier.Copy(&cards, &rawCards)

	// 找出手牌中所有的Joker
	var jokers []pkg.Card

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

	seq := []pkg.Card{cards[0]}
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

		if seqNextVal == jokerVal || seq[len(seq)-1].Suit == pkg.JokerSuit {
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
