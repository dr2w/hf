package action

import (
	"reflect"
	"testing"

	"dr2w.com/hf/model/bid"
	"dr2w.com/hf/model/seat"
	"dr2w.com/hf/model/state"
)

var bidsTests = []struct {
	name       string
	inState    state.State
	inMessage  Message
	outState   state.State
	outMessage Message
	err        bool
}{
	{
		"Selection out of range",
		state.State{},
		Message{Options: []int{100}},
		state.State{},
		Message{},
		true,
	},
	{
		"Pass => Six",
		state.State{Bids: map[seat.Seat]bid.Bid{seat.North: bid.Pass}},
		Message{Seat: seat.East, Options: []int{2}}, // B7
		state.State{Bids: map[seat.Seat]bid.Bid{seat.North: bid.Pass, seat.East: bid.B7}},
		Message{Type: Bid, Seat: seat.South, Options: []int{0, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}},
		false,
	},
	{
		"Pass, Eight => Nine",
		state.State{Bids: map[seat.Seat]bid.Bid{seat.South: bid.Pass, seat.West: bid.B8}},
		Message{Seat: seat.North, Options: []int{4}}, // B9
		state.State{Bids: map[seat.Seat]bid.Bid{seat.South: bid.Pass, seat.West: bid.B8, seat.North: bid.B9}},
		Message{Type: Bid, Seat: seat.East, Options: []int{0, 5, 6, 7, 8, 9, 10, 11, 12}},
		false,
	},
	{
		"Pass, Eight, Pass => Ten",
		state.State{Bids: map[seat.Seat]bid.Bid{seat.South: bid.Pass, seat.West: bid.B8, seat.East: bid.Pass}},
		Message{Seat: seat.North, Options: []int{5}}, // B10
		state.State{Bids: map[seat.Seat]bid.Bid{seat.South: bid.Pass, seat.West: bid.B8, seat.East: bid.Pass, seat.North: bid.B10}},
		Message{Type: Trump, Seat: seat.North, Options: []int{0, 1, 2, 3}},
		false,
	},
}

func TestBids(t *testing.T) {
	for _, test := range bidsTests {
		state, outMessage, err := bids(test.inState, test.inMessage)
		if err != nil && !test.err {
			t.Errorf("%s: unexpected error (%s)", test.name, err)
		}
		if err == nil && test.err {
			t.Errorf("%s: expected error, got none", test.name)
		}
		if !reflect.DeepEqual(state, test.outState) {
			t.Errorf("%s: want state \n%s\n,got:\n%s\n", test.name, test.outState, state)
		}
		if !reflect.DeepEqual(outMessage, test.outMessage) {
			t.Errorf("%s: want outMessage \n%s\ngot:\n%s\n", test.name, test.outMessage, outMessage)
		}
	}
}

var bidWinnerTests = []struct {
	name string
	m    map[seat.Seat]bid.Bid
	seat seat.Seat
	bid  bid.Bid
}{
	{
		"Empty Map",
		make(map[seat.Seat]bid.Bid),
		0,
		0,
	},
	{
		"One Bid",
		map[seat.Seat]bid.Bid{seat.North: bid.B6},
		seat.North,
		bid.B6,
	},
	{
		"Multiple Bids",
		map[seat.Seat]bid.Bid{
			seat.North: bid.Pass,
			seat.East:  bid.B10,
			seat.South: bid.Pass,
		},
		seat.East,
		bid.B10,
	},
}

func TestBidWinner(t *testing.T) {
	for _, test := range bidWinnerTests {
		seat, bid := bidWinner(test.m)
		if seat != test.seat {
			t.Errorf("%s: want %s for seat, got %s.", test.name, test.seat, seat)
		}
		if bid != test.bid {
			t.Errorf("%s: want %s for bid, got %s.", test.name, test.bid, bid)
		}
	}
}

var reqForNextBidTests = []struct {
	name  string
	state state.State
	want  Message
	err   bool
}{
	{
		"Empty State",
		state.State{},
		Message{},
		true,
	},
	{
		"One Pass",
		state.State{Bids: map[seat.Seat]bid.Bid{seat.North: bid.Pass}},
		Message{
			Type:    Bid,
			Seat:    seat.East,
			Options: SelectionRange(0, len(bid.Values)),
		},
		false,
	},
	{
		"Six Seven Pass",
		state.State{Bids: map[seat.Seat]bid.Bid{seat.East: bid.B6, seat.South: bid.B7, seat.West: bid.Pass}},
		Message{
			Type:    Bid,
			Seat:    seat.North,
			Options: append([]int{0}, SelectionRange(int(bid.B8), len(bid.Values))...),
		},
		false,
	},
	{
		"Everyone bid already",
		state.State{Bids: map[seat.Seat]bid.Bid{seat.North: bid.Pass, seat.East: bid.Pass, seat.South: bid.Pass, seat.West: bid.Pass}},
		Message{},
		true,
	},
}

func TestReqForNextBid(t *testing.T) {
	for _, test := range reqForNextBidTests {
		got, err := reqForNextBid(test.state)
		if err != nil && !test.err {
			t.Errorf("%s: unexpected error (%s)", test.name, err)
		}
		if err == nil && test.err {
			t.Errorf("%s: expected error got none.", test.name)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s: got %s, want %s", test.name, got, test.want)
		}
	}
}
