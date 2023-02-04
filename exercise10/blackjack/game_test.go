package blackjack

import (
	"github.com/stretchr/testify/assert"
	"gotraining/exercise9/deck"
	"strings"
	"testing"
)

type SpyPlayer struct {
	h       Hand
	actions int
}

func (s *SpyPlayer) Action(hand Hand, dealer deck.Card, actions ...Action) Action {
	s.h = hand
	s.actions = len(actions)
	return ActionHit
}

func (s *SpyPlayer) Prompt(msg string) {

}

func (s *SpyPlayer) Result(r Result) {

}

func Test_PlayerTurn(t *testing.T) {
	var s SpyPlayer
	g := NewGame(&s)
	deal(g)

	playerTurn(g)

	assert.Equal(t, 2, s.actions)
}

func Test_PlayerAction(t *testing.T) {
	var sb strings.Builder
	cli := CliPlayer{Out: &sb, In: strings.NewReader("1")}

	cli.showMenu(
		[]deck.Card{{deck.Hearts, deck.Ace}, {deck.Clubs, deck.Jack}},
		deck.Card{Suit: deck.Hearts, Value: deck.King},
		[]Action{ActionStand})

	exp := `Dealer Hand=**HIDDEN**, King of Hearts
Your Hand=Ace of Hearts, Jack of Clubs (score=21)
What do you want to do?
	1) Stand`
	assert.Equal(t, exp, sb.String())
}

func Test_PlayerInput(t *testing.T) {
	var sb strings.Builder
	cli := CliPlayer{Out: &sb, In: strings.NewReader("0\n6\n1")}

	assert.Equal(t, 1, cli.getInput(1, 5))
}

func Test_DealerHitsUntilOver16(t *testing.T) {
	g := NewGame()

	g.dealer = []deck.Card{{deck.Hearts, deck.Six}, {deck.Clubs, deck.Six}}

	dealerTurn(g)

	assert.Greater(t, g.dealer.Score(), 16)
}

func Test_DealerHitsOnSoft17(t *testing.T) {
	g := NewGame()

	g.dealer = []deck.Card{{deck.Hearts, deck.Six}, {deck.Clubs, deck.Ace}}

	dealerTurn(g)

	assert.Greater(t, g.dealer.Score(), 16)
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
			assert.Equal(t, tt.score, tt.hand.Score())
		})
	}
}
