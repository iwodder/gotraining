package main

import (
	"fmt"
	"gotraining/exercise10/blackjack"
)

type HumanPlayer struct{}

func (h *HumanPlayer) Prompt(s string) {
	fmt.Println(s)
}

func (h *HumanPlayer) NextMove() string {
	var ret string
	fmt.Printf("> ")
	_, _ = fmt.Scanf("%s\n", &ret)
	return ret
}

func main() {
	g := blackjack.New(&HumanPlayer{})
	g.Play()
}
