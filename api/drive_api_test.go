package api

import (
	"testing"
)

func TestGetSharedDrives(t *testing.T) {
	tests := []struct {
		name string
		auth string
		want []*SDrive
	}{
		{
			"Test",
			"", // use root auth
			nil,
		},
		{
			"Test",
			"WRONG_PASSWORD",
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSharedDrives(tt.auth); got == nil {
				t.Errorf("GetSharedDrives() = %v, want %v", got, tt.want)
			} else {
				t.Log("#shared drives is:", len(got))
				for i, g := range got {
					t.Log(i+1, g.Name, g.Id)
				}
			}
		})
	}
}
