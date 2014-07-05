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

var (
	c0 = card.Card{0, 0}
	c1 = card.Card{1, 0}
	c2 = card.Card{2, 0}
	c3 = card.Card{3, 0}
	c4 = card.Card{4, 0}
	c5 = card.Card{5, 0}
	c6 = card.Card{6, 0}
	c7 = card.Card{7, 0}
	c8 = card.Card{8, 0}
	c9 = card.Card{9, 0}
)

var dealWithPatternTests = []struct {
	name  string
	state state.State
	cpd   int
	dph   int
	want  state.State
	req   Message
	error bool
}{
	{
		"Bad State",
		state.State{Dealer: seat.North},
		1,
		1,
		state.State{},
		Message{},
		true,
	},
	{
		"Good State - 2 cards at time",
		state.State{
			Dealer: seat.West,
			Deck:   deck.Deck{c0, c1, c2, c3, c4, c5, c6, c7, c8, c9},
		},
		2,
		1,
		state.State{
			Dealer: seat.North,
			Deck:   deck.Deck{c8, c9},
			Hands: map[seat.Seat]*hand.Hand{
				seat.North: &hand.Hand{c6, c7},
				seat.East:  &hand.Hand{c0, c1},
				seat.South: &hand.Hand{c2, c3},
				seat.West:  &hand.Hand{c4, c5},
			},
		},
		Message{
			Type:    Bid,
			Seat:    seat.East,
			Options: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		},
		false,
	},
}

func TestDealWithPattern(t *testing.T) {
	for _, test := range dealWithPatternTests {
		got, req, err := dealWithPattern(test.state, test.cpd, test.dph)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s: want %v, got %v", test.name, test.want, got)
		}
		if !reflect.DeepEqual(req, test.req) {
			t.Errorf("%s: want req %v, got req %v", test.name, test.req, req)
		}
		if err != nil && !test.error {
			t.Errorf("unexpected error for %s: %s", test.name, err)
		}
		if err == nil && test.error {
			t.Errorf("%s: expected error, got none.", test.name)
		}
	}
}
