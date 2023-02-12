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
	{{add $index 1}}) {{$element}}{{end}}
`
	tmpl *template.Template
)

func init() {
	tmpl = template.Must(template.New("test").Funcs(template.FuncMap{
		"add": add,
	}).Parse(menu))
}

type CliPlayer struct {
	bank int
	out  io.Writer
	in   io.Reader
}

func NewCLIPlayer(in io.Reader, out io.Writer) *CliPlayer {
	return &CliPlayer{
		in:   in,
		out:  out,
		bank: 1000,
	}
}

func (c *CliPlayer) Bet(shuffled bool) int {
	_, _ = fmt.Println(c.out, "How much do you want to bet?")
	_, _ = fmt.Println(c.out, "The deck was shuffled")
	var amt int
	_, _ = fmt.Fprintln(c.out, "> ")
	for _, err := fmt.Fscan(c.in, &amt); err != nil; _, err = fmt.Fscan(c.in, &amt) {
		_, _ = fmt.Fprintf(c.out, "Error getting input.\nYou must enter a number\n> ")
	}
	return amt
}

func (c *CliPlayer) Action(hand Hand, dealer deck.Card, actions []Action) Action {
	if err := c.showMenu(hand, dealer, actions); err != nil {
		log.Println("error showing menu:", err)
	}
	opt := c.getInput(1, len(actions))
	return actions[opt-1]
}

func (c *CliPlayer) Prompt(msg string) {
	_, _ = fmt.Fprintln(c.out, msg)
}

func (c *CliPlayer) Result(r Result, winnings int) {

}

func (c *CliPlayer) showMenu(hand Hand, dealer deck.Card, actions []Action) error {
	return tmpl.Execute(c.out, struct {
		Hand    Hand
		Dealer  deck.Card
		Actions []Action
	}{
		hand, dealer, actions,
	})
}

func (c *CliPlayer) getInput(start, end int) int {
	var ret int
	_, _ = fmt.Fprintln(c.out, "> ")
	for _, err := fmt.Fscan(c.in, &ret); ret < start || ret > end; _, err = fmt.Fscan(c.in, &ret) {
		if err != nil {
			log.Println("error getting input")
		}
		_, _ = fmt.Fprintf(c.out, "You must enter a number between %d and %d\n> ", start, end)
	}
	return ret
}

type StatsPlayer struct {
	wins   int
	losses int
	draws  int
	Player
}

func (s *StatsPlayer) Result(r Result, winnings int) {
	switch r {
	case Win:
		s.wins++
	case Lose:
		s.losses++
	case Draw:
		s.draws++
	}
	s.Player.Result(r, winnings)
}

func (s *StatsPlayer) String() string {
	return fmt.Sprintf("Wins=%d, Losses=%d, Draws=%d\n", s.wins, s.losses, s.draws)
}

type AI struct {
	bank int
}

func (a *AI) Bet(shuffled bool) int {
	return 1
}

func (a *AI) Result(r Result, winnings int) {
	a.bank += winnings
}

func (a *AI) Prompt(s string) {
	//no-op
}

func (a *AI) Action(hand Hand, dealer deck.Card, actions []Action) Action {
	if hand.Score() <= 16 {
		return ActionHit
	}
	return ActionStand
}
