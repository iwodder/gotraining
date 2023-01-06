package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_parsePath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"Splits a single /",
			args{"/"},
			nil,
		},
		{
			"Splits one path single /",
			args{"/new-york"},
			[]string{"new-york"},
		},
		{
			"Handles no forward slashes",
			args{"new-york"},
			[]string{"new-york"},
		},
		{
			"Splits two part path single /",
			args{"/new-york/denver"},
			[]string{"new-york", "denver"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, parsePath(tt.args.path))
		})
	}
}
