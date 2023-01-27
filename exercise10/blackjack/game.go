package blackjack

import (
	"fmt"
	"gotraining/exercise9/deck"
	"strings"
)

type GameState func(*Game) GameState

type Game struct {
	deck       []deck.Card
	dealerHand []deck.Card
	players    []*Player
	playerIdx  int
	rounds     int
}

func NewGame(players ...*Player) *Game {
	g := Game{deck: deck.New(deck.Quantity(3), deck.Shuffle)}
	for _, v := range players {
		g.players = append(g.players, v)
	}
	return &g
}

func (g *Game) Play(rounds int) {
	g.rounds = rounds
	for state := Deal(g); state != nil; {
		state = state(g)
	}
}

func Deal(g *Game) GameState {
	for i := 0; i < 2; i++ {
		for _, p := range g.players {
			p.acceptCard(g.draw())
		}
		g.dealerHand = append(g.dealerHand, g.draw())
	}
	for _, p := range g.players {
		p.Prompt(fmt.Sprintf("Dealer Hand=%s,**HIDDEN**", g.dealerHand[0]))
		p.showHand()
	}
	return PlayerTurn
}

func PlayerTurn(g *Game) GameState {
	if p := g.nextPlayer(); p != nil {
		pScore := p.scoreHand(score)
		for stand := false; pScore < 21 && !stand; pScore = p.scoreHand(score) {
			p.Prompt("Do you want to (h)it or (s)tand?")
			switch p.Response() {
			case "h":
				p.acceptCard(g.draw())
				p.showHand()
			case "s":
				stand = true
			default:
				p.Prompt("Unknown action, must be either (h)it or (s)tand.")
			}
		}
		if pScore > 21 {
			p.bust()
		}
		return PlayerTurn
	}
	return DealerTurn
}

func DealerTurn(g *Game) GameState {
	has := func(d []deck.Card, v deck.Value) bool {
		return d[0].Value == v || d[1].Value == v
	}
	soft17 := func(d []deck.Card) bool {
		return has(d, deck.Ace) && has(d, deck.Six)
	}

	dScore := score(g.dealerHand)
	for ; dScore <= 16 || soft17(g.dealerHand); dScore = score(g.dealerHand) {
		g.dealerHand = append(g.dealerHand, g.draw())
		var dHand []string
		for _, card := range g.dealerHand {
			dHand = append(dHand, card.String())
		}
		hand := strings.Join(dHand, ", ")
		for _, p := range g.players {
			p.Prompt(fmt.Sprintf("Dealer Hand=%s", hand))
		}
	}
	return DetermineWinners(dScore)
}

func DetermineWinners(dealerScore int) GameState {
	return func(g *Game) GameState {
		for _, player := range g.players {
			playerScore := player.scoreHand(score)
			switch {
			case playerScore > 21:
			case dealerScore > 21:
				player.Prompt("Dealer busted! You win!")
			case playerScore > dealerScore:
				player.winner()
			case playerScore < dealerScore:
				player.loser()
			case playerScore == dealerScore:
				player.Prompt("Draw")
			}
		}
		return PlayAgain
	}
}

func PlayAgain(g *Game) GameState {
	g.rounds--
	if g.rounds > 1 {
		g.reset()
		return Deal
	}
	return nil
}

func (g *Game) reset() {
	g.playerIdx = 0
	g.dealerHand = nil
	for _, p := range g.players {
		p.clearHand()
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
