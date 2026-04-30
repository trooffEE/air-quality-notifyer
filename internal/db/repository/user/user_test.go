package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroupObservedIDs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		rows []userObservedDistrict
		want map[int64][]int64
	}{
		{
			name: "groups rows by telegram id",
			rows: []userObservedDistrict{
				{TelegramID: 101, DistrictID: 1},
				{TelegramID: 202, DistrictID: 3},
				{TelegramID: 101, DistrictID: 2},
			},
			want: map[int64][]int64{
				101: {1, 2},
				202: {3},
			},
		},
		{
			name: "preserves duplicate observed ids",
			rows: []userObservedDistrict{
				{TelegramID: 101, DistrictID: 1},
				{TelegramID: 101, DistrictID: 1},
			},
			want: map[int64][]int64{
				101: {1, 1},
			},
		},
		{
			name: "empty rows",
			rows: nil,
			want: map[int64][]int64{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, groupObservedIDs(tt.rows))
		})
	}
}
