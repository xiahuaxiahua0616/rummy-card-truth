package internal

import "rummy-card-truth/pkg"

func (p *Planner) pureWithJokerSetup1(rawCards []pkg.Card) (cards [][]pkg.Card, overCards []pkg.Card) {
	overCards, jokers := getJokers(rawCards, p.jokerVal)

	// 找带Joker的顺子
	pureWithJoker, overCards := p.pureWithJokerSOP(overCards, jokers)
	for _, pc := range pureWithJoker {
		jokers = pkg.SliceDifferent(jokers, pc)
	}

	setCards, overCards := getSet(overCards)

	setWithJoker, overCards := getSetWithJoker(overCards, p.jokerVal)

	cards = append(cards, pureWithJoker...)
	cards = append(cards, setCards...)
	cards = append(cards, setWithJoker...)

	var result []pkg.Card
	for _, c := range cards {
		result = append(result, c...)
	}

	overCards = append(overCards, jokers...)
	result = append(result, overCards...)
	return cards, overCards
}

func (p *Planner) pureWithJokerSetup2(rawCards []pkg.Card) (cards [][]pkg.Card, overCards []pkg.Card) {
	setCards, overCards := getSet(rawCards)

	overCards, jokers := getJokers(overCards, p.jokerVal)

	// 找带Joker的顺子
	pureWithJoker, overCards := p.pureWithJokerSOP(overCards, jokers)
	setWithJoker, overCards := getSetWithJoker(overCards, p.jokerVal)

	cards = append(cards, pureWithJoker...)
	cards = append(cards, setCards...)
	cards = append(cards, setWithJoker...)
	var result []pkg.Card
	for _, c := range cards {
		result = append(result, c...)
	}
	for _, pc := range pureWithJoker {
		jokers = pkg.SliceDifferent(jokers, pc)
	}
	overCards = append(overCards, jokers...)

	return cards, overCards
}

func (p *Planner) pureWithJokerSetup3(rawCards []pkg.Card) (cards [][]pkg.Card, overCards []pkg.Card) {
	setCards, overCards := getSet(rawCards)
	setWithJoker, overCards := getSetWithJoker(overCards, p.jokerVal)

	overCards, jokers := getJokers(overCards, p.jokerVal)

	// 找带Joker的顺子
	pureWithJoker, overCards := p.pureWithJokerSOP(overCards, jokers)

	cards = append(cards, pureWithJoker...)
	cards = append(cards, setCards...)
	cards = append(cards, setWithJoker...)
	var result []pkg.Card
	for _, c := range cards {
		result = append(result, c...)
	}
	for _, pc := range pureWithJoker {
		jokers = pkg.SliceDifferent(jokers, pc)
	}
	overCards = append(overCards, jokers...)

	return cards, overCards
}
