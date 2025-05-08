package internal

import (
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
			cards:        []byte{1, 2, 3, 5, 6, 7, 9, 10, 11},
			joker:        0,
			wantStraight: [][]byte{{1, 2, 3}, {5, 6, 7}, {9, 10, 11}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStraights, _ := GetStraight(tt.cards, tt.joker)
			if !reflect.DeepEqual(gotStraights, tt.wantStraight) {
				t.Errorf("expected %v, got %v", tt.wantStraight, gotStraights)
			}
		})
	}
}
