package exercise7

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCamelCase(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"Zero length string has zero words",
			args{
				"",
			},
			0,
		},
		{
			"All lowercase strings is one word",
			args{
				"o",
			},
			1,
		},
		{
			"String length 2, lowercase",
			args{
				"oo",
			},
			1,
		},
		{
			"Uppercase word marks second word",
			args{
				"oO",
			},
			2,
		},
		{
			"Uppercase word marks third word",
			args{
				"oOoO",
			},
			3,
		},
		{
			"Uppercase word marks third word, consecutive capital letters",
			args{
				"oOOO",
			},
			4,
		},
		{
			"Uppercase word marks third word",
			args{
				"ooOooOooOoo",
			},
			4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, CamelCase(tt.args.s), tt.want)
		})
	}
}

func TestCeasarCipher(t *testing.T) {
	type args struct {
		s     string
		shift int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Zero length string can't be shifted",
			args{
				"",
				0,
			},
			"",
		},
		{
			"Zero length string can't be shifted",
			args{
				"",
				1,
			},
			"",
		},
		{
			"Shift 1 lowercase",
			args{
				"a",
				1,
			},
			"b",
		},
		{
			"Shift 1 lowercase z",
			args{
				"z",
				1,
			},
			"a",
		},
		{
			"Shift 1 uppercase z",
			args{
				"Z",
				1,
			},
			"A",
		},
		{
			"Shift 1 lowercase",
			args{
				"a",
				3,
			},
			"d",
		},
		{
			"Shift 1 uppercase letter",
			args{
				"A",
				1,
			},
			"B",
		},
		{
			"Shift 1 lowercase letter",
			args{
				"a",
				27,
			},
			"b",
		},
		{
			"Shift 2 lowercase when shift is over alphabet length",
			args{
				"aa",
				27,
			},
			"bb",
		},
		{
			"Shifts 2 lowercase letters",
			args{
				"aa",
				25,
			},
			"zz",
		},
		{
			"Doesn't shift non-letters",
			args{
				"a-a",
				2,
			},
			"c-c",
		},
		{
			"Sample test case",
			args{
				"There's-a-starman-waiting-in-the-sky.",
				3,
			},
			"Wkhuh'v-d-vwdupdq-zdlwlqj-lq-wkh-vnb.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, CaesarCipher(tt.args.s, tt.args.shift))
		})
	}
}
