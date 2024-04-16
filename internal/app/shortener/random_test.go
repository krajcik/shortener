package shortener

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_randomString(t *testing.T) {
	type args struct {
		leng int
	}
	tests := []struct {
		name string
		args args
	}{
		{"success", args{1000}},
		{"success", args{0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := randomString(tt.args.leng)
			assert.Lenf(t, s, tt.args.leng, "randomString(%v)", tt.args.leng)
		})
	}
}
