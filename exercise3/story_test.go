package main

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseJsonStory(t *testing.T) {
	gopherStory := make(map[string]Chapter)
	gopherStory["intro"] = Chapter{
		Title: "The Little Blue Gopher",
		Story: []string{
			"Once upon a time, long long ago, there was a little blue gopher. Our little blue friend wanted to go on an adventure, but he wasn't sure where to go. Will you go on an adventure with him?",
			"One of his friends once recommended going to New York to make friends at this mysterious thing called \"GothamGo\". It is supposed to be a big event with free swag and if there is one thing gophers love it is free trinkets. Unfortunately, the gopher once heard a campfire story about some bad fellas named the Sticky Bandits who also live in New York. In the stories these guys would rob toy stores and terrorize young boys, and it sounded pretty scary.",
			"On the other hand, he has always heard great things about Denver. Great ski slopes, a bad hockey team with cheap tickets, and he even heard they have a conference exclusively for gophers like himself. Maybe Denver would be a safer place to visit.",
		},
		Options: []Option{
			{
				Text:    "That story about the Sticky Bandits isn't real, it is from Home Alone 2! Let's head to New York.",
				Chapter: "new-york",
			},
			{
				Text:    "Gee, those bandits sound pretty real to me. Let's play it safe and try our luck in Denver.",
				Chapter: "denver",
			},
		},
	}
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Can create story from default gopher.json",
			args: args{
				filename: "gopher.json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			story, err := CreateStory(tt.args.filename)
			if err != nil {
				assert.Failf(t, "Received unexpected error", err.Error())
			}
			assert.Equal(t, gopherStory["intro"], (*story)["intro"])
		})
	}
}

func TestParseJsonStoryErrors(t *testing.T) {
	type args struct {
		filename string
		err      error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Missing file returns an error",
			args: args{
				filename: "missing_story.json",
				err:      errors.New("story: Unable to open file 'missing_story.json'"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateStory(tt.args.filename)
			switch err.(type) {
			case error:
				assert.Equal(t, tt.args.err, err)
			case nil:
				assert.Failf(t, "Expected to recieve an error", "")
			}
		})
	}
}
