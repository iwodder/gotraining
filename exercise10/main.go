package main

import (
	"fmt"
	"gotraining/exercise10/blackjack"
	"os"
)

func main() {
	//ai := blackjack.AI{}
	ai := blackjack.CliPlayer{Out: os.Stdout, In: os.Stdin}
	stats := blackjack.StatsPlayer{Player: &ai}
	g := blackjack.NewGame(&stats)
	g.Play(10)
	fmt.Printf("===Results===\n\t%s", &stats)
}
