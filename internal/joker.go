package internal

import "rummy-card-truth/pkg"

// getJokers 获取卡牌中的Joker并且返回剩余牌和joker牌
func getJokers(rawCards []pkg.Card, jokerVal int) (cards []pkg.Card, jokers []pkg.Card) {
	cards = make([]pkg.Card, 0, len(rawCards))
	jokers = make([]pkg.Card, 0, len(rawCards))
	for _, card := range rawCards {
		if card.Value == jokerVal || card.Suit == pkg.JokerSuit {
			jokers = append(jokers, card)
		} else {
			cards = append(cards, card)
		}
	}
	return cards, jokers
}
