package blackjack

import (
	"fmt"
	"gotraining/exercise9/deck"
)

type Result string

const (
	Win  Result = "w"
	Lose Result = "l"
	Draw Result = "d"
)

var (
	ActionHit   GameAction
	ActionStand GameAction
	ActionSplit GameAction
)

func init() {
	ActionHit = GameAction{
		Name: "Hit",
		gs: func(g *Game) GameState {
			p := g.players[g.playerIdx]
			p.hands[g.handIdx] = append(p.hands[g.handIdx], g.draw())
			return playerTurn
		},
	}
	ActionStand = GameAction{
		Name: "Stand",
		gs: func(g *Game) GameState {
			p := g.players[g.playerIdx]
			if g.handIdx < len(p.hands)-1 {
				g.handIdx++
			} else {
				g.handIdx = 0
				g.playerIdx++
			}
			return playerTurn
		},
	}
	ActionSplit = GameAction{
		Name: "Split",
		gs: func(g *Game) GameState {
			p := g.players[g.playerIdx]
			if len(p.hands) < 2 {
				p.hands = append(p.hands, Hand{p.hands[0][1]})
				p.hands[0] = Hand{p.hands[0][0]}
				p.hands[0] = append(p.hands[0], g.draw())
				p.hands[1] = append(p.hands[1], g.draw())
			}
			return playerTurn
		},
	}
}

type GameAction struct {
	Name string
	gs   GameState
}

func (a GameAction) String() string {
	return a.Name
}

// Player describes the basic interactions required for blackjack
type Player interface {

	// Bet allows the player to place a wager for the round. Returning 0 skips
	// the player for round. Implementations are responsible for tracking
	// their account balance.
	Bet(shuffled bool) int

	// Action allows the player to hit or stand
	Action(hand Hand, dealer deck.Card, actions []GameAction) GameAction

	// Result pays out any winnings and indicates whether the hand was won,
	// lost, or a draw.
	Result(r Result, winnings int)

	// Prompt provides information about the game for supporting UIs. Non-interactive
	// players can safely no-op this method without loss of functionality.
	Prompt(s string)
}

type GameState func(*Game) GameState

type PlayerData struct {
	hands []Hand
	bets  []int
	Player
}

type Game struct {
	deck      []deck.Card
	dealer    Hand
	players   []*PlayerData
	playerIdx int
	handIdx   int
	state     GameState
	shuffled  bool
}

func NewGame(players ...Player) *Game {
	g := Game{deck: deck.New(deck.Quantity(3), deck.Shuffle), state: shuffle}
	for _, v := range players {
		g.players = append(g.players, &PlayerData{Player: v, hands: make([]Hand, 1, 2), bets: make([]int, 0, 2)})
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
		p.bets = append(p.bets, p.Bet(g.shuffled))
	}
	return deal
}

func deal(g *Game) GameState {
	for i := 0; i < 2; i++ {
		for _, p := range g.players {
			p.hands[0] = append(p.hands[0], g.draw())
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
	determineActions := func(h Hand) []GameAction {
		if len(h) < 2 {
			panic("Player cannot play with less than two cards in their hand")
		}
		if h[0].Value == h[1].Value {
			return []GameAction{ActionHit, ActionStand, ActionSplit}
		}
		return []GameAction{ActionHit, ActionStand}
	}

	p := g.players[g.playerIdx]
	if p.hands[g.handIdx].Score() > 21 {
		g.handIdx++
	}
	if g.handIdx < len(p.hands) {
		return p.Action(copyHand(p.hands[g.handIdx]), g.dealer[0], determineActions(p.hands[g.handIdx])).gs
	}
	g.playerIdx++
	g.handIdx = 0
	if g.playerIdx < len(g.players) {
		p := g.players[g.playerIdx]
		return p.Action(copyHand(p.hands[g.handIdx]), g.dealer[0], determineActions(p.hands[g.handIdx])).gs
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
		for i, h := range player.hands {
			playerScore := h.Score()
			switch {
			case playerScore > 21:
				player.Result(Lose, 0)
			case dealerScore > 21 || playerScore > dealerScore:
				player.Result(Win, player.bets[i]*2)
			case playerScore < dealerScore:
				player.Result(Lose, 0)
			case playerScore == dealerScore:
				player.Result(Draw, player.bets[i])
			}
		}
	}
	return nil
}

func (g *Game) reset() {
	g.playerIdx = 0
	g.dealer = nil
	g.state = shuffle
	for _, p := range g.players {
		p.hands = nil
	}
}

func (g *Game) draw() deck.Card {
	var ret deck.Card
	ret, g.deck = g.deck[0], g.deck[1:]
	return ret
}
