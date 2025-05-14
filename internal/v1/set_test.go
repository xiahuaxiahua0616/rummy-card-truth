package v1

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSet(t *testing.T) {
	tests := []struct {
		name         string
		cards        []byte
		joker        byte
		wantStraight [][]byte
	}{
		{
			name:         "001",
			cards:        []byte{0x1, 0x2, 0x3, 0x11, 0x21},
			joker:        0,
			wantStraight: [][]byte{{0x01, 0x11, 0x21}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSet, duplicate := GetSetV2(tt.cards)
			if !reflect.DeepEqual(gotSet, tt.wantStraight) {
				fmt.Println("剩余牌...", duplicate)
				t.Errorf("expected %v, got %v", tt.wantStraight, gotSet)
			}
		})
	}
}

func TestSetWithJoker(t *testing.T) {
	tests := []struct {
		name         string
		cards        []byte
		joker        byte
		wantStraight [][]byte
		duplicate    []byte
	}{
		{
			name:         "001",
			cards:        []byte{0x1, 0x21, 0x07},
			joker:        0x07,
			wantStraight: [][]byte{{0x1, 0x21, 0x07}},
			duplicate:    []byte{},
		},
		{
			name:         "002",
			cards:        []byte{0x1, 0x21, 0x07, 0x07, 0x02, 0x12},
			joker:        0x07,
			wantStraight: [][]byte{{0x1, 0x21, 0x07}, {0x02, 0x12, 0x07}},
			duplicate:    []byte{},
		},
		{
			name:         "003",
			cards:        []byte{0x01, 0x07, 0x07},
			joker:        0x07,
			wantStraight: [][]byte{{0x01, 0x07, 0x07}},
			duplicate:    []byte{},
		},
		{
			name:         "004",
			cards:        []byte{0x02, 0x12, 0x01, 0x07, 0x07},
			joker:        0x07,
			wantStraight: [][]byte{{0x02, 0x12, 0x07}},
			duplicate:    []byte{0x01, 0x07},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSet, duplicate := GetSetWithJokerV2(tt.cards, tt.joker)
			if !reflect.DeepEqual(gotSet, tt.wantStraight) {
				t.Errorf("expected %v, got %v", tt.wantStraight, gotSet)
			}
			if len(tt.duplicate) > 0 && !reflect.DeepEqual(duplicate, tt.duplicate) {
				t.Errorf("expected %v, got %v", tt.duplicate, duplicate)
			}
		})
	}
}
