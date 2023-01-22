package blackjack

import (
	"fmt"
	"gotraining/exercise9/deck"
)

type Game struct {
	deck      []deck.Card
	dealer    []deck.Card
	players   []*player
	playerIdx int
}

type player struct {
	won   bool
	hand  []deck.Card
	score int
	Playable
}

type Playable interface {
	HandCard(c deck.Card)
	TellScore(score int)
	Prompt(s string)
	NextMove() string
}

func New(players ...Playable) *Game {
	g := Game{deck: deck.New()}
	for i := range players {
		g.players = append(g.players, &player{
			Playable: players[i],
		})
	}
	return &g
}

func (g *Game) Play() {
	for state := Deal(g); state != nil; {
		state = state(g)
	}
}

func (g *Game) draw() deck.Card {
	ret := g.deck[0]
	g.deck = g.deck[1:]
	return ret
}

func (g *Game) nextPlayer() *player {
	if g.playerIdx >= len(g.players) {
		return nil
	}
	next := g.players[g.playerIdx]
	g.playerIdx++
	return next
}

type GameState func(*Game) GameState

func Deal(g *Game) GameState {
	for i := 0; i < 2; i++ {
		for _, p := range g.players {
			p.HandCard(g.draw())
		}
		g.dealer = append(g.dealer, g.draw())
	}
	return PlayerTurn(g.nextPlayer())
}

func PlayerTurn(p *player) GameState {
	return func(g *Game) GameState {
		stand := false
		for pScore := score(p.hand); !stand && pScore < 21; pScore = score(p.hand) {
			p.Prompt(fmt.Sprintf("Current score = %d\nDo you want to (h)it or (s)tand?\n", pScore))
			action := p.NextMove()
			switch action {
			case "h":
				p.HandCard(g.draw())
			case "s":
				stand = true
			case "default":
				p.Prompt("Unknown action, must be either (h)it or (s)tand.")
			}
		}
		if next := g.nextPlayer(); next != nil {
			return PlayerTurn(next)
		}
		return DealerTurn()
	}
}

func DealerTurn() GameState {
	return func(g *Game) GameState {
		has := func(d []deck.Card, v deck.Value) bool {
			return d[0].Value == v || d[1].Value == v
		}
		soft17 := func(d []deck.Card) bool {
			return has(d, deck.Ace) && has(d, deck.Six)
		}

		dScore := score(g.dealer)
		fmt.Printf("Dealer has %d\n", dScore)
		if dScore <= 16 || soft17(g.dealer) {
			g.dealer = append(g.dealer, g.draw())
		}
		dScore = score(g.dealer)
		fmt.Printf("Dealer has %d", dScore)
		return nil
	}
}

func score(hand []deck.Card) int {
	return scoreIter(hand, 0)
}

func scoreIter(hand []deck.Card, acc int) int {
	if len(hand) == 0 {
		return acc
	}

	card := hand[0]
	switch card.Value {
	case deck.Ace:
		s := scoreIter(hand[1:], acc+11)
		if s < 21 {
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
