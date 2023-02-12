package main

import (
	"fmt"
	"gotraining/exercise10/blackjack"
)

func main() {
	ai := blackjack.AI{}
	//p := blackjack.NewCLIPlayer(os.Stdout, os.Stdin)
	stats := blackjack.StatsPlayer{Player: &ai}
	g := blackjack.NewGame(&stats)
	g.Play(1000)
	fmt.Printf("===Results===\n\t%s", &stats)
}
