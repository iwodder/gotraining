package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

func CreateStory(filename string) (*Story, error) {
	result, err := parseJsonFile(filename)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func parseJsonFile(filename string) (Story, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("story: Unable to open file '%s'", filename))
	}
	d := json.NewDecoder(file)
	var story Story
	if err := d.Decode(&story); err != nil {
		fmt.Println(err)
		return nil, errors.New(fmt.Sprintf("story: Unable to create decoder."))
	}
	return story, nil
}

type Story map[string]Chapter

type Chapter struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}
