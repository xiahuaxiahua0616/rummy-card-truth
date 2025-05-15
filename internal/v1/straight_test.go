package v1

import (
	"fmt"
	"reflect"
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
			wantStraight: nil, // 没有顺子
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
			gotStraights, duplicate := GetStraight(tt.cards, tt.joker)
			if !reflect.DeepEqual(gotStraights, tt.wantStraight) {
				fmt.Println("剩余牌...", duplicate)
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
			duplicate:    []byte{1, 2},
			wantStraight: [][]byte{{7, 11, 13}},
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
			name:         "003",
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
		{
			name:         "007",
			cards:        []byte{0x03, 0x05, 0x28, 0x29, 0x4f},
			joker:        0x29,
			duplicate:    []byte{0x4f, 0x28},
			wantStraight: [][]byte{{0x03, 0x05, 0x29}},
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
			if !reflect.DeepEqual(gotStraights, tt.wantStraight) {
				t.Errorf("expected %v, got %v", tt.wantStraight, gotStraights)
			}

			if !reflect.DeepEqual(duplicate, tt.duplicate) {
				t.Errorf("duplicate expected %v, got %v", tt.duplicate, duplicate)
			}
		})
	}
}
