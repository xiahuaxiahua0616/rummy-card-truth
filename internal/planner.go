package internal

import (
	"fmt"
	"sync"

	"github.com/jinzhu/copier"
	"github.com/xiahua/ifonly/pkg"
	"github.com/xiahuaxiahua0616/ifonlyutils/ifonlyutils"
)

type Planner struct {
	cards    []pkg.Card
	jokerVal int
}

func (p *Planner) Run() [][]pkg.Card {
	pureCards, _ := getBasePure(p.cards, p.jokerVal)
	if !pkg.JudgeIsHave1Seq(pureCards) {
		return [][]pkg.Card{p.cards}
	}

	var result [][]pkg.Card
	var score int

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, pure := range pureCards {
		possibleCards := pkg.GetSeqAllPossible(pure)
		for _, cards := range possibleCards {
			wg.Add(1)
			go func(cards []pkg.Card) {
				defer wg.Done()
				overCards := pkg.SliceDifferent(p.cards, cards)

				var pureOverCards []pkg.Card
				var pureWithJokerOverCards []pkg.Card
				_ = copier.Copy(&pureOverCards, &overCards)
				_ = copier.Copy(&pureWithJokerOverCards, &overCards)

				pure2Cards, pureOverCards := getBasePure(pureOverCards, p.jokerVal)

				if len(pure2Cards) >= 1 {
					nextCards, pureOverCards := p.pureSetup(pureOverCards)
					newScore := pkg.CalculateScore(pureOverCards, p.jokerVal)

					mu.Lock()
					if score == 0 || newScore < score {
						score = newScore
						result = [][]pkg.Card{}
						result = append(result, cards)
						result = append(result, pure2Cards...)
						result = append(result, nextCards...)
						result = append(result, pureOverCards)
					}
					mu.Unlock()
				}

				allPossible := p.pureWithJokerAllPossible(pureWithJokerOverCards)
				for _, pure2WithJoker := range allPossible {
					pureOverCards = pkg.SliceDifferent(pureWithJokerOverCards, pure2WithJoker)
					if len(pure2WithJoker) >= 1 {
						nextCards, pureOverCards := p.pureWithJokerSetup(pureOverCards)
						newScore := pkg.CalculateScore(pureOverCards, p.jokerVal)

						mu.Lock()
						if score == 0 || newScore < score {
							score = newScore
							result = [][]pkg.Card{}
							result = append(result, cards)
							result = append(result, pure2WithJoker)
							result = append(result, nextCards...)
							result = append(result, pureOverCards)
						}
						mu.Unlock()
					}
				}

				if len(pure2Cards) == 0 && len(allPossible) == 0 {
					mu.Lock()
					result = [][]pkg.Card{}
					result = append(result, cards)
					result = append(result, pureWithJokerOverCards)
					mu.Unlock()
				}
			}(cards)
		}
	}

	wg.Wait()
	return result
}

func NewPlanner(cards []pkg.Card, jokerVal int) *Planner {
	return &Planner{
		cards, jokerVal,
	}
}

type PlannerV2 struct {
	cards []byte
	joker byte
}

func (p *PlannerV2) Run() {
	fmt.Println(p.cards)
	straights, _ := GetStraight(p.cards, p.joker)

	// for _, straight := range straights {

	// }
	fmt.Println(straights)

	fmt.Println(getStraightAllPossible([]byte{1, 2, 3, 4}))
}

func NewPlannerV2(cards []byte, joker byte) *PlannerV2 {
	return &PlannerV2{
		cards: cards,
		joker: joker,
	}
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
				result = append(result, straight[start:end])
			}
		}
	}
	return
}
