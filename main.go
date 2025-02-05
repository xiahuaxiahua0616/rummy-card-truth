package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rummy-card-truth/internal"
	"rummy-card-truth/pkg"
	"sort"
)

func main() {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	})
	r.GET("/api/v1/hand/range", func(c *gin.Context) {
		cards := []pkg.Card{
			{Suit: pkg.D, Value: 3},
			{Suit: pkg.D, Value: 5},
			{Suit: pkg.JokerSuit, Value: 0},

			{Suit: pkg.C, Value: 11},
			{Suit: pkg.C, Value: 12},
			{Suit: pkg.C, Value: 13},

			{Suit: pkg.B, Value: 2},
			{Suit: pkg.B, Value: 3},
			{Suit: pkg.B, Value: 4},
			{Suit: pkg.B, Value: 8},
			{Suit: pkg.B, Value: 9},
			{Suit: pkg.B, Value: 9},
			{Suit: pkg.B, Value: 10},

			{Suit: pkg.A, Value: 10},
		}

		result := internal.NewPlanner(cards, 9).Run()

		c.JSON(http.StatusOK, SuccessResponse{
			Success: true,
			Data: gin.H{
				"myCards": GetResponse([][]pkg.Card{cards})[0],
				"result":  GetResponse(result),
				"sysJoker": GetResponse([][]pkg.Card{
					{
						{Suit: pkg.A, Value: 9},
					},
				}),
			},
		})
	})
	r.Run(":8009") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func GetCardsResult(cards []pkg.Card) []int {
	var myCards []int
	for _, c := range cards {
		if c.Suit == pkg.A {
			myCards = append(myCards, c.Value+48)
		} else if c.Suit == pkg.B {
			myCards = append(myCards, c.Value+32)
		} else if c.Suit == pkg.C {
			myCards = append(myCards, c.Value+16)
		} else if c.Suit == pkg.D {
			myCards = append(myCards, c.Value)
		} else if c.Suit == pkg.JokerSuit {
			myCards = append(myCards, 79)
		} else if c.Suit == pkg.JokerSuit {
			myCards = append(myCards, 78)
		}
	}

	if len(myCards) == 0 {
		return []int{0}
	}
	return myCards
}

func GetResponse(cards ...[][]pkg.Card) [][]int {
	var res [][]int

	for _, cardDim := range cards {
		for _, card := range cardDim {
			if len(card) > 0 {
				res = append(res, GetCardsResult(card))
			}
		}
	}

	for _, card := range res {
		sort.Ints(card)
	}
	return res
}

type SuccessResponse struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
}
