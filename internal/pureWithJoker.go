package internal

import "rummy-card-truth/pkg"

type SetupWithJoker func(rawCards []pkg.Card, jokerVal int) (cards []pkg.Card, overCards []pkg.Card)

func (p *Planner) pureWithJokerSetup(rawCards []pkg.Card) (cards [][]pkg.Card, overCards []pkg.Card) {
	setupFuncs := [][]PureSetup{
		{p.getPureWithJokerSetup, p.getSetSetup, p.getSetWithJokerSetup},
		{p.getSetSetup, p.getPureWithJokerSetup, p.getSetWithJokerSetup},
		{p.getSetSetup, p.getSetWithJokerSetup, p.getPureWithJokerSetup},
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
