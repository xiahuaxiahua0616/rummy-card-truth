package internal

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
