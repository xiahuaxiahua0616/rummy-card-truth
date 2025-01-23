package internal

import (
	"github.com/jinzhu/copier"
	"rummy-card-truth/pkg"
	"sort"
)

func (p *Planner) pureSetup(rawCards []pkg.Card) (cards [][]pkg.Card, overCards []pkg.Card) {
	result1, overCards1 := p.pureSetup1(rawCards)
	result2, overCards2 := p.pureSetup2(rawCards)
	result3, overCards3 := p.pureSetup3(rawCards)

	score1 := pkg.CalculateScore(overCards1, p.jokerVal)
	score2 := pkg.CalculateScore(overCards2, p.jokerVal)
	score3 := pkg.CalculateScore(overCards3, p.jokerVal)

	minScore := min(score1, score2, score3)
	if score1 == minScore {
		cards = result1
		overCards = overCards1

	} else if score2 == minScore {
		cards = result2
		overCards = overCards2
	} else {
		cards = result3
		overCards = overCards3
	}
	return cards, overCards
}

func (p *Planner) getBasePure(cards []pkg.Card) (pureCards [][]pkg.Card, overCards []pkg.Card) {
	// 找到牌中所有的顺子（不带joker的）
	suitCards := pkg.SuitGroup(cards)

	for _, suit := range pkg.SuitQueen {
		cards = suitCards[suit]
		if len(cards) < 3 || suit == pkg.JokerSuit {
			// 该花色没有顺子，继续找下一个花色
			// 鬼牌一定不是顺子
			overCards = append(overCards, cards...)
			continue
		}

		// 对牌进行升序降序的比分，找出最优的解
		// 升序
		pureAsc, overCardsAsc := getPure(cards, true)

		// 降序
		pureDesc, overCardsDesc := getPure(cards, false)

		ascScore := pkg.CalculateScore(overCardsAsc, p.jokerVal)
		descScore := pkg.CalculateScore(overCardsDesc, p.jokerVal)

		if descScore < ascScore {
			pureCards = append(pureCards, pureDesc...)
			overCards = append(overCards, overCardsDesc...)
			continue
		}

		pureCards = append(pureCards, pureAsc...)
		overCards = append(overCards, overCardsAsc...)
	}

	return
}

func getPure(rawCards []pkg.Card, isAsc bool) (pure [][]pkg.Card, overCard []pkg.Card) {
	// 我们接收一组牌，在这组牌当中找到顺子，返回结果和剩余牌
	// 为什么要复制一份？因为切片是指针类型，如果直接操作会影响外面的数据
	// 这里的职责是找到顺子，然后返回结果和剩余其他的不是该函数的要点。
	if len(rawCards) < 3 {
		return nil, rawCards
	}

	var cards []pkg.Card
	_ = copier.Copy(&cards, &rawCards)

	if !isAsc {
		pkg.CardValue1To14(cards)
	}

	sort.Slice(cards, func(i, j int) bool {
		if isAsc {
			return cards[i].Value < cards[j].Value
		}
		return cards[i].Value > cards[j].Value
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

		if len(seq) < 3 {
			// 如果找不到顺子，当前卡牌加入剩余牌中
			overCard = append(overCard, cards[0])
			cards = cards[1:]
			continue
		}

		pure = append(pure, seq)
		cards = pkg.SliceDifferent(cards, seq)
	}

	// 将卡牌值14转换成值1
	if !isAsc {
		pkg.CardValue14To1(overCard)
		for j, p := range pure {
			sort.Slice(p, func(i, j int) bool {
				return p[i].Value < p[j].Value
			})
			var tempPure []pkg.Card
			_ = copier.Copy(&tempPure, p)
			pkg.CardValue14To1(tempPure)
			pure[j] = tempPure
		}
	}

	return pure, overCard
}

func getPureWithJoker(rawCards []pkg.Card, rawJokers []pkg.Card, jokerVal int, isAsc bool) (pureWithJoker [][]pkg.Card, overCard []pkg.Card) {
	// 我们接收一组牌，在这组牌当中找到顺子，返回结果和剩余牌
	// 为什么要复制一份？因为切片是指针类型，如果直接操作会影响外面的数据
	// 这里的职责是找到顺子，然后返回结果和剩余其他的不是该函数的要点。
	if len(rawCards) < 2 || len(rawJokers) < 1 {
		return nil, append(rawCards, rawJokers...)
	}

	var cards, jokers []pkg.Card
	_ = copier.Copy(&cards, &rawCards)
	_ = copier.Copy(&jokers, &rawJokers)

	if !isAsc {
		pkg.CardValue1To14(cards)
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

		if (currentVal-seqNextVal > 2 || currentVal-seqNextVal < -2) && len(seq) == 1 {
			// 什么都不是的情况
			overCard = append(overCard, seq[0])
			seq = seq[1:]
			seq = append(seq, cards[i])
			continue
		}

		if (currentVal-seqNextVal > 2 || currentVal-seqNextVal < -2) && len(seq) > 1 {
			// 什么都不是的情况
			overCard = append(overCard, cards[i])
			continue
		}

		if !isAsc {
			if currentVal-seqNextVal < 2*factors {
				// 什么都不是的情况
				overCard = append(overCard, seq[0])
				seq = seq[1:]
				seq = append(seq, cards[i])
				continue
			}
		} else {
			if currentVal-seqNextVal > 2*factors && len(seq) < 2 {
				// 什么都不是的情况
				overCard = append(overCard, seq[0])
				seq = seq[1:]
				seq = append(seq, cards[i])
				continue
			}
		}

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

	if len(seq) == 2 && !isUsed {
		isUsed = true
		seq = append(seq, jokers[0])
		jokers = jokers[1:]
	}

	if len(seq) >= 3 {
		pureWithJoker = append(pureWithJoker, seq)
		// 从 cards 中移除 seq
		overCard = pkg.SliceDifferent(cards, seq)
	} else {
		overCard = cards
	}
	overCard = append(overCard, jokers...)

	// 将卡牌值14转换成值1
	if !isAsc {
		pkg.CardValue14To1(overCard)
		for j, p := range pureWithJoker {
			sort.Slice(p, func(i, j int) bool {
				return p[i].Value < p[j].Value
			})
			var tempPure []pkg.Card
			_ = copier.Copy(&tempPure, p)
			pkg.CardValue14To1(tempPure)
			pureWithJoker[j] = tempPure
		}
	}

	return pureWithJoker, overCard
}
