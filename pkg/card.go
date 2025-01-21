package pkg

type Card struct {
	Value int
	Suit  SuitVal // 0方片 1梅花 2红桃 3黑桃 4 Joker
}

func NewCard(value int, suit SuitVal) Card {
	return Card{value, suit}
}

type SuitVal string

var D SuitVal = "D"
var C SuitVal = "C"
var B SuitVal = "B"
var A SuitVal = "A"
var JokerSuit SuitVal = "Joker"
