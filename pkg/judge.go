package pkg

func JudgeIsHave1Seq(cards [][]Card) bool {
	// 该函数调用应该在第一轮找顺子的时候判断
	if len(cards) >= 1 {
		return true
	}
	return false
}
