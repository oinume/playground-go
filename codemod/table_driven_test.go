package main

import (
	"slices"
	"testing"
)

func TestContains(t *testing.T) {
	type args struct {
		ss []string
		s  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty: false",
			args: args{[]string{}, ""},
			want: false,
		},
		{
			name: "found: true",
			args: args{[]string{"a", "b", "c"}, "b"},
			want: true,
		},
		{
			name: "not found: false",
			args: args{[]string{"a", "b"}, "c"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := contains(tt.args.ss, tt.args.s); got != tt.want {
				t.Errorf("contains(): want=%v, got=%v", got, tt.want)
			}
		})
	}
}

func contains(ss []string, s string) bool {
	return slices.Contains(ss, s)
}
