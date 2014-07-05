package action

import (
	"reflect"
	"testing"

	"dr2w.com/hf/model/card"
	"dr2w.com/hf/model/deck"
	"dr2w.com/hf/model/hand"
	"dr2w.com/hf/model/seat"
	"dr2w.com/hf/model/state"
)

var trumpTests = []struct {
	name       string
	inState    state.State
	inMessage  Message
	outState   state.State
	outMessage Message
	err        bool
}{
	{
		"Empty Test",
		state.State{},
		Message{},
		state.State{},
		Message{},
		true,
	},
	{
		"End to End",
		state.State{
			Deck: deck.Deck{
				card.Card{card.Seven, card.Clubs},
				card.Card{card.Eight, card.Clubs},
				card.Card{card.Joker, card.NoSuit},
				card.Card{card.Ace, card.Diamonds},
				card.Card{card.King, card.Clubs},
			},
			Hands: map[seat.Seat]*hand.Hand{
				seat.North: &hand.Hand{card.Card{card.Ten, card.Clubs}, card.Card{card.Five, card.Diamonds}, card.Card{card.Ten, card.Hearts}},
				seat.East:  &hand.Hand{card.Card{card.Five, card.Clubs}, card.Card{card.Five, card.Spades}},
			},
		},
		Message{
			Seat:    seat.North,
			Options: []int{1}, // Clubs
		},
		state.State{
			Deck: deck.Deck{
				card.Card{card.Seven, card.Clubs},
				card.Card{card.Eight, card.Clubs},
				card.Card{card.Joker, card.Clubs},
				card.Card{card.Ace, card.Diamonds},
				card.Card{card.King, card.Clubs},
			},
			Hands: map[seat.Seat]*hand.Hand{
				seat.North: &hand.Hand{card.Card{card.Ten, card.Clubs}, card.Card{card.Five, card.Diamonds}, card.Card{card.Ten, card.Hearts}},
				seat.East:  &hand.Hand{card.Card{card.Five, card.Clubs}, card.Card{card.OffFive, card.Clubs}},
			},
            Trump: card.Clubs,
		},
		Message{
			Type:    ReDeal,
			Seat:    seat.None,
			Options: []int{0},
		},
		false,
	},
}

func TestTrump(t *testing.T) {
	for _, test := range trumpTests {
		s, m, e := trump(test.inState, test.inMessage)
		if e != nil && !test.err {
			t.Errorf("%s: unexpected error (%s)", test.name, e)
		}
		if e == nil && test.err {
			t.Errorf("%s: expected error.", test.name)
		}
		if s.String() != test.outState.String() {
			t.Errorf("%s: want state %s, got %s.", test.name, test.outState, s)
		}
		if !reflect.DeepEqual(m, test.outMessage) {
			t.Errorf("%s: want message %s, got %s.", test.name, test.outMessage, m)
		}
	}
}
