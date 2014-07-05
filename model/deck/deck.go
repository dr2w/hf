// Package deck implements a deck of 52 cards and one Joker.
package deck

import (
	"fmt"
	"math/rand"
	"strings"

	"dr2w.com/hf/model/card"
)

type Deck card.Set

// String returns a string representation of the deck - a single line of comma-separated cards.
func (d Deck) String() string {
	s := make([]string, len(d))
	for i, c := range d {
		s[i] = c.String()
	}
	return strings.Join(s, ",")
}

// Shuffle takes a deck and reorders the card in a semi-random fashion. The integer parameter affects
// how thorough the shuffling is.
func (d Deck) Shuffle(n int) {
	size := len(d)
	if size == 0 {
		return
	}
	for i := 0; i < n; i++ {
		a := rand.Intn(size)
		b := rand.Intn(size)
		d[a], d[b] = d[b], d[a]
	}
}

// Deal "pops" the given number of cards off the top of the deck and returns them to the caller.
// Returns an error if too many cards are requested.
func (d *Deck) Deal(n int) (card.Set, error) {
    if n == 0 {
        return card.Set([]card.Card{}), nil
    }
	if n > len(*d) || n < 0 {
		return card.Set{}, fmt.Errorf("Cannot deal %d cards from a deck with only %d cards.", n, len(*d))
	}
	ret := (*d)[:n]
	*d = (*d)[n:]
	return card.Set(ret), nil
}

// New returns a fully formed deck of cards with a determinstic sorted order.
func New() Deck {
	d := make(Deck, len(card.Suits)*len(card.SuitedValues)+1)
	for s, suit := range card.Suits {
		for v, value := range card.SuitedValues {
			d[v+s*len(card.SuitedValues)] = card.Card{value, suit}
		}
	}
	d[len(card.Suits)*len(card.SuitedValues)] = card.Card{card.Joker, card.NoSuit}
	return d
}

// Shuffled returns a new Deck shuffled once for every card in the deck.
func Shuffled() Deck {
	d := New()
	d.Shuffle(len(d))
	return d
}
