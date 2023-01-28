package blackjack

import (
	"github.com/stretchr/testify/assert"
	"gotraining/exercise9/deck"
	"testing"
)

type TestPrompter struct {
	prompts []string
}

func (m *TestPrompter) Prompt(s string) {
	m.prompts = append(m.prompts, s)
}

func (m *TestPrompter) Response() string {
	return "s"
}

func Test_CreateNewGame(t *testing.T) {
	g := NewGame(NewPlayer(&TestPrompter{}))

	assert.NotNil(t, g)
	assert.Equal(t, 1, len(g.players))
}

func Test_DealsCards(t *testing.T) {
	p := NewPlayer(&TestPrompter{})
	g := NewGame(p)
	_ = Deal(g)

	assert.Equal(t, 2, len(p.hand))
}

func Test_PlayerTurn(t *testing.T) {
	spy := &TestPrompter{}
	p := NewPlayer(spy)
	g := NewGame(p)

	PlayerTurn(g)

	assert.Contains(t, spy.prompts, "Do you want to (h)it or (s)tand?")
}

func Test_DetermineWinners(t *testing.T) {

	spy := &TestPrompter{}
	p := NewPlayer(spy)
	g := NewGame(p)
	g.dealer = []deck.Card{{deck.Clubs, deck.Two}, {deck.Hearts, deck.Jack}}
	p.hand = []deck.Card{{deck.Clubs, deck.Ace}, {deck.Hearts, deck.Jack}}

	DetermineWinners(g)

	assert.Contains(t, spy.prompts, "You won!")
}

func Test_Scoring(t *testing.T) {
	tests := []struct {
		name  string
		hand  Hand
		score int
	}{
		{
			name:  "Hand of zero is zero",
			hand:  []deck.Card{},
			score: 0,
		},
		{
			name:  "2 and 3",
			hand:  []deck.Card{{deck.Clubs, deck.Two}, {deck.Spades, deck.Three}},
			score: 5,
		},
		{
			name:  "5 and 2",
			hand:  []deck.Card{{deck.Clubs, deck.Five}, {deck.Spades, deck.Two}},
			score: 7,
		},
		{
			name:  "Face cards are worth 10",
			hand:  []deck.Card{{deck.Clubs, deck.King}, {deck.Spades, deck.Two}},
			score: 12,
		},
		{
			name:  "Face cards are worth 10",
			hand:  []deck.Card{{deck.Clubs, deck.King}, {deck.Spades, deck.Jack}},
			score: 20,
		},
		{
			name:  "Two Aces is 12",
			hand:  []deck.Card{{deck.Clubs, deck.Ace}, {deck.Spades, deck.Ace}},
			score: 12,
		},
		{
			name:  "King, Queen, Ace is 21",
			hand:  []deck.Card{{deck.Spades, deck.King}, {deck.Spades, deck.Queen}, {deck.Clubs, deck.Ace}},
			score: 21,
		},
		{
			name:  "Ace, Queen, King is 21",
			hand:  []deck.Card{{deck.Clubs, deck.Ace}, {deck.Spades, deck.Queen}, {deck.Spades, deck.King}},
			score: 21,
		},
		{
			name:  "Ace, Ace, King is 12",
			hand:  []deck.Card{{deck.Spades, deck.Ace}, {deck.Spades, deck.Ace}, {deck.Clubs, deck.King}},
			score: 12,
		},
		{
			name:  "Ace, Ace, King, Five is 17",
			hand:  []deck.Card{{deck.Spades, deck.Ace}, {deck.Spades, deck.Ace}, {deck.Clubs, deck.King}, {deck.Clubs, deck.Five}},
			score: 17,
		},
		{
			name:  "Ace, Jack is 21",
			hand:  []deck.Card{{deck.Spades, deck.Ace}, {deck.Clubs, deck.Jack}},
			score: 21,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.score, tt.hand.score())
		})
	}
}
