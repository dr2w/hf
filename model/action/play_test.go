package action

import (
	"reflect"
	"testing"

	"dr2w.com/hf/model/card"
	"dr2w.com/hf/model/hand"
	"dr2w.com/hf/model/seat"
	"dr2w.com/hf/model/state"
	"dr2w.com/hf/model/trick"
)

var (
	c7d = card.Card{card.Seven, card.Diamonds}
	c9d = card.Card{card.Nine, card.Diamonds}
	c3c = card.Card{card.Three, card.Clubs}
	c9h = card.Card{card.Nine, card.Hearts}
)

var playCardTests = []struct {
	name  string
	state state.State
	seat  seat.Seat
	card  card.Card
	want  []trick.Trick
	err   bool
}{
	{
		"First Play",
		state.State{},
		seat.North,
		c7d,
		[]trick.Trick{trick.New(c7d)},
		false,
	},
	{
		"Second Play",
		state.State{Played: []trick.Trick{trick.New(c7d)}},
		seat.East,
		c3c,
		[]trick.Trick{trick.New(c7d, c3c)},
		false,
	},
	{
		"Next Trick",
		state.State{Played: []trick.Trick{trick.New(c7d, c7d, c7d, c7d)}},
		seat.North,
		c3c,
		[]trick.Trick{trick.New(c7d, c7d, c7d, c7d), trick.New(c3c)},
		false,
	},
	{
		"Already Played",
		state.State{Played: []trick.Trick{trick.New(c7d)}},
		seat.North,
		c7d,
		[]trick.Trick{},
		true,
	},
}

func TestPlayCard(t *testing.T) {
	for _, test := range playCardTests {
		got, err := playCard(test.state, test.seat, test.card)
		if !reflect.DeepEqual(test.want, got) {
			t.Errorf("%s: want %v, got %v", test.name, test.want, got)
		}
		if err != nil && !test.err {
			t.Errorf("%s: got unexpected error (%s)", test.name, err)
		}
		if err == nil && test.err {
			t.Errorf("%s: wanted error, got none", test.name)
		}
	}
}

var validCardsTests = []struct {
	name  string
	hand  *hand.Hand
	trump card.Suit
	trick trick.Trick
	want  []int
}{
	{
		"Empty Hand",
		&hand.Hand{},
		card.Hearts,
		trick.Trick{},
		nil,
	},
	{
		"Empty Trick",
		&hand.Hand{c7d, c9d, c3c, c9h},
		card.Diamonds,
		trick.New(),
		[]int{0, 1, 2, 3},
	},
	{
		"Follow Lead",
		&hand.Hand{c7d, c9d, c3c, c9h},
		card.Hearts,
		trick.New(c3c, c7d),
		[]int{2, 3},
	},
	{
		"Can't Follow Lead",
		&hand.Hand{c7d, c9d, c9h},
		card.Hearts,
		trick.New(c3c, c7d),
		[]int{0, 1, 2},
	},
}

func TestValidCards(t *testing.T) {
	for _, test := range validCardsTests {
		if got := validCards(test.hand, test.trump, test.trick); !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s: want %v, got %v", test.name, test.want, got)
		}
	}
}

var nextTests = []struct {
	name  string
	state state.State
	want  Message
	err   bool
}{
	{
		"Empty State",
		state.State{},
		Message{Type: Score, Seat: seat.None, Options: []int{0}},
		false,
	},
	{
		"Second Play",
		state.State{
			Hands:  map[seat.Seat]*hand.Hand{seat.East: &hand.Hand{c3c, c7d}},
			Played: []trick.Trick{trick.New(c9d)},
			Trump:  card.Diamonds,
		},
		Message{
			Type:    Play,
			Seat:    seat.East,
			Options: []int{1},
		},
		false,
	},
	{
		"Second Trick",
		state.State{
			Hands:  map[seat.Seat]*hand.Hand{seat.South: &hand.Hand{c3c, c7d}},
			Played: []trick.Trick{trick.New(c9d, c7d, c3c, c9h)},
			Trump:  card.Clubs,
		},
		Message{
			Type:    Play,
			Seat:    seat.South,
			Options: []int{0, 1},
		},
		false,
	},
}

func TestNext(t *testing.T) {
	for _, test := range nextTests {
		got, err := next(test.state)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s: want %v, got %v", test.name, test.want, got)
		}
		if err != nil && !test.err {
			t.Errorf("%s: unexpected error: %s", test.name, err)
		}
		if err == nil && test.err {
			t.Errorf("%s: expected error, got none.", test.name)
		}
	}
}

