package blackjack

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gotraining/exercise9/deck"
	"testing"
)

type StringPrompter struct {
	prompt string
}

func (sp *StringPrompter) Prompt(s string) {
	sp.prompt = s
}

func (sp *StringPrompter) Response() string {
	return ""
}

func Test_PlayerLoses(t *testing.T) {
	sp := &StringPrompter{}
	p := Player{
		Prompter: sp,
	}

	p.Lost()

	assert.Equal(t, "You lose, better luck next time!", sp.prompt)
}

func Test_PlayerDraw(t *testing.T) {
	sp := &StringPrompter{}
	p := Player{
		Prompter: sp,
	}

	p.Draw()

	assert.Equal(t, "Draw", sp.prompt)
}

func Test_acceptCardDisplaysHand(t *testing.T) {
	sp := &StringPrompter{}
	p := Player{
		Prompter: sp,
	}

	hand := Hand([]deck.Card{{Suit: deck.Hearts, Value: deck.Ace}})
	p.giveCard(deck.Card{Suit: deck.Hearts, Value: deck.Ace})

	assert.Equal(t, fmt.Sprintf("Your hand=%s", hand), sp.prompt)
}
