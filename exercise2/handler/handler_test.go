package handler

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var yml = `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`

func Test_parseYaml(t *testing.T) {
	type args struct {
		yml []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Handles YAML Format",
			args: args{
				yml: []byte(yml),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := parseYaml(tt.args.yml)
			assert.NotNil(t, result)
			assert.Equal(t, "https://github.com/gophercises/urlshort", result["/urlshort"])
			assert.Equal(t, "https://github.com/gophercises/urlshort/tree/solution", result["/urlshort-final"])
		})
	}
}
