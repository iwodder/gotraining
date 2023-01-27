package blackjack

import (
	"fmt"
	"gotraining/exercise9/deck"
	"strings"
)

type Player struct {
	wonHand bool
	hand    []deck.Card
	Prompter
}

type Prompter interface {
	Prompt(s string)
	Response() string
}

func NewPlayer(p Prompter) *Player {
	ret := Player{
		Prompter: p,
	}
	return &ret
}

func (p *Player) acceptCard(c deck.Card) {
	p.hand = append(p.hand, c)
}

func (p *Player) clearHand() {
	p.hand = nil
}

func (p *Player) showHand() {
	var ret []string
	for _, v := range p.hand {
		ret = append(ret, v.String())
	}
	p.Prompt(fmt.Sprintf("Your hand=%s", strings.Join(ret, ", ")))
}

func (p *Player) winner() {
	p.wonHand = true
	p.Prompt("You won!")
}

func (p *Player) loser() {
	p.Prompt("You lose, better luck next time!")
}

func (p *Player) bust() {
	p.Prompt("Bust, you lose!")
}

func (p *Player) scoreHand(scorer func([]deck.Card) int) int {
	return scorer(p.hand)
}
