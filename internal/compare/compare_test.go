package compare

import (
	"testing"

	"github.com/mathiswitte/flat-notifier-go/internal/model"
)

func TestFindNewFlats(t *testing.T) {
	tests := []struct {
		name        string
		found       []model.FlatInfo
		existingIDs []string
		wantCount   int
		wantIDs     []string
	}{
		{
			name: "empty DB - all new",
			found: []model.FlatInfo{
				{FlatID: "1", FlatURL: "url1"},
				{FlatID: "2", FlatURL: "url2"},
			},
			existingIDs: nil,
			wantCount:   2,
			wantIDs:     []string{"1", "2"},
		},
		{
			name: "all known - none new",
			found: []model.FlatInfo{
				{FlatID: "1", FlatURL: "url1"},
				{FlatID: "2", FlatURL: "url2"},
			},
			existingIDs: []string{"1", "2"},
			wantCount:   0,
		},
		{
			name: "mixed - some new",
			found: []model.FlatInfo{
				{FlatID: "1", FlatURL: "url1"},
				{FlatID: "2", FlatURL: "url2"},
				{FlatID: "3", FlatURL: "url3"},
			},
			existingIDs: []string{"1", "3"},
			wantCount:   1,
			wantIDs:     []string{"2"},
		},
		{
			name:        "empty found list",
			found:       nil,
			existingIDs: []string{"1", "2"},
			wantCount:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindNewFlats(tt.found, tt.existingIDs)
			if len(got) != tt.wantCount {
				t.Errorf("got %d new flats, want %d", len(got), tt.wantCount)
			}
			for i, wantID := range tt.wantIDs {
				if i < len(got) && got[i].FlatID != wantID {
					t.Errorf("got[%d].FlatID = %q, want %q", i, got[i].FlatID, wantID)
				}
			}
		})
	}
}
