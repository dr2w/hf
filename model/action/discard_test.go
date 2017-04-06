package action

import (
	"reflect"
	"testing"

	"dr2w.com/hf/model/bid"
	"dr2w.com/hf/model/card"
	"dr2w.com/hf/model/hand"
	"dr2w.com/hf/model/seat"
	"dr2w.com/hf/model/state"
)

var nextMessageTests = []struct {
	name    string
	state   state.State
	seat    seat.Seat
	want    Message
}{
	{
		"Winner",
		state.State{
			Bids: map[seat.Seat]bid.Bid{
				seat.North: bid.B8,
			},
			Hands: map[seat.Seat]*hand.Hand{
				seat.North: &hand.Hand{card.Card{}},
			},
		},
		seat.North,
		Message{
			Type: Play,
			Seat: seat.North,
			Options: []int{0},
			Expect: 1,
		},
	},
	{
		"NextIsWinner",
		state.State{
			Bids: map[seat.Seat]bid.Bid{
				seat.East: bid.B8,
			},
			Dealer: seat.South,
		},
		seat.North,
		Message{
			Type: ReDeal,
			Seat: seat.South,
			Options: []int{0},
			Expect: 1,
		},
	},
	{
		"NoWinnerYet",
		state.State{
			Bids: map[seat.Seat]bid.Bid{
				seat.South: bid.B8,
			},
			Trump: card.Diamonds,
			Hands: map[seat.Seat]*hand.Hand{
				seat.East: &hand.Hand{
					card.Card{card.Ace, card.Clubs},
					card.Card{card.King, card.Diamonds},
					card.Card{card.Queen, card.Hearts},
					card.Card{card.Jack, card.Diamonds},
					card.Card{card.Ten, card.Clubs},
				},
			},
		},
		seat.North,
		Message{
			Type: Discard,
			Seat: seat.East,
			Options: []int{0,2,4},
			Expect: 3,
		},
	},
}
			
			

func TestNextMessage(t *testing.T) {
	for _, test := range nextMessageTests {
		got := nextMessage(test.state, test.seat)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s: want message %s, got %s.", test.name, test.want, got)
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
