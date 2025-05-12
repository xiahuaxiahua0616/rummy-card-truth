package internal

import "github.com/xiahua/ifonly/pkg"

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

func getJokerV2(cards []byte, joker byte) ([]byte, []byte) {
	var jokers []byte
	var filterCards []byte

	// 提取Joker
	jokerVal := joker & 0x0F
	for _, card := range cards {
		val := card & 0x0F
		if val == jokerVal {
			jokers = append(jokers, card)
		} else {
			filterCards = append(filterCards, card)
		}
	}

	return filterCards, jokers
}
