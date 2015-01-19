package bidding

import "testing"
import "dr2w.com/hf/model/action"
import "dr2w.com/hf/model/bid"
import "dr2w.com/hf/model/card"
import "dr2w.com/hf/model/hand"
import "dr2w.com/hf/model/seat"
import "dr2w.com/hf/model/state"

var drwSixBidTests = []struct {
	name string
	hand card.Set
	want bool
}{
	{"No Cards", card.Set{}, false},
	{"No Coverage", card.WorstHand, false},
	{"Nothing Good 1", card.BadHand1, false},
	{"Nothing Good 2", card.OkHand1, false},
	{"One good suit", card.GoodHand1, false},
	{"Almost", card.WeakSixHand, false},
	{"Crazy Good", card.BestSixHand, true},
	{"Solid", card.SolidSixHand, true},
}

func TestDRWSixBid(t *testing.T) {
	for _, test := range drwSixBidTests {
		if got := drwSixBid(test.hand); got != test.want {
			t.Errorf("%s: want %v, got %v", test.name, test.want, got)
		}
	}
}

var drwSuitBidTests = []struct {
	name    string
	hand    card.Set
	wantMin bid.Bid
	wantMax bid.Bid
}{
	{
		"Nuttin",
		card.Set{},
		bid.Pass,
		bid.Pass,
	},
	{
		"WorstHand",
		card.WorstHand,
		bid.Pass,
		bid.Pass,
	},
	{
		"BestHand",
		card.BestHand,
		bid.B1530,
		bid.B1530,
	},
	{
		"SixBid",
		card.WeakSixHand,
		bid.Pass,
		bid.Pass,
	},
	{
		"BadHand1",
		card.BadHand1,
		bid.Pass,
		bid.Pass,
	},
	{
		"GreatHand1",
		card.GreatHand1,
		bid.B10,
		bid.B1428,
	},
	{
		"GoodHand1",
		card.GoodHand1,
		bid.B8,
		bid.B9,
	},
	{
		"GoodHand2",
		card.GoodHand2,
		bid.B9,
		bid.B10,
	},
	{
		"GoodHand3",
		card.GoodHand3,
		bid.B8,
		bid.B9,
	},
	{
		"GoodHand4",
		card.GoodHand4,
		bid.B8,
		bid.B9,
	},
}

func TestDRWSuitBid(t *testing.T) {
	for _, test := range drwSuitBidTests {
		min, max := bid.Pass, bid.Pass
		for _, suit := range card.Suits {
			cards := test.hand.TrumpCards(suit).AsTrump(suit)
			minB, maxB := drwSuitBid(cards)
			if minB > min {
				min = minB
				max = maxB
			}
		}
		if min != test.wantMin || max != test.wantMax {
			t.Errorf("%s: want %s,%s - got %s,%s", test.name,
				test.wantMin, test.wantMax, min, max)
		}
	}
}

var fromBiddersTests = []struct {
	name    string
	forSuit func(cards card.Set) (min, max bid.Bid)
	isSix   func(cards card.Set) bool
	state   state.State
	message action.Message
	want    []int
}{
	{
		"Basic Suit",
		func(cards card.Set) (min, max bid.Bid) {
			return bid.B1530, bid.B1530
		},
		func(cards card.Set) bool {
			return false
		},
		state.State{
			Hands: map[seat.Seat]*hand.Hand{
				seat.North: &hand.Hand{card.Card{card.Ace, card.Spades}},
			},
		},
		action.Message{
			Seat: seat.North,
		},
		[]int{int(bid.B1530)},
	},
	{
		"Basic Six",
		func(cards card.Set) (min, max bid.Bid) {
			return bid.Pass, bid.B1530
		},
		func(cards card.Set) bool {
			return true
		},
		state.State{
			Hands: map[seat.Seat]*hand.Hand{
				seat.North: &hand.Hand{card.Card{card.Ace, card.Spades}},
			},
		},
		action.Message{
			Seat: seat.North,
		},
		[]int{int(bid.B6)},
	},
}

func TestFromBidders(t *testing.T) {
	for _, test := range fromBiddersTests {
		f := FromBidders(test.forSuit, test.isSix)
		if got := f(test.state, test.message); got[0] != test.want[0] {
			t.Errorf("%s: want %s, got %s", test.name, test.want[0], got[0])
		}
	}
}
