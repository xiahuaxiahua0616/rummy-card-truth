package v1

import (
	"github.com/xiahuaxiahua0616/ifonlyutils/ifonlyutils"
)

type Planner struct {
	cards []byte
	joker byte
}

func (p *Planner) Run(data *[][]byte) {
	// 第一个顺子，如果一个都没有直接返回结果
	firstStraights, _ := GetStraight(p.cards, p.joker)
	var score int
	var result [][]byte
	for _, firstStraight := range firstStraights {
		// 获取全部可能性
		dataVersion := [][]byte{
			firstStraight,
			ifonlyutils.Conv1to14(firstStraight),
		}
		var datas [][]byte
		for _, dv := range dataVersion {
			datas = append(datas, getStraightAllPossible(dv)...)
		}

		for _, data := range datas {
			leftover := SliceDiffWithDup(p.cards, data)
			leftoverV2 := SliceDiffWithDup(p.cards, data)
			// 找第二个顺子。
			secondStraights, leftover := GetStraight(leftover, p.joker)
			if len(secondStraights) >= 1 {
				nextCards, leftover := p.PlannerChain(leftover)
				tempScore := ifonlyutils.CalcScore(leftover, p.joker)
				if score == 0 || tempScore < score {
					score = tempScore
					result = [][]byte{data}
					result = append(result, secondStraights...)
					result = append(result, nextCards...)
					result = append(result, leftover)
				}
			}

			// 找带joker的顺子
			datasWithJoker := getStraightWithJokerAllPossible(leftoverV2, p.joker)

			for _, dataWithJoker := range datasWithJoker {
				leftoverV2 := SliceDiffWithDup(leftoverV2, dataWithJoker)
				nextCards, leftoverV2 := p.PlannerJokerChain(leftoverV2)
				tempScore := ifonlyutils.CalcScore(leftoverV2, p.joker)
				if score == 0 || tempScore < score {
					score = tempScore
					result = [][]byte{data}
					result = append(result, dataWithJoker)
					result = append(result, nextCards...)
					result = append(result, leftoverV2)
				}
			}
			if len(secondStraights) == 0 && len(datasWithJoker) == 0 {
				result = [][]byte{data}
				result = append(result, leftover)
			}
		}
	}
	*data = result
}

type PlannerChainFunc func(cards []byte) (result [][]byte, leftover []byte)

func (p *Planner) StraightChain(cards []byte) (result [][]byte, leftover []byte) {
	return GetStraight(cards, p.joker)
}
func (p *Planner) StraightWithJokerChain(cards []byte) (result [][]byte, leftover []byte) {
	return GetStraightWithJoker(cards, p.joker)
}
func (p *Planner) SetChain(cards []byte) (result [][]byte, leftover []byte) {
	return GetSetV2(cards)
}
func (p *Planner) SetChainWithJoker(cards []byte) (result [][]byte, leftover []byte) {
	return GetSetWithJokerV2(cards, p.joker)
}

// 执行链
func (p *Planner) PlannerChain(cards []byte) (result [][]byte, leftover []byte) {
	chainFuncs := [][]PlannerChainFunc{
		{p.StraightChain, p.StraightWithJokerChain, p.SetChain, p.SetChainWithJoker},
		{p.SetChain, p.StraightChain, p.StraightWithJokerChain, p.SetChainWithJoker},
		{p.SetChain, p.StraightChain, p.SetChainWithJoker, p.StraightWithJokerChain},
	}
	var score = int(^uint(0) >> 1)

	for _, chainFunc := range chainFuncs {
		var tempResult [][]byte
		var tempLeftover []byte = make([]byte, len(cards))
		var tempScore int

		copy(tempLeftover, cards)
		for _, funcs := range chainFunc {
			var tempChainResult [][]byte
			tempChainResult, tempLeftover = funcs(tempLeftover)
			if tempChainResult != nil {
				tempResult = append(tempResult, tempChainResult...)
			}
		}
		tempScore = ifonlyutils.CalcScore(tempLeftover, p.joker)
		if tempScore < score {
			result = tempResult
			leftover = tempLeftover
			score = tempScore
		}
	}
	return
}

// 带Joker的执行链
func (p *Planner) PlannerJokerChain(cards []byte) (result [][]byte, leftover []byte) {
	chainFuncs := [][]PlannerChainFunc{
		{p.StraightWithJokerChain, p.SetChain, p.SetChainWithJoker},
		{p.SetChain, p.StraightWithJokerChain, p.SetChainWithJoker},
		{p.SetChain, p.SetChainWithJoker, p.StraightWithJokerChain},
	}
	var score = int(^uint(0) >> 1)

	for _, chainFunc := range chainFuncs {
		var tempResult [][]byte
		var tempLeftover []byte = make([]byte, len(cards))
		var tempScore int

		copy(tempLeftover, cards)

		for _, funcs := range chainFunc {
			var tempChainResult [][]byte
			tempChainResult, tempLeftover = funcs(tempLeftover)
			if tempChainResult != nil {
				tempResult = append(tempResult, tempChainResult...)
			}
		}

		tempScore = ifonlyutils.CalcScore(tempLeftover, p.joker)
		if tempScore < score {
			result = tempResult
			leftover = tempLeftover
			score = tempScore
		}
	}
	return
}

func NewPlannerV2(cards []byte, joker byte) *Planner {
	return &Planner{
		cards: cards,
		joker: joker,
	}
}

// SliceDiffWithDup 要被裁剪的放在前面，裁剪数
func SliceDiffWithDup(a, b []byte) []byte {
	countB := make(map[byte]int)
	for _, val := range b {
		countB[val]++
	}

	diff := []byte{}
	for _, val := range a {
		if countB[val] > 0 {
			countB[val]--
		} else {
			diff = append(diff, val)
		}
	}
	return diff
}

func getStraightAllPossible(straight []byte) (result [][]byte) {
	if len(straight) == 3 {
		return [][]byte{
			straight,
		}
	}

	for start := range straight {
		for end := start + 2; end <= len(straight); end++ {
			if ifonlyutils.IsContinuous(straight[start:end]) {
				result = append(result, ifonlyutils.Conv14to1(straight[start:end]))
			}
		}
	}
	return
}

func getStraightWithJokerAllPossible(datas []byte, joker byte) (cards [][]byte) {
	_, jokers := getJokerV2(datas, joker)
	for {
		straights, _ := GetStraightWithJoker(datas, joker)
		if straights == nil {
			break
		}

		for _, straight := range straights {
			if len(straight) > 3 {
				t := SliceDiffWithDup(straight, jokers)
				for i := 1; i < len(t); i++ {
					tempOvCard := make([]byte, len(datas))
					copy(tempOvCard, datas)

					tempOvCard = SliceDiffWithDup(tempOvCard, t[:i])
					otherStright, _ := GetStraightWithJoker(tempOvCard, joker)
					if otherStright == nil {
						continue
					}
					cards = append(cards, otherStright...)
				}
				for i := 1; i < len(t); i++ {
					tempOvCard := make([]byte, len(datas))
					copy(tempOvCard, datas)
					tempOvCard = SliceDiffWithDup(tempOvCard, t[i:])
					otherStright, _ := GetStraightWithJoker(tempOvCard, joker)
					if otherStright == nil {
						continue
					}
					cards = append(cards, otherStright...)
				}
			}
			cards = append(cards, straight)
			t := SliceDiffWithDup(straight, jokers)
			datas = SliceDiffWithDup(datas, t)
		}
	}
	return
}
