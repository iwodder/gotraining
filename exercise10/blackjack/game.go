package blackjack

import (
	"fmt"
	"gotraining/exercise9/deck"
)

type GameState func(*Game) GameState

type Game struct {
	deck      []deck.Card
	dealer    Hand
	players   []*Player
	playerIdx int
}

func NewGame(players ...*Player) *Game {
	g := Game{deck: deck.New(deck.Quantity(3), deck.Shuffle)}
	for _, v := range players {
		g.players = append(g.players, v)
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

func Deal(g *Game) GameState {
	for i := 0; i < 2; i++ {
		for _, p := range g.players {
			p.giveCard(g.draw())
		}
		g.dealer = append(g.dealer, g.draw())
	}
	for _, p := range g.players {
		p.Prompt(fmt.Sprintf("Dealer Hand=**HIDDEN**, %s", g.dealer[0]))
	}
	return PlayerTurn
}

func PlayerTurn(g *Game) GameState {
	if p := g.nextPlayer(); p != nil {
		stand := false
		for s := p.Score(); s < 21 && !stand; s = p.Score() {
			p.Prompt("Do you want to (h)it or (s)tand?")
			switch p.Response() {
			case "h":
				p.giveCard(g.draw())
			case "s":
				stand = true
			default:
				p.Prompt("Unknown action, must be either (h)it or (s)tand.")
			}
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
		playerScore := player.Score()
		switch {
		case playerScore > 21:
			player.Bust()
		case dealerScore > 21:
			player.Prompt("Dealer busted! You win!")
		case playerScore > dealerScore:
			player.Won()
		case playerScore < dealerScore:
			player.Lost()
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
