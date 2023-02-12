package blackjack

import (
	"fmt"
	"gotraining/exercise9/deck"
	"strings"
)

type Result string

const (
	Win  Result = "w"
	Lose Result = "l"
	Draw Result = "d"
)

var (
	ActionHit   Action
	ActionStand Action
)

func init() {
	ActionHit = Action{
		Name: "Hit",
		gs: func(g *Game) GameState {
			p := g.players[g.playerIdx]
			p.hand = append(p.hand, g.draw())
			if p.hand.Score() > 21 {
				g.playerIdx++
			}
			return playerTurn
		},
	}
	ActionStand = Action{
		Name: "Stand",
		gs: func(g *Game) GameState {
			g.playerIdx++
			return playerTurn
		},
	}
}

type Action struct {
	Name string
	gs   GameState
}

func (a Action) String() string {
	return a.Name
}

// Player describes the basic interactions required for blackjack
type Player interface {

	// Bet allows the player to place a wager for the round. Returning 0 skips
	// the player for round. Implementations are responsible for tracking
	// their account balance.
	Bet(shuffled bool) int

	// Action allows the player to hit or stand
	Action(hand Hand, dealer deck.Card, actions []Action) Action

	// Result pays out any winnings and indicates whether the hand was won,
	// lost, or a draw.
	Result(r Result, winnings int)

	// Prompt provides information about the game for supporting UIs. Non-interactive
	// players can safely no-op this method without loss of functionality.
	Prompt(s string)
}

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

type PlayerData struct {
	hand Hand
	bet  int
	Player
}

type GameState func(*Game) GameState

type Game struct {
	deck      []deck.Card
	dealer    Hand
	players   []*PlayerData
	playerIdx int
	state     GameState
	shuffled  bool
}

func NewGame(players ...Player) *Game {
	g := Game{deck: deck.New(deck.Quantity(3), deck.Shuffle), state: shuffle}
	for _, v := range players {
		g.players = append(g.players, &PlayerData{Player: v})
	}
	return &g
}

func (g *Game) Play(rounds int) {
	for ; rounds > 0; rounds-- {
		for g.state != nil {
			g.state = g.state(g)
		}
		g.reset()
	}
}

func shuffle(g *Game) GameState {
	switch {
	case len(g.deck) < 10:
		g.deck = append(g.deck, deck.New(deck.Quantity(3), deck.Shuffle)...)
		g.shuffled = true
	default:
		g.shuffled = false
	}
	return bet
}

func bet(g *Game) GameState {
	for _, p := range g.players {
		p.bet = p.Bet(g.shuffled)
	}
	return deal
}

func deal(g *Game) GameState {
	for i := 0; i < 2; i++ {
		for _, p := range g.players {
			p.hand = append(p.hand, g.draw())
		}
		g.dealer = append(g.dealer, g.draw())
	}
	return playerTurn
}

func playerTurn(g *Game) GameState {
	copyHand := func(h Hand) Hand {
		ph := make([]deck.Card, len(h))
		copy(ph, h)
		return ph
	}
	determineActions := func(h Hand) []Action {
		return []Action{ActionHit, ActionStand}
	}

	if g.playerIdx < len(g.players) {
		p := g.players[g.playerIdx]
		return p.Action(copyHand(p.hand), g.dealer[0], determineActions(p.hand)).gs
	}
	return dealerTurn
}

func dealerTurn(g *Game) GameState {
	has := func(d []deck.Card, v deck.Value) bool {
		return d[0].Value == v || d[1].Value == v
	}
	soft17 := func(d []deck.Card) bool {
		return has(d, deck.Ace) && has(d, deck.Six)
	}
	for _, p := range g.players {
		p.Prompt(fmt.Sprintf("Dealer Hand=%s", g.dealer))
	}
	if soft17(g.dealer) {
		g.dealer = append(g.dealer, g.draw())
	}
	for dScore := g.dealer.Score(); dScore <= 16; dScore = g.dealer.Score() {
		g.dealer = append(g.dealer, g.draw())
		for _, p := range g.players {
			p.Prompt(fmt.Sprintf("Dealer Hand=%s", g.dealer))
		}
	}
	return determineWinners
}

func determineWinners(g *Game) GameState {
	dealerScore := g.dealer.Score()
	for _, player := range g.players {
		playerScore := player.hand.Score()
		switch {
		case playerScore > 21:
			player.Result(Lose, 0)
		case dealerScore > 21 || playerScore > dealerScore:
			player.Result(Win, player.bet*2)
		case playerScore < dealerScore:
			player.Result(Lose, 0)
		case playerScore == dealerScore:
			player.Result(Draw, player.bet)
		}
	}
	return nil
}

func (g *Game) reset() {
	g.playerIdx = 0
	g.dealer = nil
	g.state = shuffle
	for _, p := range g.players {
		p.hand = nil
	}
}

func (g *Game) draw() deck.Card {
	var ret deck.Card
	ret, g.deck = g.deck[0], g.deck[1:]
	return ret
}
