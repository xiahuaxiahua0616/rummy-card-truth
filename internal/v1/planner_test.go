package v1

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	// internalV1 "github.com/xiahua/ifonly/internal/v1"
	"slices"
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
				{0x32, 0x33, 0x34}, {0x11, 0x1d, 0x35}, {0x19, 0x25, 0x39}, {0x05, 0x12, 0x22}, {0x37},
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
		{
			name: "003",
			cards: []byte{
				0x11, 0x1b, 0x1c, 0x1d, 0x25, 0x3b, 0x3c, 0x09, 0x19, 0x39, 0x12, 0x24, 0x34,
			},
			joker: 0x25,
			wantStraight: [][]byte{
				{0x11, 0x1b, 0x1c, 0x1d},
				{0x25, 0x3b, 0x3c},
				{0x09, 0x19, 0x39},
				{0x12, 0x24, 0x34},
			},
			duplicate: []byte{},
		},
		{
			name: "004",
			cards: []byte{
				0x04, 0x05, 0x0b, 0x0c, 0x0d, 0x15, 0x27, 0x28, 0x29, 0x2a, 0x2c, 0x37, 0x4f,
			},
			joker: 0x28,
			wantStraight: [][]byte{
				{0x0b, 0x0c, 0x0d},
				{0x29, 0x2a, 0x2c, 0x4f},
				{0x27, 0x28, 0x37},
				{0x04, 0x05, 0x15},
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
				slices.Sort(slice)
			}
			for _, slice := range tt.wantStraight {
				slices.Sort(slice)
			}
			if !reflect.DeepEqual(result, tt.wantStraight) {
				wantStraightStr := formatSlice(tt.wantStraight)
				resultStr := formatSlice(result)
				t.Errorf("expected %v, got %v", wantStraightStr, resultStr)
			}
		})
	}
}

func formatSlice(slice [][]byte) string {
	var result []string
	for _, subSlice := range slice {
		// 格式化每个子切片为带逗号的大括号字符串
		subSliceStr := fmt.Sprintf("{%s}", strings.Join(formatSubSlice(subSlice), ","))
		result = append(result, subSliceStr)
	}
	return fmt.Sprintf("{%s}", strings.Join(result, ", "))
}

// formatSubSlice 将单个切片格式化为带逗号的字符串
func formatSubSlice(slice []byte) []string {
	var result []string
	for _, val := range slice {
		result = append(result, fmt.Sprintf("%d", val))
	}
	return result
}
