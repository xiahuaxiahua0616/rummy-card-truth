package v1

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
