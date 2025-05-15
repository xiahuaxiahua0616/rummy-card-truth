package v1

import (
	"reflect"
	"sort"
	"testing"
	// internalV1 "github.com/xiahua/ifonly/internal/v1"
)

func TestPlanner(t *testing.T) {
	tests := []struct {
		name         string
		cards        []byte
		joker        byte
		wantStraight [][]byte
		duplicate    []byte
	}{
		{
			name: "001",
			cards: []byte{
				0x33, 0x34, 0x35, 0x11, 0x1d, 0x25, 0x12, 0x22, 0x32, 0x05, 0x19, 0x39, 0x37,
			},
			joker: 0x05,
			wantStraight: [][]byte{
				{0x33, 0x34, 0x35}, {0x11, 0x1d, 0x25}, {0x12, 0x22, 0x32}, {0x19, 0x39, 0x05}, {0x37},
			},
			duplicate: []byte{},
		},
		{
			name: "002",
			cards: []byte{
				0x03, 0x05, 0x4f,
				0x1b, 0x1c, 0x1d,
				0x22, 0x23, 0x24, 0x28, 0x29, 0x29, 0x2a,
				0x3a,
			},
			joker: 0x29,
			wantStraight: [][]byte{
				{0x1b, 0x1c, 0x1d},
				{0x22, 0x23, 0x24},
				{0x28, 0x29, 0x2a},
				{0x3a, 0x29, 0x4f},
				{0x03, 0x05},
			},
			duplicate: []byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result [][]byte
			NewPlannerV2(tt.cards, tt.joker).Run(&result)

			// 对每个子切片进行排序
			for _, slice := range result {
				sort.Slice(slice, func(i, j int) bool {
					return slice[i] < slice[j]
				})
			}
			for _, slice := range tt.wantStraight {
				sort.Slice(slice, func(i, j int) bool {
					return slice[i] < slice[j]
				})
			}
			if !reflect.DeepEqual(result, tt.wantStraight) {
				t.Errorf("expected %v, got %v", tt.wantStraight, result)
			}
		})
	}
}
