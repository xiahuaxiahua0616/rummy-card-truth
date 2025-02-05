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
	pureCards, overCards := getBasePure(p.cards, p.jokerVal)
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
			pure2Cards, pureOverCards := getBasePure(pureOverCards, p.jokerVal)

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

func NewPlanner(cards []pkg.Card, jokerVal int) *Planner {
	return &Planner{
		cards, jokerVal,
	}
}
