package main

import (
	"fmt"
	"gotraining/exercise10/blackjack"
)

func main() {
	ai := blackjack.AI{}
	stats := blackjack.StatsPlayer{Player: &ai}
	g := blackjack.NewGame(&stats)
	g.Play(100)
	fmt.Printf("===Results===\n\t%s", &stats)
}
