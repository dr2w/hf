// Package hand describes a single hand of cards.
package hand

import (
	"fmt"
	"strings"

	"dr2w.com/hf/model/card"
)

// Hand represents a set of cards to be owned and played by a player.
type Hand card.Set

// String returns a human readable string representation of the Hand.
func (h *Hand) String() string {
	var s []string
	for _, card := range *h {
		s = append(s, card.String())
	}
	return strings.Join(s, ",")
}

// Length returns the number of cards in the Hand.
func (h *Hand) Length() int {
	return len(*h)
}

// Add adds the given Card to this Hand.
func (h *Hand) Add(c ...card.Card) {
	*h = append(*h, c...)
}

// Remove removes the Card with the given index from the Hand and
// returns it.
func (h *Hand) Remove(i int) (card.Card, error) {
	if i >= len(*h) {
		return card.Card{}, fmt.Errorf("attempted to access card %d of hand: %v", i, *h)
	}
	c := (*h)[i]
	*h = append((*h)[:i], (*h)[i+1:]...)
	return c, nil
}

// HasSuit returns true iff the Hand has at least one card of the given Suit.
func (h *Hand) HasSuit(s card.Suit) bool {
	for _, c := range *h {
		if c.Suit == s {
			return true
		}
	}
	return false
}
