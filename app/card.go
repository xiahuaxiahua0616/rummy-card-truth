package app

type Card struct {
	Value int
	Suit  SuitVal // 0方片 1梅花 2红桃 3黑桃 4 Joker
}

func NewCard(value int, suit SuitVal) Card {
	return Card{value, suit}
}

type SuitVal int

var D SuitVal = 0
var C SuitVal = 1
var B SuitVal = 2
var A SuitVal = 3
var JokerSuit SuitVal = 4
