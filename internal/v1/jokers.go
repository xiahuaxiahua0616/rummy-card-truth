package v1

func getJokerV2(cards []byte, joker byte) ([]byte, []byte) {
	var jokers []byte
	var filterCards []byte

	jokerLow := joker & 0x0F

	// 提取Joker
	for _, card := range cards {
		if card >= byte(0x4e) || (card&0x0F) == jokerLow {
			jokers = append(jokers, card)
		} else {
			filterCards = append(filterCards, card)
		}
	}

	return filterCards, jokers
}
