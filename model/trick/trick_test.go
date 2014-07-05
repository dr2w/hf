package trick

import (
	"testing"

	"dr2w.com/hf/model/card"
	"dr2w.com/hf/model/seat"
)

var (
	c3h  = card.Card{card.Three, card.Hearts}
	c7d  = card.Card{card.Seven, card.Diamonds}
	c9d  = card.Card{card.Nine, card.Diamonds}
	c10c = card.Card{card.Ten, card.Clubs}
)

var winnerTests = []struct {
	name  string
	trick Trick
	trump card.Suit
	seat  seat.Seat
	card  card.Card
}{
	{
		"Empty Trick",
		Trick{},
		card.NoSuit,
		seat.None,
		card.Card{},
	},
	{
		"Single Card",
		New(c7d),
		card.Diamonds,
		seat.North,
		c7d,
	},
	{
		"Trump Lead",
		New(c7d, c9d, c10c),
		card.Diamonds,
		seat.East,
		c9d,
	},
	{
		"Offsuit Lead",
		New(c7d, c9d, c10c),
		card.Clubs,
		seat.South,
		c10c,
	},
	{
		"No Trump",
		New(c7d, c9d, c10c, c3h),
		card.Spades,
		seat.East,
		c9d,
	},
}

func TestWinner(t *testing.T) {
	for _, test := range winnerTests {
		s, c := test.trick.Winner(test.trump)
		if s != test.seat {
			t.Errorf("%s: want seat %v, got %s", test.name, test.seat, s)
		}
		if c != test.card {
			t.Errorf("%s: want card %v, got %s", test.name, test.card, c)
		}
	}
}

var nextSeatTests = []struct {
	name  string
	trick Trick
	trump card.Suit
	want  seat.Seat
}{
	{
		"Empty Trick",
		Trick{},
		card.Hearts,
		seat.None,
	},
	{
		"Single Card",
		New(c7d),
		card.Hearts,
		seat.East,
	},
	{
		"Three Cards",
		New(c7d, c9d, c3h),
		card.Hearts,
		seat.West,
	},
	{
		"Full Trick",
		New(c7d, c9d, c3h, c10c),
		card.Diamonds,
		seat.East,
	},
}

func TestNextSeat(t *testing.T) {
	for _, test := range nextSeatTests {
		if got := test.trick.NextSeat(test.trump); got != test.want {
			t.Errorf("%s: want %v, got %v", test.name, test.want, got)
		}
	}
}
