package blackjack

import (
	"github.com/stretchr/testify/assert"
	"gotraining/exercise9/deck"
	"strings"
	"testing"
)

type SpyPlayer struct {
	h               Hand
	possibleActions []GameAction
	action          GameAction
}

func (s *SpyPlayer) Bet(shuffled bool) int {
	return 1
}

func (s *SpyPlayer) Action(hand Hand, dealer deck.Card, actions []GameAction) GameAction {
	s.h = hand
	s.possibleActions = actions
	return s.action
}

func (s *SpyPlayer) Prompt(msg string) {

}

func (s *SpyPlayer) Result(r Result, winnings int) {

}

func Test_PlayerTurn(t *testing.T) {
	var s SpyPlayer
	g := NewGame(&s)
	deal(g)

	playerTurn(g)

	assert.Equal(t, 2, len(s.possibleActions))
}

func Test_PlayerAction(t *testing.T) {
	var sb strings.Builder
	cli := CliPlayer{out: &sb, in: strings.NewReader("1")}

	if err := cli.showMenu(
		[]deck.Card{{deck.Hearts, deck.Ace}, {deck.Clubs, deck.Jack}},
		deck.Card{Suit: deck.Hearts, Value: deck.King},
		[]GameAction{ActionStand}); err != nil {
		assert.Fail(t, "Unexpected error showing menu: %s", err)
	}

	exp := `Dealer Hand=**HIDDEN**, King of Hearts
Your Hand=Ace of Hearts, Jack of Clubs (score=21)
What do you want to do?
	1) Stand
`
	assert.Equal(t, exp, sb.String())
}

func Test_PlayerInput(t *testing.T) {
	var sb strings.Builder
	cli := CliPlayer{out: &sb, in: strings.NewReader("0\n6\n1")}

	assert.Equal(t, 1, cli.getInput(1, 5))
}

func Test_Bet(t *testing.T) {
	s := SpyPlayer{}
	g := NewGame(&s)
	bet(g)

	if g.players[0].bets[0] != 1 {
		t.Errorf("Should've set the player bet, wanted 1, got %d", g.players[0].bets)
	}
}

func Test_CLIPlayerBetWhenShuffled(t *testing.T) {
	var sb strings.Builder
	p := NewCLIPlayer(strings.NewReader("1\n"), &sb)
	amt := p.Bet(false)

	assert.Equal(t, 1, amt)
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

func Test_SplitHand(t *testing.T) {
	g := NewGame()
	g.players = append(g.players, &PlayerData{
		hands:  []Hand{{deck.Card{Suit: deck.Hearts, Value: deck.King}, deck.Card{Suit: deck.Clubs, Value: deck.King}}},
		bets:   []int{1},
		Player: &SpyPlayer{},
	})

	ActionSplit.gs(g)

	assert.Equal(t, 2, len(g.players[0].hands))
	assert.Equal(t, 2, len(g.players[0].hands[0]))
	assert.Equal(t, 2, len(g.players[0].hands[1]))
}

func Test_PlayerPromptedToSplit(t *testing.T) {
	spy := &SpyPlayer{}
	g := NewGame()
	g.dealer = Hand{deck.Card{Suit: deck.Hearts, Value: deck.King}, deck.Card{Suit: deck.Clubs, Value: deck.King}}
	g.players = append(g.players, &PlayerData{
		hands:  []Hand{{deck.Card{Suit: deck.Hearts, Value: deck.King}, deck.Card{Suit: deck.Clubs, Value: deck.King}}},
		bets:   []int{1},
		Player: spy,
	})

	playerTurn(g)

	for _, v := range spy.possibleActions {
		if v.Name == ActionSplit.Name {
			return
		}
	}
	t.Fatal("Expected to find split action")
}

func Test_PlaysBothHands(t *testing.T) {
	s := &SpyPlayer{}
	s.action = ActionHit
	g := NewGame()
	g.deck = []deck.Card{{Suit: deck.Hearts, Value: deck.Three}, {Suit: deck.Hearts, Value: deck.Four}}
	g.dealer = Hand{deck.Card{Suit: deck.Hearts, Value: deck.King}, deck.Card{Suit: deck.Clubs, Value: deck.King}}
	g.players = append(g.players, &PlayerData{
		hands: []Hand{
			{deck.Card{Suit: deck.Hearts, Value: deck.King}, deck.Card{Suit: deck.Clubs, Value: deck.King}},
			{deck.Card{Suit: deck.Spades, Value: deck.Jack}, deck.Card{Suit: deck.Clubs, Value: deck.Six}},
		},
		bets:   []int{1},
		Player: s,
	})

	nextState := playerTurn(g)
	nextState = nextState(g)
	assert.Equal(t, 3, len(g.players[0].hands[0]))
	assert.Equal(t, 23, g.players[0].hands[0].Score())

	s.action = ActionStand
	nextState = nextState(g)
	nextState(g)
	assert.Equal(t, 2, len(g.players[0].hands[1]))
	assert.Equal(t, 16, g.players[0].hands[1].Score())
	assert.Equal(t, 1, g.playerIdx)
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
