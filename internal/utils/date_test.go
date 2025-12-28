package utils

import (
	"testing"
	"time"
)

func Test_formatDateTime(t *testing.T) {
	type args struct {
		time_ time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test case 1: Standard date and time",
			args: args{time_: time.Date(2025, 4, 3, 14, 30, 0, 0, time.UTC)},
			want: "2025-04-03(Thu) 14:30",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatDateTime(tt.args.time_); got != tt.want {
				t.Errorf("formatDateTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