var advanceEmptyHandsTests = []struct {
    name string
    state state.State
    want state.State
    err bool
} {
    {
        "No Empty Hand",
        state.State{
            Hands: map[seat.Seat]*hand.Hand{
                seat.North: &hand.Hand{c7d, c9d},
                seat.East: &hand.Hand{c3c, c9h},
            },
            Played: []trick.Trick{trick.New(c3c)},
            Trump: card.Hearts,
        },
        state.State{
            Hands: map[seat.Seat]*hand.Hand{
                seat.North: &hand.Hand{c7d, c9d},
                seat.East: &hand.Hand{c3c, c9h},
            },
            Played: []trick.Trick{trick.New(c3c)},
            Trump: card.Hearts,
        },
        false,
    },
    {
        "Two Empty Hands",
        state.State{
            Hands: map[seat.Seat]*hand.Hand{
                seat.North: &hand.Hand{c7d, c9d},
                seat.East: &hand.Hand{},
                seat.South: &hand.Hand{},
                seat.West: &hand.Hand{c3c, c9h},
            },
            Played: []trick.Trick{trick.New(c3c)},
            Trump: card.Hearts,
        },
        state.State{
            Hands: map[seat.Seat]*hand.Hand{
                seat.North: &hand.Hand{c7d, c9d},
                seat.East: &hand.Hand{},
                seat.South: &hand.Hand{},
                seat.West: &hand.Hand{c3c, c9h},
            },
            Played: []trick.Trick{trick.New(c3c, card.Card{}, card.Card{})},
            Trump: card.Hearts,
        },
        false,
    },
}

func TestAdvanceEmptyHands(t *testing.T) {
    for _, test := range advanceEmptyHandsTests {
        got, err := advanceEmptyHands(test.state)
        if !reflect.DeepEqual(got, test.want) {
            t.Errorf("%s: want %v, got %v", test.name, test.want, got)
        }
        if err != nil && !test.err {
            t.Errorf("%s: unexpected error - %s", test.name, err)
        }
        if err == nil && test.err {
            t.Errorf("%s: expected error, got none", test.name)
        }
    }
}

var playTests = []struct {
    name string
    inState state.State
    inMsg Message
    outState state.State
    outMsg Message
    err bool
} {
    {
        "Bad Selection",
        state.State{},
        Message{Play, seat.North, []int{0,1,2}},
        state.State{},
        Message{},
        true,
    },
    {
        "Bad Selection Index",
        state.State{},
        Message{Play, seat.North, []int{0}},
        state.State{},
        Message{},
        true,
    },
    {
        "Valid First Play",
        state.State{
            Hands: map[seat.Seat]*hand.Hand{
                seat.North: &hand.Hand{c7d, c3c},
                seat.East: &hand.Hand{c9h, c9d},
            },
            Trump: card.Diamonds,
        },
        Message{Play, seat.North, []int{0}},
        state.State{
            Hands: map[seat.Seat]*hand.Hand{
                seat.North: &hand.Hand{c3c},
                seat.East: &hand.Hand{c9h, c9d},
            },
            Trump: card.Diamonds,
            Played: []trick.Trick{trick.New(c7d)},
        },
        Message{Play, seat.East, []int{1}},
        false,
    },
    {
        "Valid Later Play",
        state.State{
            Hands: map[seat.Seat]*hand.Hand{
                seat.South: &hand.Hand{c7d, c3c},
                seat.West: &hand.Hand{},
                seat.North: &hand.Hand{c9h, c9d},
                seat.East: &hand.Hand{c9h, c7d, c9d},
            },
            Trump: card.Hearts,
            Played: []trick.Trick{
                trick.New(c7d, c7d, c9h, c7d),
                trick.New(c7d, c9h),
            },
        },
        Message{Play, seat.South, []int{1}},
        state.State{
            Hands: map[seat.Seat]*hand.Hand{
                seat.South: &hand.Hand{c7d},
                seat.West: &hand.Hand{},
                seat.North: &hand.Hand{c9h, c9d},
                seat.East: &hand.Hand{c9h, c7d, c9d},
            },
            Trump: card.Hearts,
            Played: []trick.Trick{
                trick.New(c7d, c7d, c9h, c7d),
                trick.New(c7d, c9h, c3c, card.Card{}),
            },
        },
        Message{Play, seat.East, []int{0, 1, 2}},
        false,
    },
}

func TestPlay(t *testing.T) {
    for _, test := range playTests {
        outState, outMsg, err := play(test.inState, test.inMsg)
        if !reflect.DeepEqual(outState, test.outState) {
            t.Errorf("%s: want %v, got %v", test.name, test.outState, outState)
        }
        if !reflect.DeepEqual(outMsg, test.outMsg) {
            t.Errorf("%s: want %v, got %v", test.name, test.outMsg, outMsg)
        }
        if err != nil && !test.err {
            t.Errorf("%s: unexpected error - %s", test.name, err)
        }
        if err == nil && test.err {
            t.Errorf("%s: expected error, got none", test.name)
        }
    }
}
