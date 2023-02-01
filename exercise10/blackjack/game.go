package blackjack

import (
	"fmt"
	"gotraining/exercise9/deck"
	"strings"
)

func init() {
	Hit = Action{Name: "Hit", gs: hit}
	Stand = Action{Name: "Stand", gs: stand}
}

var (
	Hit   Action
	Stand Action
)

type Action struct {
	Name string
	gs   GameState
}

func (a Action) String() string {
	return a.Name
}

type Player interface {
	ShowHand(h Hand)
	Action(actions ...Action) Action
	Prompt(msg string)
	Win()
	Lose()
	Draw()
	Bust()
}

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
		for state := Deal(g); state != nil; {
			state = state(g)
		}
		g.reset()
	}
}

func Shuffle(g *Game) GameState {
	g.deck = append(g.deck, deck.New(deck.Quantity(3), deck.Shuffle)...)
	return Deal
}

func Deal(g *Game) GameState {
	if len(g.deck) < 10 {
		return Shuffle
	}
	for i := 0; i < 2; i++ {
		for _, p := range g.players {
			p.hand = append(p.hand, g.draw())
		}
		g.dealer = append(g.dealer, g.draw())
	}
	for _, p := range g.players {
		p.ShowHand(p.hand)
		p.Prompt(fmt.Sprintf("Dealer Hand=**HIDDEN**, %s", g.dealer[0]))
	}
	return PlayerTurn
}

func PlayerTurn(g *Game) GameState {
	if g.playerIdx < len(g.players) {
		p := g.players[g.playerIdx]
		return p.Action(Hit, Stand).gs
	}
	return DealerTurn

}

func hit(g *Game) GameState {
	p := g.players[g.playerIdx]
	p.hand = append(p.hand, g.draw())
	p.ShowHand(p.hand)
	if p.hand.score() > 21 {
		g.playerIdx++
	}
	return PlayerTurn
}

func stand(g *Game) GameState {
	g.playerIdx++
	return PlayerTurn
}

func DealerTurn(g *Game) GameState {
	has := func(d []deck.Card, v deck.Value) bool {
		return d[0].Value == v || d[1].Value == v
	}
	soft17 := func(d []deck.Card) bool {
		return has(d, deck.Ace) && has(d, deck.Six)
	}

	for _, p := range g.players {
		p.Prompt(fmt.Sprintf("Dealer Hand=%s", g.dealer))
	}
	for dScore := g.dealer.score(); dScore <= 16 || soft17(g.dealer); dScore = g.dealer.score() {
		g.dealer = append(g.dealer, g.draw())
		for _, p := range g.players {
			p.Prompt(fmt.Sprintf("Dealer Hand=%s", g.dealer))
		}
	}
	return DetermineWinners
}

func DetermineWinners(g *Game) GameState {
	dealerScore := g.dealer.score()
	for _, player := range g.players {
		playerScore := player.hand.score()
		switch {
		case playerScore > 21:
			player.Bust()
		case dealerScore > 21:
			player.Prompt("Dealer busted! You win!")
			player.Win()
		case playerScore > dealerScore:
			player.Win()
		case playerScore < dealerScore:
			player.Lose()
		case playerScore == dealerScore:
			player.Draw()
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

func (g *Game) nextPlayer() *PlayerHand {
	if g.playerIdx >= len(g.players) {
		return nil
	}
	next := g.players[g.playerIdx]
	g.playerIdx++
	return next
}
