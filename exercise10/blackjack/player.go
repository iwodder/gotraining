package blackjack

import (
	"fmt"
	"io"
)

type CliPlayer struct {
	out io.Writer
	in  io.Reader
}

func (c *CliPlayer) ShowHand(h Hand) {
	_, _ = fmt.Fprintf(c.out, "Your hand=%s, (Score=%d)\n", h, h.score())
}

func (c *CliPlayer) Action(action ...Action) Action {
	_, _ = fmt.Fprintln(c.out, "What do you want to do?")
	for i, v := range action {
		_, _ = fmt.Fprintf(c.out, "\t %d) %s\n", i+1, v)
	}
	var opt int
	_, _ = fmt.Fprintf(c.out, "> ")
	_, _ = fmt.Fscanf(c.in, "%d\n", &opt)
	return action[opt-1]
}

func (c *CliPlayer) Prompt(msg string) {
	_, _ = fmt.Fprintln(c.out, msg)
}

func (c *CliPlayer) Win() {
	_, _ = fmt.Fprintln(c.out, "You won!")
}

func (c *CliPlayer) Lose() {
	_, _ = fmt.Fprintln(c.out, "You lost.")
}

func (c *CliPlayer) Draw() {
	_, _ = fmt.Fprintln(c.out, "Draw.")
}

func (c *CliPlayer) Bust() {
	_, _ = fmt.Fprintln(c.out, "Bust, you lose.")
}
