package app

type Card struct {
	value int
	suit  int // 0方片 1梅花 2红桃 3黑桃 4 Joker
}

func (c *Card) Value() int {
	return c.value
}

func (c *Card) Suit() int {
	return c.suit
}

func NewCard(value, suit int) Card {
	return Card{value, suit}
}
