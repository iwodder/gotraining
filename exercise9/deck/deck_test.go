package deck

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var defaultCardDeck = []Card{
	{Clubs, Ace},
	{Clubs, Two},
	{Clubs, Three},
	{Clubs, Four},
	{Clubs, Five},
	{Clubs, Six},
	{Clubs, Seven},
	{Clubs, Eight},
	{Clubs, Nine},
	{Clubs, Ten},
	{Clubs, Jack},
	{Clubs, Queen},
	{Clubs, King},
	{Diamonds, Ace},
	{Diamonds, Two},
	{Diamonds, Three},
	{Diamonds, Four},
	{Diamonds, Five},
	{Diamonds, Six},
	{Diamonds, Seven},
	{Diamonds, Eight},
	{Diamonds, Nine},
	{Diamonds, Ten},
	{Diamonds, Jack},
	{Diamonds, Queen},
	{Diamonds, King},
	{Hearts, Ace},
	{Hearts, Two},
	{Hearts, Three},
	{Hearts, Four},
	{Hearts, Five},
	{Hearts, Six},
	{Hearts, Seven},
	{Hearts, Eight},
	{Hearts, Nine},
	{Hearts, Ten},
	{Hearts, Jack},
	{Hearts, Queen},
	{Hearts, King},
	{Spades, Ace},
	{Spades, Two},
	{Spades, Three},
	{Spades, Four},
	{Spades, Five},
	{Spades, Six},
	{Spades, Seven},
	{Spades, Eight},
	{Spades, Nine},
	{Spades, Ten},
	{Spades, Jack},
	{Spades, Queen},
	{Spades, King},
}

func Test_PrintsString(t *testing.T) {
	c := Card{Hearts, Ace}

	assert.Equal(t, "Ace of Hearts", fmt.Sprintf("%s", c))
}

func Test_DefaultDeckIsOrdered(t *testing.T) {
	assert.Equal(t, defaultCardDeck, New())
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

	assert.NotEqual(t, d, defaultCardDeck)
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
