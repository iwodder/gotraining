package blackjack

import (
	"fmt"
	"gotraining/exercise9/deck"
	"html/template"
	"io"
	"log"
)

var (
	add  = func(a, b int) int { return a + b }
	menu = `Dealer Hand=**HIDDEN**, {{.Dealer}}
Your Hand={{.Hand}} (score={{.Hand.Score}})
What do you want to do?
{{- range $index, $element := .Actions}}
	{{add $index 1}}) {{$element}}{{end}}`
	tmpl *template.Template
)

func init() {
	tmpl = template.Must(template.New("test").Funcs(template.FuncMap{
		"add": add,
	}).Parse(menu))
}

type CliPlayer struct {
	Out io.Writer
	In  io.Reader
}

func (c *CliPlayer) Action(hand Hand, dealer deck.Card, actions ...Action) Action {
	if err := c.showMenu(hand, dealer, actions); err != nil {
		log.Println("error showing menu:", err)
	}
	opt := c.getInput(1, len(actions))
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

func (c *CliPlayer) showMenu(hand Hand, dealer deck.Card, actions []Action) error {
	return tmpl.Execute(c.Out, struct {
		Hand    Hand
		Dealer  deck.Card
		Actions []Action
	}{
		hand, dealer, actions,
	})
}

func (c *CliPlayer) getInput(start, end int) int {
	var ret int
	_, _ = fmt.Fprintln(c.Out, "> ")
	for _, err := fmt.Fscan(c.In, &ret); ret < start || ret > end; _, err = fmt.Fscan(c.In, &ret) {
		if err != nil {
			log.Println("error getting input")
		}
		_, _ = fmt.Fprintf(c.Out, "You must enter a number between %d and %d\n> ", start, end)
	}
	return ret
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
	if hand.Score() <= 16 {
		return ActionHit
	}
	return ActionStand
}

func (a *AI) Prompt(msg string) {
	//pass
}

func (a *AI) Result(r Result) {
	//pass
}
