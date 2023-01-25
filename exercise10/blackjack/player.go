package blackjack

import (
	"gotraining/exercise9/deck"
	"strings"
)

type Player struct {
	won   bool
	hand  []deck.Card
	score int
	Playable
}

type Playable interface {
	Prompt(s string)
	NextMove() string
}

func (p *Player) showHand() string {
	var ret []string
	for _, v := range p.hand {
		ret = append(ret, v.String())
	}
	return strings.Join(ret, ", ")
}
