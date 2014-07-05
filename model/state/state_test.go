package state

import (
	"testing"

	"dr2w.com/hf/model/card"
	"dr2w.com/hf/model/deck"
	"dr2w.com/hf/model/hand"
	"dr2w.com/hf/model/seat"
	"dr2w.com/hf/model/trick"
)

var findCardTests = []struct {
	name  string
	seek  card.Card
	state State
	err   bool
}{
	{
		"No card",
		card.Card{card.Ace, card.Spades},
		State{},
		true,
	},
	{
		"In Deck",
		card.Card{card.Joker, card.NoSuit},
		State{Deck: deck.Deck{card.Card{card.Joker, card.NoSuit}, card.Card{card.Ace, card.Spades}}},
		false,
	},
	{
		"In Hand",
		card.Card{card.Five, card.Hearts},
		State{Hands: map[seat.Seat]*hand.Hand{seat.North: &hand.Hand{card.Card{card.Five, card.Hearts}}}},
		false,
	},
	{
		"In Played",
		card.Card{card.Seven, card.Clubs},
		State{Played: []trick.Trick{trick.New(card.Card{card.Seven, card.Clubs})}},
		true,
	},
}

func TestFindCard(t *testing.T) {
	for _, test := range findCardTests {
		if got := test.state.FindCard(test.seek); got == nil && !test.err {
			t.Errorf("%s: unable to find card %s in state %s", test.name, test.seek, test.state)
		} else if got != nil && test.err {
			t.Errorf("%s: found a card when we shouldn't have! (%s in %s)", test.name, test.seek, test.state)
		} else if got != nil && *got != test.seek {
			t.Errorf("%s: found the wrong card! Looking for %s, found %s", test.name, test.seek, *got)
		}
	}
}
