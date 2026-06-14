package main

import "testing"

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name  string
		left  string
		right string
		want  int
	}{
		{name: "same", left: "v0.2.0", right: "0.2.0", want: 0},
		{name: "patch newer", left: "v0.2.1", right: "0.2.0", want: 1},
		{name: "minor newer", left: "v0.3.0", right: "0.2.9", want: 1},
		{name: "numeric compare", left: "v0.2.10", right: "0.2.9", want: 1},
		{name: "older", left: "v0.1.9", right: "0.2.0", want: -1},
		{name: "suffix ignored", left: "v1.0.0+build", right: "1.0.0", want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := compareVersions(tt.left, tt.right)
			if got != tt.want {
				t.Fatalf("compareVersions(%q, %q) = %d, want %d", tt.left, tt.right, got, tt.want)
			}
		})
	}
}
