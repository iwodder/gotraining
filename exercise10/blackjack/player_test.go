package blackjack

import (
	"github.com/stretchr/testify/assert"
	"gotraining/exercise9/deck"
	"strings"
	"testing"
)

func TestCliPlayer_ShowHand(t *testing.T) {
	var sb strings.Builder
	c := CliPlayer{out: &sb}

	c.ShowHand([]deck.Card{{deck.Hearts, deck.Ace}})

	assert.Equal(t, "Your hand=Ace of Hearts, (Score=11)\n", sb.String())
}

func TestCliPlayer_Prompt(t *testing.T) {
	var sb strings.Builder
	c := CliPlayer{out: &sb}

	c.Prompt("")

	assert.Equal(t, "\n", sb.String())
}

func TestCliPlayer_Win(t *testing.T) {
	var sb strings.Builder
	c := CliPlayer{out: &sb}

	c.Win()

	assert.Equal(t, "You won!\n", sb.String())
}

func TestCliPlayer_Lose(t *testing.T) {
	var sb strings.Builder
	c := CliPlayer{out: &sb}

	c.Lose()

	assert.Equal(t, "You lost.\n", sb.String())
}

func TestCliPlayer_Draw(t *testing.T) {
	var sb strings.Builder
	c := CliPlayer{out: &sb}

	c.Draw()

	assert.Equal(t, "Draw.\n", sb.String())
}

func TestCliPlayer_Bust(t *testing.T) {
	var sb strings.Builder
	c := CliPlayer{out: &sb}

	c.Bust()

	assert.Equal(t, "Bust, you lose.\n", sb.String())
}
