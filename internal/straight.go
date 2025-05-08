package internal

import (
	"slices"

	"github.com/xiahua/ifonly/pkg"
)

func GetStraight(cards []byte, joker byte) (straights [][]byte, overCards []byte) {
	// 找到牌中所有的顺子（不带joker的）
	groupBySuitCards := pkg.GroupBySuit(cards)
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
		cards, _ := UniqueAndDuplicates(cards)

		// tempBytes 临时存储使用

		var straight []byte
		for {
			currentCards := make([]byte, len(cards))
			copy(currentCards, cards)
			slices.Sort(cards)

			straight, cards = GetStraightOnyByOne(cards)
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

// GetStraightOnyByOne 一个接一个的获取
func GetStraightOnyByOne(cards []byte) (straight []byte, overCards []byte) {
	tempBytes := []byte{}
	for {
		if len(cards) <= 0 {
			break
		}

		if isContinuous(cards) {
			// 是顺子
			return cards, tempBytes
		} else {
			tempBytes = append(tempBytes, cards[len(cards)-1])
			cards = cards[:len(cards)-1]
		}
	}
	return []byte{}, tempBytes
}

// isContinuous 判断数据是否连续
func isContinuous(nums []byte) bool {
	if len(nums) < 2 {
		return false
	}
	// TIPS: 这里无法判断Joker、是否连续
	for i := 1; i < len(nums); i++ {
		if nums[i] != nums[i-1]+1 {
			return false
		}
	}
	return true
}

// UniqueAndDuplicates 对byte切片进行去重，返回去重后的切片和重复的元素
func UniqueAndDuplicates(input []byte) (unique, duplicates []byte) {
	seen := make(map[byte]bool)
	addedDup := make(map[byte]bool)
	for _, b := range input {
		if !seen[b] {
			seen[b] = true
			unique = append(unique, b)
		} else if !addedDup[b] {
			duplicates = append(duplicates, b)
			addedDup[b] = true
		}
	}
	return
}
