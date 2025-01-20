package pkg

var SuitQueen = []SuitVal{
	A, B, C, D, JokerSuit,
}

// SuitGroup 花色分组
func SuitGroup(cards []Card) map[SuitVal][]Card {
	suitCards := map[SuitVal][]Card{}
	for _, card := range cards {
		suitCards[card.Suit] = append(suitCards[card.Suit], card)
	}
	return suitCards
}
