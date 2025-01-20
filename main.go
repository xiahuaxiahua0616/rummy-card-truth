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
		app.NewCard(14, app.D), // 方片 1
		//app.NewCard(2, 0), // 方片 2
		app.NewCard(3, app.D),  // 方片 3
		app.NewCard(4, app.D),  // 方片 4
		app.NewCard(13, app.D), // 方片 4
		app.NewCard(11, app.D), // 方片 4
		app.NewCard(11, app.C), // 方片 4
		app.NewCard(11, app.B), // 方片 4
		app.NewCard(11, app.B), // 方片 4

		app.NewCard(9, app.D), // 方片 4
	}

	cards3 := []app.Card{
		app.NewCard(3, app.D), // 方片 3
		app.NewCard(3, app.D), // 方片 4
		app.NewCard(3, app.C), // 方片 4

		app.NewCard(5, app.C), // 方片 4
		app.NewCard(5, app.B), // 方片 4

		app.NewCard(9, app.D), // 方片 4
	}

	cards4 := []app.Card{

		app.NewCard(5, app.C), // 方片 4

		app.NewCard(9, app.D), // 方片 4
		app.NewCard(9, app.A), // 方片 4
	}

	//fmt.Println(app.GetPure(cards))
	//fmt.Println(app.GetPure(cards2, false))
	//fmt.Println(app.GetPureWithJoker(cards2, 9, true))
	//fmt.Println(app.GetPureWithJoker(cards2, 9, false))
	fmt.Println(app.GetSet(cards2))

	fmt.Println(app.GetSetWithJoker(cards3, 9))
	fmt.Println(app.GetSetWithJoker(cards4, 9))
}
