package internal

import (
	"github.com/xiahua/ifonly/pkg"

	"github.com/jinzhu/copier"
)

type PureSetup func(rawCards []pkg.Card) (cards [][]pkg.Card, overCards []pkg.Card)

func (p *Planner) getPureSetup(rawCards []pkg.Card) (cards [][]pkg.Card, overCards []pkg.Card) {
	return p.processSetup(rawCards, func(cards []pkg.Card) ([][]pkg.Card, []pkg.Card) {
		return getBasePure(cards, p.jokerVal)
	})
}

func (p *Planner) getPureWithJokerSetup(rawCards []pkg.Card) (cards [][]pkg.Card, overCards []pkg.Card) {
	return p.processSetup(rawCards, func(cards []pkg.Card) ([][]pkg.Card, []pkg.Card) {
		var jokers []pkg.Card
		overCards, jokers = getJokers(cards, p.jokerVal)
		pureWithJoker, overCards, jokers := pureWithJokerSOP(overCards, jokers, p.jokerVal)
		overCards = append(overCards, jokers...)
		return pureWithJoker, overCards
	})
}

func (p *Planner) getSetSetup(rawCards []pkg.Card) (cards [][]pkg.Card, overCards []pkg.Card) {
	return p.processSetup(rawCards, func(cards []pkg.Card) ([][]pkg.Card, []pkg.Card) {
		return getSet(cards)
	})
}

func (p *Planner) getSetWithJokerSetup(rawCards []pkg.Card) (cards [][]pkg.Card, overCards []pkg.Card) {
	return p.processSetup(rawCards, func(cards []pkg.Card) ([][]pkg.Card, []pkg.Card) {
		return getSetWithJoker(cards, p.jokerVal)
	})
}

func (p *Planner) processSetup(rawCards []pkg.Card, processFunc func([]pkg.Card) ([][]pkg.Card, []pkg.Card)) (cards [][]pkg.Card, overCards []pkg.Card) {
	processedCards, remainingCards := processFunc(rawCards)
	cards = append(cards, processedCards...)
	return cards, remainingCards
}

func setupChain(rawCards []pkg.Card, setups ...PureSetup) (cards [][]pkg.Card, overCards []pkg.Card) {
	for _, setup := range setups {
		var processedCards [][]pkg.Card
		processedCards, rawCards = setup(rawCards)
		if len(processedCards) > 0 {
			cards = append(cards, processedCards...)
		}
	}
	overCards = rawCards
	return cards, overCards
}

func (p *Planner) pureWithJokerAllPossible(rawCards []pkg.Card) (cards [][]pkg.Card) {
	_, jokers := getJokers(rawCards, p.jokerVal)

	for {
		pure2WithJoker, _ := p.getBasePureWithJoker(rawCards)
		if pure2WithJoker == nil {
			break
		}

		for _, pc := range pure2WithJoker {
			if len(pc) > 3 {
				// 找出所有可能性
				// 顺子 9 11 12 joker，可能 9是用来组成三条的
				t := pkg.SliceDifferent(pc, jokers)
				for i := 1; i < len(t); i++ {
					var tempOvCard []pkg.Card
					_ = copier.Copy(&tempOvCard, rawCards)
					tempOvCard = pkg.SliceDifferent(tempOvCard, t[:i])
					otherPure, _ := p.getBasePureWithJoker(tempOvCard)
					if otherPure == nil {
						continue
					}
					cards = append(cards, otherPure...)

				}
				for i := 1; i < len(t); i++ {
					var tempOvCard []pkg.Card
					_ = copier.Copy(&tempOvCard, rawCards)
					tempOvCard = pkg.SliceDifferent(tempOvCard, t[i:])
					otherPure, _ := p.getBasePureWithJoker(tempOvCard)
					if otherPure == nil {
						continue
					}
					cards = append(cards, otherPure...)
				}
			}
			cards = append(cards, pc)
			t := pkg.SliceDifferent(pc, jokers)
			rawCards = pkg.SliceDifferent(rawCards, t)
		}
	}

	return cards
}

func (p *Planner) getBasePureWithJoker(rawCards []pkg.Card) (pureWithJokerCards [][]pkg.Card, overCards []pkg.Card) {
	overCards, jokers := getJokers(rawCards, p.jokerVal)
	suitCards := pkg.SuitGroup(overCards)

	overCards = nil
	for _, suit := range pkg.SuitQueen {
		cards := suitCards[suit]
		if len(cards) == 0 {
			continue
		}
		if len(jokers) < 1 || suit == pkg.JokerSuit {
			// 该花色没有顺子，继续找下一个花色
			// 鬼牌一定不是顺子
			overCards = append(overCards, cards...)
			continue
		}

		result1, overCards1 := getPureWithJoker(cards, jokers, p.jokerVal, true)
		result2, overCards2 := getPureWithJoker(cards, jokers, p.jokerVal, false)

		score1, score2 := pkg.CalculateScore(overCards1, p.jokerVal), pkg.CalculateScore(overCards2, p.jokerVal)

		if score1 < score2 {
			pureWithJokerCards, overCards = appendResult(pureWithJokerCards, overCards, result1, overCards1, jokers)
		} else {
			pureWithJokerCards, overCards = appendResult(pureWithJokerCards, overCards, result2, overCards2, jokers)
		}
	}

	overCards = append(overCards, jokers...)

	return pureWithJokerCards, overCards
}

