package blackjack

import (
	"fmt"
	"io"
)

type CliPlayer struct {
	Out io.Writer
	In  io.Reader
}

func (c *CliPlayer) ShowHand(h Hand) {
	_, _ = fmt.Fprintf(c.Out, "Your hand=%s, (Score=%d)\n", h, h.score())
}

func (c *CliPlayer) Action(action ...Action) Action {
	_, _ = fmt.Fprintln(c.Out, "What do you want to do?")
	for i, v := range action {
		_, _ = fmt.Fprintf(c.Out, "\t%d) %s\n", i+1, v)
	}
	var opt int
	_, _ = fmt.Fprintf(c.Out, "> ")
	for _, err := fmt.Fscan(c.In, &opt); err != nil; _, err = fmt.Fscan(c.In, &opt) {
		_, _ = fmt.Fprintf(c.Out, "You must enter a number between 1 and %d\n> ", len(action))
	}
	return action[opt-1]
}

func (c *CliPlayer) Prompt(msg string) {
	_, _ = fmt.Fprintln(c.Out, msg)
}

func (c *CliPlayer) Win() {
	_, _ = fmt.Fprintln(c.Out, "You won!")
}

func (c *CliPlayer) Lose() {
	_, _ = fmt.Fprintln(c.Out, "You lost.")
}

func (c *CliPlayer) Draw() {
	_, _ = fmt.Fprintln(c.Out, "Draw.")
}

func (c *CliPlayer) Bust() {
	_, _ = fmt.Fprintln(c.Out, "Bust, you lose.")
}

type StatsPlayer struct {
	wins   int
	losses int
	draws  int
	Player
}

func (s *StatsPlayer) Win() {
	s.wins++
	s.Player.Win()
}

func (s *StatsPlayer) Lose() {
	s.losses++
	s.Player.Lose()
}

func (s *StatsPlayer) Draw() {
	s.draws++
	s.Player.Draw()
}

func (s *StatsPlayer) Bust() {
	s.losses++
	s.Player.Draw()
}

func (s *StatsPlayer) String() string {
	return fmt.Sprintf("Wins=%d, Losses=%d, Draws=%d\n", s.wins, s.losses, s.draws)
}

type AI struct {
	h     Hand
	score int
}

func (a *AI) ShowHand(h Hand) {
	a.h = h
	a.score = h.score()
}

func (a *AI) Action(action ...Action) Action {
	if a.score >= 17 {
		return Hit
	}
	return Stand
}

func (a *AI) Prompt(msg string) {
	//pass
}

func (a *AI) Win() {
	//pass
}

func (a *AI) Lose() {
	//pass
}

func (a *AI) Draw() {
	//pass
}

func (a *AI) Bust() {
	//pass
}
