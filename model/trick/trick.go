// Package trick describes a single trick of seat.Seats mapped to Cards played.
package trick

import (
	"dr2w.com/hf/model/card"
	"dr2w.com/hf/model/seat"
)

// Size defines the number of cards which make up a complete Trick.
const Size = 4

// Trick represents the set of cards played in order against each other in a single round of play.
type Trick struct {
	Cards map[seat.Seat]card.Card
	First seat.Seat
}

// SuitLead returns the suit that was lead in this trick.
func (t Trick) SuitLead() card.Suit {
	return t.Cards[t.First].Suit
}

// Points returns the total number of points contained within the trick.
func (t Trick) Points(trump card.Suit) (p int) {
    for _, c := range t.Cards {
        p += c.Points(trump)
    }
    return p
}

// Winner returns the winner of the specified Trick and the winning Card.
func (t Trick) Winner(trump card.Suit) (seat.Seat, card.Card) {
	var (
		maxCard card.Card
		maxSeat seat.Seat
	)
	for s, c := range t.Cards {
		if c.Beats(maxCard, trump, t.SuitLead()) {
			maxCard = c
			maxSeat = s
		}
	}
	return maxSeat, maxCard
}

// NextSeat takes the most recently played Trick and returns the Seat that should play next.
func (t Trick) NextSeat(trump card.Suit) seat.Seat {
	if len(t.Cards) >= Size {
		seat, _ := t.Winner(trump)
		return seat
	}
	for seat, _ := range t.Cards {
		if _, ok := t.Cards[seat.Next()]; !ok {
			return seat.Next()
		}
	}
	return seat.None
}

// Full returns true iff all cards for this Trick have been played.
func (t Trick) Full() bool {
	return len(t.Cards) == Size
}

// Empty returns true iff no cards have yet been played on this Trick.
func (t Trick) Empty() bool {
	return len(t.Cards) == 0
}

// New is a convenience functio which constructs a trick from the first four given Cards.
func New(cards ...card.Card) Trick {
	t := Trick{First: seat.North, Cards: map[seat.Seat]card.Card{}}
	for i, c := range cards {
		if i < len(seat.Order) {
			t.Cards[seat.Order[i]] = c
		}
	}
	return t
}
