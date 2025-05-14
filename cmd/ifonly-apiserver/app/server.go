package app

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/xiahua/ifonly/internal"
	internalV1 "github.com/xiahua/ifonly/internal/v1"
	"github.com/xiahua/ifonly/pkg"
)

var configFile string

var mode string

func NewIfOnlyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ifonly",
		Short: "ifonly is hlep robot to be a winner",
		Long:  "We use ifonly to let the robot win, because ifonly is the best hand generator algorithm",
		RunE: func(cmd *cobra.Command, args []string) error {
			if mode == "release" {
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

				r.GET("/api/v2/hand/range", func(c *gin.Context) {
					cards := []byte{
						0x03, 0x05, 0x4e,
						0x1b, 0x1c, 0x1d,
						0x22, 0x23, 0x24, 0x28, 0x29, 0x29, 0x2a,
						0x3a,
					}
					internalV1.NewPlannerV2(cards, 0x39).Run()
				})
				r.Run(":8009") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
			} else {
				DoSomething()
			}

			return nil
		},
		SilenceUsage: true,
		Args:         cobra.NoArgs,
	}

	// 这应该是一个异步调用
	cobra.OnInitialize(onInitialize)

	cmd.PersistentFlags().StringVarP(&mode, "mode", "m", "release", "Running mode: debug or release")
	cmd.PersistentFlags().StringVarP(&configFile, "config", "c", filePath(), "Path to the miniblog configuration file.")
	return cmd
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

func DoSomething() {
	cards := []byte{
		0x03, 0x05, 0x4e,
		0x1b, 0x1c, 0x1d,
		0x22, 0x23, 0x24, 0x28, 0x29, 0x29, 0x2a,
		0x3a,
	}
	internalV1.NewPlannerV2(cards, 9).Run()
}
