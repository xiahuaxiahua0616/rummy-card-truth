package v1

import (
	"reflect"
	"slices"
	"testing"
)

func TestGetStraight(t *testing.T) {
	tests := []struct {
		name         string
		cards        []byte
		joker        byte
		wantStraight [][]byte
	}{
		{
			name:         "simple straight 2,3,4",
			cards:        []byte{2, 3, 4},
			joker:        0,
			wantStraight: [][]byte{{2, 3, 4}},
		},
		{
			name:         "multiple straights",
			cards:        []byte{2, 3, 4, 6, 7, 8},
			joker:        0,
			wantStraight: [][]byte{{2, 3, 4}, {6, 7, 8}},
		},
		{
			name:         "non-consecutive cards",
			cards:        []byte{2, 4, 6, 8},
			joker:        0,
			wantStraight: [][]byte{}, // 没有顺子
		},
		{
			name:         "duplicate values",
			cards:        []byte{2, 2, 3, 3, 4, 5},
			joker:        0,
			wantStraight: [][]byte{{2, 3, 4, 5}},
		},
		{
			name:         "multiple straights v2",
			cards:        []byte{0x1, 0x2, 0x3, 0x5, 0x6, 0x7, 0x9, 0xa, 0xb},
			joker:        0,
			wantStraight: [][]byte{{0x1, 0x2, 0x3}, {0x5, 0x6, 0x7}, {0x9, 0xa, 0xb}},
		},
		{
			name:         "score",
			cards:        []byte{1, 2, 3, 12, 13},
			joker:        0,
			wantStraight: [][]byte{{1, 12, 13}},
		},
		{
			name:         "007",
			cards:        []byte{3, 5, 78, 27, 28, 29, 34, 35, 36, 40, 41, 41, 42, 58},
			joker:        0x39,
			wantStraight: [][]byte{{27, 28, 29}, {34, 35, 36}, {40, 41, 42}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStraights, _ := GetStraight(tt.cards, tt.joker)

			for _, slice := range gotStraights {
				slices.Sort(slice)
			}
			for _, slice := range tt.wantStraight {
				slices.Sort(slice)
			}
			if gotStraights == nil {
				gotStraights = [][]byte{}
			}
			if !reflect.DeepEqual(gotStraights, tt.wantStraight) {
				t.Errorf("expected %v, got %v", tt.wantStraight, gotStraights)
			}
		})
	}
}

func TestGetStraightWithJoker(t *testing.T) {
	tests := []struct {
		name         string
		cards        []byte
		duplicate    []byte
		joker        byte
		wantStraight [][]byte
	}{
		{
			name:         "001",
			cards:        []byte{1, 2, 11, 13, 7},
			joker:        7,
			duplicate:    []byte{2},
			wantStraight: [][]byte{{1, 7, 11, 13}},
		},
		{
			name:         "002",
			cards:        []byte{2, 3, 12, 13, 7},
			joker:        7,
			duplicate:    []byte{2, 3},
			wantStraight: [][]byte{{7, 12, 13}},
		},
		{
			name:         "003",
			cards:        []byte{1, 2, 9, 10, 7},
			joker:        7,
			duplicate:    []byte{1, 2},
			wantStraight: [][]byte{{7, 9, 10}},
		},
		{
			name:         "004",
			cards:        []byte{5, 4, 8, 9, 1},
			joker:        1,
			duplicate:    []byte{4, 5},
			wantStraight: [][]byte{{1, 8, 9}},
		},
		{
			name:         "005",
			cards:        []byte{3, 5, 7},
			joker:        7,
			duplicate:    []byte{},
			wantStraight: [][]byte{{3, 5, 7}},
		},
		{
			name:         "006",
			cards:        []byte{2, 6, 7},
			joker:        7,
			duplicate:    []byte{2, 6, 7},
			wantStraight: [][]byte{},
		},
		// {
		// 	name:      "007",
		// 	cards:     []byte{0x03, 0x05, 0x28, 0x29, 0x4f},
		// 	joker:     0x29,
		// 	duplicate: []byte{},
		// 	wantStraight: [][]byte{
		// 		{0x03, 0x05, 0x29},
		// 		{0x28, 0x29, 0x4f}, // 这个不要开了
		// 	},
		// },
		{
			name:         "008",
			cards:        []byte{0x25, 0x39, 0x3b, 0x3c},
			joker:        0x25,
			duplicate:    []byte{},
			wantStraight: [][]byte{{0x39, 0x25, 0x3b, 0x3c}},
		},
		{
			name:      "009",
			cards:     []byte{0x21, 0x28, 0x2a, 0x14, 0x15},
			joker:     0x01,
			duplicate: []byte{0x14, 0x15},
			wantStraight: [][]byte{
				{0x21, 0x28, 0x2a},
			},
		},
		{
			name:      "010",
			cards:     []byte{0x21, 0x1a, 0x14, 0x15},
			joker:     0x01,
			duplicate: []byte{0x1a},
			wantStraight: [][]byte{
				{0x21, 0x14, 0x15},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStraights, duplicate := GetStraightWithJoker(tt.cards, tt.joker)
			if gotStraights == nil {
				gotStraights = [][]byte{}
			}
			if duplicate == nil {
				duplicate = []byte{}
			}

			for _, slice := range gotStraights {
				slices.Sort(slice)
			}
			for _, slice := range tt.wantStraight {
				slices.Sort(slice)
			}
			slices.Sort(duplicate)
			slices.Sort(tt.duplicate)
			if !reflect.DeepEqual(gotStraights, tt.wantStraight) {
				t.Errorf("expected %v, got %v", tt.wantStraight, gotStraights)
			}

			if !reflect.DeepEqual(duplicate, tt.duplicate) {
				t.Errorf("duplicate expected %v, got %v", tt.duplicate, duplicate)
			}
		})
	}
}
