package godeck_test

import (
	"testing"

	"github.com/mbaum0/euchrego/godeck"
	"github.com/stretchr/testify/assert"
)

func TestNewDeckHas52Cards(t *testing.T) {
	deck, err := godeck.NewDeck()
	assert.Nil(t, err)
	assert.Equal(t, 52, deck.Length())
}

func TestRandomShuffleSeed(t *testing.T) {
	deck1, err := godeck.NewDeck(godeck.RandomShuffleSeed())
	assert.Nil(t, err)
	deck2, err := godeck.NewDeck(godeck.RandomShuffleSeed())
	assert.Nil(t, err)
	assert.NotEqual(t, deck1, deck2)
}

func TestShuffleSeed(t *testing.T) {
	deck1, err := godeck.NewDeck(godeck.Range(godeck.Two, godeck.King), godeck.ShuffleSeed(1))
	assert.Nil(t, err)
	deck2, err := godeck.NewDeck(godeck.Range(godeck.Two, godeck.King), godeck.ShuffleSeed(1))
	assert.Nil(t, err)
	assert.Equal(t, deck1, deck2)

	deck3, err := godeck.NewDeck(godeck.Range(godeck.Two, godeck.King), godeck.ShuffleSeed(2))
	assert.Nil(t, err)
	assert.NotEqual(t, deck1, deck3)
}

func TestAceHigh(t *testing.T) {
	deck, err := godeck.NewDeck(godeck.AceHigh())
	assert.Nil(t, err)
	assert.Equal(t, 14, deck.RankValue(godeck.Ace))
}

func TestAceLow(t *testing.T) {
	deck, err := godeck.NewDeck(godeck.AceLow())
	assert.Nil(t, err)
	assert.Equal(t, 1, deck.RankValue(godeck.Ace))
}

func TestPreShuffled(t *testing.T) {
	deck1, err := godeck.NewDeck(godeck.PreShuffled())
	assert.Nil(t, err)
	deck2, err := godeck.NewDeck(godeck.PreShuffled())
	assert.Nil(t, err)
	assert.NotEqual(t, deck1, deck2)
}

func TestInvalidRange(t *testing.T) {
	_, err := godeck.NewDeck(godeck.Range(godeck.Ace, godeck.Two))
	assert.NotNil(t, err)
	assert.Equal(t, "start must be less than or equal to end", err.Error())
}

func TestDrawCard(t *testing.T) {
	deck, err := godeck.NewDeck()
	assert.Nil(t, err)
	assert.Equal(t, 52, deck.Length())
	deck.DrawCard()
	assert.Equal(t, 51, deck.Length())
}

func TestDrawingWhenNoCardsAreInDeck(t *testing.T) {
	deck, err := godeck.NewDeck()
	assert.Nil(t, err)
	assert.Equal(t, 52, deck.Length())
	for i := 0; i < 52; i++ {
		_, err := deck.DrawCard()
		assert.Nil(t, err)
	}
	_, err = deck.DrawCard()
	assert.NotNil(t, err)
	assert.Equal(t, "no cards left in deck", err.Error())
}

func TestDrawingCardsWhenNotEnoughCardsInDeck(t *testing.T) {
	deck, err := godeck.NewDeck()
	assert.Nil(t, err)
	assert.Equal(t, 52, deck.Length())
	_, err = deck.DrawCards(53)
	assert.NotNil(t, err)
	assert.Equal(t, "not enough cards in deck, tried to draw 53 cards when only 52 are in the deck", err.Error())
}

func TestReturnCards(t *testing.T) {
	deck, err := godeck.NewDeck()
	assert.Nil(t, err)
	assert.Equal(t, 52, deck.Length())
	cards, _ := deck.DrawCards(2)
	assert.Equal(t, 50, deck.Length())
	deck.ReturnCards(cards)
	assert.Equal(t, 52, deck.Length())
}

func TestReturnCard(t *testing.T) {
	deck, err := godeck.NewDeck()
	assert.Nil(t, err)
	assert.Equal(t, 52, deck.Length())
	card, _ := deck.DrawCard()
	assert.Equal(t, 51, deck.Length())
	deck.ReturnCard(card)
	assert.Equal(t, 52, deck.Length())
}

func TestSetRankValue(t *testing.T) {
	deck, err := godeck.NewDeck()
	assert.Nil(t, err)
	assert.Equal(t, 52, deck.Length())
	deck.SetRankValue(godeck.Ace, 10)
	assert.Equal(t, 10, deck.RankValue(godeck.Ace))
}

func TestSetRankValues(t *testing.T) {
	deck, err := godeck.NewDeck()
	assert.Nil(t, err)
	assert.Equal(t, 52, deck.Length())
	deck.SetRankValues(map[godeck.Rank]int{
		godeck.Ace: 10,
	})
	assert.Equal(t, 10, deck.RankValue(godeck.Ace))
}

func TestResetRankValues(t *testing.T) {
	deck, err := godeck.NewDeck()
	assert.Nil(t, err)
	assert.Equal(t, 52, deck.Length())
	deck.SetRankValue(godeck.Ace, 10)
	assert.Equal(t, 10, deck.RankValue(godeck.Ace))
	deck.ResetRankValues()
	assert.Equal(t, 14, deck.RankValue(godeck.Ace))
}