func pureWithJokerSOP(rawCards []pkg.Card, jokers []pkg.Card, jokerVal int) (pureWithJoker [][]pkg.Card, overCards, ovJoker []pkg.Card) {
	if len(jokers) < 1 {
		return nil, rawCards, jokers
	}
	suitCards := pkg.SuitGroup(rawCards)
	overCards = nil
	for _, suit := range pkg.SuitQueen {
		suitCard := suitCards[suit]
		if len(suitCard) == 0 {
			continue
		}
		if len(jokers) < 1 || suit == pkg.JokerSuit {
			// 该花色没有顺子，继续找下一个花色
			// 鬼牌一定不是顺子
			overCards = append(overCards, suitCard...)
			continue
		}

		// 找带Joker的顺子
		pwjAsc, ovCardsAsc := getPureWithJoker(suitCard, jokers, jokerVal, true)

		// 降序
		pwjDesc, ovCardsDesc := getPureWithJoker(suitCard, jokers, jokerVal, false)

		ascScore := pkg.CalculateScore(ovCardsAsc, jokerVal)
		descScore := pkg.CalculateScore(ovCardsDesc, jokerVal)

		if descScore < ascScore {
			for _, pw := range pwjDesc {
				jokers = pkg.SliceDifferent(jokers, pw)
			}
			ovCardsDesc = pkg.SliceDifferent(ovCardsDesc, jokers)
			pureWithJoker = append(pureWithJoker, pwjDesc...)
			overCards = append(overCards, ovCardsDesc...)
			continue
		}

		for _, pw := range pwjAsc {
			jokers = pkg.SliceDifferent(jokers, pw)
		}
		ovCardsAsc = pkg.SliceDifferent(ovCardsAsc, jokers)
		pureWithJoker = append(pureWithJoker, pwjAsc...)
		overCards = append(overCards, ovCardsAsc...)
	}

	return pureWithJoker, overCards, jokers
}

func appendResult(pureWithJokerCards [][]pkg.Card, overCards []pkg.Card, result [][]pkg.Card, overCardsPart []pkg.Card, jokers []pkg.Card) ([][]pkg.Card, []pkg.Card) {
	overCardsPart = pkg.SliceDifferent(overCardsPart, jokers)
	pureWithJokerCards = append(pureWithJokerCards, result...)
	overCards = append(overCards, overCardsPart...)
	for _, pure := range result {
		jokers = pkg.SliceDifferent(jokers, pure)
	}
	return pureWithJokerCards, overCards
}

// type StraightSetup func(datas []byte) (cards [][]byte, left []byte)

// func (p *PlannerV2) getStrightSetup(datas []byte) (cards [][]byte, leftover []byte) {
// 	return p.processSetup(datas, func(b []byte) ([][]byte, []byte) {
// 		return GetStraight(datas, p.joker)
// 	})
// }

// func (p *PlannerV2) getStraightWithJokerSetup(datas []byte) (cards [][]byte, leftover []byte) {
// 	return p.processSetup(datas, func(b []byte) ([][]byte, []byte) {
// 		return GetStraightWithJoker(datas, p.joker)
// 	})
// }

// func (p *PlannerV2) getSetSetup(datas []byte) (cards [][]byte, leftover []byte) {
// 	return GetSetV2(datas)
// }

// func (p *PlannerV2) getSetWithJokerSetup(datas []byte) (cards [][]byte, leftover []byte) {
// 	return GetSetWithJokerV2(datas, p.joker)
// }

// func (p *PlannerV2) processSetup(datas []byte, processFunc func([]byte) ([][]byte, []byte)) ([][]byte, []byte) {
// 	return processFunc(datas)
// }

// func (p *PlannerV2) straightSetup(datas []byte) ([][]byte, []byte) {
// 	setupFunc := [][]StraightSetup{
// 		{p.getStrightSetup, p.getStraightWithJokerSetup, p.getSetSetup, p.getSetWithJokerSetup},
// 		{p.getSetSetup, p.getStrightSetup, p.getStraightWithJokerSetup, p.getSetWithJokerSetup},
// 		{p.getSetSetup, p.getStrightSetup, p.getSetWithJokerSetup, p.getStraightWithJokerSetup},
// 	}

// 	var result [][]byte
// 	var leftover []byte
// 	var score = int(^uint(0) >> 1)

// 	for _, funcs := range setupFunc {
// 		result2, leftover2 := setupChainV2(datas, funcs...)
// 		score2 := ifonlyutils.CalcScore(leftover, p.joker)
// 		if score < score2 {
// 			score = score2
// 			result = result2
// 			leftover = leftover2
// 		}
// 	}
// 	return result, leftover
// }

// func setupChainV2(rawCards []byte, setups ...StraightSetup) (cards [][]byte, overCards []byte) {
// 	for _, setup := range setups {
// 		var processedCards [][]byte
// 		processedCards, rawCards = setup(rawCards)
// 		if len(processedCards) > 0 {
// 			cards = append(cards, processedCards...)
// 		}
// 	}
// 	overCards = rawCards
// 	return cards, overCards
// }

// type SetupWithJokerV2 func(datas []byte) ([][]byte, []byte)

// func (p *PlannerV2) straightWithJokerSetup(datas []byte) ([][]byte, []byte) {
// 	setupFuncs := [][]StraightSetup{
// 		{p.getStraightWithJokerSetup, p.getSetSetup, p.getSetWithJokerSetup},
// 		{p.getSetSetup, p.getStraightWithJokerSetup, p.getSetWithJokerSetup},
// 		{p.getSetSetup, p.getSetWithJokerSetup, p.getStraightWithJokerSetup},
// 	}

// 	var result [][]byte
// 	var leftover []byte
// 	var score = int(^uint(0) >> 1)

// 	for _, funcs := range setupFuncs {
// 		result2, leftover2 := setupChainV2(datas, funcs...)
// 		score2 := ifonlyutils.CalcScore(leftover2, p.joker)
// 		if score2 < score {
// 			result = result2
// 			score = score2
// 			leftover = leftover2
// 		}
// 	}
// 	return result, leftover
// }
