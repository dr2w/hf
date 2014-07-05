package action

import (
	"reflect"
	"testing"

	"dr2w.com/hf/model/card"
	"dr2w.com/hf/model/hand"
	"dr2w.com/hf/model/seat"
	"dr2w.com/hf/model/state"
)

var discardTests = []struct {
	name       string
	inState    state.State
	inMessage  Message
	outState   state.State
	outMessage Message
	err        bool
}{
	{
		"Bad Counting",
		state.State{
			Hands: map[seat.Seat]*hand.Hand{
				seat.North: &hand.Hand{card.Card{}, card.Card{}, card.Card{}},
			},
		},
		Message{
			Seat:    seat.North,
			Options: []int{0, 1},
		},
		state.State{},
		Message{},
		true,
	},
	{
		"End to End",
		state.State{
			Hands: map[seat.Seat]*hand.Hand{
				seat.North: &hand.Hand{
					card.Card{card.Seven, card.Clubs},
					card.Card{card.Eight, card.Clubs},
					card.Card{card.Joker, card.NoSuit},
					card.Card{card.Ace, card.Diamonds},
					card.Card{card.King, card.Clubs},
				},
				seat.East: &hand.Hand{
					card.Card{card.Seven, card.Spades},
					card.Card{card.Eight, card.Spades},
					card.Card{card.Ace, card.Hearts},
					card.Card{card.King, card.Spades},
					card.Card{card.Seven, card.Clubs},
					card.Card{card.Eight, card.Clubs},
					card.Card{card.Joker, card.NoSuit},
					card.Card{card.Ace, card.Diamonds},
					card.Card{card.King, card.Clubs},
				},
			},
		},
		Message{
			Seat:    seat.East,
			Options: []int{1, 3, 5},
		},
		state.State{
			Hands: map[seat.Seat]*hand.Hand{
				seat.North: &hand.Hand{
					card.Card{card.Seven, card.Clubs},
					card.Card{card.Eight, card.Clubs},
					card.Card{card.Joker, card.NoSuit},
					card.Card{card.Ace, card.Diamonds},
					card.Card{card.King, card.Clubs},
				},
				seat.East: &hand.Hand{
					card.Card{card.Seven, card.Spades},
					card.Card{card.Ace, card.Hearts},
					card.Card{card.Seven, card.Clubs},
					card.Card{card.Joker, card.NoSuit},
					card.Card{card.Ace, card.Diamonds},
					card.Card{card.King, card.Clubs},
				},
			},
		},
		Message{
			Type:    Play,
			Seat:    seat.East,
			Options: []int{0, 1, 2, 3, 4, 5},
		},
		false,
	},
}

func TestDiscard(t *testing.T) {
	for _, test := range discardTests {
		s, m, e := discard(test.inState, test.inMessage)
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

var omitTests = []struct {
	name    string
	hand    hand.Hand
	indices []int
	want    hand.Hand
}{
	{
		"Empty Test",
		hand.Hand{},
		[]int{},
		nil,
	},
	{
		"Exclude All",
		hand.Hand{card.Card{card.Ace, card.Hearts}, card.Card{card.Ten, card.Spades}},
		[]int{0, 1},
		nil,
	},
	{
		"Exclude Some",
		hand.Hand{card.Card{card.Ace, card.Hearts}, card.Card{card.Ten, card.Spades}, card.Card{card.Five, card.Hearts}},
		[]int{1},
		hand.Hand{card.Card{card.Ace, card.Hearts}, card.Card{card.Five, card.Hearts}},
	},
}

func TestOmit(t *testing.T) {
	for _, test := range omitTests {
		if got := omit(test.hand, test.indices); !reflect.DeepEqual(*got, test.want) {
			t.Errorf("%s: want %v, got %v", test.name, test.want, *got)
		}
	}
}
