package link

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestExtractLinks(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want []Link
	}{
		{
			"HTML1 should have links extracted",
			args{
				"../html1.html",
			},
			[]Link{
				{
					"/other-page",
					[]string{"A link to another page"},
				},
			},
		},
		{
			"HTML2 should have links extracted",
			args{
				"../html2.html",
			},
			[]Link{
				{
					"https://www.twitter.com/joncalhoun",
					[]string{"Check me out on twitter"},
				},
				{
					"https://github.com/gophercises",
					[]string{"Gophercises is on", "Github", "!"},
				},
			},
		},
		{
			"HTML3 should have links extracted",
			args{
				"../html3.html",
			},
			[]Link{
				{
					"#",
					[]string{"Login"},
				},
				{
					"/lost",
					[]string{"Lost? Need help?"},
				},
				{
					"https://twitter.com/marcusolsson",
					[]string{"@marcusolsson"},
				},
			},
		},
		{
			"HTML4 should have links extracted",
			args{
				"../html4.html",
			},
			[]Link{
				{
					"/dog-cat",
					[]string{"dog cat"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(tt.args.filename)
			if err != nil {
				assert.FailNow(t, fmt.Sprintf("Unable to open the file, %s, for parsing.", tt.args.filename))
			}
			links, err := Parse(f)
			if err != nil {
				assert.FailNow(t, "Unexpected error when processing file.")
			}
			assert.Equal(t, links, tt.want)
		})
	}
}
