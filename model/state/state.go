// Package state describes all card-related state associated with a game of High Five.
package state

import (
	"bytes"
	"fmt"

	"dr2w.com/hf/model/bid"
	"dr2w.com/hf/model/card"
	"dr2w.com/hf/model/deck"
	"dr2w.com/hf/model/hand"
	"dr2w.com/hf/model/seat"
	"dr2w.com/hf/model/trick"
)

// HandSize determines the number of cards left in the hand after discarding.
const HandSize = 6

// State encapsulates the entire card-related state of a high five game.
type State struct {
	Score  map[seat.Seat]int
	Dealer seat.Seat
	Deck   deck.Deck
	Bids   map[seat.Seat]bid.Bid
	Trump  card.Suit
	Hands  map[seat.Seat]*hand.Hand
	Played []trick.Trick
    Rounds int
}

// WinningBid returns the seat with the winning Bid.
func (s State) WinningBid() (maxSeat seat.Seat, maxBid bid.Bid) {
    for st, b := range s.Bids {
        if b > maxBid {
            maxSeat = st
            maxBid = b
        }
    }
    return maxSeat, maxBid
}

// FindCard finds an instance of the 'seek' card and returns a pointer to it.
// NOTE: This only searches through the Deck and Hands, not the Played pile.
func (s State) FindCard(seek card.Card) *card.Card {
	for i, c := range s.Deck {
		if c == seek {
			return &s.Deck[i]
		}
	}
	for seat, h := range s.Hands {
		for i, c := range *h {
			if c == seek {
				return &((*s.Hands[seat])[i])
			}
		}
	}
	return nil
}

// LastPlayed returns the Trick most recently played, including incomplete
// (current) Tricks.
func (s State) LastPlayed() trick.Trick {
	if len(s.Played) == 0 {
		return trick.Trick{}
	}
	return s.Played[len(s.Played)-1]
}

// ToPlay returns the Seat and Hand of the next Seat to play.
func (s State) ToPlay() (seat.Seat, *hand.Hand) {
	st := s.LastPlayed().NextSeat(s.Trump)
	h, ok := s.Hands[st]
	if !ok {
		return st, nil
	}
	return st, h
}

// String returns a human readable representation of the State.
func (s State) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("\n---\nScore: %v\n", s.Score))
	buffer.WriteString(fmt.Sprintf("Dealer: %v\n", s.Dealer))
	buffer.WriteString(fmt.Sprintf("Deck: %v\n", s.Deck))
	buffer.WriteString(fmt.Sprintf("Bids: %v\n", s.Bids))
    buffer.WriteString(fmt.Sprintf("Trump: %v\n", s.Trump))
	buffer.WriteString("Hands:\n")
	for _, seat := range seat.Order {
		if hand, ok := s.Hands[seat]; ok {
			buffer.WriteString(fmt.Sprintf("\t%s: %s\n", seat, hand.String()))
		}
	}
	buffer.WriteString(fmt.Sprintf("Played: %v\n", s.Played))
	return buffer.String()
}

// Initial returns the starting state for a game of High Five based on the "Dealer" provided.
func Initial(dealer seat.Seat) State {
	return State{
		Deck:   deck.Shuffled(),
		Dealer: dealer,
        Score:  map[seat.Seat]int{seat.North: 0, seat.East: 0, seat.South: 0, seat.West: 0},
	}
}
