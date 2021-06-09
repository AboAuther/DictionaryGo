package youdao

import "testing"

func Test_truncateString(t *testing.T) {
	type args struct {
		q string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"simple", args{q: "xx"}, "xx"},
		{"long", args{q: "111111111101234567892222222222"}, "1111111111302222222222"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := truncateString(tt.args.q); got != tt.want {
				t.Errorf("truncateString() = %v, want %v", got, tt.want)
			}
		})
	}
}
