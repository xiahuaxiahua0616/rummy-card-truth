package internal

import (
	"slices"

	"github.com/xiahua/ifonly/pkg"
	"github.com/xiahuaxiahua0616/ifonlyutils"
)

func GetStraight(cards []byte, joker byte) (straights [][]byte, overCards []byte) {
	// 找到牌中所有的顺子（不带joker的）
	groupBySuitCards := ifonlyutils.GroupBySuit(cards)
	// 根据方片-梅花-红桃-黑桃-王的顺序进行遍历
	for i, cards := range groupBySuitCards {
		if len(cards) < 3 || i == pkg.JokerSuitV2 {
			// 该花色没有顺子，继续找下一个花色
			// 鬼牌一定不是顺子
			continue
		}
		// 排序升序
		slices.Sort(cards)

		// 去重数据
		// duplicates
		cards, _ := ifonlyutils.UniqueAndDuplicates(cards)

		// tempBytes 临时存储使用

		var straight []byte
		for {
			currentCards := make([]byte, len(cards))
			copy(currentCards, cards)
			slices.Sort(cards)

			straight, cards = ifonlyutils.GetStraightOnyByOne(cards)
			if len(straight) > 0 {
				straights = append(straights, straight)
			}

			if len(currentCards) == len(cards) {
				// fmt.Println(currentCards, cards)
				break
			}

			if len(cards) < 3 {
				break
			}
			// fmt.Println("执行...", cards)
		}

		// fmt.Println("顺子", straights, "重复", duplicates)
	}

	// fmt.Println("原始数据...", groupBySuitCards)
	return
}
