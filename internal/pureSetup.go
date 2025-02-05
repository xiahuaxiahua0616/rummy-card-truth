package internal

import (
	"github.com/jinzhu/copier"
	"rummy-card-truth/pkg"
)

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

type PureSetup func(rawCards []pkg.Card, jokerVal int) (cards [][]pkg.Card, overCards []pkg.Card)

func getPureSetup(rawCards []pkg.Card, jokerVal int) (cards [][]pkg.Card, overCards []pkg.Card) {
	pure, overCards := getBasePure(rawCards, jokerVal)

	cards = append(cards, pure...)

	return cards, overCards
}

func getPureWithJokerSetup(rawCards []pkg.Card, jokerVal int) (cards [][]pkg.Card, overCards []pkg.Card) {
	var jokers []pkg.Card
	overCards, jokers = getJokers(rawCards, jokerVal)

	// 找带Joker的顺子
	pureWithJoker, overCards, jokers := pureWithJokerSOP(overCards, jokers, jokerVal)

	cards = append(cards, pureWithJoker...)

	overCards = append(overCards, jokers...)
	return cards, overCards
}

func getSetSetup(rawCards []pkg.Card, jokerVal int) (cards [][]pkg.Card, overCards []pkg.Card) {
	setCards, overCards := getSet(rawCards)
	cards = append(cards, setCards...)

	return cards, overCards
}

func getSetWithJokerSetup(rawCards []pkg.Card, jokerVal int) (cards [][]pkg.Card, overCards []pkg.Card) {
	setWithJoker, overCards := getSetWithJoker(rawCards, jokerVal)

	cards = append(cards, setWithJoker...)

	return cards, overCards
}

func setupChain(rawCards []pkg.Card, jokerVal int, setups ...PureSetup) (cards [][]pkg.Card, overCards []pkg.Card) {
	for _, setup := range setups {
		var card [][]pkg.Card
		card, rawCards = setup(rawCards, jokerVal)
		if len(card) > 0 {
			cards = append(cards, card...)
		}
	}
	return cards, rawCards
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

		// getPureWithJoker 该函数会把joker再次返回，这个是为了解决重复joker的问题
		score1, score2 := pkg.CalculateScore(overCards1, p.jokerVal), pkg.CalculateScore(overCards2, p.jokerVal)
		if score1 < score2 {
			overCards1 = pkg.SliceDifferent(overCards1, jokers)
			pureWithJokerCards = append(pureWithJokerCards, result1...)
			overCards = append(overCards, overCards1...)
			for _, pure := range result1 {
				jokers = pkg.SliceDifferent(jokers, pure)
			}
		} else {
			overCards2 = pkg.SliceDifferent(overCards2, jokers)
			pureWithJokerCards = append(pureWithJokerCards, result2...)
			overCards = append(overCards, overCards2...)

			for _, pure := range result2 {
				jokers = pkg.SliceDifferent(jokers, pure)
			}
		}
	}

	overCards = append(overCards, jokers...)

	return pureWithJokerCards, overCards
}
