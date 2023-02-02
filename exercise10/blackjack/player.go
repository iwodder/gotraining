package blackjack

import (
	"fmt"
	"gotraining/exercise9/deck"
	"io"
)

type CliPlayer struct {
	Out io.Writer
	In  io.Reader
}

func (c *CliPlayer) Action(hand Hand, dealer deck.Card, actions ...Action) Action {
	_, _ = fmt.Fprintf(c.Out, "Dealer Hand=**HIDDEN**, %s\n", dealer)
	_, _ = fmt.Fprintf(c.Out, "Your Hand=%s (score=%d)\n", hand, hand.Score())
	_, _ = fmt.Fprintln(c.Out, "What do you want to do?")
	for i, v := range actions {
		_, _ = fmt.Fprintf(c.Out, "\t%d) %s\n", i+1, v)
	}
	var opt int
	_, _ = fmt.Fprintf(c.Out, "> ")
	for _, err := fmt.Fscan(c.In, &opt); err != nil; _, err = fmt.Fscan(c.In, &opt) {
		_, _ = fmt.Fprintf(c.Out, "You must enter a number between 1 and %d\n> ", len(actions))
	}
	return actions[opt-1]
}

func (c *CliPlayer) Prompt(msg string) {
	_, _ = fmt.Fprintln(c.Out, msg)
}

func (c *CliPlayer) Result(r Result) {
	switch r {
	case Win:
		_, _ = fmt.Fprintf(c.Out, "You're a winner!\n")
	case Lose:
		_, _ = fmt.Fprintf(c.Out, "You lose, try again!\n")
	case Draw:
		_, _ = fmt.Fprintf(c.Out, "---Draw----\n")
	}
}

type StatsPlayer struct {
	wins   int
	losses int
	draws  int
	Player
}

func (s *StatsPlayer) Result(r Result) {
	switch r {
	case Win:
		s.wins++
	case Lose:
		s.losses++
	case Draw:
		s.draws++
	}
	s.Player.Result(r)
}

func (s *StatsPlayer) String() string {
	return fmt.Sprintf("Wins=%d, Losses=%d, Draws=%d\n", s.wins, s.losses, s.draws)
}

type AI struct{}

func (a *AI) Action(hand Hand, dealer deck.Card, actions ...Action) Action {
	if hand.Score() >= 17 {
		return ActionHit
	}
	return ActionStand
}

func (a *AI) Result(r Result) {
	//pass
}
