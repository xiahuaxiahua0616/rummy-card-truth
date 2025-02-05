package internal

import "rummy-card-truth/pkg"

type SetupWithJoker func(rawCards []pkg.Card, jokerVal int) (cards []pkg.Card, overCards []pkg.Card)

func (p *Planner) pureWithJokerSetup(rawCards []pkg.Card) (cards [][]pkg.Card, overCards []pkg.Card) {
	result1, overCards1 := setupChain(rawCards, p.jokerVal, getPureWithJokerSetup, getSetSetup, getSetWithJokerSetup)
	result2, overCards2 := setupChain(rawCards, p.jokerVal, getSetSetup, getPureWithJokerSetup, getSetWithJokerSetup)
	result3, overCards3 := setupChain(rawCards, p.jokerVal, getSetSetup, getSetWithJokerSetup, getPureWithJokerSetup)

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
