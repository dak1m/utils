package utils

import (
	"testing"
	"time"
)

func TestIsOnlySet(t *testing.T) {
	type args struct {
		obj     any
		field   string
		exclude []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "TestIsOnlySet",
			args: args{
				obj: struct {
					Id    string     `json:"id"`
					Name  string     `json:"name"`
					Age   int        `json:"age"`
					Start *time.Time `json:"start"`
				}{
					Id:    "1",
					Start: timePtr(time.Now()),
				},
				field:   "start",
				exclude: []string{"id"},
			},
			want: true,
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsOnlySet(tt.args.obj, tt.args.field, tt.args.exclude...); got != tt.want {
				t.Errorf("IsOnlySet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func timePtr(t time.Time) *time.Time {
	return &t
}
