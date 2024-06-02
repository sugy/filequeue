package filequeue

import (
	"testing"
)

func TestSelectCopyCmd(t *testing.T) {
	tests := []struct {
		name string
		os   string
		want string
	}{
		{"Darwin", "darwin", "pbcopy"},
		{"Linux", "linux", "cat"},
		{"OtherOS", "windows", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := clipboardCopyCmd(tt.os)
			if got != tt.want {
				t.Errorf("selectCopyCmd(%v) = %v, want %v", tt.os, got, tt.want)
			}
		})
	}
}
