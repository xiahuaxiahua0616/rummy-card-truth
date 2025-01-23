package internal

import "rummy-card-truth/pkg"

func (p *Planner) pureSetup1(rawCards []pkg.Card) (cards [][]pkg.Card, overCards []pkg.Card) {
	// 找顺子
	pure, overCards := p.getBasePure(rawCards)

	overCards, jokers := getJokers(overCards, p.jokerVal)

	// 找带Joker的顺子
	pureWithJoker, overCards := p.pureWithJokerSOP(overCards, jokers)

	setCards, overCards := getSet(overCards)

	setWithJoker, overCards := getSetWithJoker(overCards, p.jokerVal)

	cards = append(pure, pureWithJoker...)
	cards = append(cards, setCards...)
	cards = append(cards, setWithJoker...)
	for _, pc := range pureWithJoker {
		jokers = pkg.SliceDifferent(jokers, pc)
	}
	overCards = append(overCards, jokers...)

	return cards, overCards
}

func (p *Planner) pureSetup2(rawCards []pkg.Card) (cards [][]pkg.Card, overCards []pkg.Card) {
	setCards, overCards := getSet(overCards)

	// 找顺子
	pure, overCards := p.getBasePure(rawCards)

	overCards, jokers := getJokers(overCards, p.jokerVal)

	// 找带Joker的顺子
	pureWithJoker, overCards := p.pureWithJokerSOP(overCards, jokers)

	setWithJoker, overCards := getSetWithJoker(overCards, p.jokerVal)

	cards = append(pure, pureWithJoker...)
	cards = append(cards, setCards...)
	cards = append(cards, setWithJoker...)
	for _, pc := range pureWithJoker {
		jokers = pkg.SliceDifferent(jokers, pc)
	}
	overCards = append(overCards, jokers...)

	return cards, overCards
}

func (p *Planner) pureSetup3(rawCards []pkg.Card) (cards [][]pkg.Card, overCards []pkg.Card) {
	setCards, overCards := getSet(overCards)

	// 找顺子
	pure, overCards := p.getBasePure(rawCards)

	setWithJoker, overCards := getSetWithJoker(overCards, p.jokerVal)

	overCards, jokers := getJokers(overCards, p.jokerVal)

	// 找带Joker的顺子
	pureWithJoker, overCards := p.pureWithJokerSOP(overCards, jokers)

	cards = append(pure, pureWithJoker...)
	cards = append(cards, setCards...)
	cards = append(cards, setWithJoker...)
	for _, pc := range pureWithJoker {
		jokers = pkg.SliceDifferent(jokers, pc)
	}
	overCards = append(overCards, jokers...)

	return cards, overCards
}
