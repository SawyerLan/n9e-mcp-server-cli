package app

import (
	"testing"
)

func TestSlicePage(t *testing.T) {
	items := []int{1, 2, 3, 4, 5, 6, 7}

	tests := []struct {
		name      string
		page      int
		limit     int
		wantItems []int
		wantTotal int64
	}{
		{"defaults", 0, 0, []int{1, 2, 3, 4, 5, 6, 7}, 7},
		{"page 1 limit 3", 1, 3, []int{1, 2, 3}, 7},
		{"page 2 limit 3", 2, 3, []int{4, 5, 6}, 7},
		{"page 3 limit 3", 3, 3, []int{7}, 7},
		{"page beyond", 4, 3, []int{}, 7},
		{"empty slice", 1, 10, []int{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := items
			if tt.name == "empty slice" {
				input = []int{}
			}
			got, total := SlicePage(input, tt.page, tt.limit)
			if total != tt.wantTotal {
				t.Errorf("total = %d, want %d", total, tt.wantTotal)
			}
			if len(got) != len(tt.wantItems) {
				t.Errorf("len = %d, want %d", len(got), len(tt.wantItems))
				return
			}
			for i, v := range got {
				if v != tt.wantItems[i] {
					t.Errorf("item[%d] = %d, want %d", i, v, tt.wantItems[i])
				}
			}
		})
	}
}
