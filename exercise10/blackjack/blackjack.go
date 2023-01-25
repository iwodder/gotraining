package blackjack

import (
	"fmt"
	"gotraining/exercise9/deck"
	"strings"
)

type GameState func(*Game) GameState

type Game struct {
	deck      []deck.Card
	dealer    []deck.Card
	players   []*Player
	playerIdx int
}

func New(players ...Playable) *Game {
	g := Game{deck: deck.New(deck.Quantity(3), deck.Shuffle)}
	for i := range players {
		g.players = append(g.players, &Player{
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

func Deal(g *Game) GameState {
	for i := 0; i < 2; i++ {
		for _, p := range g.players {
			p.hand = append(p.hand, g.draw())
		}
		g.dealer = append(g.dealer, g.draw())
	}
	for _, p := range g.players {
		p.Prompt(fmt.Sprintf("Dealer Hand=%s,**HIDDEN**\nYour Hand=%s", g.dealer[0], p.showHand()))
	}
	return PlayerTurn(g.nextPlayer())
}

func PlayerTurn(p *Player) GameState {
	return func(g *Game) GameState {
		pScore := score(p.hand)
	loop:
		for ; pScore < 21; pScore = score(p.hand) {
			p.Prompt("Do you want to (h)it or (s)tand?")
			switch p.NextMove() {
			case "h":
				p.hand = append(p.hand, g.draw())
			case "s":
				break loop
			case "default":
				p.Prompt("Unknown action, must be either (h)it or (s)tand.")
			}
		}
		if pScore > 21 {
			p.won = false
			p.Prompt("Bust, you lose!")
		}
		if next := g.nextPlayer(); next != nil {
			return PlayerTurn(next)
		}
		return DealerTurn()
	}
}

func DealerTurn() GameState {
	return func(g *Game) GameState {
		ch := make(chan []deck.Card)
		go func() {
			for s := range ch {
				var dHand []string
				for _, card := range s {
					dHand = append(dHand, card.String())
				}
				hand := strings.Join(dHand, ", ")
				for _, p := range g.players {
					p.Prompt(fmt.Sprintf("Dealer Hand=%s", hand))
				}
			}
		}()
		has := func(d []deck.Card, v deck.Value) bool {
			return d[0].Value == v || d[1].Value == v
		}
		soft17 := func(d []deck.Card) bool {
			return has(d, deck.Ace) && has(d, deck.Six)
		}

		ch <- g.dealer
		dScore := score(g.dealer)
		for ; dScore <= 16 || soft17(g.dealer); dScore = score(g.dealer) {
			g.dealer = append(g.dealer, g.draw())
			ch <- g.dealer
		}
		close(ch)
		return nil
	}
}

func (g *Game) draw() deck.Card {
	var ret deck.Card
	ret, g.deck = g.deck[0], g.deck[1:]
	return ret
}

func (g *Game) nextPlayer() *Player {
	if g.playerIdx >= len(g.players) {
		return nil
	}
	next := g.players[g.playerIdx]
	g.playerIdx++
	return next
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
