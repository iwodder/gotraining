package blackjack

import (
	"fmt"
	"gotraining/exercise9/deck"
	"strings"
)

type Hand []deck.Card

func (h Hand) String() string {
	var ret []string
	for _, v := range h {
		ret = append(ret, v.String())
	}
	return strings.Join(ret, ", ")
}

func (h Hand) score() int {
	return scoreIter(h, 0)
}

func scoreIter(hand []deck.Card, acc int) int {
	if len(hand) == 0 {
		return acc
	}
	card := hand[0]
	switch card.Value {
	case deck.Ace:
		s := scoreIter(hand[1:], acc+11)
		if s <= 21 {
			return s
		} else {
			return scoreIter(hand[1:], acc+1)
		}
	case deck.Jack, deck.Queen, deck.King:
		return scoreIter(hand[1:], acc+10)
	default:
		return scoreIter(hand[1:], acc+int(card.Value))
	}
}

type Player struct {
	hand Hand
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

func (p *Player) giveCard(c deck.Card) {
	p.hand = append(p.hand, c)
	p.ShowHand()
}

func (p *Player) clearHand() {
	p.hand = nil
}

func (p *Player) ShowHand() {
	p.Prompt(fmt.Sprintf("Your hand=%s", p.hand))
}

func (p *Player) Won() {
	p.Prompt("You won!")
}

func (p *Player) Lost() {
	p.Prompt("You lose, better luck next time!")
}

func (p *Player) Bust() {
	p.Prompt("Bust, you lose!")
}

func (p *Player) Score() int {
	return p.hand.score()
}

func (p *Player) Draw() {
	p.Prompt("Draw")
}
