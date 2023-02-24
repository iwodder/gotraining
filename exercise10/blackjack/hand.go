package blackjack

import (
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

func (h Hand) Score() int {
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
