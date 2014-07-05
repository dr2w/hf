package deck

import (
	"reflect"
	"testing"

	"dr2w.com/hf/model/card"
)

func TestNew(t *testing.T) {
	d := New()
	// Assert total # cards == 53
	if len(d) != 53 {
		t.Errorf("Deck Size != 53 (%d)", len(d))
	}
	// Assert # of cards per suit == 13
	for _, suit := range card.Suits {
		var count int
		for _, c := range d {
			if c.Suit == suit {
				count++
			}
		}
		if count != 13 {
			t.Errorf("Found %d cards, not 13 for suit %d", count, suit)
		}
	}
	// Assert # of cards per value = 4
	for _, value := range card.SuitedValues {
		var count int
		for _, c := range d {
			if c.Value == value {
				count++
			}
		}
		if count != 4 {
			t.Errorf("Found %d cards, not 4 for value %d", count, value)
		}
	}
}

var dealTests = []struct {
	name  string
	deck  Deck
	n     int
	dealt card.Set
	left  int
	err   bool
}{
	{
		"Empty Deck",
		Deck{},
		0,
		card.Set{},
		0,
		false,
	},
	{
		"Removing from an Empty Deck",
		Deck{},
		1,
		card.Set{},
		0,
		true,
	},
	{
		"Removing all cards",
		Deck{card.Card{0, 0}, card.Card{1, 1}},
		2,
		card.Set{card.Card{0, 0}, card.Card{1, 1}},
		0,
		false,
	},
	{
		"Removing some cards",
		Deck{card.Card{0, 0}, card.Card{1, 1}, card.Card{2, 2}},
		1,
		card.Set{card.Card{0, 0}},
		2,
		false,
	},
}

func TestDeal(t *testing.T) {
	for _, test := range dealTests {
		dealt, err := test.deck.Deal(test.n)
		switch {
		case err == nil && test.err:
			t.Errorf("%s: Expected an error, but got none.", test.name)
		case err != nil && !test.err:
			t.Errorf("%s: Unexpected error (%s)", err.Error())
		case len(test.deck) != test.left:
			t.Errorf("%s: Want %d remaining cards, got %d", test.left, len(test.deck))
		case !reflect.DeepEqual(dealt, test.dealt):
			t.Errorf("%s: Want %v, got %v", dealt, test.dealt)
		}
	}
}

var shuffleTests = []struct {
	name string
	deck Deck
	n    int
}{
	{
		"Empty Deck",
		Deck{},
		0,
	},
	{
		"Single Card",
		Deck{card.Card{0, 0}},
		1,
	},
	{
		"Multiple Cards",
		Deck{card.Card{0, 0}, card.Card{1, 1}, card.Card{2, 2}, card.Card{3, 3}},
		50,
	},
}

func TestShuffle(t *testing.T) {
	for _, test := range shuffleTests {
		deckCopy := test.deck
		deckCopy.Shuffle(test.n)
		if len(deckCopy) != len(test.deck) {
			t.Errorf("Deck length changed during shuffle (%d => %d)", len(test.deck), len(deckCopy))
		}
		if test.n > 1 && !reflect.DeepEqual(test.deck, deckCopy) {
			t.Errorf("Deck wasn't shuffled: %v", test.deck)
		}
	}
}
