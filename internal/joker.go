package internal

import "rummy-card-truth/pkg"

// GetJokers 获取卡牌中的Joker并且返回剩余牌和joker牌
func GetJokers(rawCards []pkg.Card, jokerVal int) (cards []pkg.Card, jokers []pkg.Card) {
	for i := 0; i < len(rawCards); i++ {
		if rawCards[i].Value == jokerVal || rawCards[i].Suit == pkg.JokerSuit {
			jokers = append(jokers, rawCards[i])
		} else {
			cards = append(cards, rawCards[i])
		}
	}
	return cards, jokers
}
