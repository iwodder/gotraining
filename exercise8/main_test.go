package main

import (
	"reflect"
	"testing"
)

func Test_normalizePhoneNumbers(t *testing.T) {
	type args struct {
		numbers []string
	}
	tests := []struct {
		name string
		args args
		want []phoneNumber
	}{
		{
			name: "Normalizing nil returns empty map",
			args: args{
				nil,
			},
			want: []phoneNumber{},
		},
		{
			name: "Normalizing removes dashes",
			args: args{
				[]string{"123-456-8765"},
			},
			want: []phoneNumber{{oldFmt: "123-456-8765", newFmt: "1234568765"}},
		},
		{
			name: "Doesn't modify correctly formatted number",
			args: args{
				[]string{"1234568765"},
			},
			want: []phoneNumber{{oldFmt: "1234568765", newFmt: "1234568765"}},
		},
		{
			name: "Normalizes list of numbers",
			args: args{
				[]string{"1234567890", "123 456 7891", "(123) 456 7892", "(123) 456-7893", "123-456-7894", "123-456-7890", "1234567892", "(123)456-7892)"},
			},
			want: []phoneNumber{
				{oldFmt: "1234567890", newFmt: "1234567890"},
				{oldFmt: "123 456 7891", newFmt: "1234567891"},
				{oldFmt: "(123) 456 7892", newFmt: "1234567892"},
				{oldFmt: "(123) 456-7893", newFmt: "1234567893"},
				{oldFmt: "123-456-7894", newFmt: "1234567894"},
				{oldFmt: "123-456-7890", newFmt: "1234567890"},
				{oldFmt: "1234567892", newFmt: "1234567892"},
				{oldFmt: "(123)456-7892)", newFmt: "1234567892"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizePhoneNumbers(tt.args.numbers); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("normalizePhoneNumbers() = %v, want %v", got, tt.want)
			}
		})
	}
}
