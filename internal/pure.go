package internal

import (
	"github.com/xiahua/ifonly/pkg"
	"sort"
)

func (p *Planner) pureSetup(rawCards []pkg.Card) (cards [][]pkg.Card, overCards []pkg.Card) {
	setupFuncs := [][]PureSetup{
		{p.getPureSetup, p.getPureWithJokerSetup, p.getSetSetup, p.getSetWithJokerSetup},
		{p.getSetSetup, p.getPureSetup, p.getPureWithJokerSetup, p.getSetWithJokerSetup},
		{p.getSetSetup, p.getPureSetup, p.getSetWithJokerSetup, p.getPureWithJokerSetup},
	}

	var bestCards [][]pkg.Card
	var bestOverCards []pkg.Card
	bestScore := int(^uint(0) >> 1) // 初始化为最大整数

	for _, funcs := range setupFuncs {
		result, overCards := setupChain(rawCards, funcs...)
		score := pkg.CalculateScore(overCards, p.jokerVal)
		if score < bestScore {
			bestScore = score
			bestCards = result
			bestOverCards = overCards
		}
	}

	return bestCards, bestOverCards
}

func getBasePure(cards []pkg.Card, jokerVal int) (pureCards [][]pkg.Card, overCards []pkg.Card) {
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

		ascScore := pkg.CalculateScore(overCardsAsc, jokerVal)
		descScore := pkg.CalculateScore(overCardsDesc, jokerVal)

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
	if len(rawCards) < 3 {
		return nil, rawCards
	}

	cards := make([]pkg.Card, len(rawCards))
	copy(cards, rawCards)

	if !isAsc {
		pkg.CardValue1To14(cards)
	}

	sort.Slice(cards, func(i, j int) bool {
		if isAsc {
			return cards[i].Value < cards[j].Value
		}
		return cards[i].Value > cards[j].Value
	})

	factors := 1
	if !isAsc {
		factors = -1
	}

	overCard = make([]pkg.Card, 0, len(cards))
	pure = make([][]pkg.Card, 0, len(cards)/3)

	for len(cards) >= 3 {
		seq := []pkg.Card{cards[0]}
		for i := 1; i < len(cards); i++ {
			if cards[i].Value-seq[len(seq)-1].Value == 1*factors {
				seq = append(seq, cards[i])
			} else {
				break
			}
		}

		if len(seq) < 3 {
			overCard = append(overCard, cards[0])
			cards = cards[1:]
			continue
		}

		pure = append(pure, seq)
		cards = pkg.SliceDifferent(cards, seq)
	}

	overCard = append(overCard, cards...)

	if !isAsc {
		pkg.CardValue14To1(overCard)
		for j, p := range pure {
			sort.Slice(p, func(i, j int) bool {
				return p[i].Value < p[j].Value
			})
			tempPure := make([]pkg.Card, len(p))
			copy(tempPure, p)
			pkg.CardValue14To1(tempPure)
			pure[j] = tempPure
		}
	}

	return pure, overCard
}

func getPureWithJoker(rawCards []pkg.Card, rawJokers []pkg.Card, jokerVal int, isAsc bool) (pureWithJoker [][]pkg.Card, overCard []pkg.Card) {
	if len(rawCards) < 2 || len(rawJokers) < 1 {
		return nil, append(rawCards, rawJokers...)
	}

	cards := make([]pkg.Card, len(rawCards))
	copy(cards, rawCards)
	jokers := make([]pkg.Card, len(rawJokers))
	copy(jokers, rawJokers)

	if !isAsc {
		pkg.CardValue1To14(cards)
	}

	sort.Slice(cards, func(i, j int) bool {
		if isAsc {
			return cards[i].Value < cards[j].Value
		}
		return cards[i].Value > cards[j].Value
	})

	factors := 1
	if !isAsc {
		factors = -1
	}

	overCard = make([]pkg.Card, 0, len(cards))
	pureWithJoker = make([][]pkg.Card, 0, len(cards)/3)

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
		overCard = pkg.SliceDifferent(cards, seq)
	} else {
		overCard = cards
	}
	overCard = append(overCard, jokers...)

	if !isAsc {
		pkg.CardValue14To1(overCard)
		for j, p := range pureWithJoker {
			sort.Slice(p, func(i, j int) bool {
				return p[i].Value < p[j].Value
			})
			tempPure := make([]pkg.Card, len(p))
			copy(tempPure, p)
			pkg.CardValue14To1(tempPure)
			pureWithJoker[j] = tempPure
		}
	}

	return pureWithJoker, overCard
}
