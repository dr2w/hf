// Package card describes a standard deck of 52 playing cards with the addition of a single "Joker" card.
package card

import (
	"fmt"
)

type Suit int

func (s Suit) String() string {
    return suitNames[s]
}

const (
	NoSuit Suit = iota
	Diamonds
	Clubs
	Hearts
	Spades
)

type Value int

func (v Value) String() string {
    return valueNames[v]
}

const (
	NoValue Value = iota
	Deuce
	Three
	Four
	OffFive
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Joker
	Jack
	Queen
	King
	Ace
)

var (
	Suits        = []Suit{Diamonds, Clubs, Hearts, Spades}
	SuitedValues = []Value{Deuce, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King, Ace}
	Values       = []Value{Deuce, Three, Four, OffFive, Five, Six, Seven, Eight, Nine, Ten, Joker, Jack, Queen, King, Ace}
)

var valueNames = map[Value]string{
	NoValue: "X",
	Deuce:   "2",
	Three:   "3",
	Four:    "4",
	OffFive: "o5",
	Five:    "5",
	Six:     "6",
	Seven:   "7",
	Eight:   "8",
	Nine:    "9",
	Ten:     "10",
	Joker:   "Jo",
	Jack:    "Ja",
	Queen:   "Qu",
	King:    "Ki",
	Ace:     "Ac",
}

var suitNames = map[Suit]string{
	NoSuit:   "X",
	Diamonds: "D",
	Hearts:   "H",
	Spades:   "S",
	Clubs:    "C",
}

type Card struct {
	Value Value
	Suit  Suit
}

// String returns a humna readable string representation of the card.
func (c Card) String() string {
	return fmt.Sprintf("[%s of %s]", c.Value, c.Suit)
}

// Beats returns true if this card outranks the given card given the trump and lead suits.
func (c Card) Beats(o Card, trump Suit, lead Suit) bool {
	if c.Suit == o.Suit {
		return c.Value > o.Value
	}
	return c.Suit == trump || c.Suit == lead && o.Suit != trump
}

// Points returns the number of scoring points this card is worth given the trump suit.
func (c Card) Points(s Suit) int {
	if c.Suit != s {
		return 0
	}
	switch c.Value {
	case Ace, Jack, Joker, Ten, Deuce:
		return 1
	case Five, OffFive:
		return 5
	}
	return 0
}

// Set describes a slice of cards with no particular ordering.
type Set []Card

// PointCards returns the subset of cards with point values > 0 given the specified trump.
func (s Set) PointCards(trump Suit) (p Set) {
	for _, c := range s {
		if c.Points(trump) > 0 {
			p = append(p, c)
		}
	}
	return p
}

// TrumpCards returns the subset of cards which match the trump suit specified.
func (s Set) TrumpCards(trump Suit) (t Set) {
	for _, c := range s {
		if c.Suit == trump {
			t = append(t, c)
		}
	}
	return t
}

// SameColorSuit takes a suit and returns the other suit of the same color.
func SameColorSuit(s Suit) Suit {
	m := map[Suit]Suit{
		NoSuit:   NoSuit,
		Diamonds: Hearts,
		Hearts:   Diamonds,
		Clubs:    Spades,
		Spades:   Clubs,
	}
	return m[s]
}
