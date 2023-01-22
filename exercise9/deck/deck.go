package deck

import (
	"fmt"
	"math/rand"
	"sort"
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
	Ace Value = iota + 1
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
	DefaultSort = func(c []Card) func(i, j int) bool {
		return func(i, j int) bool {
			if c[i].Value == c[j].Value {
				return c[i].Suit < c[j].Suit
			}
			return c[i].Value < c[j].Value
		}
	}
)

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
	if c.Value == Joker {
		return "Joker"
	}
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

func Sort(f func([]Card) func(i, j int) bool) Opt {
	return func(cards []Card) []Card {
		sort.Slice(cards, f(cards))
		return cards
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
