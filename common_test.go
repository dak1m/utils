package utils

import (
	"testing"
	"time"
)

func TestIsOnlySet(t *testing.T) {
	type args struct {
		obj   any
		field string
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
					Name  string     `json:"name"`
					Age   int        `json:"age"`
					Start *time.Time `json:"start"`
				}{
					Start: timePtr(time.Now()),
				},
				field: "start",
			},
			want: true,
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsOnlySet(tt.args.obj, tt.args.field); got != tt.want {
				t.Errorf("IsOnlySet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func timePtr(t time.Time) *time.Time {
	return &t
}
