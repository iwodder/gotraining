package main

import (
	"fmt"
	"gotraining/exercise10/blackjack"
)

type CliPrompt struct{}

func (h *CliPrompt) Prompt(s string) {
	fmt.Println(s)
}

func (h *CliPrompt) Response() string {
	fmt.Printf("> ")
	var ret string
	fmt.Scanf("%s\n", &ret)
	return ret
}

func main() {
	g := blackjack.NewGame(blackjack.NewPlayer(&CliPrompt{}))
	g.Play(1)
}
