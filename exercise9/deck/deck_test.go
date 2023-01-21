package deck

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_PrintsString(t *testing.T) {
	tests := []struct {
		name string
		arg  Card
	}{
		{
			"Ace of Hearts",
			Card{Hearts, Ace},
		},
		{
			"Two of Clubs",
			Card{Clubs, Two},
		},
		{
			"Joker",
			Card{Any, Joker},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.name, fmt.Sprintf("%s", tt.arg))
		})
	}
}

func Test_DefaultDeckIsOrdered(t *testing.T) {
	assertDeckHasDefaultOrder := func(t *testing.T, d []Card) {
		deckIdx := 0
		for _, s := range suits {
			for _, v := range values {
				assert.Equal(t, d[deckIdx], Card{Suit: s, Value: v}, "New deck should have default order")
				deckIdx++
			}
		}
	}
	d := New()

	assert.Equal(t, 52, len(d), "New deck should have 52 cards")
	assertDeckHasDefaultOrder(t, d)
}

func Test_CanFilterCards(t *testing.T) {
	d := New(Remove(Ten, Nine))

	for _, v := range d {
		assert.NotEqual(t, Ten, v.Value)
		assert.NotEqual(t, Nine, v.Value)
	}
}

func Test_CanShuffleStartingDeck(t *testing.T) {
	d := New(Shuffle)
	d1 := New()

	assert.NotEqual(t, d, d1)
}

func Test_MakeMultipleDecks(t *testing.T) {
	d := New(Quantity(3))

	assert.Equal(t, 156, len(d))
}

func Test_CanAddJokers(t *testing.T) {
	d := New(Jokers(2))

	exp := Card{
		Any,
		Joker,
	}

	assert.Equal(t, exp, d[len(d)-1])
	assert.Equal(t, exp, d[len(d)-2])
}

func Test_DefaultSort(t *testing.T) {
	d := New(Sort(DefaultSort))

	assert.Equal(t, Card{Suit: Clubs, Value: Ace}, d[0])
}
