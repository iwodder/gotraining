package deck

import (
	"fmt"
	"math/rand"
	"time"
)

type Suit int

const (
	Any Suit = iota
	Clubs
	Diamonds
	Hearts
	Spades
)

type Value int

const (
	Ace Value = iota
	One
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Joker
)

var (
	suits   = []Suit{Clubs, Diamonds, Hearts, Spades}
	values  = []Value{Ace, Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King}
	Shuffle = func(d []Card) []Card {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(d), func(i, j int) { d[i], d[j] = d[j], d[i] })
		return d
	}
)

type Reverse []Card

func (r Reverse) Less(i, j int) bool {
	if r[i].Value < r[j].Value {
		return true
	}
	return r[i].Suit < r[j].Suit
}

func (r Reverse) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

func (r Reverse) Len() int { return len(r) }

type Card struct {
	Suit
	Value
}

type Opt func([]Card) []Card

func New(opts ...Opt) []Card {
	ret := defaultDeck()
	for _, opt := range opts {
		ret = opt(ret)
	}
	return ret
}

func defaultDeck() []Card {
	ret := make([]Card, 0, 52)
	for _, v := range suits {
		for _, vv := range values {
			ret = append(ret, Card{Suit: v, Value: vv})
		}
	}
	return ret
}

func (c Card) String() string {
	return fmt.Sprintf("%s of %s", c.Value, c.Suit)
}

func Remove(values ...Value) Opt {
	return func(d []Card) []Card {
		ret := make([]Card, 0, len(d)-len(values))
		for _, card := range d {
			if contains(values, card.Value) {
				continue
			}
			ret = append(ret, card)
		}
		return ret
	}
}

func Quantity(n int) Opt {
	return func(c []Card) []Card {
		for i := 1; i < n; i++ {
			c = append(c, defaultDeck()...)
		}
		return c
	}
}

func Jokers(n int) Opt {
	return func(c []Card) []Card {
		for i := 0; i < n; i++ {
			c = append(c, Card{Any, Joker})
		}
		return c
	}
}

func contains(values []Value, v Value) bool {
	for _, vv := range values {
		if vv == v {
			return true
		}
	}
	return false
}
