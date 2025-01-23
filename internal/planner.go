package internal

import (
	"github.com/jinzhu/copier"
	"rummy-card-truth/pkg"
)

type Planner struct {
	cards    []pkg.Card
	jokerVal int
}

func (p *Planner) Run() [][]pkg.Card {
	pureCards, overCards := p.getBasePure(p.cards)
	if !pkg.JudgeIsHave1Seq(pureCards) {
		return [][]pkg.Card{p.cards}
		// todo:: 返回分数
	}

	var result [][]pkg.Card
	var score int

	// 循环所有纯顺子
	for _, pure := range pureCards {
		possibleCards := pkg.GetSeqAllPossible(pure)
		for _, cards := range possibleCards {
			overCards = pkg.SliceDifferent(p.cards, cards)

			var pureOverCards []pkg.Card
			var pureWithJokerOverCards []pkg.Card
			_ = copier.Copy(&pureOverCards, &overCards)
			_ = copier.Copy(&pureWithJokerOverCards, &overCards)

			// 找纯顺子的流程
			pure2Cards, pureOverCards := p.getBasePure(pureOverCards)

			if len(pure2Cards) >= 1 {
				// 其他流程
				nextCards, pureOverCards := p.pureSetup(pureOverCards)

				// 当前结果的分数
				newScore := pkg.CalculateScore(pureOverCards, p.jokerVal)

				if score == 0 || newScore < score {
					score = pkg.CalculateScore(pureOverCards, p.jokerVal)
					result = [][]pkg.Card{}
					result = append(result, cards)
					result = append(result, pure2Cards...)
					result = append(result, nextCards...)
					result = append(result, pureOverCards)
				}
			}

			// 从剩余的牌中找pureWithJoker，同时要把他所有可能行找到

			// 找带Joker的顺子流程(全部可能）
			allPossible := p.pureWithJokerAllPossible(pureWithJokerOverCards)
			for _, pure2WithJoker := range allPossible {
				pureOverCards = pkg.SliceDifferent(pureWithJokerOverCards, pure2WithJoker)
				if len(pure2WithJoker) >= 1 {
					// 其他流程
					nextCards, pureOverCards := p.pureWithJokerSetup(pureOverCards)
					// 当前结果的分数
					newScore := pkg.CalculateScore(pureOverCards, p.jokerVal)

					if score == 0 || newScore < score {
						score = pkg.CalculateScore(pureOverCards, p.jokerVal)
						result = [][]pkg.Card{}
						result = append(result, cards)
						result = append(result, pure2WithJoker)
						result = append(result, nextCards...)
						result = append(result, pureOverCards)
					}
				}
			}

			if len(pure2Cards) == 0 && len(allPossible) == 0 {
				result = [][]pkg.Card{}
				result = append(result, cards)
				result = append(result, pureWithJokerOverCards)
			}
		}
	}

	return result
}

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

func (p *Planner) pureWithJokerSetup(rawCards []pkg.Card) (cards [][]pkg.Card, overCards []pkg.Card) {
	result1, overCards1 := p.pureWithJokerSetup1(rawCards)
	result2, overCards2 := p.pureWithJokerSetup2(rawCards)
	result3, overCards3 := p.pureWithJokerSetup3(rawCards)

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

func (p *Planner) getBasePureWithJoker(rawCards []pkg.Card) (pureWithJokerCards [][]pkg.Card, overCards []pkg.Card) {
	//fmt.Println("开始的over_card", rawCards, len(rawCards))
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

func (p *Planner) pureWithJokerSOP(rawCards []pkg.Card, jokers []pkg.Card) (pureWithJoker [][]pkg.Card, overCards []pkg.Card) {
	if len(jokers) < 1 {
		return nil, rawCards
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
		pwjAsc, ovCardsAsc := getPureWithJoker(suitCard, jokers, p.jokerVal, true)

		// 降序
		pwjDesc, ovCardsDesc := getPureWithJoker(suitCard, jokers, p.jokerVal, false)

		ascScore := pkg.CalculateScore(ovCardsAsc, p.jokerVal)
		descScore := pkg.CalculateScore(ovCardsDesc, p.jokerVal)

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
	return pureWithJoker, overCards
}

func NewPlanner(cards []pkg.Card, jokerVal int) *Planner {
	return &Planner{
		cards, jokerVal,
	}
}
