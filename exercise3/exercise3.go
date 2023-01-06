package main

import (
	"fmt"
	"log"
	"net/http"
)

var defaultStory = "gopher.json"

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	story, err := CreateStory(defaultStory)
	if err != nil {
		return err
	}
	fmt.Printf("Running with story %s...\n", defaultStory)
	s := newServer()
	s.story = story
	return http.ListenAndServe(":8080", s)
}
