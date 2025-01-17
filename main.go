package main

import (
	"fmt"
	"rummy-card-truth/app"
)

func main() {
	// 不考虑其他花色，只小范围开始做
	// 当前例子采用方片，1，2，3，5进行示范，找顺子

	// 创建牌对象
	//cards := []app.Card{
	//	app.NewCard(1, 0), // 方片 1
	//	app.NewCard(2, 0), // 方片 2
	//	app.NewCard(3, 0), // 方片 3
	//	app.NewCard(4, 0), // 方片 4
	//	app.NewCard(5, 0), // 方片 5
	//}

	cards2 := []app.Card{
		app.NewCard(1, 0), // 方片 1
		//app.NewCard(2, 0), // 方片 2
		app.NewCard(3, 0), // 方片 3
		app.NewCard(4, 0), // 方片 4
		//app.NewCard(5, 0), // 方片 5
		app.NewCard(6, 0),  // 方片 5
		app.NewCard(7, 0),  // 方片 5
		app.NewCard(8, 0),  // 方片 5
		app.NewCard(9, 0),  // 方片 5
		app.NewCard(10, 0), // 方片 5
		//app.NewCard(11, 0), // 方片 5
		app.NewCard(12, 0), // 方片 5
	}

	//fmt.Println(app.GetPure(cards))
	fmt.Println(app.GetPure(cards2, false))
}
