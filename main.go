package main

import (
	"rummy-card-truth/internal"
	"rummy-card-truth/pkg"
)

func main() {
	//cards := []pkg.Card{
	//	pkg.NewCard(1, pkg.D),
	//	pkg.NewCard(2, pkg.D),
	//	pkg.NewCard(3, pkg.D),
	//	pkg.NewCard(4, pkg.D),
	//	pkg.NewCard(5, pkg.D),
	//	pkg.NewCard(6, pkg.D),
	//	pkg.NewCard(11, pkg.D),
	//	pkg.NewCard(12, pkg.D),
	//	pkg.NewCard(13, pkg.D),
	//}

	cards := []pkg.Card{
		pkg.NewCard(1, pkg.D),
		pkg.NewCard(2, pkg.D),
		pkg.NewCard(12, pkg.D),
		pkg.NewCard(13, pkg.D),
		pkg.NewCard(9, pkg.D),
		pkg.NewCard(10, pkg.D),
		pkg.NewCard(11, pkg.D),

		pkg.NewCard(9, pkg.C),

		pkg.NewCard(2, pkg.B),
		pkg.NewCard(3, pkg.B),
		pkg.NewCard(4, pkg.B),
		pkg.NewCard(5, pkg.B),
		pkg.NewCard(9, pkg.B),
	}
	internal.NewPlanner(cards, 8).Run()
}
