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

type Player interface {
	Action(hand Hand, dealer deck.Card, actions ...Action) Action
	Result(r Result)
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

type PlayerHand struct {
	hand Hand
	Player
}

type GameState func(*Game) GameState

type Game struct {
	deck      []deck.Card
	dealer    Hand
	players   []*PlayerHand
	playerIdx int
}

func NewGame(players ...Player) *Game {
	g := Game{deck: deck.New(deck.Quantity(3), deck.Shuffle)}
	for _, v := range players {
		g.players = append(g.players, &PlayerHand{Player: v})
	}
	return &g
}

func (g *Game) Play(rounds int) {
	for ; rounds > 0; rounds-- {
		for state := deal(g); state != nil; {
			state = state(g)
		}
		g.reset()
	}
}

func shuffle(g *Game) GameState {
	g.deck = append(g.deck, deck.New(deck.Quantity(3), deck.Shuffle)...)
	return deal
}

func deal(g *Game) GameState {
	if len(g.deck) < 10 {
		return shuffle
	}
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
		return p.Action(copyHand(p.hand), g.dealer[0], determineActions(p.hand)...).gs
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
	for dScore := g.dealer.Score(); dScore <= 16 || soft17(g.dealer); dScore = g.dealer.Score() {
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
			player.Result(Lose)
		case dealerScore > 21:
			player.Result(Win)
		case playerScore > dealerScore:
			player.Result(Win)
		case playerScore < dealerScore:
			player.Result(Lose)
		case playerScore == dealerScore:
			player.Result(Draw)
		}
	}
	return nil
}

func (g *Game) reset() {
	g.playerIdx = 0
	g.dealer = nil
	for _, p := range g.players {
		p.hand = nil
	}
}

func (g *Game) draw() deck.Card {
	var ret deck.Card
	ret, g.deck = g.deck[0], g.deck[1:]
	return ret
}
