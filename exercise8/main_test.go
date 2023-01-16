package main

import (
	"github.com/stretchr/testify/assert"
	"gotraining/exercise8/model"
	"testing"
)

func Test_normalizePhoneNumbers(t *testing.T) {
	type args struct {
		numbers []model.PhoneNumber
	}
	tests := []struct {
		name string
		args args
		want []model.PhoneNumber
	}{
		{
			name: "Normalizing nil returns empty slice",
			args: args{
				nil,
			},
			want: []model.PhoneNumber{},
		},
		{
			name: "Normalizing removes dashes",
			args: args{
				[]model.PhoneNumber{{Number: "123-456-8765"}},
			},
			want: []model.PhoneNumber{{Number: "1234568765"}},
		},
		{
			name: "Doesn't modify correctly formatted number",
			args: args{
				[]model.PhoneNumber{{Number: "1234568765"}},
			},
			want: []model.PhoneNumber{{Number: "1234568765"}},
		},
		{
			name: "Normalizes list of numbers",
			args: args{
				[]model.PhoneNumber{{Number: "1234567890"}, {Number: "123 456 7891"}, {Number: "(123) 456 7892"}, {Number: "(123) 456-7893"}, {Number: "123-456-7894"}, {Number: "(123)456-7892)"}},
			},
			want: []model.PhoneNumber{{Number: "1234567890"}, {Number: "1234567891"}, {Number: "1234567892"}, {Number: "1234567893"}, {Number: "1234567894"}, {Number: "1234567892"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizePhoneNumbers(tt.args.numbers)
			assert.Equal(t, tt.want, got)
		})
	}
}
